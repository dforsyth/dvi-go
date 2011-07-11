package main

import (
	"bufio"
	"os"
	"utf8"
)

type Line struct {
	text    []byte // string
	displen uint
	next    *Line
	prev    *Line
	dirty   bool
	// attributes
}

func (l *Line) length() int {
	return utf8.RuneCount(l.text)
}

func NewLine(text []byte) *Line {
	dst := make([]byte, len(text))
	copy(dst, text)
	return &Line{
		text: dst,
		next: nil,
		prev: nil,
	}
}

const (
	ADD = iota
	REMOVE
)

type Action struct {
	op   int // the action that happened, do the opposite to undo
	s, e *Position
	text []byte
	char *Action // multi op actions?
}

type Buffer struct {
	name              string
	dirty             bool
	temp              bool
	first, last, disp *Line
	pos               *Position
	line              bool // line mode buffer, for yank/put
	next, prev        *Buffer
	db                []*Line
}

func newBuffer() *Buffer {
	l := NewLine([]byte{})
	return &Buffer{
		first: l,
		last:  l,
		pos: &Position{
			line: l,
			off:  0,
		},
	}
}

func (b *Buffer) resetPos() {
	b.pos.line = b.first
	b.pos.off = 0
}

func (b *Buffer) rename(name string) {
	b.name = name
}

func insertLineBelow(l, n *Line) {
	if l.next != nil {
		l.next.prev = n
		n.next = l.next
	}
	n.prev = l
	l.next = n
}

func insertLineAbove(l, n *Line) {
	if l.prev != nil {
		l.prev.next = n
		n.prev = l.prev
	}
	n.next = l
	l.prev = n
}

func (b *Buffer) insertLineAbove(l, n *Line) {
	insertLineAbove(l, n)
	if l == b.first {
		b.first = n
	}
}

func (b *Buffer) insertLineBelow(l, n *Line) {
	insertLineBelow(l, n)
	if l == b.last {
		b.last = n
	}
}

func (b *Buffer) add(p Position, text []byte) *Position {
	for _, c := range text {
		if c != '\n' && c != '\r' {
			l := p.line
			l.text = append(l.text[:p.off], append([]byte{c}, l.text[p.off:]...)...)
			p.off++
		} else {
			linetext := p.line.text
			l := NewLine(linetext[p.off:])
			p.line.text = linetext[:p.off]
			/*
				if p.line.next != nil {
					p.line.next.prev = l
					l.next = p.line.next
				}
				l.prev = p.line
				p.line.next = l
			*/
			b.insertLineBelow(p.line, l)
			p.line = l
			p.off = 0
		}
	}
	return &p
}

func (b *Buffer) lineCount() int {
	i := 0
	for l := b.first; l != b.last.next; l = l.next {
		i++
	}
	return i
}

func (b *Buffer) lineNumber(l *Line) int {
	for i, bl := 1, b.first; bl != nil && l != nil; i, bl = i+1, bl.next {
		if bl == l {
			return i
		}
	}
	return -1
}

func (b *Buffer) getLine(lno int) *Line {
	l := b.first
	for i := 1; i < lno; i, l = i+1, l.next {
		if l.next == nil || l == b.last {
			break
		}
	}
	return l
}

func (b *Buffer) remove(start, end Position, line bool) {
	// XXX This function returns b.pos.  It should actually just return the first safe 
	// position after (or before) the removed chunk.
	if start.line == end.line && !line {
		start.line.text = append(start.line.text[:start.off], start.line.text[end.off:]...)
		/*return &Position{
			line: start.line,
			off:  start.off,
		}*/
	} else {
		if !line {
			// If we aren't in line mode, check if only the 0th char is marked by end.
			// If it is, move back to the end of the prev, because we don't actually
			// want to delete the line.
			if end.off == 0 {
				end.line = end.line.prev
				end.off = end.line.length()
			}
			for l := start.line; l != end.line; l = l.next {
				start.line.next = l.next
				start.line.next.prev = start.line
			}
			start.line.text = append(start.line.text[:start.off], end.line.text[end.off:]...)
			if end.line.next != nil {
				end.line.next.prev = start.line
			}
			start.line.next = end.line.next
		} else {
			for l := start.line; l != end.line.next; l = l.next {
				if l != b.first /* l.prev != nil */ {
					l.prev.next = l.next
				} else {
					b.first = l.next
				}
				if l != b.last /* l.next != nil */ {
					l.next.prev = l.prev
				} else {
					b.last = l.prev
				}
				if b.disp == l {
					if l.next != nil {
						b.disp = l.next
					} else {
						b.disp = l.prev
					}
				}
			}
			// we deleted the entire file
			if b.first == nil {
				b.first = NewLine([]byte{})
				b.last = b.first
				b.disp = b.first
				b.pos.line = b.first
				b.pos.off = 0
			} else if !b.contains(b.pos.line) {
				b.pos.line = b.first
				b.pos.off = 0
			}
		}
	}
	// return b.pos
}

func (b *Buffer) removeLine(l *Line) *Line {
	return nil
}

func (b *Buffer) contains(l *Line) bool {
	for c := b.first; c != b.last.next; c = c.next {
		if l == c {
			return true
		}
	}
	return false
}

func get(a, b *Position) []byte {
	// doesn't need to be on a buffer
	text := []byte{}
	for l, s := a.line, a.off; l != b.line.next; l, s = l.next, 0 {
		if l == b.line {
			text = append(text, l.text[s:b.off]...)
		} else {
			text = append(text, l.text[s:]...)
			if l.next != nil {
				text = append(text, '\n')
			}
		}
	}
	return text
}

func (b *Buffer) getAll() []byte {
	return get(&Position{b.first, 0}, &Position{b.last, b.last.length()})
}

func (b *Buffer) clear() {
	b.first = NewLine([]byte{})
	b.last = b.first
}

func (b *Buffer) loadFile(f *os.File) os.Error {
	rdr := bufio.NewReader(f)
	buf := make([]byte, 4096)
	var e os.Error
	var n int
	for n, e = rdr.Read(buf); n > 0 && e == nil; n, e = rdr.Read(buf) {
		b.pos = b.add(*b.pos, buf[:n])
	}
	if e != os.EOF {
		return e
	}
	if b.last.length() == 0 && b.last != b.first {
		p := b.last.prev
		p.next = nil
		b.last.prev = nil
		b.last = p
	}
	return nil
}

func (b *Buffer) writeFile() os.Error {
	wf, e := os.Create(b.name)
	if e != nil {
		return e
	}
	defer wf.Close()
	for l := b.first; l != nil; l = l.next {
		if _, e := wf.Write(append(l.text, '\n')); e != nil {
			return e
		}
	}
	return nil
}
