package rnd

import (
	"math/rand"
	"time"
)

// Returns sum of given array
func float64Sum(l []float64) float64 {
	sum := float64(0)
	for _, n := range l {
		sum += n
	}
	return sum
}

// Returns index of `w` array based on its probability weights. For `w`
// equals to {.2, .3, .5} is 50% change that function returns `2`, 30% `1`
// and 20% `1`. `r` - uniformly distributed random number
func Weight(w []float64) int {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := rand.Float64() * float64Sum(w)
	for i, w := range w {
		rnd -= w
		if rnd < 0 {
			return i
		}
	}
	return 0
}
