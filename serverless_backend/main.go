package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jyotiprakashh/Spotify-Playlist-Maker/serverless_backend/functions"
)

func main() {
	http.HandleFunc("/generate-playlist", functions.GeneratePlaylistHandler)
	http.HandleFunc("/add-playlist", functions.AddPlaylistHandler)
	http.HandleFunc("/user-profile", functions.UserProfileHandler)
	http.HandleFunc("/login", functions.LoginHandler)
	http.HandleFunc("/callback", functions.CompleteAuth)
	http.HandleFunc("/logout", functions.LogoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
