package fm_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddTrack(spotify_uri string, token string) {

	url := "https://api.thisissoon.fm/player/queue"

	user := &QueueTrack{Uri: spotify_uri}
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Access-Token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
