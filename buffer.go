package main

import (
	"container/list"
	"os"
)

// todo: lockable
type EditBuffer struct {
	dirty bool
	lines *list.List
	line  *list.Element
	lno	  int
	st    *os.FileInfo
	title string
}

func NewEditBuffer(title string) *EditBuffer {

	b := new(EditBuffer)
	b.lines = list.New()
	b.line = nil
	b.lno = 0
	b.st = nil
	b.title = title

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
		n := b.line.Next()
		if n != nil {
			b.line = n
		}
	}
}

func (b *EditBuffer) MoveCursorUp() {
	if b.line != nil {
		p := b.line.Prev()
		if p != nil {
			b.line = p
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
}

func (b *EditBuffer) AppendLine() {
	b.InsertLine(NewGapBuffer([]byte("")))
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
