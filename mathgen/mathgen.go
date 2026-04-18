// Package mathgen generates math-operations crossword puzzles.
package mathgen

import (
	"errors"
	"math/rand"
	"sort"
	"strconv"
)

type Op byte

const (
	OpAdd Op = '+'
	OpSub Op = '-'
	OpMul Op = '*'
	OpDiv Op = '/'
)

type CellIdx struct{ R, C int }

type Equation struct {
	R, C       int
	Horizontal bool
	Op         Op
}

func (e Equation) Cells() [3]CellIdx {
	if e.Horizontal {
		return [3]CellIdx{{e.R, e.C}, {e.R, e.C + 2}, {e.R, e.C + 4}}
	}
	return [3]CellIdx{{e.R, e.C}, {e.R + 2, e.C}, {e.R + 4, e.C}}
}

func (e Equation) OpCell() CellIdx {
	if e.Horizontal {
		return CellIdx{e.R, e.C + 1}
	}
	return CellIdx{e.R + 1, e.C}
}

func (e Equation) EqCell() CellIdx {
	if e.Horizontal {
		return CellIdx{e.R, e.C + 3}
	}
	return CellIdx{e.R + 3, e.C}
}

func apply(a, b int, op Op) (int, bool) {
	switch op {
	case OpAdd:
		return a + b, true
	case OpSub:
		r := a - b
		return r, r >= 1
	case OpMul:
		return a * b, true
	case OpDiv:
		if b == 0 || a%b != 0 {
			return 0, false
		}
		r := a / b
		return r, r >= 1
	}
	return 0, false
}

func verify(a, b, c int, op Op) bool {
	r, ok := apply(a, b, op)
	return ok && r == c
}

type Puzzle struct {
	Size      int
	MaxVal    int
	Equations []Equation
	Values    map[CellIdx]int
}

func (p *Puzzle) UniqueCells() []CellIdx {
	set := map[CellIdx]bool{}
	for _, eq := range p.Equations {
		for _, c := range eq.Cells() {
			set[c] = true
		}
	}
	cells := make([]CellIdx, 0, len(set))
	for c := range set {
		cells = append(cells, c)
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].R != cells[j].R {
			return cells[i].R < cells[j].R
		}
		return cells[i].C < cells[j].C
	})
	return cells
}

func (p *Puzzle) Solve(rng *rand.Rand) bool {
	cells := p.UniqueCells()
	p.Values = make(map[CellIdx]int, len(cells))
	return p.backtrack(cells, 0, rng)
}

func (p *Puzzle) backtrack(cells []CellIdx, idx int, rng *rand.Rand) bool {
	if idx == len(cells) {
		return true
	}
	cell := cells[idx]
	for _, v0 := range rng.Perm(p.MaxVal) {
		p.Values[cell] = v0 + 1
		if p.consistent(cell) && p.backtrack(cells, idx+1, rng) {
			return true
		}
	}
	delete(p.Values, cell)
	return false
}

func (p *Puzzle) consistent(changed CellIdx) bool {
	for _, eq := range p.Equations {
		cs := eq.Cells()
		involves := false
		for _, c := range cs {
			if c == changed {
				involves = true
				break
			}
		}
		if !involves {
			continue
		}
		a, aOk := p.Values[cs[0]]
		b, bOk := p.Values[cs[1]]
		c, cOk := p.Values[cs[2]]
		switch {
		case aOk && bOk && cOk:
			if !verify(a, b, c, eq.Op) {
				return false
			}
		case aOk && bOk:
			r, ok := apply(a, b, eq.Op)
			if !ok || r < 1 || r > p.MaxVal {
				return false
			}
		case aOk && cOk:
			switch eq.Op {
			case OpAdd:
				x := c - a
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpSub:
				x := a - c
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpMul:
				if a == 0 || c%a != 0 {
					return false
				}
				x := c / a
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpDiv:
				if c == 0 || a%c != 0 {
					return false
				}
				x := a / c
				if x < 1 || x > p.MaxVal {
					return false
				}
			}
		case bOk && cOk:
			switch eq.Op {
			case OpAdd:
				x := c - b
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpSub:
				x := c + b
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpMul:
				if b == 0 || c%b != 0 {
					return false
				}
				x := c / b
				if x < 1 || x > p.MaxVal {
					return false
				}
			case OpDiv:
				if b == 0 {
					return false
				}
				x := c * b
				if x < 1 || x > p.MaxVal {
					return false
				}
			}
		}
	}
	return true
}

