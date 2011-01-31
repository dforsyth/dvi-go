package main

import (
	"os"
	"fmt"
)

// todo: lockable
type EditBuffer struct {
	lno	  int
	st    *os.FileInfo
	title string

	lines *Line
	line *Line

	prev, next *EditBuffer // roll ourselves because type assertions are pointless in this case.
}

func NewEditBuffer(title string) *EditBuffer {

	b := new(EditBuffer)

	b.lines = new(Line)
	b.line = nil
	b.lno = 0
	b.st = nil
	b.title = title
	b.next = nil
	b.prev = nil

	return b
}

func (b *EditBuffer) InsertChar(ch byte) {
	b.line.insertCharacter(ch)
}

func (b *EditBuffer) BackSpace() {
	if b.line == nil {
		Debug = "nothing to backspace"
		return
	}


	if b.line.cursor == 0 {
		if b.line.size != 0 {
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
	if b.line.cursor == 0 {
		Beep()
	} else {
		b.line.moveCursor(b.line.cursor - 1)
	}
}

func (b *EditBuffer) MoveCursorRight() {
	if b.line.CursorIsMax() {
		Beep()
		return
	}
	b.line.moveCursor(b.line.cursor + 1)
}

func (b *EditBuffer) MoveCursorDown() {
	if b.line != nil {
		if n := b.line.next; n != nil {
			c := b.line.cursor
			b.line = n
			b.line.moveCursor(c)
			b.lno++
			Message = fmt.Sprintf("down %d", b.lno, b.line.bytes())
		}
	}
}

func (b *EditBuffer) MoveCursorUp() {
	if b.line != nil {
		if p := b.line.prev; p != nil {
			c := b.line.cursor
			b.line = p
			b.line.moveCursor(c)
			b.lno--
			Message = fmt.Sprintf("up %d", b.lno, b.line.bytes())
		}
	}
}

func (b *EditBuffer) DeleteSpan(p, l int) {
	b.line.delete(p, l)
}

func (b *EditBuffer) FirstLine() {
	b.line = b.lines
}

func (b *EditBuffer) InsertLine(line *Line) {
	if b.line == nil {
		b.lines = line
	} else {
		b.line.next = line
	}
	b.line = line
	b.lno++
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewLine([]byte("")))
}


func (b *EditBuffer) NewLine(nlchar byte) {

	// newbuf := b.line.bytes()[b.line.cursor:]
	// b.Line().InsertChar(nlchar)
	// b.Line().DeleteAfterGap()
	// b.InsertLine(NewGapBuffer([]byte(newbuf)))
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
