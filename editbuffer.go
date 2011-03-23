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

type EditBuffer struct {
	fi         *os.FileInfo
	Name       string
	Lines      *list.List
	Line       *list.Element
	Anchor     *list.Element
	Window     *Window
	X, Y       int
	ScreenMap  []string
	CurX, CurY int
	Pathname   string
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)
	eb.Name = name
	eb.Lines = list.New()
	eb.Lines.Init()
	eb.Line = nil
	eb.Anchor = eb.Line

	eb.Window = gs.Window
	eb.ScreenMap = make([]string, eb.Window.Rows-1)
	eb.CurX, eb.CurY = 0, 0
	eb.X, eb.Y = eb.Window.Cols, eb.Window.Rows-1

	return eb
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
			eb.Backspace()
		} else if k == 0xd || k == 0xa {
			eb.NewLine(byte('\n'))
		} else {
			eb.InsertChar(byte(k))
		}
	}
}

func (eb *EditBuffer) RunRoutine(fn func(Interacter)) {
	go fn(eb)
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

func (eb *EditBuffer) InsertChar(c byte) {
	if eb.Line == nil {
		panic(NilLine)
	}

	eb.Line.Value.(*EditLine).InsertChar(c)
	eb.MapToScreen()
}

func (eb *EditBuffer) MapToScreen() {
	i := 0
	for l := eb.Anchor; l != nil && i < eb.Y; l = l.Next() {
		e := l.Value.(*EditLine)
		// XXX: screen Lines code for wrap
		row := make([]byte, eb.X)
		// panic(fmt.Sprintf("len of e.raw is %d", len(e.raw())))
		for i, _ := range row {
			row[i] = ' '
		}
		copy(row, e.raw())
		eb.ScreenMap[i] = string(row)
		if l == eb.Line {
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

func (eb *EditBuffer) Backspace() {
	if eb.Line == nil {
		panic(NilLine)
	}

	l := eb.Line.Value.(*EditLine)
	if l.b.gs == 0 {
		if prev := eb.Line.Prev(); prev != nil {
			eb.DeleteLine()
			eb.Line = prev
		} else {
			Beep()
		}
	} else {
		l.Delete(1)
	}
	eb.MapToScreen()
}

func (eb *EditBuffer) InsertLine(e *EditLine) *list.Element {
	if eb.Line == nil {
		eb.Line = eb.Lines.PushFront(e)
		eb.Anchor = eb.Line
	} else {
		eb.Line = eb.Lines.InsertAfter(e, eb.Line)
	}
	return eb.Line
}

func (eb *EditBuffer) AppendLine() *list.Element {
	return eb.InsertLine(NewEditLine([]byte("")))
}

func (eb *EditBuffer) DeleteLine() {
	eb.Lines.Remove(eb.Line)
}

func (eb *EditBuffer) NewLine(d byte) {
	// XXX This is pretty wrong lol
	if eb.Line != nil {
		eb.Line.Value.(*EditLine).InsertChar(d)
		eb.Line = eb.AppendLine()
		eb.MapToScreen()
	}
}

func (b *EditBuffer) top() {
	b.Line = b.Lines.Front()
	b.Anchor = b.Line
}

func (b *EditBuffer) moveLeft() {
	if !b.Line.Value.(*EditLine).moveCursor(-1) {
		Beep()
	}
	b.MapToScreen()
}

func (b *EditBuffer) moveRight() {
	if !b.Line.Value.(*EditLine).moveCursor(1) {
		Beep()
	}
	b.MapToScreen()
}

func (b *EditBuffer) moveUp() {
	if p := b.Line.Prev(); p != nil {
		b.Line = p
		b.MapToScreen()
	} else {
		Beep()
	}
}

func (b *EditBuffer) moveDown() {
	if n := b.Line.Next(); n != nil {
		b.Line = n
		b.MapToScreen()
	} else {
		Beep()
	}
}
