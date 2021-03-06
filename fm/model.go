package fm

import (
	"github.com/thisissoon/fm-deepmind/math"
	"time"
)

func Map(vs []AudioSummary, f func(AudioSummary) float64) []float64 {
	vsm := make([]float64, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Raw AudioSummary object from a track
type AudioSummary struct {
	Key              int       `json:"key",csv:"key"`
	Tempo            float64   `json:"tempo",csv:"tempo"`
	Energy           float64   `json:"energy",csv:"energy"`
	Liveness         float64   `json:"liveness",csv:"liveness"`
	AnalysisUrl      string    `json:"analysis_url",csv:"-"`
	Speechiness      float64   `json:"speechiness",csv:"speechiness"`
	Acousticness     float64   `json:"acousticness",csv:"acousticness"`
	Instrumentalness float64   `json:"instrumentalness",csv:"instrumentalness"`
	Mode             float64   `json:"mode",csv:"mode"`
	TimeSignature    float64   `json:"time_signature",csv:"time_signature"`
	Duration         float64   `json:"duration",csv:"duration"`
	Loudness         float64   `json:"loudness",csv:"loudness"`
	Valence          float64   `json:"valence",csv:"valence"`
	Danceability     float64   `json:"danceability",csv:"danceability"`
	Timestamp        time.Time `csv:"timestamp"`
	IsValid          bool
}

// Quartiles of AudioSumary from set of tracks
type AudioSummaryQuartile struct {
	Tempo            math.Quartile
	Energy           math.Quartile
	Liveness         math.Quartile
	Speechiness      math.Quartile
	Acousticness     math.Quartile
	Instrumentalness math.Quartile
	Loudness         math.Quartile
	Valence          math.Quartile
	Danceability     math.Quartile
}

// Weights for AudioSumary from set of tracks
type AudioSummaryWeights struct {
	Tempo            float64
	Energy           float64
	Liveness         float64
	Speechiness      float64
	Acousticness     float64
	Instrumentalness float64
	Loudness         float64
	Valence          float64
	Danceability     float64
}

func (a *AudioSummaryWeights) Populate(l []AudioSummary) {
	a.Tempo = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Tempo
	}))
	a.Energy = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Energy
	}))
	a.Energy = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Energy
	}))
	a.Liveness = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Liveness
	}))
	a.Speechiness = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Speechiness
	}))
	a.Acousticness = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Acousticness
	}))
	a.Instrumentalness = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Instrumentalness
	}))
	a.Loudness = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Loudness
	}))
	a.Valence = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Valence
	}))
	a.Danceability = math.StandardDeviation(Map(l, func(v AudioSummary) float64 {
		return v.Danceability
	}))
	a.normalise()
}

func (a *AudioSummaryWeights) normalise() {
	l := []float64{
		// a.Tempo,
		a.Energy,
		a.Liveness,
		a.Speechiness,
		a.Acousticness,
		a.Instrumentalness,
		// a.Loudness,
		a.Valence,
		a.Danceability,
	}
	s := math.Sum(l)

	// a.Tempo = a.Tempo / s
	a.Energy = a.Energy / s
	a.Liveness = a.Liveness / s
	a.Speechiness = a.Speechiness / s
	a.Acousticness = a.Acousticness / s
	a.Instrumentalness = a.Instrumentalness / s
	// a.Loudness = a.Loudness / s
	a.Valence = a.Valence / s
	a.Danceability = a.Danceability / s
}

type DynamicList struct {
	L []AudioSummary
}

func (dl *DynamicList) Add(i AudioSummary) {
	dl.L = append(dl.L, i)
}

type DataObject struct {
	Id      string
	Label   string
	Total   int
	Weight  float64
	Endorse float64
	Meta    interface{}
}

type DataSet struct {
	D     []DataObject
	Total int
}

func (ds *DataSet) Endorse(f func(DataObject) float64) {
	for i := 0; i < len(ds.D); i++ {
		ds.D[i].Endorse = f(ds.D[i])
	}
}

func (ds *DataSet) Append(d DataObject) {
	ds.D = append(ds.D, d)
	ds.Total += d.Total
	ds.CalculateWeights()
}

func (ds *DataSet) CalculateWeights() {
	for i := 0; i < len(ds.D); i++ {
		ds.D[i].Weight = float64(ds.D[i].Total) / float64(ds.Total)
	}
}

func (ds *DataSet) GetWeights() []float64 {
	weights := []float64{}
	for _, d := range ds.D {
		weights = append(weights, d.Weight+d.Weight*d.Endorse)
	}
	return weights
}

func (ds *DataSet) Get(d int) DataObject {
	return ds.D[d]
}

// Holds AudioSummary objects in a weekday and hour matrix to easily pick a
// certain AudioSummary quartile in a certain time and day
type AudioSumaryMatrix struct {
	M [][]AudioSummaryQuartile
}

// Populate AudioSumary matrix from list of AudioSummary objects
func (m *AudioSumaryMatrix) Populate(a []AudioSummary) {
	week := make([][]DynamicList, 7)
	for i := range week {
		week[i] = make([]DynamicList, 24)
	}

	for _, audioSummary := range a {
		time := audioSummary.Timestamp
		weekday := time.Weekday()
		hour := time.Hour()
		week[weekday][hour].Add(audioSummary)
	}

	m.M = make([][]AudioSummaryQuartile, 7)
	for i := range m.M {
		m.M[i] = make([]AudioSummaryQuartile, 24)
	}

	for d, hours := range week {
		for h, hour := range hours {
			m.M[d][h] = AudioSummaryQuartile{
				Tempo: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Tempo
				})),
				Energy: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Energy
				})),
				Liveness: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Liveness
				})),
				Speechiness: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Speechiness
				})),
				Acousticness: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Acousticness
				})),
				Instrumentalness: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Instrumentalness
				})),
				Loudness: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Loudness
				})),
				Valence: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Valence
				})),
				Danceability: math.GetQuartile(Map(hour.L, func(v AudioSummary) float64 {
					return v.Danceability
				})),
			}
		}
	}
}

func (m *AudioSumaryMatrix) GetQuartile(t time.Time) AudioSummaryQuartile {
	return m.M[t.Weekday()][t.Hour()]
}
