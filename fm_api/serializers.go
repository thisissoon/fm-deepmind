package fm_api

type QueueTrack struct {
	Uri string `json:"uri"`
}

type User struct {
	FamilyName       string `json:"family_name"`
	DisplayName      string `json:"display_name"`
	AvatarUrl        string `json:"avatar_url"`
	SpotifyPlaylists string `json:"spotify_playlists"`
	GivenName        string `json:"given_name"`
	Id               string `json:"id"`
}

type Track struct {
	Name      string `json:"Name"`
	Uri       string `json:"Uri"`
	PlayCount int    `json:"Play_count"`
	Duration  int    `json:"Duration"`
}

type QueueItem struct {
	Track Track `json:"track"`
	User  User  `json:"user"`
}
