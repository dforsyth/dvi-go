package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	// "strings"
)

const (
	NilLine  = "nil line"
	CacheMax = 5
)

type EditBuffer struct {
	gs *GlobalState

	fi       *os.FileInfo
	pathname string
	lines    []*EditLine
	lno      int
	col      int
	dirty    bool // This should become an int, so that updates are just after a given line
	temp     bool

	undoList list.List

	// buffer settings
	tabs     bool
	tabwidth int
	tabstop  int

	// yank buffer
	yb []*EditLine

	cmdbuff *GapBuffer

	// Stuff for painting
	redraw bool

	head int
	tail int

	X, Y       int
	CurX, CurY int
}

func NewEditBuffer(gs *GlobalState, name string) *EditBuffer {
	eb := new(EditBuffer)

	eb.gs = gs

	eb.pathname = name
	eb.lines = make([]*EditLine, 0)
	eb.lno = 0
	eb.col = 0
	eb.redraw = true
	eb.temp = false

	eb.cmdbuff = NewGapBuffer([]byte(""))

	eb.head = eb.lno
	eb.CurX, eb.CurY = 0, 0
	eb.X, eb.Y = eb.gs.Window.Cols, eb.gs.Window.Rows-1

	return eb
}

func (eb *EditBuffer) line() *EditLine {
	return eb.lines[eb.lno]
}

func (eb *EditBuffer) RunRoutine(fn func(Buffer)) {
	go fn(eb)
}

func (eb *EditBuffer) getWindow() *Window {
	return eb.gs.Window
}

func (eb *EditBuffer) mapScreen() {
	if eb.redraw {
		eb.MapToScreen()
		eb.redraw = false
	}
}

func (eb *EditBuffer) SetDimensions(x, y int) {
	eb.X, eb.Y = x, y
}

func (eb *EditBuffer) getCursor() (int, int) {
	return eb.CurX, eb.CurY
}

func (eb *EditBuffer) ident() string {
	return eb.pathname
}

func (eb *EditBuffer) insertChar(c byte) {
	eb.lines[eb.lno].insert(c)
}

func (eb *EditBuffer) isDirty() bool {
	return eb.dirty
}

func (eb *EditBuffer) screenLines(el *EditLine) int {
	raw := el.raw()
	l := len(raw)
	if l > 0 && raw[l-1] == '\n' {
		l--
	}

	// XXX It would be better if I was getting the screen width from the smap
	if sl := int(math.Ceil(float64(l) / float64(eb.gs.Window.Cols))); sl > 0 {
		return sl
	}
	return 1
}

func (eb *EditBuffer) MapToScreen() {
	var i int
	smap := eb.gs.Window.screenMap
	for lno, e := range eb.lines[eb.head:] {
		if i >= eb.Y {
			break
		}

		cnt := eb.screenLines(e)
		raw := e.raw()
		wrap := 0
		for lim := i + cnt; i < lim; i++ {
			beg := wrap * eb.X
			end := beg + eb.X
			s := ""
			if end < len(raw) {
				s = string(raw[beg:end])
			} else {
				end = len(raw)
				s = string(raw[beg:end])
			}
			if lno+eb.head == eb.lno && (e.Cursor() >= beg && e.Cursor() <= end) {
				eb.CurY = i
				eb.CurX = e.Cursor() - beg
				// lol.  it needs to automaticaly realize its at the end of a line
				// and puch the cursor to the next screen line.
				if eb.X == eb.CurX {
					eb.CurY++
					eb.CurX = 0
				}
			}
			smap[i] = s
			wrap++
		}
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
func (eb *EditBuffer) insertLn(e *EditLine, lno int) {
	if lno < 0 || lno > len(eb.lines) {
		panic(fmt.Sprintf("Unable to insert line at %d in buffer of %d lines", lno,
			len(eb.lines)))
	}

	eb.lines = append(eb.lines[:lno], append([]*EditLine{e}, eb.lines[lno:]...)...)
	e.eb = eb
}

func (eb *EditBuffer) insertEmptyLine(lno int) {
	eb.insertLn(NewEditLine([]byte("")), lno)
}

func (eb *EditBuffer) AppendEmptyLine() {
	eb.insertLn(NewEditLine([]byte("")), eb.lno+1)
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
		eb.insertLn(NewEditLine([]byte("")), 0)
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
		eb.yb[i] = NewEditLine(ln.raw())
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

	if lno < eb.head {
		eb.head = lno
	} else {
		dist := 0
		for _, e := range eb.lines[eb.head:lno] {
			dist += eb.screenLines(e)
			if dist > eb.Y-1 {
				s := eb.screenLines(eb.lines[eb.head])
				dist -= s
				eb.head++
			}
		}
	}

	eb.col = eb.lines[eb.lno].Cursor()
	eb.lno = lno
	if l := eb.lines[eb.lno]; len(l.raw()) > eb.col {
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
		eb.insertLn(NewEditLine(ln.raw()), lno)
		lno++
	}
	eb.lno = lno - 1
}

func (eb *EditBuffer) undo(c int) {
	eb.gs.queueMessage(&Message{
		"Got undo",
		false,
	})
}

// Reads a file at pathname into EditBuffer eb
// Returns the number of lines read or error
func (eb *EditBuffer) readFile(f *os.File, mark int) (int, os.Error) {
	rdr := bufio.NewReader(f)
	// XXX fix this loop
	lno := mark
	for {
		if ln, err := rdr.ReadBytes('\n'); err == nil {
			eb.insertLn(NewEditLine(ln), lno)
		} else {
			if err != os.EOF {
				return -1, err
			} else {
				if len(ln) > 0 {
					eb.insertLn(NewEditLine(ln), lno)
				}
				return lno - mark, nil
			}
		}
		lno++
	}

	return -1, nil // NOT REACHED
}

func (eb *EditBuffer) writeFile(f *os.File) (int, os.Error) {
	wb := 0
	for _, ln := range eb.lines {
		wl := ln.raw()
		if len(wl) > 0 && wl[len(wl)-1] != '\n' {
			wl = append(wl, '\n')
		}

		if b, e := f.Write(wl); e != nil {
			return 0, e
		} else {
			wb += b
		}
	}

	eb.gs.queueMessage(&Message{
		fmt.Sprintf("\"%s\"\t%dL\t%db", eb.pathname, len(eb.lines), wb),
		false,
	})

	return wb, nil
}
