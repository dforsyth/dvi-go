package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
)

const (
	NilLinePanicString = "current EditBuffer line is nil"
	CacheMax           = 5
)

// XXX once a better drawing interface is figured out, there should be an lru
// cache of computed lines.  the way things are mapped now (rangeless) doesn't
// really let me cache the way i want to.

// XXX editbuffers are editable text buffers that happen to also be a screen.

// todo: lockable
type EditBuffer struct {
	fi     *os.FileInfo /* EditBuffer info. */
	name   string
	dirty  bool
	rdonly bool
	lines  *list.List    /* List of lines in the file. */
	line   *list.Element /* Current line in the file. */
	edits  *list.List    /* Edit history. */

	anchor         *list.Element // first line to draw
	view           *Screen       // view this buffer draws to
	curs_x, curs_y int           // cursor position
}


func NewEditBuffer(name string, view *Screen) *EditBuffer {

	f := new(EditBuffer)
	f.name = name
	// XXX it would be better to init this to nil and add the line from another
	// place so that the file loader can use this
	f.lines = list.New()
	// b.InsertLine(NewLine([]byte("")))
	f.anchor = f.line
	f.dirty = false

	f.view = view
	return f
}

func (b *EditBuffer) Insert(ch byte) {
	if l, ok := b.line.Value.(*EditLine); ok {
		l.insertCharacter(ch)
	}
	b.dirty = true
	for b.ScreenRange(b.anchor, b.line) > b.view.Rows-1 && b.anchor != b.line {
		b.anchor = b.anchor.Next()
	}
	b.Map()
}

func (b *EditBuffer) BackSpace() {
	if b.line == nil {
		Beep()
		return
	}

	if l, ok := b.line.Value.(*EditLine); ok {
		if (l.cursor == 0 && !l.hasNewLine) || (l.cursor == 1 && l.hasNewLine) {
			if l.DisplayLength() != 0 && b.line.Prev() != nil {
				// combine this line and the previous
			} else {

				if b.line.Prev() != nil {
					b.DeleteCurrLine()
				} else {
					Beep()
				}
			}
		} else {
			l.backspace()
		}
	}
	for b.ScreenRange(b.anchor, b.line) > b.view.Rows-1 && b.anchor != b.line {
		b.anchor = b.anchor.Next()
	}
	b.Map()
}

func (b *EditBuffer) MoveLeft() {
	b.MoveCursor(-1)
}

func (b *EditBuffer) MoveRight() {
	b.MoveCursor(1)
}

func (b *EditBuffer) MoveCursor(d int) {
	if l, ok := b.line.Value.(*EditLine); ok && l.moveCursor(l.cursor+d) < 0 {
		Beep()
	}
}

func (b *EditBuffer) ScreenRange(a, l *list.Element) int {
	// aln, bln := a.Value.(*EditLine), l.Value.(*EditLine)
	cnt := 0
	for c := a; c != nil && c != l.Next(); c = c.Next() {
		// XXX hacks on hacks
		cnt += b.ScreenLines(c.Value.(*EditLine))
	}
	return cnt
}

func (b *EditBuffer) MoveDown() {
	if b.line == nil {
		panic(NilLinePanicString)
	}

	if n := b.line.Next(); n != nil {
		ln := n.Value.(*EditLine)
		if ln.moveCursor(b.line.Value.(*EditLine).cursor) < 0 {
			ln.moveCursor(ln.cursorMax())
		}
		b.line = n
		// We are now at line n.  We need to adjust anchor properly so that we can
		// remap the buffer.
		for b.ScreenRange(b.anchor, b.line) > b.view.Rows-1 && b.anchor != b.line {
			b.anchor = b.anchor.Next()
		}
		curr.Map()
	} else {
		Beep()
	}
}

func (b *EditBuffer) MoveUp() {
	if b.line == nil {
		panic(NilLinePanicString)
	}

	if p := b.line.Prev(); p != nil {
		lp := p.Value.(*EditLine)
		if lp.moveCursor(b.line.Value.(*EditLine).cursor) < 0 {
			lp.moveCursor(lp.cursorMax())
		}
		b.line = p
		// We are now at line p.  We need to adjust the anchor in case we've moved
		// above it.  If we have, there is a need for an entire screen remap :(
		for b.line.Value.(*EditLine).lno < b.anchor.Value.(*EditLine).lno && b.anchor != b.line {
			b.anchor = b.anchor.Prev()
		}
		curr.Map()
	} else {
		Beep()
	}
}

