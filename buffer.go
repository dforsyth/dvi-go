package main

import (
	"os"
	"fmt"
)

// todo: lockable
type EditBuffer struct {
	lno, lco int
	st       *os.FileInfo
	title    string
	dirty bool

	lines *Line
	line  *Line

	prev, next *EditBuffer // roll ourselves because type assertions are pointless in this case.
}

func NewEditBuffer(title string) *EditBuffer {

	b := new(EditBuffer)

	// XXX it would be better to init this to nil and add the line from another
	// place so that the file loader can use this
	b.lines = NewLine([]byte(""))
	b.line = nil
	b.lno = 0
	b.lco = 0
	b.st = nil
	b.title = title
	b.next = nil
	b.prev = nil
	b.dirty = false

	return b
}

func (b *EditBuffer) InsertChar(ch byte) {
	b.line.insertCharacter(ch)
	b.dirty = true
}

func (b *EditBuffer) BackSpace() {
	if b.line == nil {
		Debug = "nothing to backspace"
		return
	}

	Message = fmt.Sprintf("%d", b.line.cursor)
	if b.line.cursor == 0 {
		if b.line.size != 0 && b.line.prev != nil {
			// combine this line and the previous
		} else {

			if b.line.prev != nil {
				b.DeleteCurrLine()
			} else {
				Beep()
			}
		}
	} else {
		b.line.backspace()
	}
}

func (b *EditBuffer) MoveCursorLeft() {
	if b.line.moveCursor(b.line.cursor-1) < 0 {
		Beep()
	}
}

func (b *EditBuffer) MoveCursorRight() {
	if b.line.moveCursor(b.line.cursor+1) < 0 {
		Beep()
	}
}

func (b *EditBuffer) MoveCursorDown() {
	if n := b.line.next; n != nil {
		if n.moveCursor(b.line.cursor) < 0 {
			n.moveCursor(b.line.cursorMax())
		}
		b.line = n
		b.lno++
	} else {
		Beep()
	}
}

func (b *EditBuffer) MoveCursorUp() {
	if p := b.line.prev; p != nil {
		if p.moveCursor(b.line.cursor) < 0 {
			p.moveCursor(b.line.cursorMax())
		}
		b.line = p
		b.lno--
	} else {
		Beep()
	}
}

func (b *EditBuffer) DeleteSpan(p, l int) {
	b.line.delete(p, l)
	b.dirty = true
}

func (b *EditBuffer) FirstLine() {
	b.line = b.lines
}

func (b *EditBuffer) InsertLine(line *Line) {
	if b.line == nil {
		b.lines = line
	} else {
		line.prev = b.line
		line.next = b.line.next
		if b.line.next != nil {
			b.line.next.prev = line
		}
		b.line.next = line
	}
	b.line = line
	b.lno++
	b.dirty = true
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewLine([]byte("")))
}


func (b *EditBuffer) NewLine(nlchar byte) {

	newbuf := b.line.bytes()[b.line.cursor:]
	b.line.insertCharacter(nlchar)
	b.line.ClearAfterCursor()
	b.line.size -= len(newbuf)
	b.InsertLine(NewLine(newbuf))
}

func (b *EditBuffer) DeleteCurrLine() {
	p, n := b.line.prev, b.line.next
	if p != nil {
		p.next = n
		b.line = p
		b.lno--
	} else if n != nil {
		n.prev = p
		b.line = n
		// line number doesn't change
	}
	b.dirty = true
}

// Move to line p
func (b *EditBuffer) MoveLine(p int) {
	i := 0
	for l := b.lines; l != nil; l = l.next {
		if i == p {
			b.line = l
			return
		}
		i++
	}
}

func (b *EditBuffer) MoveLineNext() {
	n := b.line.next
	if n != nil {
		b.line = n
	}
}

func (b *EditBuffer) MoveLinePrev() {
	p := b.line.prev
	if p != nil {
		b.line = p
	}
}

func (b *EditBuffer) Lines() *Line {
	return b.lines
}

func (b *EditBuffer) Title() string {
	return b.title
}
