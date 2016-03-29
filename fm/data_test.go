package fm

import (
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
