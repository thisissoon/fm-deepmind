package fm

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"fmt"
)

type DataObject struct {
	Name  string
	Total int
}

type DataSet struct {
	D     []DataObject
	Total int
}

func (ds *DataSet) append(d DataObject) {
	ds.D = append(ds.D, d)
	ds.Total += d.Total
}

func (ds *DataSet) GetWeights() []float64 {
	weights := []float64{}
	for _, d := range ds.D {
		weights = append(weights, float64(d.Total)/float64(ds.Total))
	}
	return weights
}

type DataAdapter struct {
	Db *sql.DB
}

func (d *DataAdapter) Conn(s string) error {
	var err error
	d.Db, err = sql.Open("postgres", s)
	return err
}

func (d *DataAdapter) Tracks() []Track {
	rows, err := d.Db.Query("SELECT id, name, audio_summary, spotify_uri FROM track")
	if err != nil {
		log.Fatal("Tracks: %e", err)
	}
	defer rows.Close()

	tracks := []Track{}
	for rows.Next() {
		var id sql.NullString
		var name sql.NullString
		var audioSummary sql.NullString
		var spotifyUri sql.NullString

		err = rows.Scan(&id, &name, &audioSummary, &spotifyUri)
		if err != nil {
			log.Fatal("Scan: %e", err)
		}

		track := Track{Id: id.String, Name: name.String, AudioSumary: audioSummary.String, SpotifyUri: spotifyUri.String}
		tracks = append(tracks, track)
	}

	err = rows.Err() // get any error encountered during iteration
	if err != nil {
		log.Fatal("Scan: %e", err)
	}

	return tracks
}

func (d *DataAdapter) GetGenreWeight(days int) {
	query := fmt.Sprintf(`
	SELECT
	  (
	    SELECT (
	      SELECT array_length(array_agg(users), 1)
	      FROM (
	        SELECT DISTINCT users
	        FROM (
	          SELECT unnest(array_agg(playlist_history.user_id)) AS users
	        ) AS t1
	      ) AS t2
	    ) AS the_array
	  ),
	  genre.name
	FROM genre
	INNER JOIN artist_genre ON genre.id = artist_genre.grenre_id
	INNER JOIN artist_track ON artist_track.artist_id = artist_genre.artist_id
	INNER JOIN track ON track.id = artist_track.track_id
	INNER JOIN playlist_history ON playlist_history.track_id = track.id
	WHERE playlist_history.created > now() - interval '%d days'
	GROUP BY genre.id`, days)

	rows, err := d.Db.Query(query)
	if err != nil {
		log.Fatal("Tracks: %e", err)
	}
	defer rows.Close()
	dataset := DataSet{}
	for rows.Next() {
		var weight int
		var name string
		err = rows.Scan(&weight, &name)
		if err != nil {
			log.Fatal("Scan: %e", err)
		}
		dataset.append(DataObject{Name: name, Total: weight})
	}

	weights := dataset.GetWeights()
	for i, a := range dataset.D {
		log.Printf("%10d %10f %30s", a.Total, weights[i], a.Name)
	}

}

func (d *DataAdapter) GetArtist(uri string) (Album, error) {
	var id sql.NullString
	query := fmt.Sprintf("SELECT id FROM artist WHERE spotify_uri = '%s' LIMIT 1", uri)
	if err := d.Db.QueryRow(query).Scan(&id); err != nil {
		log.Printf("GetArtist: %e", err)
		return Album{}, err
	}
	return Album{Id: id.String}, nil
}

func (d *DataAdapter) Close() {
	d.Close()
}
