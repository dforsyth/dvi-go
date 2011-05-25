package main

// A simple gap buffer implementation on slices.

const (
	size = 64
	max  = 4096
)

type GapBuffer struct {
	buf    []byte
	gs, ge int
}

// Create a new gap buffer
func NewGapBuffer(t []byte) *GapBuffer {
	g := new(GapBuffer)
	g.buf = make([]byte, len(t)+size)

	copy(g.buf[:len(t)], t)

	// gs is the first index the gap, ge is the first index after the gap
	g.gs = len(t)
	g.ge = len(t) + size

	g.moveGap(0)
	return g
}

// Insert a character at the first position of the gap
func (g *GapBuffer) insertChar(c byte) {
	g.buf[g.gs] = c
	g.gs++
	if g.gs == g.ge {
		g.GrowGap(size)
	}
}

func (g *GapBuffer) InsertString(s string) {
	for _, c := range s {
		g.insertChar(byte(c))
	}
}

func (g *GapBuffer) DeleteSpan(p, s int) {
	g.moveGap(p + s)
	for i := 0; i < s; i++ {
		if g.gs == 0 {
			return
		}
		g.gs--
	}
}

func (g *GapBuffer) DeleteAfterGap() {
	g.ge = len(g.buf)
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

// move the gap to p.  p does not take the gap byteo account.
func (g *GapBuffer) moveGap(p int) {
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

func (g *GapBuffer) Buffer() []byte {
	return g.buf
}

func (g *GapBuffer) BeforeGap() []byte {
	return g.buf[:g.gs]
}

func (g *GapBuffer) AfterGap() []byte {
	return g.buf[g.ge:]
}

func (g *GapBuffer) GaplessBuffer() []byte {
	b := make([]byte, len(g.buf[:g.gs])+len(g.buf[g.ge:]))
	copy(b, g.BeforeGap())
	copy(b[g.gs:], g.AfterGap())
	return b
	// return []byte(string(g.buf[:g.gs]) + string(g.buf[g.ge:]))
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
