package main

import (
	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/rnd"
	"log"
)

func main() {
	connString := ""
	data := fm.DataAdapter{}
	if err := data.Conn(connString); err != nil {
		log.Fatal("Scan: %e", err)
		return
	}

	gengers := data.GetGenreDataSet(14)
	d := gengers.Get(rnd.Weight(gengers.GetWeights()))

	tracks := data.GetTrackDataSet(d.Id)
	track := tracks.Get(rnd.Weight(tracks.GetWeights()))
	log.Printf("Picked track `%s` based on genre `%s`", track.Label, d.Label)
}
