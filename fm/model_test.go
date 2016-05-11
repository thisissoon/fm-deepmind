package fm

import (
	"github.com/stretchr/testify/assert"
	"github.com/thisissoon/fm-deepmind/math"
	"testing"
	"time"
)

func TestDataSetAppend(t *testing.T) {
	ds := DataSet{}

	assert.Equal(t, 0, len(ds.D))
	ds.Append(DataObject{Label: "abc", Total: 4})
	assert.Equal(t, 1, len(ds.D))
}

func TestDatasetKeepsCorrectTotal(t *testing.T) {
	ds := DataSet{}
	assert.Equal(t, 0, ds.Total)

	ds.Append(DataObject{Label: "abc", Total: 4})
	assert.Equal(t, 4, ds.Total)

	ds.Append(DataObject{Label: "abc", Total: 5})
	assert.Equal(t, 9, ds.Total)
}

func TestDatasetReturnsCorrectWeights(t *testing.T) {
	ds := DataSet{}
	ds.Append(DataObject{Label: "abc", Total: 4})
	ds.Append(DataObject{Label: "abc", Total: 5})

	assert.Equal(t, []float64{4.0 / 9.0, 5.0 / 9.0}, ds.GetWeights())
}

func TestSumOfWeighShouleBeZero(t *testing.T) {
	ds := DataSet{}
	ds.Append(DataObject{Label: "abc", Total: 4})
	ds.Append(DataObject{Label: "abc", Total: 5})
	ds.Append(DataObject{Label: "abc", Total: 6})

	weights := ds.GetWeights()

	assert.Equal(t, float64(1), weights[0]+weights[1]+weights[2])
}

func TestAudioSumaryMatrixTempoMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Tempo:     0.1,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Tempo:     0.3,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Tempo:     0.4,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Tempo:     0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Tempo)
}

func TestAudioSumaryMatrixLivenessMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Liveness:  0.1,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Liveness:  0.3,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Liveness:  0.4,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Liveness:  0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Liveness)
}

func TestAudioSumaryMatrixSpeechinessMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp:   time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Speechiness: 0.1,
		},
		AudioSummary{
			Timestamp:   time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Speechiness: 0.3,
		},
		AudioSummary{
			Timestamp:   time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Speechiness: 0.4,
		},
		AudioSummary{
			Timestamp:   time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Speechiness: 0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Speechiness)
}

func TestAudioSumaryMatrixAcousticnessMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Acousticness: 0.1,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Acousticness: 0.3,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Acousticness: 0.4,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Acousticness: 0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Acousticness)
}

func TestAudioSumaryMatrixInstrumentalnessMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp:        time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Instrumentalness: 0.1,
		},
		AudioSummary{
			Timestamp:        time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Instrumentalness: 0.3,
		},
		AudioSummary{
			Timestamp:        time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Instrumentalness: 0.4,
		},
		AudioSummary{
			Timestamp:        time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Instrumentalness: 0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Instrumentalness)
}

func TestAudioSumaryMatrixLoudnessMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Loudness:  0.1,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Loudness:  0.3,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Loudness:  0.4,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Loudness:  0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Loudness)
}

func TestAudioSumaryMatrixValenceMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Valence:   0.1,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Valence:   0.3,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Valence:   0.4,
		},
		AudioSummary{
			Timestamp: time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Valence:   0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Valence)
}

func TestAudioSumaryMatrixDanceabilityMediam(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Danceability: 0.1,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Danceability: 0.3,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Danceability: 0.4,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Danceability: 0.4,
		},
	}
	gm.Populate(audioSummaries)
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), gm.M[1][8].Danceability)
}

func TestAudioSumaryMatrixGetExistingQuartile(t *testing.T) {
	gm := AudioSumaryMatrix{}

	audioSummaries := []AudioSummary{
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 0, 0, 0, time.UTC),
			Danceability: 0.1,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.April, 25, 8, 30, 0, 0, time.UTC),
			Danceability: 0.3,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 25, 9, 0, 0, 0, time.UTC),
			Danceability: 0.4,
		},
		AudioSummary{
			Timestamp:    time.Date(2016, time.May, 3, 8, 0, 0, 0, time.UTC),
			Danceability: 0.4,
		},
	}
	gm.Populate(audioSummaries)
	quartile := gm.GetQuartile(time.Date(2016, time.April, 25, 8, 15, 0, 0, time.UTC))
	assert.Equal(t, math.GetQuartile([]float64{0.1, 0.3}), quartile.Danceability)
}

func TestAudioSumaryMatrixGetNonExistingQuartile(t *testing.T) {
	gm := AudioSumaryMatrix{}
	gm.Populate([]AudioSummary{})

	quartile := gm.GetQuartile(time.Now())
	assert.Equal(t, math.Quartile{0, 0, 0}, quartile.Danceability)
}

func TestAudioSummaryWeightsPopulateTempo(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Tempo: 1},
		AudioSummary{Tempo: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Tempo)
}

func TestAudioSummaryWeightsPopulateEnergy(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Energy: 1},
		AudioSummary{Energy: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Energy)
}

func TestAudioSummaryWeightsPopulateLiveness(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Liveness: 1},
		AudioSummary{Liveness: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Liveness)
}

func TestAudioSummaryWeightsPopulateSpeechiness(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Speechiness: 1},
		AudioSummary{Speechiness: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Speechiness)
}

func TestAudioSummaryWeightsPopulateAcousticness(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Acousticness: 1},
		AudioSummary{Acousticness: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Acousticness)
}

func TestAudioSummaryWeightsPopulateInstrumentalness(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Instrumentalness: 1},
		AudioSummary{Instrumentalness: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Instrumentalness)
}

func TestAudioSummaryWeightsPopulateLoudness(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Loudness: 1},
		AudioSummary{Loudness: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Loudness)
}

func TestAudioSummaryWeightsPopulateValence(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Valence: 1},
		AudioSummary{Valence: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Valence)
}

func TestAudioSummaryWeightsPopulateDanceability(t *testing.T) {
	a := []AudioSummary{
		AudioSummary{Danceability: 1},
		AudioSummary{Danceability: 5},
	}

	w := AudioSummaryWeights{}
	w.Populate(a)

	assert.Equal(t, float64(4), w.Danceability)
}
