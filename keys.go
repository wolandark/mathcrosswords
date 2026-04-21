package main

import (
	// "bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	// "os"
)

func setupKeys([][]string) {
	gridView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyUp || event.Rune() == 'k':
			oldRow, oldCol := focusRow, focusCol
			moveFocusCell(app, cells, focusRow-1, focusCol)
			resetCellColor(oldRow, oldCol)
			cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

		case event.Key() == tcell.KeyDown || event.Rune() == 'j':
			oldRow, oldCol := focusRow, focusCol
			moveFocusCell(app, cells, focusRow+1, focusCol)
			resetCellColor(oldRow, oldCol)
			cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

		case event.Key() == tcell.KeyLeft || event.Rune() == 'h':
			oldRow, oldCol := focusRow, focusCol
			moveFocusCell(app, cells, focusRow, focusCol-1)
			resetCellColor(oldRow, oldCol)
			cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

		case event.Key() == tcell.KeyRight || event.Rune() == 'l':
			oldRow, oldCol := focusRow, focusCol
			moveFocusCell(app, cells, focusRow, focusCol+1)
			resetCellColor(oldRow, oldCol)
			cells[focusRow][focusCol].SetBackgroundColor(tcell.ColorBlack)

		case event.Key() == tcell.KeyEnter:
			insertMode = true
			return nil

		case insertMode && event.Key() == tcell.KeyEsc:
			insertMode = false
			return nil
			// mkForm()

		// case insertMode && event.Rune() != 0:
		case insertMode && event.Rune() >= '0' && event.Rune() <= '9':
			cell := cells[focusRow][focusCol]
			cell.SetText(cell.GetText(false) + string(event.Rune()))
			return nil

		case event.Rune() == 'S': //Submti
			res := checkAnswers()
			fmt.Println(res)

		case insertMode && event.Key() == tcell.KeyBackspace:
			cell := cells[focusRow][focusCol]
			t := cell.GetText(false)
			if len(t) > 0 {
				cell.SetText(t[:len(t)-1])
			}
			return nil

		case event.Rune() == 'q':
			// app.Suspend(func() {
			// sheet := answerSheet()
			// fmt.Println(sheet)
			// fmt.Println("Press Enter To Continue...")
			// bufio.NewReader(os.Stdin).ReadString('\n')
			// })
			app.Stop()
		}

		return nil
	})
}
