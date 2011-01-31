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

	v := d.view

	v.win.Clear()

	ln := d.buf.lines
	for i := 1; i < v.rows-2; i++ {
		if ln != nil {
			v.win.Mvwaddnstr(i, 0, string(ln.bytes()), v.cols)
			ln = ln.next
		} else {
			v.win.Mvwaddnstr(i, 0, NaL, v.cols)
		}
	}

	v.win.Mvwaddnstr(0, 0, d.buf.Title(), v.cols)
	v.win.Mvwaddnstr(v.rows-2, 0, statusLine(), v.cols)

	UpdateMessageLine()

	if d.buf.line != nil {
		DrawCursor()
	}

	v.win.Refresh()
}

// update line l with str and refresh
func UpdateLine(l int, str string) {

	v := d.view

	v.win.Move(l, 0)
	v.win.Clrtoeol()
	v.win.Mvwaddnstr(l, 0, str, v.cols)
	v.win.Refresh()
}

func UpdateStatusLine() {
	d.view.win.Mvwaddnstr(d.view.rows-2, 0, statusLine(), d.view.cols)
	d.view.win.Refresh()
}

func UpdateMessageLine() {
	d.view.win.Mvwaddnstr(d.view.rows-1, 0, Message, d.view.cols)
	d.view.win.Refresh()
}

func (v *View) DrawLine(lno int, line *Line) {
	for _, c := range line.bytes() {
		switch c {
		case '\t':
		case '\n':
		default:
		}
	}
}

func DrawCursor() {
	d.view.win.Move(d.buf.lno, d.buf.line.cursor)
}
