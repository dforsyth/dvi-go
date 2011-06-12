package main

import (
	"curses"
	"os"
)

func initscreen(s *State) {
	curses.Initscr()
	curses.Cbreak()
	curses.Noecho()
	curses.Nonl()
	curses.Start_color()
	curses.Stdwin.Keypad(true)

	s.w = curses.Stdwin
	curses.Init_pair(1, curses.COLOR_WHITE, curses.COLOR_BLUE)
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
	// TODO this fn
	return 1
}

func draw(s *State) os.Error {

	f := s.f
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
	for y, l := 0, f.disp; y < *curses.Rows-1 && l != nil; y, l = y+1, l.next {
		x := 0
		s.w.Move(y, 0)
		s.w.Clrtoeol()

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
				s.w.Move(y, 0)
				s.w.Clrtoeol()
			}
			for i := x; i < x+charlen(c); i++ {
				s.w.Mvwaddch(y, i, int32(c), colors)
			}
			x += charlen(c)
		}
		if colors != 0 {
			for ; x < *curses.Cols; x++ {
				s.w.Mvwaddch(y, x, int32(' '), colors)
			}
		}
	}

	s.currx = cursorx
	s.curry = cursory

	msg := ""
	if s.msg == nil {
		msg = message(s) + " " + str
	} else {
		msg = ":" + string(*s.msg)
	}
	s.w.Move(*curses.Rows-1, 0)
	s.w.Clrtoeol()
	s.w.Mvwaddnstr(*curses.Rows-1, 0, msg, *curses.Cols)

	s.w.Move(cursory, cursorx)
	s.w.Refresh()

	return nil
}
