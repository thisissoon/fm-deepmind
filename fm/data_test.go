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
				"key": 5,
				"tempo": 88.088,
				"energy": 0.4246929822221389,
				"liveness": 0.11306411657392115,
				"analysis_url": "",
				"speechiness": 0.06528033930143728,
				"acousticness": 0.32661968843315015,
				"instrumentalness": 1.7090341899105053e-07,
				"mode": 0,
				"time_signature": 4,
				"duration": 284.73333,
				"loudness": -12.406,
				"valence": 0.5644159063204205,
				"danceability": 0.7072561410157304
			}
		`,
		Valid: true,
	}

	da := DataAdapter{}
	as, _ := da.parseAudioSummary(v)

	assert.Equal(t, 88.088, as.Tempo)
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