func (p *Puzzle) ToGrid(hidden map[CellIdx]bool) [][]string {
	g := make([][]string, p.Size)
	for i := range g {
		g[i] = make([]string, p.Size)
		for j := range g[i] {
			g[i][j] = "0"
		}
	}
	for _, eq := range p.Equations {
		cs := eq.Cells()
		opc := eq.OpCell()
		eqc := eq.EqCell()
		g[opc.R][opc.C] = string(eq.Op)
		g[eqc.R][eqc.C] = "="
		for _, c := range cs {
			if hidden[c] {
				g[c.R][c.C] = "999"
			} else {
				g[c.R][c.C] = strconv.Itoa(p.Values[c])
			}
		}
	}
	return g
}

// BuildLayout returns the default 15x15 crossword skeleton.
func BuildLayout() *Puzzle {
	return &Puzzle{
		Size:   15,
		MaxVal: 99,
		Equations: []Equation{
			{R: 1, C: 1, Horizontal: true, Op: OpAdd},
			{R: 3, C: 1, Horizontal: true, Op: OpSub},
			{R: 1, C: 1, Horizontal: false, Op: OpAdd},
			{R: 1, C: 5, Horizontal: false, Op: OpAdd},
			{R: 1, C: 9, Horizontal: true, Op: OpMul},
			{R: 3, C: 9, Horizontal: true, Op: OpAdd},
			{R: 3, C: 11, Horizontal: false, Op: OpAdd},
			{R: 7, C: 1, Horizontal: true, Op: OpSub},
			{R: 7, C: 3, Horizontal: false, Op: OpMul},
			{R: 11, C: 1, Horizontal: true, Op: OpAdd},
			{R: 13, C: 1, Horizontal: true, Op: OpSub},
			{R: 9, C: 7, Horizontal: true, Op: OpDiv},
			{R: 13, C: 9, Horizontal: true, Op: OpAdd},
		},
	}
}

// Result bundles the solved puzzle, the printable grids, and the hidden mask.
type Result struct {
	Layout   *Puzzle
	Hidden   map[CellIdx]bool
	Puzzle   [][]string
	Solution [][]string
}

// isSolvable returns true if every hidden cell can be deduced by repeatedly
// applying the rule: if 2 of 3 cells in an equation are known, the 3rd is too.
func isSolvable(p *Puzzle, hidden map[CellIdx]bool) bool {
	known := map[CellIdx]bool{}
	for _, c := range p.UniqueCells() {
		if !hidden[c] {
			known[c] = true
		}
	}
	for {
		progress := false
		for _, eq := range p.Equations {
			cs := eq.Cells()
			var unknowns []CellIdx
			for _, c := range cs {
				if !known[c] {
					unknowns = append(unknowns, c)
				}
			}
			if len(unknowns) == 1 {
				known[unknowns[0]] = true
				progress = true
			}
		}
		if !progress {
			break
		}
	}
	for _, c := range p.UniqueCells() {
		if !known[c] {
			return false
		}
	}
	return true
}

// Generate solves the default layout, hides cells with probability hideProb,
// and returns a ready-to-render Result.
func Generate(rng *rand.Rand, hideProb float64) (*Result, error) {
	p := BuildLayout()
	if !p.Solve(rng) {
		return nil, errors.New("mathgen: failed to solve layout")
	}
	cells := p.UniqueCells()
	var hidden map[CellIdx]bool
	for attempt := 0; attempt < 200; attempt++ {
		hidden = map[CellIdx]bool{}
		for _, c := range cells {
			if rng.Float64() < hideProb {
				hidden[c] = true
			}
		}
		// ensure every equation has at least one hidden cell
		for _, eq := range p.Equations {
			cs := eq.Cells()
			hasHidden := false
			for _, c := range cs {
				if hidden[c] {
					hasHidden = true
					break
				}
			}
			if !hasHidden {
				hidden[cs[rng.Intn(3)]] = true
			}
		}
		if isSolvable(p, hidden) {
			break
		}
	}
	// fallback: reveal cells one by one until solvable
	if !isSolvable(p, hidden) {
		order := rng.Perm(len(cells))
		for _, i := range order {
			if hidden[cells[i]] {
				delete(hidden, cells[i])
				if isSolvable(p, hidden) {
					break
				}
			}
		}
	}
	return &Result{
		Layout:   p,
		Hidden:   hidden,
		Puzzle:   p.ToGrid(hidden),
		Solution: p.ToGrid(nil),
	}, nil
}
