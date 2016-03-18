package fm

type TotalStruct struct {
	Total int `json:"total"`
}

type TotalTrack struct {
	TotalStruct
	Item Track `json:"track"`
}

type TotalGenres struct {
	TotalStruct
	Item string `json:"name"`
}

type TotalArtist struct {
	TotalStruct
	Item Artist `json:"artist"`
}

type Album struct {
	Id         string `json:"id"`
	SpotifyUri string `json:"spotify_uri"`
	Name       string `json:"name"`
}

type Artist struct {
	Id         string `json:"id"`
	SpotifyUri string `json:"spotify_uri"`
	Name       string `json:"name"`
}

type Track struct {
	Name       string   `json:"name"`
	SpotifyUri string   `json:"spotify_uri"`
	Artists    []Artist `json:"artists"`
	Id         string   `json:"id"`
	Albums     Album    `json:"album"`
}

type Stats struct {
	MostPlayedTracks  []TotalTrack  `json:"most_played_tracks"`
	MostPlayedGenres  []TotalGenres `json:"most_played_genres"`
	MostPlayedArtists []TotalArtist `json:"most_played_artists"`
	TotalPlays        int           `json:"total_plays"`
}
