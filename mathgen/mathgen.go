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

type CellPosition struct{ Row, Col int }

type Equation struct {
	Row        int
	Col        int
	Horizontal bool
	Op         Op
}

func (eq Equation) Cells() [3]CellPosition {
	if eq.Horizontal {
		return [3]CellPosition{{eq.Row, eq.Col}, {eq.Row, eq.Col + 2}, {eq.Row, eq.Col + 4}}
	}
	return [3]CellPosition{{eq.Row, eq.Col}, {eq.Row + 2, eq.Col}, {eq.Row + 4, eq.Col}}
}

func (eq Equation) OpCell() CellPosition {
	if eq.Horizontal {
		return CellPosition{eq.Row, eq.Col + 1}
	}
	return CellPosition{eq.Row + 1, eq.Col}
}

func (eq Equation) EqCell() CellPosition {
	if eq.Horizontal {
		return CellPosition{eq.Row, eq.Col + 3}
	}
	return CellPosition{eq.Row + 3, eq.Col}
}

func applyOp(left, right int, op Op) (int, bool) {
	switch op {
	case OpAdd:
		return left + right, true
	case OpSub:
		result := left - right
		return result, result >= 1
	case OpMul:
		return left * right, true
	case OpDiv:
		if right == 0 || left%right != 0 {
			return 0, false
		}
		result := left / right
		return result, result >= 1
	}
	return 0, false
}

func equationHolds(left, right, result int, op Op) bool {
	computed, ok := applyOp(left, right, op)
	return ok && computed == result
}

type Puzzle struct {
	Size      int
	MaxVal    int
	Equations []Equation
	Values    map[CellPosition]int
}

func (puzzle *Puzzle) UniqueCells() []CellPosition {
	seen := map[CellPosition]bool{}
	for i := 0; i < len(puzzle.Equations); i++ {
		positions := puzzle.Equations[i].Cells()
		for j := 0; j < len(positions); j++ {
			seen[positions[j]] = true
		}
	}
	cells := make([]CellPosition, 0, len(seen))
	for pos := range seen {
		cells = append(cells, pos)
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Row != cells[j].Row {
			return cells[i].Row < cells[j].Row
		}
		return cells[i].Col < cells[j].Col
	})
	return cells
}

func (puzzle *Puzzle) Solve(rng *rand.Rand) bool {
	cells := puzzle.UniqueCells()
	puzzle.Values = make(map[CellPosition]int, len(cells))
	return puzzle.backtrack(cells, 0, rng)
}

func (puzzle *Puzzle) backtrack(cells []CellPosition, index int, rng *rand.Rand) bool {
	if index == len(cells) {
		return true
	}
	pos := cells[index]
	perm := rng.Perm(puzzle.MaxVal)
	for i := 0; i < len(perm); i++ {
		puzzle.Values[pos] = perm[i] + 1
		if puzzle.consistent(pos) && puzzle.backtrack(cells, index+1, rng) {
			return true
		}
	}
	delete(puzzle.Values, pos)
	return false
}

func (puzzle *Puzzle) consistent(changedPos CellPosition) bool {
	for i := 0; i < len(puzzle.Equations); i++ {
		eq := puzzle.Equations[i]
		positions := eq.Cells()
		involves := false
		for j := 0; j < len(positions); j++ {
			if positions[j] == changedPos {
				involves = true
				break
			}
		}
		if !involves {
			continue
		}
		left, leftKnown := puzzle.Values[positions[0]]
		right, rightKnown := puzzle.Values[positions[1]]
		result, resultKnown := puzzle.Values[positions[2]]
		switch {
		case leftKnown && rightKnown && resultKnown:
			if !equationHolds(left, right, result, eq.Op) {
				return false
			}
		case leftKnown && rightKnown:
			computed, ok := applyOp(left, right, eq.Op)
			if !ok || computed < 1 || computed > puzzle.MaxVal {
				return false
			}
		case leftKnown && resultKnown:
			switch eq.Op {
			case OpAdd:
				missing := result - left
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpSub:
				missing := left - result
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpMul:
				if left == 0 || result%left != 0 {
					return false
				}
				missing := result / left
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpDiv:
				if result == 0 || left%result != 0 {
					return false
				}
				missing := left / result
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			}
		case rightKnown && resultKnown:
			switch eq.Op {
			case OpAdd:
				missing := result - right
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpSub:
				missing := result + right
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpMul:
				if right == 0 || result%right != 0 {
					return false
				}
				missing := result / right
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			case OpDiv:
				if right == 0 {
					return false
				}
				missing := result * right
				if missing < 1 || missing > puzzle.MaxVal {
					return false
				}
			}
		}
	}
	return true
}

func (puzzle *Puzzle) ToGrid(hidden map[CellPosition]bool) [][]string {
	grid := make([][]string, puzzle.Size)
	for row := 0; row < puzzle.Size; row++ {
		grid[row] = make([]string, puzzle.Size)
		for col := 0; col < puzzle.Size; col++ {
			grid[row][col] = "0"
		}
	}
	for i := 0; i < len(puzzle.Equations); i++ {
		eq := puzzle.Equations[i]
		positions := eq.Cells()
		opPos := eq.OpCell()
		eqPos := eq.EqCell()
		grid[opPos.Row][opPos.Col] = string(eq.Op)
		grid[eqPos.Row][eqPos.Col] = "="
		for j := 0; j < len(positions); j++ {
			pos := positions[j]
			if hidden[pos] {
				grid[pos.Row][pos.Col] = "999"
			} else {
				grid[pos.Row][pos.Col] = strconv.Itoa(puzzle.Values[pos])
			}
		}
	}
	return grid
}

