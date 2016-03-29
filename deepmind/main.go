package main

import (
	// "encoding/json"
	"github.com/thisissoon/fm-deepmind/fm"
	// "github.com/thisissoon/fm-deepmind/rnd"
	"log"

	"encoding/json"
	"fmt"
	"net/http"
)

// func totalSum(t []fm.TotalGenres) int {
// 	sum := 0
// 	for _, e := range t {
// 		sum += e.Total
// 	}
// 	return sum
// }

// type Dataset struct {
// 	Name      string
// 	Prob      float64
// 	TotalHits int
// }

type SpotifyArtist struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}

type SpotifyTrack struct {
	Artists []SpotifyArtist `json:"artists"`
}

func getSpotifyTrack(uri string, target interface{}) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", uri)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	// log.Printf("fas")

	// connString := "postgresql://deepmind:deepmind@fm-production.cnwmyjremvfe.eu-west-1.rds.amazonaws.com:5432/production"

	connString := "postgres://deepmind:deepmind@fm-production.cnwmyjremvfe.eu-west-1.rds.amazonaws.com:5432/production"
	data := fm.DataAdapter{}
	if err := data.Conn(connString); err != nil {
		log.Fatal("Scan: %e", err)
		return
	}

	data.GetGenreWeight(14)

	// for _, t := range data.Tracks() {
	// 	track := &SpotifyTrack{}
	// 	getSpotifyTrack(t.GetSpotifyUri(), track)
	// 	for _, a := range track.Artists {
	// 		if dbArtist, err := data.GetArtist(a.Uri); err == nil {
	// 			log.Printf("SaveArtistToTrack: [%s, %s]", t.Id, dbArtist.Id)
	// 			data.SaveArtistToTrack(t.Id, dbArtist.Id)
	// 		}
	// 		// log.Printf("Artists: %s, %s", data.GetArtist(a.Uri).Id, t.Id)
	// 	}
	// }

	// db, _ := sql.Open("postgres", connString)
	// rows, err := db.Query("SELECT * FROM track LIMIT 10")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	_ = rows.Scan(&id, &name)

	// 	log.Printf("Name: %d %s", id, name)
	// }
	// err = rows.Err() // get any error encountered during iteration

	// url := "https://api.thisissoon.fm/player/stats"
	// res, err := http.Get(url)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // read body
	// body, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if res.StatusCode != 200 {
	// 	log.Fatal("Unexpected status code", res.StatusCode)
	// }

	// user := fm.Stats{}

	// err = json.Unmarshal(body, &user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// crippled := map[string]float64{
	// 	"indie r&b": .001,
	// 	"pop":       2,
	// }

	// gengers := user.MostPlayedGenres
	// for index, element := range gengers {
	// 	total := element.Total
	// 	if val, exists := crippled[element.Item]; exists {
	// 		total = int(float64(element.Total) * val)
	// 	}
	// 	gengers[index].Total = total
	// }

	// totalSum := totalSum(gengers)
	// datasets := []Dataset{}
	// weights := []float64{}
	// for _, element := range gengers {
	// 	prob := float64(element.Total) / float64(totalSum)
	// 	datasets = append(datasets, Dataset{Name: element.Item, Prob: prob, TotalHits: 0})
	// 	weights = append(weights, prob)
	// }

	// runs := 100000
	// for i := 0; i < runs; i++ {
	// 	datasetIndex := rnd.Weight(weights)
	// 	datasets[datasetIndex].TotalHits += 1
	// }

	// for index, value := range datasets {
	// 	log.Printf("%2d %-20s %20f %20f", index, value.Name, value.Prob, float64(value.TotalHits)/float64(runs))
	// }
}
