package main

import (
	"container/list"
	"curses"
	// "fmt"
	// "math"
	"os"
	"strings"
)

const (
	NilLine  = "nil line"
	CacheMax = 5
)

type EditBuffer struct {
	fi       *os.FileInfo
	Name     string
	Pathname string
	Lines    *list.List
	Line     *list.Element
	Column   int
	// Stuff for painting
	Anchor     *list.Element
	Window     *Window
	X, Y       int
	ScreenMap  []string
	CurX, CurY int
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)
	eb.Pathname = name
	eb.Lines = list.New()
	eb.Lines.Init()
	eb.Line = nil
	eb.Column = 0

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
		switch k {
		case curses.KEY_BACKSPACE, 127:
			eb.Backspace()
		case 0xd, 0xa:
			eb.NewLine(byte('\n'))
		default:
			eb.InsertChar(byte(k))
		}
	case NORMAL:
		switch k {
		case 'j':
			eb.MoveLeft()
		case 'k':
			eb.MoveUp()
		case 'l':
			eb.MoveDown()
		case ';':
			eb.MoveRight()
		}
	case COMMAND: // XXX How did you get here?
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
		rs := string(row)
		t := strings.Count(rs, "\t")
		s := strings.Replace(rs, "\t", "        ", -1)
		eb.ScreenMap[i] = s
		if l == eb.Line {
			eb.CurY = i
			eb.CurX = e.b.gs + (t * 7)
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

func (eb *EditBuffer) Top() {
	eb.Line = eb.Lines.Front()
	eb.Anchor = eb.Line
}

// TODO If the column is the length of a line, set b.Column to -1 so that moving
// vertically will put the cursor at the end of the new line.
func (eb *EditBuffer) MoveHorizontal(dir int) {
	if l := eb.Line.Value.(*EditLine); !l.MoveCursor(l.Cursor() + dir) {
		Beep()
	} else {
		eb.Column = l.Cursor()
		eb.MapToScreen()
	}
}

func (eb *EditBuffer) MoveLeft() {
	eb.MoveHorizontal(-1)
}

func (eb *EditBuffer) MoveRight() {
	eb.MoveHorizontal(1)
}

func (b *EditBuffer) MoveUp() {
	if p := b.Line.Prev(); p != nil {
		b.Line = p
		if l := b.Line.Value.(*EditLine); len(l.raw()) > b.Column {
			l.MoveCursor(b.Column)
		}
		b.MapToScreen()
	} else {
		Beep()
	}
}

func (b *EditBuffer) MoveDown() {
	if n := b.Line.Next(); n != nil {
		b.Line = n
		if l := b.Line.Value.(*EditLine); len(l.raw()) > b.Column {
			l.MoveCursor(b.Column)
		}
		b.MapToScreen()
	} else {
		Beep()
	}
}
