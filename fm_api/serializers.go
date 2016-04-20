package fm_api

type QueueTrack struct {
	Uri string `json:"uri"`
}

type QueueItem struct {
	Track string `json:"track"`
	User  string `json:"user"`
}
