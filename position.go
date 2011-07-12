package main

import (
	"os"
)

type Position struct {
	line *Line
	off  int
}

func posEq(p, q *Position) bool {
	return p.line == q.line && p.off == q.off
}

func orderPos(a, b *Position) (*Position, *Position) {
	if a.line != b.line {
		for l := b.line; l != nil; l = l.next {
			if l == a.line {
				return b, a
			}
		}
		return a, b
	}
	if b.off < a.off {
		return b, a
	}
	return a, b
}

func (p *Position) getChar() (int, os.Error) {
	if p.line.length() > 0 && p.off < p.line.length() {
		// XXX cast this until we go full rune up in this bitch
		return int(p.line.text[p.off]), nil
	}
	return -1, &DviError{}
}

func (p *Position) setChar(c int) os.Error {
	if p.line.length() > 0 && p.off < p.line.length() {
		p.line.text[p.off] = byte(c)
		return nil
	}
	return &DviError{}
}

// XXX These should really be renamed nextPos and prevPos
func prevChar(p Position) *Position {
	if p.off > 0 {
		p.off--
	}
	return &p
}

func prevChar2(p Position) *Position {
	if p.off > 0 {
		p.off--
	} else if p.line.prev != nil {
		p.line = p.line.prev
		p.off = p.line.length()
	}
	return &p
}

func nextChar(p Position) *Position {
	if p.off < p.line.length() {
		p.off++
	}
	return &p
}

func nextChar2(p Position) *Position {
	if p.off < p.line.length() {
		p.off++
	} else if p.line.next != nil {
		p.line = p.line.next
		p.off = 0
	}
	return &p
}

func prevWord(p Position) *Position {
	return &p
}

func nextWord(p Position) *Position {
	return &p
}

func prevLine(p Position) *Position {
	if p.line.prev != nil {
		p.line = p.line.prev
		// TODO utf8-itize this
		if p.off > p.line.length() {
			p.off = p.line.length()
		}
	}
	return &p
}

func nextLine(p Position) *Position {
	if p.line.next != nil {
		p.line = p.line.next
		if p.off > p.line.length() {
			p.off = p.line.length()
		}
	}
	return &p
}

func eol(p Position) *Position {
	p.off = p.line.length() - 1
	return &p
}

func bol(p Position) *Position {
	p.off = 0
	return &p
}
