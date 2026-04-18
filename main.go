package main

import (
	// "fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

var focusRow, focusCol int
var cells [][]*tview.TextView
var grid = GenGrid()
var ans = GenSolution()
var size = len(grid)
var app = tview.NewApplication()
var gridView = tview.NewGrid()
var insertMode bool

func moveFocusCell(app *tview.Application, cells [][]*tview.TextView, row, col int) {
	size := len(cells)

	if row >= 0 && row < size && col >= 0 && col < size {
		focusRow = row
		focusCol = col
		app.SetFocus(cells[row][col])
	}
}

func holdCells() [][]*tview.TextView {
	cells := make([][]*tview.TextView, size)
	for r := 0; r < size; r++ {
		cells[r] = make([]*tview.TextView, size)
	}
	return cells
}

func readCells() [][]string {
	out := make([][]string, size)
	for r := 0; r < size; r++ {
		out[r] = make([]string, size)
		for c := 0; c < size; c++ {
			out[r][c] = strings.TrimSpace(cells[r][c].GetText(true))
		}
	}
	return out
}

func compareSubAns(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if len(a[i]) != len(b[i]) {
			return false
		}

		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func test() bool {
	out := readCells()
	v := compareSubAns(out, ans)
	return v
}

func PopulateGrid(gridView *tview.Grid, grid [][]string) {
	size := len(grid)

	rows := make([]int, size)
	cols := make([]int, size)

	for i := 0; i < size; i++ {
		rows[i] = 2
		cols[i] = 6
	}

	gridView.SetRows(rows...)
	gridView.SetColumns(cols...)

	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			val := grid[r][c]

			cell := tview.NewTextView().
				SetTextAlign(tview.AlignCenter)

			switch val {
			case "0":
				cell.SetText("   ")
				cell.SetBackgroundColor(tcell.ColorGray)
			case "999":
				cell.SetText("     ")
				cell.SetBackgroundColor(tcell.ColorSteelBlue)
			default:
				cell.SetText(" " + val + " ")
				cell.SetBackgroundColor(tcell.ColorTeal)
				cell.SetTextColor(tcell.ColorBlack)
			}
			gridView.AddItem(cell, r, c, 1, 1, 0, 0, false)
			cells[r][c] = cell
		}
	}
}

func resetCellColor(r, c int) {
	switch grid[r][c] {
	case "0":
		cells[r][c].SetBackgroundColor(tcell.ColorGray)
	case "999":
		cells[r][c].SetBackgroundColor(tcell.ColorSteelBlue)
	default:
		cells[r][c].SetBackgroundColor(tcell.ColorTeal)
	}
}

func main() {
	cells = holdCells()

	PopulateGrid(gridView, grid)
	setupKeys(grid)

	if err := app.SetRoot(gridView, true).Run(); err != nil {
		panic(err)
	}
}

// ------------------------------------
// ------------------------------------
// Vim Keys

// gridView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 	switch event.Rune() {
// 	case 'k':
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow-1, focusCol)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)
// 	case 'j':
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow+1, focusCol)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)
// 	case 'h':
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow, focusCol-1)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)
// 	case 'l':
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow, focusCol+1)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)
// 	}
// 	return nil
// })

// ------------------------------------
// ------------------------------------
// Arrow Keys

// gridView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 	switch event.Key() {
// 	case tcell.KeyUp:
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow-1, focusCol)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

// 	case tcell.KeyDown:
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow+1, focusCol)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

// 	case tcell.KeyLeft:
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow, focusCol-1)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

// 	case tcell.KeyRight:
// 		oldRow, oldCol := focusRow, focusCol
// 		moveFocusCell(app, cells, focusRow, focusCol+1)
// 		resetCellColor(oldRow, oldCol)
// 		cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

// 	case tcell.KeyEnter:
// 		mkForm()

// 	case tcell.KeyEscape:
// 		gridView.Clear()
// 	}
// 	return nil
// })

// -------------------------------------
// -------------------------------------
// Make Form

// func mkForm() {
// 	form := tview.NewForm().
// 		AddInputField("?", "", 5, nil, nil)
// 	form.SetBorder(true)

// 	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		if event.Key() == tcell.KeyEscape {
// 			app.SetRoot(gridView, true)
// 			return nil
// 		}
// 		return event
// 	})
// 	formGrid := tview.NewGrid().
// 		SetRows(0, 10, 0).
// 		SetColumns(0, 15, 0).
// 		AddItem(form, 1, 1, 1, 1, 0, 0, true)

// 	app.SetRoot(formGrid, true)
// }
