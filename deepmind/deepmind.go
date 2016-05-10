package deepmind

import (
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/gorhill/cronexpr"

	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/fm_api"
	"github.com/thisissoon/fm-deepmind/rnd"
	"github.com/thisissoon/fm-deepmind/socket"
)

// return current london time
func getLondonTime() time.Time {
	location, _ := time.LoadLocation("Europe/London")
	now := time.Now()
	return now.In(location)
}

type Config struct {
	UserToken    string
	EventService string
	Secret       string
	Db           string
	MinTracks    int // keep min number of tracks in the queue
	LastTrackAt  int // play until
}

func NewDeepmind(c Config) *Deepmind {
	log.Printf("Creating Deepmind instance")
	return &Deepmind{
		config:            c,
		eventChannel:      make(chan []byte),
		endChannel:        make(chan bool),
		data:              fm.NewDataAdapter(c.Db),
		audioSumaryMatrix: &fm.AudioSumaryMatrix{},
	}
}

type Deepmind struct {
	config            Config
	eventChannel      chan []byte
	endChannel        chan bool
	data              *fm.DataAdapter
	audioSumaryMatrix *fm.AudioSumaryMatrix
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

	// push tracks into the queue
	fmApi := fm_api.FmApiManager{Token: d.config.UserToken}
	if false {
		go d.Listen(fmApi, d.getRandomTrackUri)
	}

	// synchronise echonest stats
	d.ScheduleAction(cronexpr.MustParse("* * * * * * *"), func() {
		d.PopulateAudioSumaryMatrix()
		log.Printf("%s", d.audioSumaryMatrix.GetQuartile(getLondonTime()))
	})
}

func (d *Deepmind) PopulateAudioSumaryMatrix() {
	d.audioSumaryMatrix.Populate(d.data.GetAudioSummary())
}

func (d *Deepmind) Listen(m fm_api.FmApiManager, r func() string) {
	for {
		<-d.endChannel
		if getLondonTime().Hour() < d.config.LastTrackAt { // dont run after LastTrackAt
			queue, err := m.GetQueue()
			if err != nil {
				log.Printf("Error when fetching queue", err)
				continue
			}

			for i := 0; i < (d.config.MinTracks - len(queue)); i++ {
				go m.AddTrack(r())
			}
		}
	}
}

func (d *Deepmind) ScheduleAction(cron *cronexpr.Expression, f func()) {
	go func() {
		// Create a ticker function
		ticker := func() *time.Ticker {
			next := cron.Next(time.Now())
			diff := next.Sub(time.Now())
			log.Infof("Next Tick: %s", next.Format(time.RFC1123))
			return time.NewTicker(diff)
		}
		// Run the Ticker
		tkr := ticker()
		for {
			<-tkr.C
			f()
			tkr = ticker()
		}
	}()
}
