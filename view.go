package main

import (
	"curses"
)

// buffer view

type View struct {
	win        *curses.Window
	cols, rows int
}

func Beep() {
	curses.Beep()
}

func UpdateDisplay() {
	Vw.win.Clear()

	UpdateTitleLine()

	ln := Eb.lines
	for i := 1; i < Vw.rows-1; i++ {
		if ln != nil {
			for j, c := range ln.raw() {
				Vw.win.Addch(i, j, int32(c), 0)
			}
			ln = ln.next
		} else {
			Vw.win.Mvwaddnstr(i, 0, NaL, Vw.cols)
		}
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
	Vw.win.Mvwaddnstr(0, 0, Eb.Title(), Vw.cols)
}

func UpdateModeLine(m Message) {
	l := Vw.rows - 1
	Vw.win.Move(l, 0)
	Vw.win.Clrtoeol()
	Vw.win.Mvwaddnstr(l, 0, m.String(), Vw.cols)
}

func DrawCursor() {
	Vw.win.Move(Eb.lno, Eb.line.cursor)
}
