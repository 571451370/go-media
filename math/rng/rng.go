package rng

import (
	"math"
	"math/rand"
)

type RNG interface {
	Int() int
	Float64() float64
}

// https://en.wikipedia.org/wiki/Poisson_distribution
// Poisson generates Poisson-distributed random variables
func Poisson(lambda float64) float64 {
	const step = 500

	l := lambda
	k := 0.0
	p := 1.0
	for {
		k++
		u := rand.Float64()
		if u == 0 {
			u += 1e-3
		}
		p *= u
		for p < 1 && l > 0 {
			if l > step {
				p *= math.Exp(step)
				l -= step
			} else {
				p *= math.Exp(l)
				l = 0
			}
		}
		if !(p > 1) {
			break
		}
	}
	return k - 1
}
