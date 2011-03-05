package main

import (
	"curses"
)

// Message lines
type Message interface {
	String() string
}

// A view displays a title bar (always), a text buffer (when available), and a
// message line (always).
type Screen struct {
	Window     *curses.Window
	Cols, Rows int
	StartRow   uint
	Lines      []string
	msg        Message
	update	chan int
}

func (s *Screen) ScreenRoutine() {
	go func() {
		<-s.update
		/*
		switch t := upd.(type) {
		case []string:
			s.Lines = t
			s.RedrawAfter(0)
		case Message:
			s.msg = t
			s.RedrawMessage()
		default:
		}
		//panic("i made it!")
		s.RedrawAfter(0)
		s.RedrawMessage()
		*/
		s.Window.Refresh()
	}()
}

func NewScreen(window *curses.Window) *Screen {
	v := new(Screen)

	v.Window = window
	v.Rows, v.Cols = *curses.Rows, *curses.Cols
	v.Lines = make([]string, v.Rows-1)
	v.update = make(chan int)
	return v
}

func Beep() {
	curses.Beep()
}

func (scr *Screen) RedrawRange(s, e int) {
	for i := s; i < e; i++ {
		scr.Window.Move(i, 0)
		scr.Window.Clrtoeol()
		scr.Window.Mvwaddnstr(i, 0, scr.Lines[i], scr.Cols)
	}
	if curr.line != nil {
		DrawCursor()
	}
}

func (scr *Screen) RedrawAfter(r int) {
	scr.RedrawRange(r, scr.Rows-1)
}

func (scr *Screen) RedrawMessage() {
	scr.Window.Move(scr.Rows-1, 0)
	scr.Window.Clrtoeol()
	scr.Window.Mvwaddnstr(scr.Rows-1, 0, scr.msg.String(), scr.Cols)
}

func (scr *Screen) SetMessage(m Message) {
	scr.msg = m
}

func (scr *Screen) RedrawCursor(y, x int) {
	scr.Window.Move(y, x)
}

func UpdateModeLine(m Message) {
	l := screen.Rows - 1
	screen.Window.Move(l, 0)
	screen.Window.Clrtoeol()
	screen.Window.Mvwaddnstr(l, 0, m.String(), screen.Cols)
}

func DrawCursor() {
	x, y := curr.CursorCoord()
	screen.Window.Move(y, x)
}
