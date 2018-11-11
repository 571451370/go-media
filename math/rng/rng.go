package rng

type RNG interface {
	Int() int
	Float64() float64
}
