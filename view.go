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
	vw.win.Clear()

	UpdateTitleLine()

	ln := eb.lines
	for i := 1; i < vw.rows-1; i++ {
		if ln != nil {
			for j, c := range ln.bytes() {
				vw.win.Addch(i, j, int32(c), 0)
			}
			ln = ln.next
		} else {
			vw.win.Mvwaddnstr(i, 0, NaL, vw.cols)
		}
	}

	UpdateModeLine(ml)

	if eb.line != nil {
		DrawCursor()
	}
}

func UpdateLine(rno int, ln *Line) {
}

func UpdateLineAndAfter(rno, ln *Line) {
}

func UpdateTitleLine() {
	vw.win.Move(0, 0)
	vw.win.Clrtoeol()
	vw.win.Mvwaddnstr(0, 0, eb.Title(), vw.cols)
}

func UpdateModeLine(m Message) {
	l := vw.rows-1
	vw.win.Move(l, 0)
	vw.win.Clrtoeol()
	vw.win.Mvwaddnstr(l, 0, m.String(), vw.cols)
}

func DrawCursor() {
	vw.win.Move(eb.lno, eb.line.cursor)
}
