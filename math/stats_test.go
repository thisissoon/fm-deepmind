package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedianOddLenghOfList(t *testing.T) {
	l := []float64{2, 1, 6, 3, 8}
	assert.Equal(t, float64(3), Median(l))
}

func TestMedianEvenLenghOfList(t *testing.T) {
	l := []float64{3, 1, 7, 5}
	assert.Equal(t, float64(4), Median(l))
}

func TestGetQuartileEvenLength(t *testing.T) {
	l := []float64{15, 36, 7, 41, 39, 40}

	q := GetQuartile(l)
	assert.Equal(t, float64(15), q.First)
	assert.Equal(t, float64(37.5), q.Second)
	assert.Equal(t, float64(40), q.Third)
}

func TestGetQuartileOddLength(t *testing.T) {
	l := []float64{36, 7, 6, 15, 39, 41, 40, 49, 42, 47, 43}

	q := GetQuartile(l)
	assert.Equal(t, float64(15), q.First)
	assert.Equal(t, float64(40), q.Second)
	assert.Equal(t, float64(42.5), q.Third)
}

func TestQuartileHasNotIn(t *testing.T) {
	q := Quartile{First: 0, Third: 10}

	assert.Equal(t, true, q.HasIn(5))
}

func TestQuartileHasIn(t *testing.T) {
	q := Quartile{First: 0, Third: 10}

	assert.Equal(t, false, q.HasIn(15))
}
