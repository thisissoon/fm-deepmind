package fm

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"

	"errors"
	"fmt"
	"time"
)

type DataAdapter struct {
	Db *sql.DB
}

func (d *DataAdapter) Conn(s string) error {
	var err error
	d.Db, err = sql.Open("postgres", s)
	return err
}

func (d *DataAdapter) populateDataSet(query string) DataSet {
	rows, err := d.Db.Query(query)
	if err != nil {
		log.Fatal("%e", err)
	}
	defer rows.Close()
	dataset := DataSet{}
	for rows.Next() {
		var total int
		var id string
		var label string
		err = rows.Scan(&total, &id, &label)
		if err != nil {
			log.Fatal("Scan: %e", err)
		}
		dataset.Append(DataObject{Id: id, Label: label, Total: total})
	}

	return dataset
}

func (d *DataAdapter) GetTrackDataSetBasedOnGenre(genreId string) DataSet {
	query := fmt.Sprintf(`
		SELECT count(track.spotify_uri), track.spotify_uri, track.name, audio_summary
		FROM track
		INNER JOIN artist_track ON artist_track.track_id = track.id
		INNER JOIN artist_genre ON artist_genre.artist_id = artist_track.artist_id
		INNER JOIN playlist_history ON playlist_history.track_id = track.id
		WHERE artist_genre.grenre_id = '%s'
		GROUP BY track.id
		HAVING count(track.spotify_uri) > 1`, genreId)

	rows, err := d.Db.Query(query)
	if err != nil {
		log.Fatal("%e", err)
	}
	defer rows.Close()
	dataset := DataSet{}

	for rows.Next() {
		var total int
		var id string
		var label string
		var as sql.NullString
		err = rows.Scan(&total, &id, &label, &as)
		if err != nil {
			log.Fatal("Scan: %e", err)
		}
		if audioSummary, err := d.parseAudioSummary(as); err != nil {
			dataset.Append(DataObject{
				Id:    id,
				Label: label,
				Total: total,
				Meta:  audioSummary,
			})
		}
	}

	return dataset
}

func (d *DataAdapter) GetGenreDataSet(days int) DataSet {
	query := fmt.Sprintf(`
		SELECT
			(
			SELECT (
				SELECT array_length(array_agg(users), 1)
				FROM (
					SELECT DISTINCT users
					FROM (
						SELECT unnest(array_agg(playlist_history.user_id)) AS users) AS t1
					) AS t2
				) AS the_array
			),
			genre.id,
			genre.name
		FROM genre
		INNER JOIN artist_genre ON genre.id = artist_genre.grenre_id
		INNER JOIN artist_track ON artist_track.artist_id = artist_genre.artist_id
		INNER JOIN track ON track.id = artist_track.track_id
		INNER JOIN playlist_history ON playlist_history.track_id = track.id
		WHERE playlist_history.created > now() - interval '%d days'
		GROUP BY genre.id`, days)

	return d.populateDataSet(query)
}

func (d *DataAdapter) parseAudioSummary(v sql.NullString) (AudioSummary, error) {
	as := AudioSummary{}
	if !v.Valid {
		return as, errors.New("Value is not valid")
	}
	labelValue, err := v.Value()
	val := labelValue.(string)
	if val != "null" {
		if err = json.Unmarshal([]byte(val), &as); err != nil {
			return AudioSummary{}, err
		}
	}
	return as, err
}

func (d *DataAdapter) GetAudioSummary() []AudioSummary {
	audioSummaries := []AudioSummary{}

	query := `
		SELECT audio_summary, created
		FROM track WHERE audio_summary is not null`
	rows, err := d.Db.Query(query)
	if err != nil {
		log.Fatal("%e", err)
	}
	defer rows.Close()
	for rows.Next() {
		var label sql.NullString
		var timestamp time.Time
		err = rows.Scan(&label, &timestamp)
		if err != nil {
			log.Fatal("Scan: %e", err)
		} else {
			labelValue, err := label.Value()
			val := labelValue.(string)
			if val != "null" {
				audioSummary := AudioSummary{Timestamp: timestamp}
				if err = json.Unmarshal([]byte(val), &audioSummary); err != nil {
					log.Println(err)
				} else {
					audioSummaries = append(audioSummaries, audioSummary)
				}
			}
		}
	}
	return audioSummaries
}

func (d *DataAdapter) Close() {
	d.Close()
}
