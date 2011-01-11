package main

import (
	"container/list"
	// "fmt"
	"os"
)

const (
	gSize int = 4
	maxLine int = 32
)

type EditBuffer struct {
	lines *list.List
	line *list.Element
	ln int
	title string
	st *os.FileInfo
}

func NewEditBuffer(title string) *EditBuffer {

	b := new(EditBuffer)
	b.lines = list.New()
	b.line = nil
	b.ln = 0
	b.title = title

	return b

}

func (b *EditBuffer) InsertChar(ch byte) {
	b.Line().InsertChar(ch)
}

func (b *EditBuffer) BackSpace() {
	if b.Line() == nil {
		Debug = "nothing to backspace"
		return
	}

	if b.Line().gs == 0 {
		if len(b.Line().buf) != (b.Line().ge - b.Line().gs) {
		} else {
			b.DeleteCurrLine()
		}
	} else {
		b.Line().DeleteSpan(b.Line().gs - 1, 1)
	}
}

func (b *EditBuffer) MoveCursorLeft() {
	b.Line().CursorLeft()
}

func (b *EditBuffer) MoveCursorRight() {
	b.Line().CursorRight()
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

