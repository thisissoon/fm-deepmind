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
	switch len(l) {
	case 0:
		return 0
	case 1:
		return l[0]
	case 2:
		return (l[0] + l[1]) / 2
	}
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

func Average(l []float64) float64 {
	total := 0.0
	for _, v := range l {
		total += v
	}
	return total / float64(len(l))
}

func SumSq(l []float64) float64 {
	sum := float64(0)
	for _, v := range l {
		sum += v * v
	}
	return sum
}

func StandardDeviation(l []float64) float64 {
	avg := Average(l)
	return SumSq(l)/float64(len(l)) - avg*avg
}
