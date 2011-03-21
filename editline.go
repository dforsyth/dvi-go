package main

type EditLine struct {
	b  *gapBuffer
	nl bool
}

func newEditLine(s []byte) *EditLine {
	e := new(EditLine)
	e.b = newGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		e.nl = true
	} else {
		e.nl = false
	}
	return e
}

func (e *EditLine) insertChar(c byte) {
	e.b.insertChar(c)
	if c == '\n' {
		e.nl = true
	}
}

func (e *EditLine) delete(d int) {
	e.b.deleteSpan(e.b.gs-1, d)
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
