package functions

import (
	"net/http"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	redirectURI = os.Getenv("SPOTIFY_REDIRECT_URI")
	auth        = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistModifyPrivate))
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
