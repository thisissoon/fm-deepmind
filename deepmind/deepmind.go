package deepmind

import (
	log "github.com/Sirupsen/logrus"

	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/fm_api"
	"github.com/thisissoon/fm-deepmind/rnd"
	"github.com/thisissoon/fm-deepmind/socket"
)

type Config struct {
	UserToken    string
	EventService string
	Secret       string
	Db           string
}

func NewDeepmind(c Config) *Deepmind {
	log.Printf("Creating Deepmind instance")
	return &Deepmind{
		config:       c,
		eventChannel: make(chan []byte),
		endChannel:   make(chan bool),
		data:         fm.NewDataAdapter(c.Db),
	}
}

type Deepmind struct {
	config       Config
	eventChannel chan []byte
	endChannel   chan bool
	data         *fm.DataAdapter
}

func (d *Deepmind) getRandomTrackUri() string {
	genres := d.data.GetGenreDataSet(14)
	genreIndex := rnd.Weight(genres.GetWeights())
	genre := genres.Get(genreIndex)

	tracks := d.data.GetTrackDataSetBasedOnGenre(genre.Id)
	trackIndex := rnd.Weight(tracks.GetWeights())
	track := tracks.Get(trackIndex)
	return track.Id
}

func (d *Deepmind) runPerceptorListener() {
	perceptor := socket.NewPerceptorService(
		d.config.EventService,
		d.config.Secret,
		d.eventChannel)
	go perceptor.Run()

	eventHandler := socket.NewHandler(d.eventChannel, d.endChannel)
	go eventHandler.Run()
}

func (d *Deepmind) Run() {
	d.runPerceptorListener()
	fmApi := fm_api.FmApiManager{Token: d.config.UserToken}
	go fmApi.Listen(d.endChannel, 6, d.getRandomTrackUri)
}
