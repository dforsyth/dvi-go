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

	ln := eb.lines
	for i := 1; i < vw.rows-1; i++ {
		if ln != nil {
			vw.win.Mvwaddnstr(i, 0, string(ln.bytes()), vw.cols)
			ln = ln.next
		} else {
			vw.win.Mvwaddnstr(i, 0, NaL, vw.cols)
		}
	}

	vw.win.Mvwaddnstr(0, 0, eb.Title(), vw.cols)
	UpdateModeLine(ml)

	if eb.line != nil {
		DrawCursor()
	}
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
