package main

import (
	"curses"
	"math"
	"os"
)

func initscreen(s *Dvi) {
	curses.Initscr()
	curses.Cbreak()
	curses.Noecho()
	curses.Nonl()
	curses.Start_color()
	curses.Stdwin.Keypad(true)

	s.w = curses.Stdwin
	curses.Init_pair(1, curses.COLOR_WHITE, curses.COLOR_BLUE)
	curses.Init_pair(2, curses.COLOR_RED, curses.COLOR_WHITE)
	curses.Init_pair(3, curses.COLOR_BLACK, curses.COLOR_YELLOW)
}

func endscreen() {
	curses.Endwin()
}

func charlen(c byte) int {
	if c == '\t' {
		return 8
	}
	return 1
}

func screenlines(l *Line) int {
	if s := int(math.Ceil(float64(l.length()) / float64(*curses.Cols))); s > 0 {
		return s
	}
	return 1
}

func draw(d *Dvi) os.Error {

	f := d.b
	if f == nil {
		return nil
	}

	// this is some slow ass scroll action, ill figure out a way to avoid this later
	for l := f.first; l != f.disp && l != nil; l = l.next {
		if l == f.pos.line {
			for l != f.disp {
				f.disp = f.disp.prev
			}
			break
		}
	}
	for i, l := 0, f.disp; l != nil; i, l = i+screenlines(l), l.next {
		if l == f.pos.line {
			for i+screenlines(l) > *curses.Rows-1 {
				f.disp = f.disp.next
				i -= screenlines(l)
			}
			break
		}
	}

	cursory := 0
	cursorx := 0
	str := ""
	y := 0
	for l := f.disp; y < *curses.Rows-1 && l != nil; y, l = y+1, l.next {
		x := 0
		d.w.Move(y, 0)
		d.w.Clrtoeol()

		if l == f.pos.line {
			cursory = y
			for i, c := range l.text {
				if i >= f.pos.off {
					break
				}
				cursorx += charlen(c)
				if cursorx > *curses.Cols-1 {
					cursory++
					cursorx = 0
				}
			}
		}

		x = 0
		colors := int32(0)
		if f.pos.line == l {
			colors = curses.Color_pair(1)
		}
		for _, c := range l.text {
			if x > *curses.Cols-1 {
				y++
				x = 0
				d.w.Move(y, 0)
				d.w.Clrtoeol()
			}
			for i := x; i < x+charlen(c); i++ {
				d.w.Mvaddch(y, i, int32(c), colors)
			}
			x += charlen(c)
		}
		if colors != 0 {
			for ; x < *curses.Cols; x++ {
				d.w.Mvaddch(y, x, int32(' '), colors)
			}
		}
	}

	for ; y < *curses.Rows-1; y++ {
		d.w.Move(y, 0)
		d.w.Clrtoeol()
		d.w.Mvaddch(y, 0, int32('~'), 0)
	}

	d.currx = cursorx
	d.curry = cursory

	msg := ""
	mcolor := 3
	beep := false
	if d.msg == nil {
		msg = message(d) + " " + str
	} else {
		msg = ":" + d.msg.Message()
		mcolor = d.msg.Color()
		beep = d.msg.Beep()
		d.msg = nil
	}
	d.w.Move(*curses.Rows-1, 0)
	d.w.Clrtoeol()

	for i := 0; i < *curses.Cols && i < len(msg); i++ {
		d.w.Mvaddch(*curses.Rows-1, i, int32(msg[i]), curses.Color_pair(mcolor))
	}
	// s.w.Mvwaddnstr(*curses.Rows-1, 0, msg, *curses.Cols)

	d.w.Move(cursory, cursorx)
	if beep {
		curses.Beep()
	}
	d.w.Refresh()

	return nil
}
