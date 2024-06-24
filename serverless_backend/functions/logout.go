package functions

import (
	"net/http"
	// "sync"
)

// var (
// 	clientMutex = &sync.Mutex{}
// 	clientStore *spotify.Client
// )

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clientMutex.Lock()
	clientStore = nil
	clientMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
