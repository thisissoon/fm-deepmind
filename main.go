package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/thisissoon/fm-deepmind/deepmind"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("USER_TOKEN", "")
	viper.SetDefault("EVENT_SERVICE", "")
	viper.SetDefault("SECRET", "")
	viper.SetDefault("DB", "")
	viper.SetDefault("MIN_TRACKS", 2)
	viper.SetDefault("LAST_TRACK_AT", 18)

	deepmind := deepmind.NewDeepmind(
		deepmind.Config{
			UserToken:    viper.GetString("USER_TOKEN"),
			EventService: viper.GetString("EVENT_SERVICE"),
			Secret:       viper.GetString("SECRET"),
			Db:           viper.GetString("DB"),
			MinTracks:    viper.GetInt("MIN_TRACKS"),
			LastTrackAt:  viper.GetInt("LAST_TRACK_AT"),
		},
	)

	deepmind.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// Run forever unless we get a signal
	for sig := range signals {
		log.Println(sig)
		os.Exit(1)
	}
}
