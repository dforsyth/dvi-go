package main

import (
	// "container/list"
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
	Lines    []*EditLine
	Line     int
	Column   int
	dirty    bool

	tabs     bool
	tabwidth int
	tabstop  int

	// Stuff for painting
	Anchor     int
	Window     *Window
	X, Y       int
	CurX, CurY int
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)
	eb.Pathname = name
	eb.Lines = make([]*EditLine, 0, 100)
	eb.Line = 0
	eb.Column = 0
	eb.dirty = true

	eb.Anchor = eb.Line
	eb.Window = gs.Window
	// eb.ScreenMap = make([]string, eb.Window.Rows-1)
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
		case ESC:
			eb.MoveLeft()
		default:
			eb.InsertChar(byte(k))
		}
		eb.dirty = true
	case NORMAL:
		switch k {
		case 'j':
			eb.MoveLeft()
		case 'k':
			eb.MoveDown()
		case 'l':
			eb.MoveUp()
		case ';':
			eb.MoveRight()
		case 'p':
			eb.PasteBelow()
		case 'P':
			eb.PasteAbove()
		case 'i':
			// Insert
		case 'a':
			// Append
			eb.MoveRight()
		}
		// XXX Until I fix mapping, mark the whole buffer as dirty on movement
		eb.dirty = true
	case COMMAND: // XXX How did you get here?
	}
}

func (eb *EditBuffer) RunRoutine(fn func(Interacter)) {
	go fn(eb)
}

func (eb *EditBuffer) GetMap() *[]string {
	if eb.dirty {
		eb.MapToScreen()
		eb.dirty = false
	}
	return eb.Window.ScreenMap
}

func (eb *EditBuffer) SetDimensions(x, y int) {
	eb.X, eb.Y = x, y
}

func (eb *EditBuffer) GetCursor() (int, int) {
	return eb.CurX, eb.CurY
}

func (eb *EditBuffer) InsertChar(c byte) {
	eb.Lines[eb.Line].InsertChar(c)
}

func (eb *EditBuffer) MapToScreen() {
	var i int
	smap := *eb.Window.ScreenMap
	for _, e := range eb.Lines[eb.Anchor:] {
		if i >= eb.Y {
			break
		}
		// XXX: screen Lines code for wrap
		row := make([]byte, eb.X)
		// panic(fmt.Sprintf("len of e.raw is %d", len(e.raw())))
		for j, _ := range row {
			row[j] = ' '
		}
		copy(row, e.GetRaw())
		rs := string(row)
		// XXX this is all sorts of wrong, but need to fix line mapping before fixing
		// this
		t := strings.Count(rs, "\t")
		s := strings.Replace(rs, "\t", "        ", -1)
		s = strings.Replace(s, "\n", "", -1)
		smap[i] = s
		if i == eb.Line {
			eb.CurY = i
			eb.CurX = e.b.gs + (t * 7)
		}
		i++
	}
	for i < eb.Y {
		smap[i] = NaL
		i++
	}
}

func (eb *EditBuffer) GoToLine(lno int) {
	if lno < 1 {
		return
	}

	if lno > len(eb.Lines) {
		eb.Line = len(eb.Lines)
	}
	eb.Line = lno - 1
}

func (eb *EditBuffer) Backspace() {
	l := eb.Lines[eb.Line]
	if l.Cursor() == 0 {
		if eb.Line > 0 {
			sav := eb.Lines[eb.Line]
			eb.DeleteLine()
			eb.Lines[eb.Line].Delete(1)
			// There *should* be a newline character to delete from prev
			if sav != nil {
			}
		} else {
			Beep()
		}
	} else {
		l.Delete(1)
	}
}

func (eb *EditBuffer) InsertLine(e *EditLine) *EditLine {
	if len(eb.Lines) > 0 {
		eb.Line += 1
	}
	eb.Lines = append(eb.Lines[:eb.Line], append([]*EditLine{e}, eb.Lines[eb.Line:]...)...)
	return eb.Lines[eb.Line]
}

func (eb *EditBuffer) AppendEmptyLine() *EditLine {
	return eb.InsertLine(NewEditLine([]byte("")))
}

func (eb *EditBuffer) DeleteLine() {
	eb.Lines = append(eb.Lines[:eb.Line], eb.Lines[eb.Line+1:]...)
	if eb.Line > 0 {
		eb.Line -= 1
	}
}

func (eb *EditBuffer) NewLine(d byte) {
	l := eb.Lines[eb.Line]
	l.InsertChar(d)
	newLine := l.AfterCursor()
	l.ClearToEOL()
	eb.InsertLine(NewEditLine(newLine))
}

func (eb *EditBuffer) Top() {
	eb.Line = 0
	eb.Anchor = eb.Line
}

// TODO If the column is the length of a line, set b.Column to -1 so that moving
// vertically will put the cursor at the end of the new line.
func (eb *EditBuffer) MoveHorizontal(dir int) {
	if l := eb.Lines[eb.Line]; !l.MoveCursor(l.Cursor() + dir) {
		Beep()
	} else {
		eb.Column = l.Cursor()
	}
}

func (eb *EditBuffer) MoveLeft() {
	eb.MoveHorizontal(-1)
}

func (eb *EditBuffer) MoveRight() {
	eb.MoveHorizontal(1)
}

func (b *EditBuffer) MoveUp() {
	if b.Line > 0 {
		b.Line -= 1
		if l := b.Lines[b.Line]; len(l.GetRaw()) > b.Column {
			l.MoveCursor(b.Column)
		}
	} else {
		Beep()
	}
}

func (b *EditBuffer) MoveDown() {
	if b.Line < len(b.Lines)-1 {
		b.Line += 1
		if l := b.Lines[b.Line]; len(l.GetRaw()) > b.Column {
			l.MoveCursor(b.Column)
		}
	} else {
		Beep()
	}
}

func (eb *EditBuffer) PasteAbove() {
}

func (eb *EditBuffer) PasteBelow() {
}
