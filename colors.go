package main

import (
	"github.com/gdamore/tcell/v2"
)

// Meh
// var answerColor = tcell.ColorTeal
// var questionColor = tcell.ColorSteelBlue
// var bgColor = tcell.ColorGray
// var cursorColor = tcell.ColorSilver

var answerColor = tcell.NewRGBColor(0, 0, 170)     // classic DOS
var questionColor = tcell.NewRGBColor(0, 170, 170) // cyan
var bgColor = tcell.NewRGBColor(0, 0, 85)          // dark blue
var fgColor = tcell.NewRGBColor(255, 255, 255)     // white
var hintText = tcell.NewRGBColor(255, 255, 255)    // white
var sheetBgColor = tcell.ColorLime
var sheetFgColor = tcell.ColorWhite
var cursorColor = tcell.NewRGBColor(255, 255, 85) // bright

type palette struct {
	answer,
	question,
	bg,
	fg,
	hint,
	sheetbg,
	sheetfg,
	cursor tcell.Color
}

var palettes = []palette{
	{ // Teal
		tcell.ColorTeal,
		tcell.ColorSteelBlue,
		tcell.ColorGray,
		tcell.ColorBlack,
		tcell.ColorWheat,
		tcell.ColorViolet,
		tcell.ColorWheat,
		tcell.ColorSilver},
	{ // Mocha
		tcell.NewRGBColor(180, 190, 254),
		tcell.NewRGBColor(137, 220, 235),
		tcell.NewRGBColor(49, 50, 68),
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(255, 255, 255),
		tcell.ColorPurple,
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(249, 226, 175),
	},
	{ // Nord
		tcell.NewRGBColor(136, 192, 208),
		tcell.NewRGBColor(143, 188, 187),
		tcell.NewRGBColor(59, 66, 82),
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(255, 255, 255),
		tcell.ColorGray,
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(235, 203, 139),
	},
	{ // Tokyo
		tcell.NewRGBColor(125, 207, 255),
		tcell.NewRGBColor(158, 206, 106),
		tcell.NewRGBColor(36, 40, 59),
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(255, 255, 255),
		tcell.ColorDarkMagenta,
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(255, 158, 100),
	},
	{ // Everforest
		tcell.NewRGBColor(167, 192, 128),
		tcell.NewRGBColor(219, 188, 127),
		tcell.NewRGBColor(61, 72, 77),
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(255, 255, 255),
		tcell.ColorLightGreen,
		tcell.NewRGBColor(0, 0, 0),
		tcell.NewRGBColor(230, 126, 128)},
}

var paletteIdx int

func toggleColors() {
	paletteIdx = (paletteIdx + 1) % len(palettes)
	p := palettes[paletteIdx]
	answerColor, questionColor, bgColor, fgColor, hintText, sheetBgColor, sheetFgColor, cursorColor = p.answer,
		p.question, p.bg, p.fg, p.hint, p.sheetbg, p.sheetfg, p.cursor

	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			resetCellColor(r, c)
		}
	}
	cells[focusRow][focusCol].SetBackgroundColor(cursorColor)
	hint.SetTextColor(hintText)
	for _, p := range sheetPanel {
		p.SetBackgroundColor(sheetBgColor)
		p.SetTextColor(sheetFgColor)
	}

}

func resetCellColor(r, c int) {
	cells[r][c].SetTextColor(fgColor)
	switch questions[r][c] {
	case "0":
		cells[r][c].SetBackgroundColor(bgColor)
	case "999":
		cells[r][c].SetBackgroundColor(questionColor)
	default:
		cells[r][c].SetBackgroundColor(answerColor)
	}
}
