package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "sync"

	"github.com/zmb3/spotify/v2"
)

// var (
// 	clientMutex = &sync.Mutex{}
// 	clientStore *spotify.Client
// )

func GeneratePlaylistHandler(w http.ResponseWriter, r *http.Request) {
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
	case "happy", "joyful", "excited":
		characteristics["energy"] = 0.8
		characteristics["tempo"] = 120
	case "sad", "bad", "blue":
		characteristics["energy"] = 0.2
		characteristics["tempo"] = 60
	case "relaxed", "calm", "chill":
		characteristics["energy"] = 0.4
		characteristics["tempo"] = 80
	case "energetic", "upbeat", "motivated":
		characteristics["energy"] = 0.9
		characteristics["tempo"] = 140
	case "angry", "aggressive", "furious":
		characteristics["energy"] = 0.7
		characteristics["tempo"] = 150
	case "romantic", "love", "affectionate":
		characteristics["energy"] = 0.3
		characteristics["tempo"] = 70
	case "party", "celebratory", "festive":
		characteristics["energy"] = 0.9
		characteristics["tempo"] = 130
	case "focused", "concentrated", "attentive":
		characteristics["energy"] = 0.5
		characteristics["tempo"] = 100
	case "nostalgic", "sentimental", "reflective":
		characteristics["energy"] = 0.4
		characteristics["tempo"] = 90
	case "melancholic", "down", "depressed":
		characteristics["energy"] = 0.2
		characteristics["tempo"] = 65
	default:
		log.Printf("Mood not recognized: %v. Using default characteristics.", mood)
	}

	return characteristics
}

func generatePlaylist(client *spotify.Client, mood string) []TrackInfo {
	topTracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Limit(40))
	if err != nil {
		log.Printf("Error fetching top tracks: %v", err)
		return []TrackInfo{}
	}

	var results []TrackInfo

	var trackIDs []spotify.ID
	for _, track := range topTracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	trackAudioFeatures, err := client.GetAudioFeatures(context.Background(), trackIDs...)
	if err != nil {
		log.Printf("Error fetching track audio features: %v", err)
		return []TrackInfo{}
	}

	moodCharacteristics := analyzeMood(mood)

	var playlist []spotify.FullTrack

	for i, features := range trackAudioFeatures {
		if features.Energy >= moodCharacteristics["energy"]-0.1 && features.Energy <= moodCharacteristics["energy"]+0.1 &&
			features.Tempo >= moodCharacteristics["tempo"]-10 && features.Tempo <= moodCharacteristics["tempo"]+10 {
			playlist = append(playlist, topTracks.Tracks[i])
		}
	}

	for _, track := range playlist {
		results = append(results, TrackInfo{Name: track.Name, Image: track.Album.Images[0].URL})
	}

	return results
}
