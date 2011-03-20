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
	update     chan int
}

func (s *Screen) ScreenRoutine() {
	for {
		<-s.update
		s.Window.Refresh()
	}
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
