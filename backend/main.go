
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2"

	// "github.com/joho/godotenv"

	"github.com/google/generative-ai-go/genai"
	"github.com/zmb3/spotify/v2"
	"google.golang.org/api/option"
)

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("Error loading .env filee")
// 	}
// }

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "https://racial-angela-jyotiprakashh-5b571460.koyeb.app/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistModifyPrivate))
	ch    = make(chan *spotify.Client)
	state = "abc123"
	clientMutex = &sync.Mutex{}
	clientStore *spotify.Client
	// These should be randomly generated for each request
	//  More information on generating these can be found here,
	// https://www.oauth.com/playground/authorization-code-with-pkce.html
	codeVerifier  = "w0HfYrKnG8AihqYHA9_XUPTIcqEXQvCQfOF2IitRgmlF43YWJ8dy2b49ZUwVUOR.YnvzVoTBL57BwIhM4ouSa~tdf0eE_OmiMC_ESCcVOe7maSLIk9IOdBhRstAxjCl7"
	codeChallenge = "ZhZJzPQXYBMjH8FlGAdYK5AndohLzFfZT-8J7biT7ig"
)

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
        if origin == "https://spotify-playlist-maker-olive.vercel.app" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
        } else {
            w.Header().Set("Access-Control-Allow-Origin", "*") // Change this later for better security
        }
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

        if r.Method == "OPTIONS" {
            w.WriteHeader(204)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/api/generate-playlist", generatePlaylistHandler)
	http.HandleFunc("/api/getUserID", getUserIDHandler)

	http.HandleFunc("/login", loginhandler)
	http.HandleFunc("/api/addPlaylist", addPlaylistHandler)
	http.HandleFunc("/api/user-profile", userProfileHandler)
	http.HandleFunc("/logout", logoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Starting server at port", port)

	handler := CORSMiddleware(http.DefaultServeMux)

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, handler))
	}()

	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	select {}

}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	clientStore = nil
	clientMutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "spotify_auth",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	fmt.Println("Logged out")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("client_id", os.Getenv("SPOTIFY_ID")),
		oauth2.SetAuthURLParam("show_dialog", "true"),
	), http.StatusTemporaryRedirect)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier), oauth2.SetAuthURLParam("client_id", os.Getenv("SPOTIFY_ID")), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))

	clientMutex.Lock()
	clientStore = client
	clientMutex.Unlock()
	fmt.Println(w, "Login Completed!")
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are now logged in as:", user.ID)
	http.Redirect(w, r, "http://localhost:5173/app", http.StatusSeeOther)

	ch <- client
}

func getUserIDHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	client := clientStore
	clientMutex.Unlock()
	if client == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"user_id": user.ID})
}

func generatePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	client := clientStore
	clientMutex.Unlock()

	if client == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var moodReq MoodRequest
	if err := json.NewDecoder(r.Body).Decode(&moodReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Println(moodReq)
	playlist := generatePlaylist(client, moodReq.Mood)

	resp := PlaylistResponse{Playlist: playlist}
	for i, track := range playlist {
		fmt.Println(i, track)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type MoodRequest struct {
	Mood string `json:"mood"`
}
type TrackInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type PlaylistResponse struct {
	Playlist    []TrackInfo `json:"playlist"`
	PlaylistURL string      `json:"playlist_url"`
}

func analyzeMood(mood string) map[string]float32 {
	characteristics := map[string]float32{
		"energy": 0.5,
		"tempo":  100,
	}

	switch mood {
	case "Happy":
		characteristics["energy"] = 0.8
		characteristics["tempo"] = 120
	case "Sad":
		characteristics["energy"] = 0.2
		characteristics["tempo"] = 60
	case "Relaxed":
		characteristics["energy"] = 0.4
		characteristics["tempo"] = 80
	case "Energetic":
		characteristics["energy"] = 0.9
		characteristics["tempo"] = 140
	case "Romantic":
		characteristics["energy"] = 0.3
		characteristics["tempo"] = 70
	case "Celebratory":
		characteristics["energy"] = 0.9
		characteristics["tempo"] = 130
	case "Focused":
		characteristics["energy"] = 0.5
		characteristics["tempo"] = 100
	case "Nostalgic":
		characteristics["energy"] = 0.4
		characteristics["tempo"] = 90
	default:
		log.Printf("Mood not recognized: %v. Using default characteristics.", mood)
	}

	return characteristics
}

func moodAnalysertoWord(mood string) string {
	ctx := context.Background()
	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	Gclient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer Gclient.Close()

	model := Gclient.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(2)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(1024)
	model.ResponseMIMEType = "text/plain"

	prompt := "You are an advanced sentiment analysis tool designed to accurately determine the mood expressed in a given text. Your expertise lies in distilling complex emotional states into a single word among happy, sad, relaxed, energetic, celebratory, focused, and nostalgic. Your task is to analyze a paragraph of text provided by the user. Here is the text you will be working with Paragraph:" + mood + " Keep in mind the nuances of language and context to make the most accurate mood assessment. Consider the choice of words, the tone of voice, and any underlying sentiments that may not be immediately apparent but contribute to the overall mood. Please present your response as a single word that best represents the mood exemplified in the paragraph above, ensuring clarity and precision in your analysis."

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}
	formattedContent := formatResponse(resp)
	fmt.Println("Formatted content: ", formattedContent)
	return formattedContent
}

func formatResponse(resp *genai.GenerateContentResponse) string {
	var formattedContent strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}
	return formattedContent.String()
}

