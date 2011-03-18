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

type editBuffer struct {
	fi     *os.FileInfo
	name   string
	lines  *list.List
	l   *list.Element
	anchor *list.Element
	cache  map[int]string
	x, y int
}

func newEditBuffer(name string) *editBuffer {
	b := new(editBuffer)
	b.name = name
	b.lines = list.New()
	b.lines.Init()
	b.l = nil
	b.anchor = b.l
	b.cache = make(map[int]string)
	return b
}

func (b *editBuffer) insertChar(c byte) {
	if b.l != nil {
		b.l.Value.(*editLine).insertChar(c)
	} else {
		panic(NilLine)
	}
	b.mapToScreen()
}

func (b *editBuffer) mapToScreen() {
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

func (b *editBuffer) backspace() {
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

func (b *editBuffer) insertLine(e *editLine) *list.Element {
	if b.l == nil {
		b.l = b.lines.PushFront(e)
		b.anchor = b.l
	} else {
		b.l = b.lines.InsertAfter(e, b.l)
	}
	return b.l
}

func (b *editBuffer) appendLine() *list.Element {
	return b.insertLine(newEditLine([]byte("")))
}

func (b *editBuffer) deleteLine() {
}

func (b *editBuffer) newLine(d byte) {
	// XXX This is pretty wrong lol
	if b.l != nil {
		b.l.Value.(*editLine).insertChar(d)
		b.l = b.appendLine()
		b.mapToScreen()
	}
}

func (b *editBuffer) top() {
	b.l = b.lines.Front()
	b.anchor = b.l
}

func (b *editBuffer) moveLeft() {
	if !b.l.Value.(*editLine).moveCursor(-1) {
		Beep()
	}
	b.mapToScreen()
}

func (b *editBuffer) moveRight() {
	if !b.l.Value.(*editLine).moveCursor(1) {
		Beep()
	}
	b.mapToScreen()
}

func (b *editBuffer) moveUp() {
	if p := b.l.Prev(); p != nil {
		b.l = p
		b.mapToScreen()
	} else {
		Beep()
	}
}

func (b *editBuffer) moveDown() {
	if n := b.l.Next(); n != nil {
		b.l = n
		b.mapToScreen()
	} else {
		Beep()
	}
}

