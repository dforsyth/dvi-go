package main

// A simple gap buffer implementation on slices.

import (
	"fmt"
)

const (
	size = 4
)

type GapBuffer struct {
	buf []byte
	gs, ge int
	c int // cursor
}

// Create a new gap buffer
func NewGapBuffer(t []byte) *GapBuffer {
	g := new(GapBuffer)
	g.buf = make([]byte, size)
	// gs is the first index the gap, ge is the first index after the gap
	g.gs = 0
	g.ge = size

	// ghetto hack to make this work easily
	for _, b := range t {
		g.InsertChar(b)
	}
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

// g.c = g.gs + (p - g.gs - (g.ge - g.gs))
// Move cursor to p.  p does not account for the gap.
func (g *GapBuffer) MoveCursor(p int) {
	if p > g.gs {
		g.c = g.gs + (p - g.gs - (g.ge - g.gs))
		// g.c = p + (g.ge - g.gs) + g.gs
	} else {
		g.c = p
	}
}

func (g *GapBuffer) CursorLeft() {
	if g.c == 0 {
		return
	}
	if g.c <= g.gs {
		g.MoveCursor(g.c - 1)
	} else {
		g.MoveCursor(g.c - (g.ge - g.gs))
	}
	Debug = fmt.Sprintf("cursor now at %d or %s", g.c, g.String())
}

func (g *GapBuffer) CursorRight() {
	if g.c == len(g.GaplessBuffer()) {
		return
	}

	if g.c <= g.gs {
		g.MoveCursor(g.c + 1)
	} else {
		g.MoveCursor(g.c + (g.ge - g.gs))
	}
	Debug = fmt.Sprintf("cursor now at %d out of %d", g.c, len(g.String()))
}

// Move the cursor to g.gs
func (g *GapBuffer) UpdateCursor() {
	g.c = g.gs
	Debug = "update cursor"
}

// Grow the size of gap by s
func (g *GapBuffer) GrowGap(s int) {
	// TODO: slices double on realloc up to 1024, then they increase by 25%.
	// we should check the cap() of our new slice after append and see if we
	// can use the new size to make a larger gap.
	b := make([]byte, s)
	for i, _ := range b {
		b[i] = byte(0)
	}

	g.buf = append(g.buf, b...)
	copy(g.buf[g.ge + s:], g.buf[g.ge:])
	g.gs = g.ge
	g.ge += s
}

// Move the gap to p.  p does not take the gap into account.
func (g *GapBuffer) MoveGap(p int) {
	if g.gs == p {
		return
	}

	if p < g.gs {
		copy(g.buf[g.ge - (g.gs - p):g.ge], g.buf[p:g.gs])
		g.ge -= (g.gs - p)
	} else {
		copy(g.buf[g.gs:p], g.buf[g.ge:g.ge + (p - g.gs)])
		g.ge += (p - g.gs)
	}
	g.gs = p

	Debug = fmt.Sprintf("moved gap to %d", g.gs)
}

func (g *GapBuffer) MoveGapToCursor() {
	if g.c <= g.gs {
		g.MoveGap(g.c)
	} else {
		g.MoveGap(g.ge + g.c)
	}
}

func (g *GapBuffer) Buffer() []byte {
	return g.buf
}

func (g *GapBuffer) GaplessBuffer() []byte {
	return append(g.buf[:g.gs], g.buf[g.ge:]...)
}

func (g *GapBuffer) String() string {
	return string(g.buf[:g.gs]) + string(g.buf[g.ge:])
}

