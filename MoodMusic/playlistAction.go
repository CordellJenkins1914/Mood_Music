package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

func createPlaylist(ctx context.Context, client *spotify.Client, user *spotify.PrivateUser, tracks []spotify.SimpleTrack) spotify.ID {
	fmt.Println("Beginning createPlaylist")

	playlistName := "Success"
	playlistDescription := "Created automatically"
	createPublicPlaylist := false
	collaborativePlaylist := false
	userId := user.ID

	createdPlaylist, err := client.CreatePlaylistForUser(ctx, userId, playlistName, playlistDescription, createPublicPlaylist, collaborativePlaylist)
	if err != nil {
		log.Fatal(err)
	}
	playlistId := createdPlaylist.ID
	//Get track IDs from list of tracks
	var trackIds []spotify.ID
	for _, track := range tracks {
		attributes, err := client.GetAudioFeatures(ctx, track.ID)
		if err != nil {
			log.Fatal(err)
		}

		mood := findMood(context.Background(), client, attributes)
		if mood == "angry" {
			trackIds = append(trackIds, track.ID)
		}
		fmt.Println(track.Name)
		fmt.Println()
	}
	_, err = client.AddTracksToPlaylist(ctx, playlistId, trackIds...)
	if err != nil {
		log.Fatal(err)
	}
	return playlistId

}

func getPlaylist(ctx context.Context, client *spotify.Client, playlistId spotify.ID) *spotify.FullPlaylist {
	//fmt.Println("Beginning getPlaylist")

	playlist, err := client.GetPlaylist(ctx, playlistId)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(playlist.ID)
	return playlist
}

func findMood(ctx context.Context, client *spotify.Client, attributes []*spotify.AudioFeatures) string {
	fmt.Println("Fetching track attributes")
	var energy float64
	var valence float64
	var danceability float64
	var tempo float64
	var mood string
	var intensity int = 0
	var sentiment int = 0

	for _, attribute := range attributes {
		energy = float64(attribute.Energy)
		valence = float64(attribute.Valence)
		danceability = float64(attribute.Danceability)
		tempo = float64(attribute.Tempo)

	}
	fmt.Println(valence)
	if valence >= 0.5 {
		sentiment++
	} else if valence < .3 {
		sentiment = sentiment - 4
	} else {
		sentiment--
	}

	if energy >= 0.5 {
		intensity = intensity + 2
	} else {
		intensity = intensity - 2
	}

	if tempo >= 130 {
		intensity++
		sentiment++
	} else {
		intensity--
		sentiment--
	}

	if danceability >= 0.5 {
		intensity++
		sentiment++
	} else {
		intensity--
		sentiment--
	}

	if intensity >= 0 {
		if sentiment >= 0 {
			mood = "excited"
		} else {
			mood = "angry"
		}
	} else {
		if sentiment <= 0 {
			mood = "sad"
		} else {
			mood = "happy"
		}
	}

	fmt.Println("Song mood is", mood)
	return mood
}
