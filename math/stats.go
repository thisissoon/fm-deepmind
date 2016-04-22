package math

import (
	"math"
	"sort"
)

type Quartile struct {
	First  float64
	Second float64
	Third  float64
}

func (q *Quartile) HasIn(f float64) bool {
	return q.First <= f && f <= q.Third
}

func Median(l []float64) float64 {
	sort.Float64s(l)
	lenght := float64(len(l))
	isOdd := math.Mod(lenght, float64(2)) == 1
	if isOdd {
		return l[int(lenght/2)]
	}

	middle := lenght / float64(2)
	return (l[int(middle-1)] + l[int(middle)]) / 2
}

func GetQuartile(l []float64) Quartile {
	sort.Float64s(l)

	lenght := float64(len(l))
	isOdd := math.Mod(lenght, float64(2)) == 1

	middle := int(lenght / 2)
	q := Quartile{Second: Median(l)}
	if isOdd {
		q.First = Median(l[:middle])
		q.Third = Median(l[middle:])
	} else {
		q.First = Median(l[:middle])
		q.Third = Median(l[middle:])
	}
	return q
}
