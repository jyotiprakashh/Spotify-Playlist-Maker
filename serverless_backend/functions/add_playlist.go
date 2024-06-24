package functions

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	// "sync"

	"github.com/zmb3/spotify/v2"
)

// var (
// 	clientMutex = &sync.Mutex{}
// 	clientStore *spotify.Client
// )

func AddPlaylistHandler(w http.ResponseWriter, r *http.Request) {
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
	playlist, err := client.CreatePlaylistForUser(context.Background(), userID, "My "+playlistName, "Playlist generated based on your mood", false, false)
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
