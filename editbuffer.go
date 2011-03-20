package main

import (
	"container/list"
	// "fmt"
	// "math"
	"os"
)

const (
	NilLine = "nil line"
	CacheMax           = 5
)

// XXX once a better drawing interface is figured out, there should be an lru
// cache of computed lines.  the way things are mapped now (rangeless) doesn't
// really let me cache the way i want to.

// XXX editbuffers are editable text buffers that happen to also be a screen.

type EditBuffer struct {
	fi     *os.FileInfo
	name   string
	lines  *list.List
	l   *list.Element
	anchor *list.Element
	x, y int

	Window *Window
	X, Y int
	ScreenMap []string
	CurX, CurY int
}

func (eb *EditBuffer) GetWindow() *Window {
	return eb.Window
}

func (eb *EditBuffer) SetWindow(w *Window) {
	eb.Window = w
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

func newEditBuffer(name string) *EditBuffer {
	b := new(EditBuffer)
	b.name = name
	b.lines = list.New()
	b.lines.Init()
	b.l = nil
	b.anchor = b.l
	return b
}

func (b *EditBuffer) insertChar(c byte) {
	if b.l != nil {
		b.l.Value.(*editLine).insertChar(c)
	} else {
		panic(NilLine)
	}
	b.mapToScreen()
}

func (b *EditBuffer) mapToScreen() {
	i := 0
	for l := b.anchor; l != nil && i < screen.Rows - 1; l = l.Next() {
		e := l.Value.(*editLine)
		// XXX: screen lines code for wrap
		row := make([]byte, screen.Cols)
		// panic(fmt.Sprintf("len of e.raw is %d", len(e.raw())))
		for i, _ := range row {
			row[i] = ' '
		}
		copy(row, e.raw())
		screen.Lines[i] = string(row)
		if l == b.l {
			b.y = i
			b.x = e.b.gs
		}
		i++
	}
	for i < screen.Rows - 1 {
		screen.Lines[i] = NaL
		i++
	}
}

func (b *EditBuffer) backspace() {
	if b.l == nil {
		panic(NilLine)
	}

	l := b.l.Value.(*editLine)
	if (l.b.gs == 0) {
		if b.l.Prev() != nil {
			// XXX
		} else {
			Beep()
		}
	} else {
		l.delete(1)
	}
	b.mapToScreen()
}

func (b *EditBuffer) insertLine(e *editLine) *list.Element {
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
		b.l.Value.(*editLine).insertChar(d)
		b.l = b.appendLine()
		b.mapToScreen()
	}
}

func (b *EditBuffer) top() {
	b.l = b.lines.Front()
	b.anchor = b.l
}

func (b *EditBuffer) moveLeft() {
	if !b.l.Value.(*editLine).moveCursor(-1) {
		Beep()
	}
	b.mapToScreen()
}

func (b *EditBuffer) moveRight() {
	if !b.l.Value.(*editLine).moveCursor(1) {
		Beep()
	}
	b.mapToScreen()
}

func (b *EditBuffer) moveUp() {
	if p := b.l.Prev(); p != nil {
		b.l = p
		b.mapToScreen()
	} else {
		Beep()
	}
}

func (b *EditBuffer) moveDown() {
	if n := b.l.Next(); n != nil {
		b.l = n
		b.mapToScreen()
	} else {
		Beep()
	}
}