func generatePlaylist(client *spotify.Client, mood string) []TrackInfo {
	// topTracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Limit(40))
	// if err != nil {
	// 	log.Printf("Error fetching top tracks: %v", err)
	// 	return []TrackInfo{}
	// }

	// // Placeholder implementation
	// var results []TrackInfo

	// // Fetch track IDs for top tracks
	// var trackIDs []spotify.ID
	// for _, track := range topTracks.Tracks {
	// 	trackIDs = append(trackIDs, track.ID)
	// }

	// // Fetch audio features for these tracks
	// audioFeatures, err := client.GetAudioFeatures(context.Background(), trackIDs...)
	// if err != nil {
	// 	log.Printf("Error fetching audio features: %v", err)
	// 	return []TrackInfo{}
	// }

	// Mood characteristics based on user input
	mood = moodAnalysertoWord(mood)
	currMood = mood
	characteristics := analyzeMood(mood)
	var results []TrackInfo
	searchQuery := "top " + mood + " songs"
	searchResult, err := client.Search(context.Background(), searchQuery, spotify.SearchTypeTrack)
	if err != nil {
		log.Printf("Error searching for tracks based on mood: %v", err)
		return []TrackInfo{}
	}

	var trackIDs []spotify.ID
	for _, track := range searchResult.Tracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	audioFeatures, err := client.GetAudioFeatures(context.Background(), trackIDs...)
	if err != nil {
		log.Printf("Error fetching audio features: %v", err)
		return []TrackInfo{}
	}

	for _, track := range searchResult.Tracks.Tracks {
		// Find corresponding audio feature for the track
		var feature *spotify.AudioFeatures
		for _, af := range audioFeatures {
			if af.ID == track.ID {
				feature = af
				break
			}
		}

		if feature == nil {
			continue
		}

		// Placeholder logic to filter tracks based on mood characteristics
		if characteristics["energy"] <= feature.Energy && feature.Danceability > 0.5 && feature.Tempo > characteristics["tempo"] {
			results = append(results, TrackInfo{
				Name:  track.Name,
				Image: track.Album.Images[0].URL,
				URL:   track.ExternalURLs["spotify"],
			})
		} else if characteristics["energy"] > feature.Energy && feature.Danceability <= 0.5 && feature.Tempo <= characteristics["tempo"] {
			results = append(results, TrackInfo{
				Name:  track.Name,
				Image: track.Album.Images[0].URL,
				URL:   track.ExternalURLs["spotify"],
			})
		}
	}

	return results
}

var currMood string

func addPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	client := clientStore
	clientMutex.Unlock()

	if client == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	var req struct {
		Playlist []string `json:"playlist"`
		Name     string   `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	req.Name = currMood

	var trackIDs []spotify.ID
	for _, trackName := range req.Playlist {
		result, err := client.Search(context.Background(), trackName, spotify.SearchTypeTrack)
		if err != nil {
			log.Printf("Error searching for track %s: %v", trackName, err)
			http.Error(w, "Error searching for tracks", http.StatusInternalServerError)
			return
		}
		if len(result.Tracks.Tracks) > 0 {
			trackIDs = append(trackIDs, result.Tracks.Tracks[0].ID)
		}
	}
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	playlistURL, err := createAndAddTracksToPlaylist(client, user.ID, req.Name, trackIDs)
	if err != nil {
		log.Printf("Error creating playlist: %v", err)
		http.Error(w, "Error creating playlist", http.StatusInternalServerError)
		return
	}

	resp := PlaylistResponse{PlaylistURL: playlistURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func createAndAddTracksToPlaylist(client *spotify.Client, userID, playlistName string, trackIDs []spotify.ID) (string, error) {
	// Create a new playlist for the user
	playlist, err := client.CreatePlaylistForUser(context.Background(), userID, "My "+playlistName + " Playlist", "Playlist generated based on your mood", false, false)
	if err != nil {
		return "", err
	}

	// Add tracks to the playlist
	_, err = client.AddTracksToPlaylist(context.Background(), playlist.ID, trackIDs...)
	if err != nil {
		return "", err
	}

	return playlist.ExternalURLs["spotify"], nil
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	client := clientStore
	clientMutex.Unlock()

	if client == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	profile := struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		ProfileURL  string `json:"profile_url"`
		Images      string `json:"images"`
	}{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		ProfileURL:  user.ExternalURLs["spotify"],
		Images:      user.Images[0].URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
