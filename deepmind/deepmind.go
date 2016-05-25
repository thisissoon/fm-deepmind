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

func Map(vs []fm.AudioSummaryQuartile, f func(fm.AudioSummaryQuartile) fm.AudioSummary) []fm.AudioSummary {
	vsm := make([]fm.AudioSummary, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
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
		config:              c,
		eventChannel:        make(chan []byte),
		endChannel:          make(chan bool),
		data:                fm.NewDataAdapter(c.Db),
		audioSumaryMatrix:   &fm.AudioSumaryMatrix{},
		audioSummaryWeights: &fm.AudioSummaryWeights{},
	}
}

type Deepmind struct {
	config              Config
	eventChannel        chan []byte
	endChannel          chan bool
	data                *fm.DataAdapter
	audioSumaryMatrix   *fm.AudioSumaryMatrix
	audioSummaryWeights *fm.AudioSummaryWeights
}

func (d *Deepmind) getRandomTrackUri() fm.DataObject {
	genres := d.data.GetGenreDataSet(14)
	genreIndex := rnd.Weight(genres.GetWeights())
	genre := genres.Get(genreIndex)

	tracks := d.data.GetTrackDataSetBasedOnGenre(genre.Id)
	tracks.Endorse(func(o fm.DataObject) float64 {
		return d.AudioSummaryEndorse(o, d.audioSumaryMatrix.GetQuartile(time.Now()))
	})

	trackIndex := rnd.Weight(tracks.GetWeights())
	return tracks.Get(trackIndex)
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
	// push tracks into the queue
	fmApi := fm_api.FmApiManager{Token: d.config.UserToken}

	d.runPerceptorListener()
	go d.Listen(fmApi, d.getRandomTrackUri)

	go d.PopulateAudioSumaryMatrix()
	d.ScheduleAction(cronexpr.MustParse("0 0 1 * *"), func() {
		d.PopulateAudioSumaryMatrix()
	})

	d.ScheduleAction(cronexpr.MustParse("0 0 * * *"), func() {
		d.PopulateAudioSummaryWeights()
	})

	d.ScheduleAction(cronexpr.MustParse("0 10 * * *"), func() { // stop for stand ups
		fmApi.Pause()
	})
}

func (d *Deepmind) PopulateAudioSumaryMatrix() {
	a := d.data.GetAudioSummary()
	d.audioSumaryMatrix.Populate(a)
	d.PopulateAudioSummaryWeights()
}

func (d *Deepmind) PopulateAudioSummaryWeights() {
	today := time.Now().Weekday()
	audioSumaryMedians := Map(d.audioSumaryMatrix.M[today], func(v fm.AudioSummaryQuartile) fm.AudioSummary {
		return fm.AudioSummary{
			Energy:           v.Energy.Second,
			Liveness:         v.Liveness.Second,
			Speechiness:      v.Speechiness.Second,
			Acousticness:     v.Acousticness.Second,
			Instrumentalness: v.Instrumentalness.Second,
			Valence:          v.Valence.Second,
			Danceability:     v.Danceability.Second,
		}
	})
	d.audioSummaryWeights.Populate(audioSumaryMedians)
	log.Infof("Today's weights :%+v", d.audioSummaryWeights)
}

func (d *Deepmind) AudioSummaryEndorse(o fm.DataObject, q fm.AudioSummaryQuartile) float64 {
	endorse := float64(.0)
	if q.Energy.HasIn(o.Meta.(fm.AudioSummary).Energy) {
		endorse += d.audioSummaryWeights.Energy
	}
	if q.Liveness.HasIn(o.Meta.(fm.AudioSummary).Liveness) {
		endorse += d.audioSummaryWeights.Liveness
	}
	if q.Speechiness.HasIn(o.Meta.(fm.AudioSummary).Speechiness) {
		endorse += d.audioSummaryWeights.Speechiness
	}
	if q.Acousticness.HasIn(o.Meta.(fm.AudioSummary).Acousticness) {
		endorse += d.audioSummaryWeights.Acousticness
	}
	if q.Instrumentalness.HasIn(o.Meta.(fm.AudioSummary).Instrumentalness) {
		endorse += d.audioSummaryWeights.Instrumentalness
	}
	if q.Valence.HasIn(o.Meta.(fm.AudioSummary).Valence) {
		endorse += d.audioSummaryWeights.Valence
	}
	if q.Danceability.HasIn(o.Meta.(fm.AudioSummary).Danceability) {
		endorse += d.audioSummaryWeights.Danceability
	}
	return endorse
}

func (d *Deepmind) Listen(m fm_api.FmApiManager, r func() fm.DataObject) {
	for {
		<-d.endChannel
		if getLondonTime().Hour() < d.config.LastTrackAt { // dont run after LastTrackAt
			queue, err := m.GetQueue()
			if err != nil {
				log.Printf("Error when fetching queue", err)
				continue
			}

			for i := 0; i < (d.config.MinTracks - len(queue)); i++ {
				go m.AddTrack(r().Id)
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
