package socket

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

const (
	END_EVENT string = "end"
)

type Event struct {
	Type string `json:"event"`
}

type Handler struct {
	eventChannel chan []byte
	endChannel   chan bool
}

// Consumes event message channel, decoding the event and passing to
// the relevant service
func (h *Handler) Run() {
	for {
		msg := <-h.eventChannel
		var event Event
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Errorf("Failed to unmarshal event", err)
			continue
		}

		switch event.Type {

		case END_EVENT:
			go func() {
				h.endChannel <- true
			}()
		}
	}
}

// Constructs a new Event Handler
func NewHandler(e chan []byte, m chan bool) *Handler {
	return &Handler{
		eventChannel: e,
		endChannel:   m,
	}
}
