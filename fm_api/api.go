package fm_api

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// return current london time
func getLondonTime() (time.Time, error) {
	location, err := time.LoadLocation("Europe/London")
	now := time.Now()
	if err != nil {
		log.Errorf("Cannot load location %v", err)
		return now, err
	}
	return now.In(location), nil
}

type FmApiManager struct {
	Token string
}

func (m *FmApiManager) AddTrack(t string) {
	url := "https://api.thisissoon.fm/player/queue"

	user := &QueueTrack{Uri: t}
	b, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
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

func (m *FmApiManager) parseQueue(c []byte) ([]QueueItem, error) {
	var err error
	var queueItem []QueueItem

	if err = json.Unmarshal(c, &queueItem); err != nil {
		log.Println("cannot unmarshal queue:", err)
	}
	return queueItem, err
}

func (m *FmApiManager) GetQueue() ([]QueueItem, error) {
	url := "https://api.thisissoon.fm/player/queue"

	response, err := http.Get(url)
	if err != nil {
		log.Errorf("%s", err)
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("Cannot get response body: %s", err)
		}
		return m.parseQueue(content)
	}

	return nil, err
}

func (m *FmApiManager) Listen(c chan bool, l int, r func() string) {
	for {
		<-c
		time, _ := getLondonTime()
		if time.Hour() < 18 { // dont run after 6pm
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
