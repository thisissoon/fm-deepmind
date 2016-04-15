package fm_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type FmApiManager struct {
	Token string
}

func (m *FmApiManager) AddTrack(t string) {
	url := "https://api.thisissoon.fm/player/queue"

	user := &QueueTrack{Uri: t}
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Access-Token", m.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusCreated {
		log.Printf("Song `%s` added to the queue", t)
	}
	defer resp.Body.Close()
}

func (m *FmApiManager) GetQueue() ([]QueueItem, error) {
	url := "https://api.thisissoon.fm/player/queue"

	var err error
	var queueItem []QueueItem

	response, err := http.Get(url)
	if err != nil {
		log.Errorf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("Cannot get response body: %s", err)
		}

		if err = json.Unmarshal(contents, &queueItem); err != nil {
			log.Println("cannot unmarshal queue:", err)
		}
	}

	return queueItem, err
}

func (m *FmApiManager) Listen(c chan bool, l int, r func() string) {
	for {
		<-c
		if time.Now().Hour() <= 18 {
			queue, err := m.GetQueue()
			if err != nil {
				log.Printf("Error when fetching queue", err)
				continue
			}

			for i := 0; i < (l - len(queue)); i++ {
				go m.AddTrack(r())
			}
		}
	}
}
