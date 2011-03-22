package main

import (
	"curses"
	"fmt"
)

func Beep() {
	curses.Beep()
}

type Window struct {
	Curses     *curses.Window
	Cols, Rows int
	gs         *GlobalState
}

func NewWindow(gs *GlobalState) *Window {
	w := new(Window)
	w.Curses = curses.Stdwin
	w.Cols = *curses.Cols
	w.Rows = *curses.Rows
	w.gs = gs
	return w
}

func (w *Window) InputRoutine(ch chan int) {
	go func() {
		for {
			ch <- w.Curses.Getch()
		}
	}()
}

func (w *Window) UpdateRoutine(ch chan int) {
	go func() {
		for {
			<-ch
			w.PaintModeliner(false)
			w.PaintMapper(0, w.Rows-1, true)
			w.Curses.Refresh()
		}
	}()
}

func (w *Window) PaintMapper(start, end int, paintCursor bool) {
	cols, rows := w.Cols, w.Rows-1

	gs := w.gs
	mapper := *gs.CurrentMapper

	if start < 0 || start > rows || end > rows {
		EndScreen()
		panic(fmt.Sprintf("Window.Paint: Bad range (%d, %d) [%d, %d]", start, end, cols, rows))
	}

	smap := mapper.GetMap()
	for i := start; i < end; i++ {
		w.Curses.Move(i, 0)
		w.Curses.Clrtoeol()
		w.Curses.Mvwaddnstr(i, 0, smap[i], cols)
	}

	cX, cY := mapper.GetCursor()
	if paintCursor {
		if cX < 0 || cY < 0 || cX > cols || cY > rows {
			EndScreen()
			panic(fmt.Sprintf("Window.Paint: Bad cursor (%d, %d) [%d, %d]", start, end, cols, rows))
		}
		w.Curses.Move(cY, cX)
	}
}

func (w *Window) PaintModeliner(paintCursor bool) {
	maxRow := w.Rows - 1
	gs := w.gs

	// XXX check for modeline until i have everything set up
	if gs.Modeline == nil {
		return
	}

	modeline := *gs.Modeline

	w.Curses.Move(maxRow, 0)
	w.Curses.Clrtoeol()
	// This needs hscroll
	w.Curses.Mvwaddnstr(maxRow, 0, modeline.String(), w.Cols)

	if paintCursor {
		w.Curses.Move(maxRow, modeline.GetCursor())
	}
}