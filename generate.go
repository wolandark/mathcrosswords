package main

import (
	"math/rand"
	"time"

	"mathgen"
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

// GenSolution returns the solution grid for the most recently generated puzzle.
func GenSolution() [][]string {
	if lastResult == nil {
		GenGrid()
	}
	return lastResult.Solution
}
