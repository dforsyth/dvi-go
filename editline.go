package main

type EditLine struct {
	b  *GapBuffer
	nl bool
}

func NewEditLine(s []byte) *EditLine {
	e := new(EditLine)
	e.b = NewGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		e.nl = true
	} else {
		e.nl = false
	}
	return e
}

func (e *EditLine) InsertChar(c byte) {
	e.b.InsertChar(c)
	if c == '\n' {
		e.nl = true
	}
}

func (e *EditLine) Delete(d int) {
	e.b.DeleteSpan(e.b.gs-1, d)
}

func (e *EditLine) raw() []byte {
	return []byte(e.b.GaplessBuffer())
}

func (e *EditLine) moveCursor(p int) bool {
	if p < 0 || p > len(e.b.GaplessBuffer()) {
		return false
	}
	e.b.MoveGap(p)
	return true
}
