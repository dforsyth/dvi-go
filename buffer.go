package main

import (
	"container/list"
	"os"
	"fmt"
)

// todo: lockable
type EditBuffer struct {
	dirty bool
	lines *list.List
	line  *list.Element
	lno	  int
	st    *os.FileInfo
	title string

	prev, next *EditBuffer // roll ourselves because type assertions are pointless in this case.
}

func NewEditBuffer(title string) *EditBuffer {

	b := new(EditBuffer)
	b.lines = list.New()
	b.line = nil
	b.lno = 0
	b.st = nil
	b.title = title
	b.next = nil
	b.prev = nil

	return b
}

func (b *EditBuffer) Dirty() bool {
	return b.dirty
}

func (b *EditBuffer) SetDirty(d bool) {
	b.dirty = d
}

func (b *EditBuffer) InsertChar(ch byte) {
	b.Line().InsertChar(ch)
	b.SetDirty(true)
}

func (b *EditBuffer) BackSpace() {
	if b.Line() == nil {
		Debug = "nothing to backspace"
		return
	}

	if b.Line().gs == 0 {
		if len(b.Line().buf) != (b.Line().ge - b.Line().gs) {
		} else {
			if b.line.Prev() != nil {
				b.DeleteCurrLine()
			} else {
				Beep()
			}
		}
	} else {
		b.Line().DeleteSpan(b.Line().gs-1, 1)
	}
}

func (b *EditBuffer) MoveCursorLeft() {
	b.Line().CursorLeft()
}

func (b *EditBuffer) MoveCursorRight() {
	b.Line().CursorRight()
}

func (b *EditBuffer) MoveCursorDown() {
	if b.line != nil {
		if n := b.line.Next(); n != nil {
			c := b.Line().c
			b.line = n
			b.Line().MoveCursor(c)
			b.lno++
			Message = fmt.Sprintf("down %d", b.lno, b.line.Value.(*GapBuffer).String())
		}
	}
}

func (b *EditBuffer) MoveCursorUp() {
	if b.line != nil {
		if p := b.line.Prev(); p != nil {
			c := b.Line().c
			b.line = p
			b.Line().MoveCursor(c)
			b.lno--
			Message = fmt.Sprintf("up %d", b.lno, b.line.Value.(*GapBuffer).String())
		}
	}
}

func (b *EditBuffer) DeleteSpan(p, l int) {
	b.Line().DeleteSpan(p, l)
}

func (b *EditBuffer) FirstLine() {
	b.line = b.lines.Front()
}

func (b *EditBuffer) InsertLine(g *GapBuffer) {
	if b.line == nil {
		b.line = b.lines.PushFront(g)
		return
	}
	b.line = b.lines.InsertAfter(g, b.line)
	b.lno++
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewGapBuffer([]byte("")))
}

func (b *EditBuffer) NewLine(nlchar byte) {

	newbuf := b.Line().String()[b.Line().c:]
	// b.Line().InsertChar(nlchar)
	b.Line().DeleteAfterGap()
	b.InsertLine(NewGapBuffer([]byte(newbuf)))
}

func (b *EditBuffer) Line() *GapBuffer {
	if b.line == nil {
		return nil
	}
	return b.line.Value.(*GapBuffer)
}

func (b *EditBuffer) DeleteCurrLine() {
	p := b.line.Prev()
	b.lines.Remove(b.line)
	b.line = p
	b.lno--
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
	n := b.line.Next()
	if n != nil {
		b.line = n
	}
}

func (b *EditBuffer) MoveLinePrev() {
	p := b.line.Prev()
	if p != nil {
		b.line = p
	}
}

func (b *EditBuffer) Lines() *list.List {
	return b.lines
}

func (b *EditBuffer) Title() string {
	return b.title
}
