package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/fm_api"
	"github.com/thisissoon/fm-deepmind/rnd"

	"github.com/spf13/viper"
	"github.com/thisissoon/fm-deepmind/socket"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("USER_TOKEN", "")
	viper.SetDefault("PERCEPTOR_ADDRESS", "perceptor.thisissoon.fm")
	viper.SetDefault("SECRET", "")
	viper.SetDefault("DB", "")

	log.Printf("%s", viper.GetString("PERCEPTOR_ADDRESS"))

	eventChannel := make(chan []byte)
	endChannel := make(chan bool)

	perceptor := socket.NewPerceptorService(
		viper.GetString("PERCEPTOR_ADDRESS"),
		viper.GetString("SECRET"),
		eventChannel)
	go perceptor.Run()

	eventHandler := socket.NewHandler(eventChannel, endChannel)
	go eventHandler.Run()

	volumeManager := fm_api.FmApiManager{
		Token: viper.GetString("USER_TOKEN"),
	}

	data := fm.DataAdapter{}
	if err := data.Conn(viper.GetString("DB")); err != nil {
		log.Fatal("Scan: %e", err)
		return
	}

	anon := func() string {
		genres := data.GetGenreDataSet(14)
		genreIndex := rnd.Weight(genres.GetWeights())
		d := genres.Get(genreIndex)

		tracks := data.GetTrackDataSetBasedOnGenre(d.Id)
		trackIndex := rnd.Weight(tracks.GetWeights())
		track := tracks.Get(trackIndex)
		return track.Id
	}
	go volumeManager.Listen(endChannel, 2, anon)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// Run forever unless we get a signal
	for sig := range signals {
		log.Println(sig)
		os.Exit(1)
	}
}
