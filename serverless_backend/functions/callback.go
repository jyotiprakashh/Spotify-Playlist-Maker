package functions

import (
	// "context"
	"fmt"
	"log"
	"net/http"
	"os"
	// "sync"

	"golang.org/x/oauth2"
	"github.com/zmb3/spotify/v2"
	// spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const RedirectURI = "http://localhost:8080/callback"

var (
	
	ch    = make(chan *spotify.Client)
	state = "abc123"
	// These should be randomly generated for each request
	//  More information on generating these can be found here,
	// https://www.oauth.com/playground/authorization-code-with-pkce.html
	codeVerifier  = "w0HfYrKnG8AihqYHA9_XUPTIcqEXQvCQfOF2IitRgmlF43YWJ8dy2b49ZUwVUOR.YnvzVoTBL57BwIhM4ouSa~tdf0eE_OmiMC_ESCcVOe7maSLIk9IOdBhRstAxjCl7"
	// codeChallenge = "ZhZJzPQXYBMjH8FlGAdYK5AndohLzFfZT-8J7biT7ig"
)

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "http://localhost:5173/app", http.StatusSeeOther)
	ch <- client
}
