package fm

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestParseAudioSumary(t *testing.T) {

	v := sql.NullString{
		String: `
			{
				"key": 2,
				"tempo": 149.999,
				"energy": 0.9,
				"liveness": 0.3,
				"analysis_url": "h",
				"speechiness": 0.1,
				"acousticness": 0.7,
				"instrumentalness": 0.9,
				"mode": 1,
				"time_signature": 4,
				"duration": 3.4,
				"loudness": -9.3,
				"valence": 0.22267801204461538,
				"danceability": 0.6063510834169058
			}
		`,
		Valid: true,
	}

	da := DataAdapter{}
	as, _ := da.parseAudioSummary(v)

	assert.Equal(t, 149.999, as.Tempo)
}

func TestParseAudioSumaryWhenValueIsNull(t *testing.T) {
	v := sql.NullString{
		String: `null`,
		Valid:  true,
	}

	da := DataAdapter{}
	as, _ := da.parseAudioSummary(v)

	assert.Equal(t, float64(0), as.Tempo)
}

func TestParseAudioSumaryWhenValueIsNill(t *testing.T) {
	v := sql.NullString{
		Valid: false,
	}

	da := DataAdapter{}
	as, _ := da.parseAudioSummary(v)

	assert.Equal(t, float64(0), as.Tempo)
}
