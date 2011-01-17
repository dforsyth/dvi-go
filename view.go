package main

import (
	"curses"
)

// buffer view

type View struct{}

func Beep() {
	curses.Beep()
}

func (d *D) UpdateDisplay() {

	win.Clear()

	ln := d.Buffer().Lines().Front()
	for i := d.s; i < d.e; i++ {
		if ln != nil {
			win.Mvwaddnstr(i, 0, ln.Value.(*GapBuffer).String(), *curses.Cols)
			ln = ln.Next()
		} else {
			win.Mvwaddnstr(i, 0, NaL, *curses.Cols)
		}
	}

	win.Mvwaddnstr(0, 0, d.Buffer().Title(), *curses.Cols)
	win.Mvwaddnstr(d.e, 0, d.StatusLine(), *curses.Cols)

	if d.Buffer().Line() != nil {
		win.Move(0, d.Buffer().Line().c)
	}

	win.Refresh()
}

// update line l with str and refresh
func UpdateLine(l int, str string) {
	win.Move(l, 0)
	win.Clrtoeol()
	win.Mvwaddnstr(l, 0, str, *curses.Cols)
	win.Refresh()
}
