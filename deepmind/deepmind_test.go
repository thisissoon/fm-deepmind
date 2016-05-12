package deepmind

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/math"
)

func TestEchonestEndorsment(t *testing.T) {
	d := Deepmind{}
	d.audioSummaryWeights = &fm.AudioSummaryWeights{
		Energy:           float64(.11),
		Liveness:         float64(.22),
		Speechiness:      float64(.33),
		Acousticness:     float64(.11),
		Instrumentalness: float64(.11),
		Valence:          float64(.11),
		Danceability:     float64(.11),
	}

	o := fm.DataObject{
		Label: "abc",
		Total: 4,
		Meta: fm.AudioSummary{
			Energy:           float64(.2),
			Liveness:         float64(.3),
			Speechiness:      float64(.4),
			Acousticness:     float64(.5),
			Instrumentalness: float64(.6),
			Valence:          float64(.7),
			Danceability:     float64(.8),
		},
	}

	q := fm.AudioSummaryQuartile{
		Energy:           math.Quartile{.1, .2, .3},
		Liveness:         math.Quartile{.1, .2, .3},
		Speechiness:      math.Quartile{.1, .2, .3},
		Acousticness:     math.Quartile{.1, .2, .3},
		Instrumentalness: math.Quartile{.1, .2, .3},
		Valence:          math.Quartile{.1, .2, .3},
		Danceability:     math.Quartile{.1, .2, .3},
	}

	assert.Equal(t, 0.33, d.AudioSummaryEndorse(o, q))
}
