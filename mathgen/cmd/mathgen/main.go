package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"mathgen"
)

func printGoGrid(w *os.File, label string, g [][]string) {
	fmt.Fprintf(w, "// %s\n", label)
	fmt.Fprintln(w, "grid := [][]string{")
	for _, row := range g {
		fmt.Fprint(w, "\t{")
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprintf(w, "%q", cell)
		}
		fmt.Fprintln(w, "},")
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func main() {
	seed := flag.Int64("seed", 0, "random seed (0 = time-based)")
	hideProb := flag.Float64("hide", 0.65, "probability a cell is hidden in puzzle")
	flag.Parse()

	s := *seed
	if s == 0 {
		s = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(s))

	r, err := mathgen.Generate(rng, *hideProb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("// seed: %d\n\n", s)
	fmt.Println("// Solved equations:")
	for _, eq := range r.Layout.Equations {
		cs := eq.Cells()
		a, b, c := r.Layout.Values[cs[0]], r.Layout.Values[cs[1]], r.Layout.Values[cs[2]]
		orient := "H"
		if !eq.Horizontal {
			orient = "V"
		}
		fmt.Printf("//  %s@(%d,%d): %d %c %d = %d\n", orient, eq.Row, eq.Col, a, eq.Op, b, c)
	}
	fmt.Println()

	printGoGrid(os.Stdout, "Puzzle", r.Puzzle)
	printGoGrid(os.Stdout, "Solution", r.Solution)
}
