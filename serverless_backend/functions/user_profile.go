package functions

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/zmb3/spotify/v2"
)

var (
	clientMutex = &sync.Mutex{}
	clientStore *spotify.Client
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	client := clientStore
	clientMutex.Unlock()

	if client == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
