package manager

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/heatxsink/go-hue/groups"
	"github.com/heatxsink/go-hue/lights"
	"github.com/heatxsink/go-hue/portal"

	"github.com/zmb3/spotify"
)

type Scene struct {
	soundtrackSpotifyID string
	hueLightColor       uint16
}

const (
	Red    = 65280
	Green  = 25500
	Blue   = 46920
	Yellow = 12750
	Purple = 46920
)

const (
	redirectURI = "http://localhost:8080/callback"
)

var (
	ch          = make(chan *spotify.Client)
	auth        = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserModifyPlaybackState)
	state       = "abc123"
	local       = len(os.Args) > 1 && os.Args[1] == "true"
	testing     = len(os.Args) > 2 && os.Args[2] == "true"
	hueBridgeIP = "192.168.1.2"
)

func main() {

	var setups = make(map[string]Scene)

	setups["forest"] = Scene{
		soundtrackSpotifyID: "spotify:track:2NMA93HoL1DkxNLbnyOOqA",
		hueLightColor:       Green,
	}
	setups["desert"] = Scene{
		soundtrackSpotifyID: "spotify:track:6A4JieXGIddpzoC56Vcp3N",
		hueLightColor:       Yellow,
	}
	setups["cave"] = Scene{
		soundtrackSpotifyID: "spotify:track:6wUQ3jIw4uM7sK8YxbV7iN",
		hueLightColor:       Purple,
	}

	auth.SetAuthInfo(clientID, clientSecret)
	userAuth()
	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		panic(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	sceneKey := os.Args[1]

	newScene := setups[sceneKey]

	setLights(newScene.hueLightColor)
	playSong(newScene.soundtrackSpotifyID)

}

func playSong(songID string) {

	client := <-ch

	songURI := spotify.URI(songID)

	err := client.PlayOpt(&spotify.PlayOptions{
		URIs: []spotify.URI{songURI},
	})

	if err != nil {
		panic(err)
	}

}

func setLights(color uint16) {
	pp, err := portal.GetPortal()
	if err != nil {
		panic(err)
	}

	gg := groups.New(pp[0].InternalIPAddress, hueAPIKey)

	livingRoom, err := gg.GetGroup(1)

	newLightsState := lights.State{
		On:  true,
		Hue: color,
	}

	res, err := gg.SetGroupState(livingRoom.ID, newLightsState)

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(err)

	fmt.Println(livingRoom.Name)

}

func userAuth() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
	ch <- &client
}
