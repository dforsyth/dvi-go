package main

import (
	// "container/list"
	"curses"
	"fmt"
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
	dirty    bool // This should become an int, so that updates are just after a given line

	tabs     bool
	tabwidth int
	tabstop  int

	cmdbuff *GapBuffer

	// Stuff for painting
	Anchor     int
	Window     *Window
	X, Y       int
	CurX, CurY int
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)
	eb.Pathname = name
	eb.Lines = make([]*EditLine, 0)
	eb.Line = 0
	eb.Column = 0
	eb.dirty = true

	eb.cmdbuff = NewGapBuffer([]byte(""))

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
		if eb.cmdbuff.String() == "" {
			b :=true
			switch k {
			case 'j':
				b = eb.MoveLeft()
			case 'k':
				b = eb.MoveDown()
			case 'l':
				b = eb.MoveUp()
			case ';':
				b = eb.MoveRight()
			case 'p':
				eb.PasteBelow()
			case 'P':
				eb.PasteAbove()
			case 'i':
				// Insert
			case 'a':
				// Append
				eb.MoveRight()
			case 'o':
				// Add a line and go to insert mode
				eb.AppendEmptyLine()
				eb.MoveDown()
			case 'd':
				eb.DeleteLine(eb.Line)
			}
			if !b {
				Beep()
			}
		} else {
			eb.cmdbuff.InsertChar(byte(k))
			eb.EvalCmdBuff()
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
	} else {
		eb.Line = lno - 1
	}
}

func (eb *EditBuffer) Backspace() {
	if l := eb.Lines[eb.Line]; l.Cursor() == 0 {
		if eb.Line > 0 {
			sav := eb.DeleteLine(eb.Line)
			eb.Lines[eb.Line].Delete(1)
			if sav != nil {
			}
		} else {
			Beep()
		}
	} else {
		l.Delete(1)
	}
}

// Insert a line at lno, 0-n, into an EditBuffer.
func (eb *EditBuffer) InsertLine(e *EditLine, lno int) {
	if lno < 0 || lno > len(eb.Lines) {
		panic(fmt.Sprintf("Unable to insert line at %d in buffer of %d lines", lno,
			len(eb.Lines)))
	}

	eb.Lines = append(eb.Lines[:lno], append([]*EditLine{e}, eb.Lines[lno:]...)...)
}

func (eb *EditBuffer) AppendEmptyLine() {
	eb.InsertLine(NewEditLine([]byte("")), eb.Line+1)
}

// Delete a line at lno, 0-n, from an EditBuffer
func (eb *EditBuffer) DeleteLine(lno int) *EditLine {
	// If we are removing the 0th line from a file with a single line,
	// after the line is removed, a new one needs to be inserted
	if len(eb.Lines) == 0 {
		// This is an error case, we're going to panic here because it
		// really should not happen
		panic("Trying to delete line 0 in a buffer with no lines")
	}

	// The line that's going away
	ln := eb.Lines[lno]

	eb.Lines = append(eb.Lines[:lno], eb.Lines[lno+1:]...)
	if len(eb.Lines) == 0 {
		eb.InsertLine(NewEditLine([]byte("")), 0)
		eb.Line = 0
		// vim would set "--no lines in buffer--" in this case
	} else if eb.Line > 0 {
		// Move up one line
		eb.Line -= 1
	}
	return ln
}

func (eb *EditBuffer) NewLine(d byte) {
	l := eb.Lines[eb.Line]
	l.InsertChar(d)
	newLine := NewEditLine(l.AfterCursor())
	l.ClearToEOL()
	eb.InsertLine(newLine, eb.Line+1)
	eb.MoveDown()
}

func (eb *EditBuffer) Top() {
	eb.Line = 0
	eb.Anchor = eb.Line
}

// TODO If the column is the length of a line, set b.Column to -1 so that moving
// vertically will put the cursor at the end of the new line.
func (eb *EditBuffer) MoveHorizontal(dir int) bool {
	if l := eb.Lines[eb.Line]; l.MoveCursor(l.Cursor() + dir) {
		eb.Column = l.Cursor()
		return true
	}
	return false
}

func (eb *EditBuffer) MoveLeft() bool {
	return eb.MoveHorizontal(-1)
}

func (eb *EditBuffer) MoveRight() bool {
	return eb.MoveHorizontal(1)
}

func (eb *EditBuffer) MoveVertical(dir int) bool {
	lidx := eb.Line + dir
	if lidx < 0 || lidx > len(eb.Lines) - 1 {
		return false
	}
	eb.Line = lidx
	if l := eb.Lines[eb.Line]; len(l.GetRaw()) > eb.Column {
		l.MoveCursor(eb.Column)
	}
	return true
}

func (eb *EditBuffer) MoveUp() bool {
	return eb.MoveVertical(-1)
}

func (eb *EditBuffer) MoveDown() bool {
	return eb.MoveVertical(1)
}

func (eb *EditBuffer) PasteAbove() {
}

func (eb *EditBuffer) PasteBelow() {
}

func (eb *EditBuffer) EvalCmdBuff() {
}
