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

type DynamicList struct {
	L []AudioSummary
}

func (dl *DynamicList) Add(i AudioSummary) {
	dl.L = append(dl.L, i)
}

type DataObject struct {
	Id    string
	Label string
	Total int
	Meta  interface{}
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

func (ds *DataSet) Get(d int) DataObject {
	return ds.D[d]
}

type GenreMatrix struct {
	M [][]AudioSummaryQuartile
}

func (m *GenreMatrix) Populate(a []AudioSummary) {
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

	for w, hours := range week {
		for h, hour := range hours {
			m.M[w][h] = AudioSummaryQuartile{
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
