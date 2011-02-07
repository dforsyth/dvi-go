package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
)

type Screen interface {
	Map() int // draw the screen in the views map
	View() *View
	SetView(*View)
}

// XXX editbuffers are editable text buffers

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

	hist *list.List

	prev, next *EditBuffer // roll ourselves because type assertions are pobyteless in this case.

	OptLineNo bool
}

func NewEditBuffer(title string) *EditBuffer {

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

	b.view = Vw
	b.OptLineNo = true

	return b
}

func (b *EditBuffer) InsertChar(ch byte) {
	if l, ok := b.line.Value.(*Line); ok {
		l.insertCharacter(ch)
	}
	b.dirty = true
}

func (b *EditBuffer) BackSpace() {
	if b.line == nil {
		Beep()
		return
	}

	if l, ok := b.line.Value.(*Line); ok {
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
	if l, ok := b.line.Value.(*Line); ok && l.moveCursor(l.cursor+d) < 0 {
		Beep()
	}
}

func (b *EditBuffer) MoveDown() {
	if l, ok := b.line.Value.(*Line); ok {
		if n := b.line.Next(); n != nil {
			ln := n.Value.(*Line)
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
	if l, ok := b.line.Value.(*Line); ok {
		if p := b.line.Prev(); p != nil {
			lp := p.Value.(*Line)
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
	if ln, ok := b.line.Value.(*Line); ok {
		ln.delete(p, l)
		b.dirty = true
	}
}

func (b *EditBuffer) FirstLine() {
	b.line = b.lines.Front()
}

// Insert a new line after line
// XXX this doesn't really work as desired (can't insert at 0, for instance)
func (b *EditBuffer) InsertLine(line *Line) {
	if b.line == nil {
		b.line = b.lines.PushFront(line)
		l := b.line.Value.(*Line)
		l.lno = 1
	} else {
		b.line = b.lines.InsertAfter(line, b.line)
		l := b.line.Value.(*Line)
		l.lno = b.line.Prev().Value.(*Line).lno + 1
		for p := b.line.Next(); p != nil; p = p.Next() {
			p.Value.(*Line).lno++
		}
	}
	b.lco++
	b.dirty = true
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewLine([]byte("")))
}


func (b *EditBuffer) NewLine(nlchar byte) {
	if l, ok := b.line.Value.(*Line); ok {
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
	// if we show line numbers, we reserve at least 3 lines.
	if b.OptLineNo {
		if dg := math.Log10(float64(b.lco)) + 1; dg > 2 {
			return int(dg) + 1
		} else {
			return 3
		}
	}
	return 0
}

func (b *EditBuffer) ScreenLines(ln *Line) int {
	offset := b.LnoOffset()
	actual := b.view.Cols - offset
	return int(math.Ceil(float64(len(ln.raw())) / float64(actual)))
}

func (b *EditBuffer) Map() int {
	offset := b.LnoOffset()
	i := 0
	for l := b.anchor; l != nil && i < b.view.Rows; l = l.Next() {
		ln := l.Value.(*Line)
		cnt := b.ScreenLines(ln)
		if cnt == 0 {
			cnt = 1
		}
		first := true
		for lmt := i + cnt; i < lmt; i++ {
			str := make([]byte, b.view.Cols)
			//for j := 0; j < offset; j++ {
			//	str[j] = ' '
			//}
			if first {
				lno := []byte(fmt.Sprintf("%*d ", offset-1, ln.lno))
				copy(str[0:offset], lno)
				first = false
			} else {
				for j := 0; j < offset; j++ {
					str[j] = ' '
				}
			}
			copy(str[offset:], ln.raw())
			b.view.Lines[i] = string(str)
		}
	}
	for ; i < b.view.Rows-1; i++ {
		b.view.Lines[i] = NaL
	}
	return 0
}

func (b *EditBuffer) CursorCoord() (int, int) {
	l := b.line.Value.(*Line)
	return int(l.cursor) + b.LnoOffset(), int(l.lno)
}

func (b *EditBuffer) Lines() *list.List {
	return b.lines
}

func (b *EditBuffer) Title() string {
	return b.title
}
