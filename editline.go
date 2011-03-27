package main

type EditLine struct {
	b  *GapBuffer
	nl bool
	raw []byte
	dirty bool
}

func NewEditLine(s []byte) *EditLine {
	e := new(EditLine)
	e.b = NewGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		e.nl = true
	} else {
		e.nl = false
	}
	e.raw = e.b.GaplessBuffer()
	e.dirty = false
	return e
}

func (e *EditLine) InsertChar(c byte) {
	e.b.InsertChar(c)
	if c == '\n' {
		e.nl = true
	}
	e.dirty = true
}

func (e *EditLine) Delete(d int) {
	e.b.DeleteSpan(e.b.gs-1, d)
	e.dirty = true
}

func (e *EditLine) GetRaw() []byte {
	if e.dirty {
		e.raw = e.b.GaplessBuffer()
		e.dirty = false
	}
	return e.raw
}

func (e *EditLine) MoveCursor(p int) bool {
	if p < 0 || p > len(e.b.GaplessBuffer()) {
		return false
	}
	e.b.MoveGap(p)
	return true
}

func (e *EditLine) Cursor() int {
	return e.b.gs
}

