package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
)


// XXX editbuffers are editable text buffers that happen to also be a screen.

// todo: lockable
type EditBuffer struct {
	st                             *os.FileInfo // file info
	fromDisk, fromArchive, fromNet bool

	lco   uint   // line count
	title string // buffer title
	dirty bool   // dirty

	// XXX text
	lines *list.List    // list of lines
	line  *list.Element // current line

	anchor *list.Element // first line to draw
	view   *View         // view this buffer draws to
	curs_x, curs_y int // cursor position

	hist *list.List

	prev, next *EditBuffer // roll ourselves because type assertions are pobyteless in this case.

	OptLineNo bool // draw line numbers
	OptHLLine bool // highlight the current line
}

func NewEditBuffer(title string, optLineNo, optHLLine bool, view *View) *EditBuffer {

	b := new(EditBuffer)

	// XXX it would be better to init this to nil and add the line from another
	// place so that the file loader can use this
	b.lines = list.New()
	// b.InsertLine(NewLine([]byte("")))
	b.anchor = b.line
	b.st = nil
	b.title = title
	b.next = nil
	b.prev = nil
	b.dirty = false

	b.view = view
	b.OptLineNo = optLineNo
	b.OptHLLine = optHLLine

	return b
}

func (b *EditBuffer) InsertChar(ch byte) {
	if l, ok := b.line.Value.(*EditLine); ok {
		l.insertCharacter(ch)
	}
	b.dirty = true
}

func (b *EditBuffer) BackSpace() {
	if b.line == nil {
		Beep()
		return
	}

	if l, ok := b.line.Value.(*EditLine); ok {
		if l.cursor == 0 {
			if l.size != 0 && b.line.Prev() != nil {
				// combine this line and the previous
			} else {

				if b.line.Prev() != nil {
					b.DeleteCurrLine()
					b.lco--
				} else {
					Beep()
				}
			}
		} else {
			l.backspace()
		}
	}
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

func (b *EditBuffer) MoveDown() {
	if l, ok := b.line.Value.(*EditLine); ok {
		if n := b.line.Next(); n != nil {
			ln := n.Value.(*EditLine)
			if ln.moveCursor(l.cursor) < 0 {
				ln.moveCursor(ln.cursorMax())
			}
			b.line = n
		} else {
			Beep()
		}
	}
}

func (b *EditBuffer) MoveUp() {
	if l, ok := b.line.Value.(*EditLine); ok {
		if p := b.line.Prev(); p != nil {
			lp := p.Value.(*EditLine)
			if lp.moveCursor(l.cursor) < 0 {
				lp.moveCursor(lp.cursorMax())
			}
			b.line = p
		} else {
			Beep()
		}
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
	} else {
		b.line = b.lines.InsertAfter(line, b.line)
		l := b.line.Value.(*EditLine)
		l.lno = b.line.Prev().Value.(*EditLine).lno + 1
		for p := b.line.Next(); p != nil; p = p.Next() {
			p.Value.(*EditLine).lno++
		}
	}
	b.lco++
	b.dirty = true
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewLine([]byte("")))
}


func (b *EditBuffer) NewLine(nlchar byte) {
	if l, ok := b.line.Value.(*EditLine); ok {
		newbuf := l.raw()[l.cursor:]
		l.insertCharacter(nlchar)
		l.ClearAfterCursor()
		l.size -= len(newbuf)
		b.InsertLine(NewLine(newbuf))
	}
}

func (b *EditBuffer) DeleteCurrLine() {
	if b.line.Prev() != nil {
		p := b.line.Prev()
		b.lines.Remove(b.line)
		b.line = p
	} else if b.line.Next() != nil {
		n := b.line.Next()
		b.lines.Remove(b.line)
		b.line = n
	} else {
		return
	}
	b.dirty = true
}

// Move to line p
func (b *EditBuffer) MoveLine(p int) {
	i := 0
	for l := b.lines.Front(); l != nil; l = l.Next() {
		if i == p {
			b.line = l
			return
		}
		i++
	}
}

func (b *EditBuffer) MoveLineNext() {
	if n := b.line.Next(); n != nil {
		b.line = n
	}
}

func (b *EditBuffer) MoveLinePrev() {
	if p := b.line.Prev(); p != nil {
		b.line = p
	}
}

func (b *EditBuffer) LnoOffset() int {
	// if we show line numbers, we reserve at least 3 columns.
	if b.OptLineNo {
		if dg := math.Log10(float64(b.lco)) + 1; dg > 2 {
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
	return int(math.Ceil(float64(ln.DisplayLength()) / float64(actual)))
}

// Maps every visible line to a position on the screen.  This is a super-slow complete refresh.
func (b *EditBuffer) Map() int {
	offset := b.LnoOffset()
	i := 0
	for l := b.anchor; l != nil && i < b.view.Rows; l = l.Next() {
		ln := l.Value.(*EditLine)
		cnt := b.ScreenLines(ln)
		if cnt == 0 {
			cnt = 1
		}
		wrap := 0
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
			copy(str[offset:], ln.raw()[start:end])

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
	return b.curs_x, b.curs_y
}

func (b *EditBuffer) Lines() *list.List {
	return b.lines
}

func (b *EditBuffer) Title() string {
	return b.title
}
