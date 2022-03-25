package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

func home(w http.ResponseWriter, r *http.Request) {
	var tmplt = template.Must(template.ParseFiles("templates/home.html"))
	tmplt.Execute(w, nil)
}

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/home", home)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	/*http.ListenAndServe(":8080", nil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}*/
	client, err := getSpotifyClient()
	log.Println("successfully created Spotify client!")

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	track, err := client.CurrentUsersTracks(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tracks := []spotify.SimpleTrack{}
	playlistTracks := track.Tracks
	for _, playlistTrack := range playlistTracks {
		tracks = append(tracks, playlistTrack.SimpleTrack)

	}

	playlistId := createPlaylist(context.Background(), client, user, tracks)
	playlist := getPlaylist(context.Background(), client, playlistId)
	fmt.Println(playlist.Name)

	fmt.Println("You are logged in as:", user.ID)
	select {}
}
