package main

import (
	"math/rand"
)

// IRandom service generating random values
type IRandom interface {
	Name() string
}

// Random service implementaion
type Random struct{}

// NewRandom is a IRandom factory
func NewRandom() IRandom {
	rand.Seed(312719581257)

	return &Random{}
}

// Name generate random name
func (r *Random) Name() string {
	names := []string{
		"Artem",
		"Maxim",
		"Andrey",
	}

	return names[rand.Intn(len(names))]
}
