package main

import (
	"flag"
	"log"

	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/fm_api"
	"github.com/thisissoon/fm-deepmind/rnd"
)

var token = flag.String("token", "", "deepmind user")
var db = flag.String("db", "", "soonfm database string")

func main() {
	data := fm.DataAdapter{}
	if err := data.Conn(*db); err != nil {
		log.Fatal("Scan: %e", err)
		return
	}

	for i := 0; i < 10; i++ {
		genres := data.GetGenreDataSet(14)
		genreWeights := genres.GetWeights()
		genreIndex := rnd.Weight(genres.GetWeights())
		d := genres.Get(genreIndex)

		tracks := data.GetTrackDataSet(d.Id)
		trackWeights := tracks.GetWeights()
		trackIndex := rnd.Weight(trackWeights)
		track := tracks.Get(trackIndex)
		log.Printf("%-9f %-50s %-9f %-10s", trackWeights[trackIndex], track.Label, genreWeights[genreIndex], d.Label)

		if false {
			go fm_api.AddTrack(track.Id, *token)
		}
	}

}
