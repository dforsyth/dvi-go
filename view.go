package main

import (
	"curses"
)

// A view displays a title bar (always), a text buffer (when available), and a
// message line (always).
type View struct {
	win        *curses.Window
	Cols, Rows int
	StartRow   uint
	Lines      []string
}

func NewView() *View {
	v := new(View)

	v.win = curses.Stdwin
	v.Rows, v.Cols = *curses.Rows, *curses.Cols
	v.Lines = make([]string, v.Rows)
	return v
}

type ScrLine struct {
	str string
	// XXX we can put color bolding information in here later...
}

func Beep() {
	curses.Beep()
}

func UpdateDisplay() {
	Vw.win.Clear()

	UpdateTitleLine()

	Eb.Map()
	xmax := Vw.Cols
	for i, row := range Vw.Lines {
		Vw.win.Move(i+1, 0)
		Vw.win.Clrtoeol()
		Vw.win.Mvwaddnstr(i+1, 0, row, xmax)
	}

	UpdateModeLine(Ml)

	if Eb.line != nil {
		DrawCursor()
	}
}


func UpdateLine(rno int, ln *Line) {
}

func UpdateLineAndAfter(rno, ln *Line) {
}

func UpdateTitleLine() {
	Vw.win.Move(0, 0)
	Vw.win.Clrtoeol()
	Vw.win.Mvwaddnstr(0, 0, Eb.Title(), Vw.Cols)
}

func UpdateModeLine(m Message) {
	l := Vw.Rows - 1
	Vw.win.Move(l, 0)
	Vw.win.Clrtoeol()
	Vw.win.Mvwaddnstr(l, 0, m.String(), Vw.Cols)
}

func DrawCursor() {
	x, y := Eb.CursorCoord()
	Vw.win.Move(y, x)
}
