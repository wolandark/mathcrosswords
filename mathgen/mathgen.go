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
	for _, eq := range puzzle.Equations {
		for _, pos := range eq.Cells() {
			seen[pos] = true
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
	for _, offset := range rng.Perm(puzzle.MaxVal) {
		puzzle.Values[pos] = offset + 1
		if puzzle.consistent(pos) && puzzle.backtrack(cells, index+1, rng) {
			return true
		}
	}
	delete(puzzle.Values, pos)
	return false
}

func (puzzle *Puzzle) consistent(changedPos CellPosition) bool {
	for _, eq := range puzzle.Equations {
		positions := eq.Cells()
		involves := false
		for _, pos := range positions {
			if pos == changedPos {
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
	for row := range grid {
		grid[row] = make([]string, puzzle.Size)
		for col := range grid[row] {
			grid[row][col] = "0"
		}
	}
	for _, eq := range puzzle.Equations {
		positions := eq.Cells()
		opPos := eq.OpCell()
		eqPos := eq.EqCell()
		grid[opPos.Row][opPos.Col] = string(eq.Op)
		grid[eqPos.Row][eqPos.Col] = "="
		for _, pos := range positions {
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
	known := map[CellPosition]bool{}
	for _, pos := range puzzle.UniqueCells() {
		if !hidden[pos] {
			known[pos] = true
		}
	}
	for {
		progress := false
		for _, eq := range puzzle.Equations {
			positions := eq.Cells()
			var unknowns []CellPosition
			for _, pos := range positions {
				if !known[pos] {
					unknowns = append(unknowns, pos)
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
	for _, pos := range puzzle.UniqueCells() {
		if !known[pos] {
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
		for _, pos := range allCells {
			if rng.Float64() < hideProb {
				hidden[pos] = true
			}
		}
		// ensure every equation has at least one hidden cell
		for _, eq := range puzzle.Equations {
			positions := eq.Cells()
			hasHidden := false
			for _, pos := range positions {
				if hidden[pos] {
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
		for _, index := range order {
			pos := allCells[index]
			if !hidden[pos] {
				continue
			}
			isSoleHidden := false
			for _, eq := range puzzle.Equations {
				positions := eq.Cells()
				involvedInEq := false
				for _, eqPos := range positions {
					if eqPos == pos {
						involvedInEq = true
						break
					}
				}
				if !involvedInEq {
					continue
				}
				hiddenCount := 0
				for _, eqPos := range positions {
					if hidden[eqPos] {
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
