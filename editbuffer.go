package main

import (
	"bufio"
	"curses"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	NilLine  = "nil line"
	CacheMax = 5
)

type EditBuffer struct {
	gs *GlobalState

	fi       *os.FileInfo
	Name     string
	Pathname string
	lines    []*EditLine
	lno      int
	col      int
	dirty    bool // This should become an int, so that updates are just after a given line

	// buffer settings
	tabs     bool
	tabwidth int
	tabstop  int

	// yank buffer
	yb []*EditLine

	cmdbuff *GapBuffer

	// Stuff for painting
	head int
	tail int

	X, Y       int
	CurX, CurY int
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)

	eb.gs = gs

	eb.Pathname = name
	eb.lines = make([]*EditLine, 0)
	eb.lno = 0
	eb.col = 0
	eb.dirty = true

	eb.cmdbuff = NewGapBuffer([]byte(""))

	eb.head = eb.lno
	eb.CurX, eb.CurY = 0, 0
	eb.X, eb.Y = eb.gs.Window.Cols, eb.gs.Window.Rows-1

	return eb
}

func (eb *EditBuffer) SendInput(k int) {
	gs := eb.gs
	switch gs.Mode {
	case INSERT:
		switch k {
		case curses.KEY_BACKSPACE, 127:
			eb.backspace()
		case 0xd, 0xa:
			eb.NewLine(byte('\n'))
		case ESC:
			eb.moveLeft()
		default:
			eb.lines[eb.lno].insertChar(byte(k))
		}
		eb.dirty = true
	case NORMAL:
		/* XXX
		if eb.cmdbuff.String() != "" {
			eb.cmdbuff.insertChar(byte(k))
			eb.EvalCmdBuff()
			return
		}
		*/

		b := true
		switch k {
		case 'j':
			b = eb.moveLeft()
		case 'k':
			b = eb.moveDown(1)
		case 'l':
			b = eb.moveUp(1)
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
			eb.moveDown(1)
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

func (eb *EditBuffer) RunRoutine(fn func(Buffer)) {
	go fn(eb)
}

func (eb *EditBuffer) getWindow() *Window {
	return eb.gs.Window
}

func (eb *EditBuffer) mapScreen() {
	if eb.dirty {
		eb.MapToScreen()
		eb.dirty = false
	}
}

func (eb *EditBuffer) SetDimensions(x, y int) {
	eb.X, eb.Y = x, y
}

func (eb *EditBuffer) getCursor() (int, int) {
	return eb.CurX, eb.CurY
}

func (eb *EditBuffer) insertChar(c byte) {
	eb.lines[eb.lno].insertChar(c)
}

func (eb *EditBuffer) screenLines(el *EditLine) int {
	raw := el.getRaw()
	l := len(raw)
	if l > 0 && raw[l-1] == '\n' {
		l--
	}

	// XXX It would be better if I was getting the screen width from the smap
	if sl := int(math.Ceil(float64(l)/float64(eb.gs.Window.Cols))); sl > 0 {
		return sl
	}
	return 1
}

func (eb *EditBuffer) MapToScreen() {
	var i int
	smap := eb.gs.Window.screenMap
	for _, e := range eb.lines[eb.head:] {
		if i >= eb.Y {
			break
		}
		// XXX: screen lines code for wrap
		rs := string(e.getRaw())
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

func (eb *EditBuffer) gotoLine(lno int) {
	if lno < 0 {
		return
	}

	if lno > len(eb.lines) {
		eb.lno = len(eb.lines)
	} else {
		eb.lno = lno - 1
	}
}

func (eb *EditBuffer) backspace() {
	if l := eb.lines[eb.lno]; l.Cursor() == 0 {
		if eb.lno > 0 {
			sav := eb.delete(eb.lno)
			eb.lines[eb.lno].Delete(1)
			// XXX This doesn't append the rest of the line that backspaced on in eb.lno...
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
	if lno < 0 || lno > len(eb.lines) {
		panic(fmt.Sprintf("Unable to insert line at %d in buffer of %d lines", lno,
			len(eb.lines)))
	}

	eb.lines = append(eb.lines[:lno], append([]*EditLine{e}, eb.lines[lno:]...)...)
}

func (eb *EditBuffer) AppendEmptyLine() {
	eb.insert(NewEditLine([]byte("")), eb.lno+1)
}

// Delete a line at lno, 0-n, from an EditBuffer
func (eb *EditBuffer) delete(lno int) *EditLine {
	// If we are removing the 0th line from a file with a single line,
	// after the line is removed, a new one needs to be inserted
	if len(eb.lines) == 0 {
		// This is an error case, we're going to panic here because it
		// really should not happen
		panic("Trying to delete line 0 in a buffer with no lines")
	}

	// The line that's going away
	ln := eb.lines[lno]

	eb.lines = append(eb.lines[:lno], eb.lines[lno+1:]...)
	if len(eb.lines) == 0 {
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
	if len(eb.lines) == 0 {
		panic("cannot yank from an empty buffer")
	}

	max := lno + cnt
	if max > len(eb.lines) {
		max = len(eb.lines)
	}

	eb.yb = make([]*EditLine, max-lno)
	for i, ln := range eb.lines[lno:max] {
		eb.yb[i] = NewEditLine(ln.getRaw())
	}

	return max - lno
}

func (eb *EditBuffer) cut(lno, cnt int) *EditLine {
	if len(eb.lines) == 0 {
		panic("Cannot cut line from empty buffer")
	}

	ln := eb.lines[eb.lno]

	return ln
}

func (eb *EditBuffer) NewLine(d byte) {
	l := eb.lines[eb.lno]
	l.insertChar(d)
	newLine := NewEditLine(l.AfterCursor())
	l.ClearToEOL()
	eb.insert(newLine, eb.lno+1)
	eb.moveDown(1)
}

func (eb *EditBuffer) TopLine() {
	eb.lno = 0
	eb.head = eb.lno
}

func (eb *EditBuffer) LastLine() {
	eb.lno = len(eb.lines) - 1
	// XXX move anchor
}

// TODO If the column is the length of a line, set b.Column to -1 so that moving
// vertically will put the cursor at the end of the new line.
func (eb *EditBuffer) moveHorizontal(dir int) bool {
	if l := eb.lines[eb.lno]; l.moveCursor(l.Cursor() + dir) {
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

// Move vertically within the buffer.
func (eb *EditBuffer) moveVertical(dir int) bool {
	lno := eb.lno + dir
	if lno < 0 || lno > len(eb.lines)-1 {
		return false
	}

	eb.col = eb.lines[eb.lno].Cursor()
	eb.lno = lno
	if l := eb.lines[eb.lno]; len(l.getRaw()) > eb.col {
		l.moveCursor(eb.col)
	}
	return true
}

func (eb *EditBuffer) moveUp(cnt int) bool {
	return eb.moveVertical(-cnt)
}

func (eb *EditBuffer) moveDown(cnt int) bool {
	return eb.moveVertical(cnt)
}

func (eb *EditBuffer) paste(lno int) {
	for _, ln := range eb.yb {
		eb.insert(NewEditLine(ln.getRaw()), lno)
		lno++
	}
	eb.lno = lno - 1
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

func (eb *EditBuffer) writeFile(f *os.File) (int, os.Error) {
	wb := 0
	for _, ln := range eb.lines {
		wl := ln.getRaw()
		if wl[len(wl)-1] != '\n' {
			wl = append(wl, '\n')
		}

		if b, e := f.Write(wl); e != nil {
			return 0, e
		} else {
			wb += b
		}
	}
	return wb, nil
}
