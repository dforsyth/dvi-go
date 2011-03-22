package main

import (
	"container/list"
	"curses"
	// "fmt"
	// "math"
	"os"
)

const (
	NilLine  = "nil line"
	CacheMax = 5
)

// XXX once a better drawing interface is figured out, there should be an lru
// cache of computed lines.  the way things are mapped now (rangeless) doesn't
// really let me cache the way i want to.

// XXX editbuffers are editable text buffers that happen to also be a screen.

type EditBuffer struct {
	fi     *os.FileInfo
	name   string
	lines  *list.List
	l      *list.Element
	anchor *list.Element
	x, y   int

	Window     *Window
	X, Y       int
	ScreenMap  []string
	CurX, CurY int
	Pathname   string
}

func (eb *EditBuffer) GetWindow() *Window {
	return eb.Window
}

func (eb *EditBuffer) SetWindow(w *Window) {
	eb.Window = w
	eb.X, eb.Y = w.Cols, w.Rows
}

func (eb *EditBuffer) SendInput(k int) {
	gs := eb.Window.gs
	switch gs.Mode {
	case INSERT:
		if k == curses.KEY_BACKSPACE || k == 127 {
			eb.backspace()
		} else if k == 0xd || k == 0xa {
			eb.newLine(byte('\n'))
		} else {
			eb.insertChar(byte(k))
		}
	}
}

func (eb *EditBuffer) GetMap() []string {
	return eb.ScreenMap
}

func (eb *EditBuffer) SetDimensions(x, y int) {
	eb.X, eb.Y = x, y
}

func (eb *EditBuffer) GetCursor() (int, int) {
	return eb.CurX, eb.CurY
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)
	eb.name = name
	eb.lines = list.New()
	eb.lines.Init()
	eb.l = nil
	eb.anchor = eb.l

	eb.Window = gs.Window
	eb.ScreenMap = make([]string, eb.Window.Rows-1)
	eb.CurX, eb.CurY = 0, 0
	eb.X, eb.Y = eb.Window.Cols, eb.Window.Rows-1

	return eb
}

func (b *EditBuffer) insertChar(c byte) {
	if b.l != nil {
		b.l.Value.(*EditLine).insertChar(c)
	} else {
		panic(NilLine)
	}
	b.MapToScreen()
}

func (eb *EditBuffer) MapToScreen() {
	i := 0
	for l := eb.anchor; l != nil && i < eb.Y; l = l.Next() {
		e := l.Value.(*EditLine)
		// XXX: screen lines code for wrap
		row := make([]byte, eb.X)
		// panic(fmt.Sprintf("len of e.raw is %d", len(e.raw())))
		for i, _ := range row {
			row[i] = ' '
		}
		copy(row, e.raw())
		eb.ScreenMap[i] = string(row)
		if l == eb.l {
			eb.CurY = i
			eb.CurX = e.b.gs
		}
		i++
	}
	for i < eb.Y {
		eb.ScreenMap[i] = NaL
		i++
	}
}

func (b *EditBuffer) backspace() {
	if b.l == nil {
		panic(NilLine)
	}

	l := b.l.Value.(*EditLine)
	if l.b.gs == 0 {
		if b.l.Prev() != nil {
			// XXX
		} else {
			Beep()
		}
	} else {
		l.delete(1)
	}
	b.MapToScreen()
}

func (b *EditBuffer) insertLine(e *EditLine) *list.Element {
	if b.l == nil {
		b.l = b.lines.PushFront(e)
		b.anchor = b.l
	} else {
		b.l = b.lines.InsertAfter(e, b.l)
	}
	return b.l
}

func (b *EditBuffer) appendLine() *list.Element {
	return b.insertLine(newEditLine([]byte("")))
}

func (b *EditBuffer) deleteLine() {
}

func (b *EditBuffer) newLine(d byte) {
	// XXX This is pretty wrong lol
	if b.l != nil {
		b.l.Value.(*EditLine).insertChar(d)
		b.l = b.appendLine()
		b.MapToScreen()
	}
}

func (b *EditBuffer) top() {
	b.l = b.lines.Front()
	b.anchor = b.l
}

func (b *EditBuffer) moveLeft() {
	if !b.l.Value.(*EditLine).moveCursor(-1) {
		Beep()
	}
	b.MapToScreen()
}

func (b *EditBuffer) moveRight() {
	if !b.l.Value.(*EditLine).moveCursor(1) {
		Beep()
	}
	b.MapToScreen()
}

func (b *EditBuffer) moveUp() {
	if p := b.l.Prev(); p != nil {
		b.l = p
		b.MapToScreen()
	} else {
		Beep()
	}
}

func (b *EditBuffer) moveDown() {
	if n := b.l.Next(); n != nil {
		b.l = n
		b.MapToScreen()
	} else {
		Beep()
	}
}