func (b *EditBuffer) DeleteSpan(p, l int) {
	if ln, ok := b.line.Value.(*EditLine); ok {
		ln.delete(p, l)
		b.dirty = true
	}
}

func (b *EditBuffer) FirstLine() {
	b.line = b.lines.Front()
}

// Insert a new line after line
// XXX this doesn't really work as desired (can't insert at 0, for instance)
func (b *EditBuffer) InsertLine(line *EditLine) {
	if b.line == nil {
		b.line = b.lines.PushFront(line)
		l := b.line.Value.(*EditLine)
		l.lno = 1
		b.anchor = b.line
	} else {
		e := b.lines.InsertAfter(line, b.line)
		l := e.Value.(*EditLine)
		l.lno = e.Prev().Value.(*EditLine).lno + 1
		for p := e.Next(); p != nil; p = p.Next() {
			p.Value.(*EditLine).lno++
		}
	}
	b.dirty = true
	b.MoveDown() // does the mapping
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewLine([]byte("")))
}


func (b *EditBuffer) NewLine(dlm byte) {
	if l, ok := b.line.Value.(*EditLine); ok {
		newbuf := l.raw()[l.cursor:]
		l.insertCharacter(dlm)
		l.ClearAfterCursor()
		l.size -= len(newbuf)
		b.InsertLine(NewLine(newbuf))
	}
}

func (b *EditBuffer) DeleteCurrLine() {
	if b.line.Prev() != nil {
		rm := b.line
		b.MoveUp()
		b.lines.Remove(rm)
	} else if b.line.Next() != nil {
		rm := b.line
		b.MoveDown()
		b.lines.Remove(rm)
	} else {
		return
	}
	b.dirty = true
	// moveup and movedown take care of mapping
}

func (b *EditBuffer) LnoOffset() int {
	// if we show line numbers, we reserve at least 3 columns.
	if optLineNo {
		if dg := math.Log10(float64(b.lines.Len())) + 1; dg > 2 {
			return int(dg) + 1
		} else {
			return 3
		}
	}
	return 0
}

func (b *EditBuffer) ScreenLines(ln *EditLine) int {
	offset := b.LnoOffset()
	actual := b.view.Cols - offset
	// XXX this is an example of why displaylength needs to be fixed.
	sz := ln.DisplayLength()
	if sz == 0 && ln.hasNewLine {
		sz = 1
	}

	l := int(math.Ceil(float64(sz) / float64(actual)))
	if l == 0 {
		l = 1
	}
	return l
}

// Maps every visible line to a position on the screen.
func (b *EditBuffer) Map() int {
	offset := b.LnoOffset()
	i := 0
	for l := b.anchor; l != nil && i < b.view.Rows-1; l = l.Next() {
		ln := l.Value.(*EditLine)
		cnt := b.ScreenLines(ln)
		wrap := 0
		raw := ln.raw()
		for lmt := i + cnt; i < lmt; i++ {
			str := make([]byte, b.view.Cols)
			// XXX this is the first part of the line, the optional
			// line number.  if this isn't the first screen line of
			// a line (for a wrapped line), then just draw empty
			// space if line numbers are on
			if wrap == 0 {
				lno := []byte(fmt.Sprintf("%*d ", offset-1, ln.lno))
				copy(str[0:offset], lno)
			} else {
				for j := 0; j < offset; j++ {
					str[j] = ' '
				}
			}

			// XXX the second part of the line, which shows actual
			// text that the user is viewing
			actual := b.view.Cols - offset
			start := actual * wrap
			end := start + actual - 1
			if end >= ln.DisplayLength() {
				end = ln.DisplayLength()
			}
			copy(str[offset:], raw[start:end])

			if b.line == l && (ln.cursor >= start || ln.cursor <= end) {
				b.curs_y = i
				b.curs_x = b.LnoOffset() + (ln.cursor - start)
			}
			b.view.Lines[i] = string(str)
			wrap++
		}
	}
	for ; i < b.view.Rows-1; i++ {
		b.view.Lines[i] = NaL
	}
	return 0
}

func (b *EditBuffer) CursorCoord() (int, int) {
	return b.curs_y, b.curs_x
}

func (b *EditBuffer) Lines() *list.List {
	return b.lines
}
