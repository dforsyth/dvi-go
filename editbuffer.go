package main

import (
	"bufio"
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
	lno      int
	Column   int
	dirty    bool // This should become an int, so that updates are just after a given line

	// buffer settings
	tabs     bool
	tabwidth int
	tabstop  int

	// yank buffer
	yb []*EditLine

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
	eb.lno = 0
	eb.Column = 0
	eb.dirty = true

	eb.cmdbuff = NewGapBuffer([]byte(""))

	eb.Anchor = eb.lno
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
			eb.moveLeft()
		default:
			eb.InsertChar(byte(k))
		}
		eb.dirty = true
	case NORMAL:
		/* XXX
		if eb.cmdbuff.String() != "" {
			eb.cmdbuff.InsertChar(byte(k))
			eb.EvalCmdBuff()
			return
		}
		*/

		b := true
		switch k {
		case 'j':
			b = eb.moveLeft()
		case 'k':
			b = eb.moveDown()
		case 'l':
			b = eb.moveUp()
		case ';':
			b = eb.moveRight()
		case 'p':
			eb.paste(eb.lno + 1)
		case 'P':
			eb.paste(eb.lno)
		case 'i':
			// Insert
		case 'a':
			// Append
			eb.moveRight()
		case 'o':
			// Add a line and go to insert mode
			eb.AppendEmptyLine()
			eb.moveDown()
		case 'd':
			eb.delete(eb.lno)
		case 'y':
			eb.yank(eb.lno, 1)
		case 'G':
			eb.LastLine()
		}
		if !b {
			Beep()
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
	eb.Lines[eb.lno].InsertChar(c)
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
		if i == eb.lno {
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

func (eb *EditBuffer) MapLine(el *EditLine) {
	w := eb.X
	o := 0 // XXX offset.  0 until I add line number support
	if w < 0 || o != 0 {
		return
	}
}

func (eb *EditBuffer) GoToLine(lno int) {
	if lno < 1 {
		return
	}

	if lno > len(eb.Lines) {
		eb.lno = len(eb.Lines)
	} else {
		eb.lno = lno - 1
	}
}

func (eb *EditBuffer) Backspace() {
	if l := eb.Lines[eb.lno]; l.Cursor() == 0 {
		if eb.lno > 0 {
			sav := eb.delete(eb.lno)
			eb.Lines[eb.lno].Delete(1)
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
func (eb *EditBuffer) insert(e *EditLine, lno int) {
	if lno < 0 || lno > len(eb.Lines) {
		panic(fmt.Sprintf("Unable to insert line at %d in buffer of %d lines", lno,
			len(eb.Lines)))
	}

	eb.Lines = append(eb.Lines[:lno], append([]*EditLine{e}, eb.Lines[lno:]...)...)
}

func (eb *EditBuffer) AppendEmptyLine() {
	eb.insert(NewEditLine([]byte("")), eb.lno+1)
}

// Delete a line at lno, 0-n, from an EditBuffer
func (eb *EditBuffer) delete(lno int) *EditLine {
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
		eb.insert(NewEditLine([]byte("")), 0)
		eb.lno = 0
		// vim would set "--no lines in buffer--" in this case
	} else if eb.lno > 0 {
		// move up one line
		eb.lno -= 1
	}
	return ln
}

func (eb *EditBuffer) yank(lno, cnt int) int {
	if len(eb.Lines) == 0 {
		panic("cannot yank from an empty buffer")
	}

	max := lno + cnt
	if max > len(eb.Lines) {
		max = len(eb.Lines)
	}

	eb.yb = make([]*EditLine, max-lno)
	for i, ln := range eb.Lines[lno:max] {
		eb.yb[i] = NewEditLine(ln.GetRaw())
	}

	return max - lno
}

func (eb *EditBuffer) cut(lno, cnt int) *EditLine {

	if len(eb.Lines) == 0 {
		panic("Cannot cut line from empty buffer")
	}

	ln := eb.Lines[eb.lno]

	return ln
}

func (eb *EditBuffer) NewLine(d byte) {
	l := eb.Lines[eb.lno]
	l.InsertChar(d)
	newLine := NewEditLine(l.AfterCursor())
	l.ClearToEOL()
	eb.insert(newLine, eb.lno+1)
	eb.moveDown()
}

func (eb *EditBuffer) TopLine() {
	eb.lno = 0
	eb.Anchor = eb.lno
}

func (eb *EditBuffer) LastLine() {
	eb.lno = len(eb.Lines) - 1
	// XXX move anchor
}

// TODO If the column is the length of a line, set b.Column to -1 so that moving
// vertically will put the cursor at the end of the new line.
func (eb *EditBuffer) moveHorizontal(dir int) bool {
	if l := eb.Lines[eb.lno]; l.moveCursor(l.Cursor() + dir) {
		eb.Column = l.Cursor()
		return true
	}
	return false
}

func (eb *EditBuffer) moveLeft() bool {
	return eb.moveHorizontal(-1)
}

func (eb *EditBuffer) moveRight() bool {
	return eb.moveHorizontal(1)
}

func (eb *EditBuffer) moveVertical(dir int) bool {
	lno := eb.lno + dir
	if lno < 0 || lno > len(eb.Lines)-1 {
		return false
	}

	eb.lno = lno
	if l := eb.Lines[eb.lno]; len(l.GetRaw()) > eb.Column {
		l.moveCursor(eb.Column)
	}
	return true
}

func (eb *EditBuffer) moveUp() bool {
	return eb.moveVertical(-1)
}

func (eb *EditBuffer) moveDown() bool {
	return eb.moveVertical(1)
}

func (eb *EditBuffer) paste(lno int) {
	for _, ln := range eb.yb {
		eb.insert(NewEditLine(ln.GetRaw()), lno)
		lno++
	}
	eb.lno = lno
}

func (eb *EditBuffer) EvalCmdBuff() {
}

// Reads a file at pathname into EditBuffer eb
// Returns the number of lines read or error
func (eb *EditBuffer) readFile(f *os.File, mark int) (int, os.Error) {
	rdr := bufio.NewReader(f)
	// XXX fix this loop
	lno := mark
	for {
		if ln, err := rdr.ReadBytes('\n'); err == nil {
			eb.insert(NewEditLine(ln), lno)
		} else {
			if err != os.EOF {
				return -1, err
			} else {
				eb.insert(NewEditLine(ln), lno)
				return lno - mark, nil
			}
		}
		lno++
	}

	return lno - mark, nil
}
