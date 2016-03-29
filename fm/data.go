package fm

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"fmt"
)

type DataObject struct {
	Id    string
	Label string
	Total int
}

type DataSet struct {
	D     []DataObject
	Total int
}

func (ds *DataSet) Append(d DataObject) {
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

func (d *DataAdapter) GetGenreDataSet(days int) DataSet {
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
	  genre.id,
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

func (d *DataAdapter) Close() {
	d.Close()
}
