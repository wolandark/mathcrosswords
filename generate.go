package main

import (
	"math/rand"
	"time"

	"mathcrossword/mathgen"
)

var lastResult *mathgen.Result

func GenGrid() [][]string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, err := mathgen.Generate(rng, 0.65)
	if err != nil {
		panic(err)
	}
	lastResult = r
	return r.Puzzle
}

func GenSolution() [][]string {
	if lastResult == nil {
		GenGrid()
	}
	return lastResult.Solution
}
