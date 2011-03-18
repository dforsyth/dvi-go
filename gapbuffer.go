package main

// A simple gap buffer implementation on slices.

const (
	size = 64
	max  = 4096
)

type gapBuffer struct {
	buf    []byte
	gs, ge int
}

// Create a new gap buffer
func newGapBuffer(t []byte) *gapBuffer {
	g := new(gapBuffer)
	g.buf = make([]byte, len(t)+size)

	copy(g.buf[:len(t)], t)

	// gs is the first index the gap, ge is the first index after the gap
	g.gs = len(t)
	g.ge = len(t) + size

	g.MoveGap(0)
	return g
}

// Insert a character at the first position of the gap
func (g *gapBuffer) insertChar(c byte) {
	g.buf[g.gs] = c
	g.gs++
	if g.gs == g.ge {
		g.GrowGap(size)
	}
}

func (g *gapBuffer) InsertString(s string) {
	for _, c := range s {
		g.insertChar(byte(c))
	}
}

func (g *gapBuffer) deleteSpan(p, s int) {
	g.MoveGap(p + s)
	for i := 0; i < s; i++ {
		if g.gs == 0 {
			return
		}
		g.gs--
	}
}

func (g *gapBuffer) DeleteAfterGap() {
	g.ge = len(g.buf)
}

// Grow the size of gap by s
func (g *gapBuffer) GrowGap(s int) {
	// TODO: slices double on realloc up to 1024, then they increase by 25%.
	// we should check the cap() of our new slice after append and see if we
	// can use the new size to make a larger gap.
	b := make([]byte, s)

	g.buf = append(g.buf, b...)
	copy(g.buf[g.ge+s:], g.buf[g.ge:])
	g.gs = g.ge
	g.ge += s
}

// Move the gap to p.  p does not take the gap byteo account.
func (g *gapBuffer) MoveGap(p int) {
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

func (g *gapBuffer) Buffer() []byte {
	return g.buf
}

func (g *gapBuffer) GaplessBuffer() []byte {
	b := make([]byte, len(g.buf[:g.gs])+len(g.buf[g.ge:]))
	copy(b, g.buf[:g.gs])
	copy(b[g.gs:], g.buf[g.ge:])
	return b
	// return []byte(string(g.buf[:g.gs]) + string(g.buf[g.ge:]))
}

func (g *gapBuffer) String() string {
	return string(g.buf[:g.gs]) + string(g.buf[g.ge:])
}

func (g *gapBuffer) DebugString() string {
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
