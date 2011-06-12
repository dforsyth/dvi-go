package main

import (
//	"fmt"
)

func isBlank(b byte) bool {
	return b < 0x30
}

// add text at p.  return the positon at the end of the inserted text.
// XXX It might be a better idea to take a *Position and just update that?
func add(p Position, text []byte) *Position {
	for _, b := range text {
		if b != '\n' && b != '\r' {
			l := p.line
			l.text = append(l.text[:p.off], append([]byte{b}, l.text[p.off:]...)...)
			p.off++
		} else {
			linetext := p.line.text
			l := NewLine(linetext[p.off:])
			p.line.text = linetext[:p.off]
			if p.line.next != nil {
				p.line.next.prev = l
				l.next = p.line.next
			}
			l.prev = p.line
			p.line.next = l
			p.line = l
			p.off = 0
		}
	}
	return &p
}

// remove text between a and b.  return position the text was remove from.
func remove(a, b Position) *Position {
	if a.line == b.line {
		a.line.text = append(a.line.text[:a.off], a.line.text[b.off:]...)
	} else {
		for l := a.line; l != b.line; l = l.next {
			a.line.next = l.next
			a.line.next.prev = a.line
		}
		a.line.text = append(a.line.text[:a.off], b.line.text[b.off:]...)
		b.line.next.prev = a.line
		a.line.next = b.line.next
	}
	return &a
}

// return the text between a and b
func get(a, b *Position) []byte {
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