// BuildLayout returns the default 15x15 crossword skeleton.
func BuildLayout() *Puzzle {
	return &Puzzle{
		Size:   15,
		MaxVal: 99,
		Equations: []Equation{
			{Row: 1, Col: 1, Horizontal: true, Op: OpAdd},
			{Row: 3, Col: 1, Horizontal: true, Op: OpSub},
			{Row: 1, Col: 1, Horizontal: false, Op: OpAdd},
			{Row: 1, Col: 5, Horizontal: false, Op: OpAdd},
			{Row: 1, Col: 9, Horizontal: true, Op: OpMul},
			{Row: 3, Col: 9, Horizontal: true, Op: OpAdd},
			{Row: 3, Col: 11, Horizontal: false, Op: OpAdd},
			{Row: 7, Col: 1, Horizontal: true, Op: OpSub},
			{Row: 7, Col: 3, Horizontal: false, Op: OpMul},
			{Row: 11, Col: 1, Horizontal: true, Op: OpAdd},
			{Row: 13, Col: 1, Horizontal: true, Op: OpSub},
			{Row: 9, Col: 7, Horizontal: true, Op: OpDiv},
			{Row: 13, Col: 9, Horizontal: true, Op: OpAdd},
		},
	}
}

// Result bundles the solved puzzle, the printable grids, and the hidden mask.
type Result struct {
	Layout   *Puzzle
	Hidden   map[CellPosition]bool
	Puzzle   [][]string
	Solution [][]string
}

// isSolvable returns true if every hidden cell can be deduced by repeatedly
// applying the rule: if 2 of 3 cells in an equation are known, the 3rd is too.
func isSolvable(puzzle *Puzzle, hidden map[CellPosition]bool) bool {
	allCells := puzzle.UniqueCells()
	known := map[CellPosition]bool{}
	for i := 0; i < len(allCells); i++ {
		if !hidden[allCells[i]] {
			known[allCells[i]] = true
		}
	}
	for {
		progress := false
		for i := 0; i < len(puzzle.Equations); i++ {
			positions := puzzle.Equations[i].Cells()
			unknownCount := 0
			lastUnknown := CellPosition{}
			for j := 0; j < len(positions); j++ {
				if !known[positions[j]] {
					unknownCount++
					lastUnknown = positions[j]
				}
			}
			if unknownCount == 1 {
				known[lastUnknown] = true
				progress = true
			}
		}
		if !progress {
			break
		}
	}
	for i := 0; i < len(allCells); i++ {
		if !known[allCells[i]] {
			return false
		}
	}
	return true
}

// Generate solves the default layout, hides cells with probability hideProb,
// and returns a ready-to-render Result.
func Generate(rng *rand.Rand, hideProb float64) (*Result, error) {
	puzzle := BuildLayout()
	if !puzzle.Solve(rng) {
		return nil, errors.New("mathgen: failed to solve layout")
	}
	allCells := puzzle.UniqueCells()
	var hidden map[CellPosition]bool
	for attempt := 0; attempt < 200; attempt++ {
		hidden = map[CellPosition]bool{}
		for i := 0; i < len(allCells); i++ {
			if rng.Float64() < hideProb {
				hidden[allCells[i]] = true
			}
		}
		// ensure every equation has at least one hidden cell
		for i := 0; i < len(puzzle.Equations); i++ {
			positions := puzzle.Equations[i].Cells()
			hasHidden := false
			for j := 0; j < len(positions); j++ {
				if hidden[positions[j]] {
					hasHidden = true
					break
				}
			}
			if !hasHidden {
				hidden[positions[rng.Intn(3)]] = true
			}
		}
		if isSolvable(puzzle, hidden) {
			break
		}
	}
	// fallback: reveal cells one by one until solvable,
	// but never un-hide a cell that is the sole hidden cell in its equation
	if !isSolvable(puzzle, hidden) {
		order := rng.Perm(len(allCells))
		for i := 0; i < len(order); i++ {
			pos := allCells[order[i]]
			if !hidden[pos] {
				continue
			}
			isSoleHidden := false
			for j := 0; j < len(puzzle.Equations); j++ {
				positions := puzzle.Equations[j].Cells()
				involvedInEq := false
				for k := 0; k < len(positions); k++ {
					if positions[k] == pos {
						involvedInEq = true
						break
					}
				}
				if !involvedInEq {
					continue
				}
				hiddenCount := 0
				for k := 0; k < len(positions); k++ {
					if hidden[positions[k]] {
						hiddenCount++
					}
				}
				if hiddenCount == 1 {
					isSoleHidden = true
					break
				}
			}
			if isSoleHidden {
				continue
			}
			delete(hidden, pos)
			if isSolvable(puzzle, hidden) {
				break
			}
		}
	}
	return &Result{
		Layout:   puzzle,
		Hidden:   hidden,
		Puzzle:   puzzle.ToGrid(hidden),
		Solution: puzzle.ToGrid(nil),
	}, nil
}
