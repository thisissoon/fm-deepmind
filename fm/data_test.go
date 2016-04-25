package fm

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
