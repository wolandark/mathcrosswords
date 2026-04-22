package main

import (
	// "fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type diff struct {
	Row, Col  int
	Got, Want string
}

var focusRow, focusCol int
var cells [][]*tview.TextView
var hint *tview.TextView
var sheetPanel *tview.TextView
var questions = GenGrid()
var size = len(questions)
var ans = GenSolution()
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

func compareSubAns(a, b [][]string) (bool, []diff) {
	if len(a) != len(b) {
		return false, nil
	}
	var diffs []diff
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false, nil
		}
		for j := range a[i] {
			if b[i][j] == "0" || b[i][j] == "999" {
				continue
			}
			if a[i][j] != b[i][j] {
				diffs = append(diffs, diff{i, j, a[i][j], b[i][j]})
			}
		}
	}
	return len(diffs) == 0, diffs
}

func answerSheet() []string {
	out := readCells()
	_, diffs := compareSubAns(out, ans)
	// fmt.Println("ok:", ok)

	var answersArr []string

	for _, d := range diffs {
		// fmt.Printf("(%d,%d) got=%q want=%q\n", d.Row, d.Col, d.Got, d.Want)
		// fmt.Printf("%q\n", d.Want)

		answersArr = append(answersArr, d.Want)
	}

	return answersArr
}

func checkAnswers() bool {
	out := readCells()
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if ans[row][col] == "0" {
				continue
			}
			if out[row][col] != ans[row][col] {
				return false
			}
		}
	}
	return true
}

func renderGrid(gridView *tview.Grid, questions [][]string) {
	size := len(questions)

	rows := make([]int, size)
	cols := make([]int, size)

	for i := 0; i < size; i++ {
		rows[i] = 2
		cols[i] = 6
	}

	rows = append(rows, 2, 2, 2) // gap row, then sheet row

	gridView.SetRows(rows...)
	gridView.SetColumns(cols...)

	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			val := questions[r][c]

			cell := tview.NewTextView().
				SetTextAlign(tview.AlignCenter)

			switch val {
			case "0":
				cell.SetText("   ")
				cell.SetBackgroundColor(bgColor)
			case "999":
				cell.SetText("     ")
				cell.SetBackgroundColor(questionColor)
			default:
				cell.SetText(" " + val + " ")
				cell.SetBackgroundColor(answerColor)
				cell.SetTextColor(fgColor)
			}
			gridView.AddItem(cell, r, c, 1, 1, 0, 0, true)
			cells[r][c] = cell
		}
	}

	hint = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(hintText)
	hint.SetText(` press enter to start typing. hjkl or arrows to move  
 S to check answers when done. if lost, choose from the answers below...`)
	gridView.AddItem(hint, size, 0, 1, size, 0, 0, false)

	// sheet row
	var sheet = answerSheet()
	for i := 0; i < size; i++ {

		a := ""
		if i < len(sheet) {
			a = sheet[i]
		}

		sheetPanel = tview.NewTextView()
		sheetPanel.SetTextColor(tcell.ColorWhite)
		sheetPanel.SetText(" " + a + " ")
		gridView.AddItem(sheetPanel, size+1, i, 1, 1, 0, 0, false)
	}

}

func main() {
	cells = holdCells()

	renderGrid(gridView, questions)

	gridSize := len(questions)

	cellH := 2
	cellW := 6

	gridHeight := gridSize*cellH + 4
	gridWidth := gridSize * cellW

	bordered := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(gridView, 0, 1, true)
	bordered.SetBorder(true).
		SetTitle("[ Math CrossWords ]").
		SetTitleColor(tcell.ColorBlack).
		SetBorderStyle(
			tcell.StyleDefault.
				Foreground(tcell.ColorBlack).
				Background(tcell.ColorTeal).
				Bold(true),
		)

	centered := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(bordered, gridHeight+2, 0, true).
				AddItem(nil, 0, 1, false),
			gridWidth+2, 0, true,
		).
		AddItem(nil, 0, 1, false)

	app.SetRoot(centered, true)

	setupKeys(questions)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
