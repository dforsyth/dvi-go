package main

// A simple gap buffer implementation on slices.

import (
	"fmt"
)

const (
	size = 64
	max  = 4096
)

type GapBuffer struct {
	buf    []byte
	gs, ge int
	c      int // cursor
}

// Create a new gap buffer
func NewGapBuffer(t []byte) *GapBuffer {
	g := new(GapBuffer)
	g.buf = make([]byte, len(t) + size)

	copy(g.buf[:len(t)], t)

	// gs is the first index the gap, ge is the first index after the gap
	g.gs = len(t)
	g.ge = len(t) + size

	// ghetto hack to make this work easily
	/*
	for _, b := range t {
		g.InsertChar(b)
	}
	*/
	g.MoveGap(0)
	return g
}

// Insert a character at the first position of the gap
func (g *GapBuffer) InsertChar(c byte) {
	g.buf[g.gs] = c
	g.gs++
	if g.gs == g.ge {
		g.GrowGap(size)
	}
}

func (g *GapBuffer) InsertString(s string) {
	for _, c := range s {
		g.InsertChar(byte(c))
	}
}

func (g *GapBuffer) DeleteSpan(p, s int) {
	g.MoveGap(p + s)
	for i := 0; i < s; i++ {
		if g.gs == 0 {
			return
		}
		g.gs--
	}
}

// Move cursor to p.  p does not account for the gap.
func (g *GapBuffer) MoveCursor(p int) {
	g.c = p
}

func (g *GapBuffer) CursorLeft() {
	if g.c == 0 {
		return
	}
	g.MoveCursor(g.c - 1)
	Debug = fmt.Sprintf("%d", g.c)
}

func (g *GapBuffer) CursorRight() {
	if g.c == len(g.GaplessBuffer()) {
		return
	}
	g.MoveCursor(g.c + 1)
}

// Move the cursor to g.gs
func (g *GapBuffer) UpdateCursor() {
	g.c = g.gs
}

// Grow the size of gap by s
func (g *GapBuffer) GrowGap(s int) {
	// TODO: slices double on realloc up to 1024, then they increase by 25%.
	// we should check the cap() of our new slice after append and see if we
	// can use the new size to make a larger gap.
	b := make([]byte, s)

	g.buf = append(g.buf, b...)
	copy(g.buf[g.ge+s:], g.buf[g.ge:])
	g.gs = g.ge
	g.ge += s
}

// Move the gap to p.  p does not take the gap into account.
func (g *GapBuffer) MoveGap(p int) {
	if g.gs == p {
		return
	}

	if p < g.gs {
		s := g.gs - p
		copy(g.buf[g.ge-s:g.ge], g.buf[p:g.gs])
		g.ge -= s
	} else {
		s := p - g.gs
		copy(g.buf[g.gs:p], g.buf[g.ge:g.ge+s])
		g.ge += s
	}
	g.gs = p
}

func (g *GapBuffer) MoveGapToCursor() {
	g.MoveGap(g.c)
}

func (g *GapBuffer) Buffer() []byte {
	return g.buf
}

func (g *GapBuffer) GaplessBuffer() []byte {
	b := make([]byte, len(g.buf[:g.gs])+len(g.buf[g.ge:]))
	copy(b, g.buf[:g.gs])
	copy(b[:g.gs], g.buf[g.ge:])
	return b
}

func (g *GapBuffer) String() string {
	return string(g.buf[:g.gs]) + string(g.buf[g.ge:])
}

func (g *GapBuffer) DebugString() string {
	s := ""
	for i, c := range g.buf {
		if i >= g.gs && i < g.ge {
			s += "_"
		} else {
			s += string(c)
		}

	}

	return s
}

func (g *GapBuffer) DebugCursor() int {
	if g.c > g.gs {
		return g.c + (g.ge - g.gs)
	}
	return g.c
}
