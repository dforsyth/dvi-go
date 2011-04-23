package main

import ()

type EditLine struct {
	eb     *EditBuffer // The buffer containing this line
	gb     *GapBuffer
	nl     bool
	dirty  bool
	indent int // index of first character
}

func newEditLine(s []byte) *EditLine {
	return NewEditLine(s)
}

func NewEditLine(s []byte) *EditLine {
	e := new(EditLine)
	e.eb = nil
	e.gb = NewGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		e.nl = true
	} else {
		e.nl = false
	}
	e.indent = 0
	e.dirty = false
	return e
}

func (el *EditLine) splitLn(pos int) *EditLine {
	if el.move(pos) {
		nl := newEditLine(el.afterCur())
		el.clearToEOL()
		return nl
	}
	return nil
}

func (e *EditLine) insert(c byte) {
	e.gb.insertChar(c)
	if c == '\n' {
		e.nl = true
	}
	e.dirty = true
}

func (e *EditLine) replace(c byte) {
}

func (e *EditLine) Delete(d int) {
	e.gb.DeleteSpan(e.gb.gs-1, d)
	e.dirty = true
}

func (e *EditLine) raw() []byte {
	return e.gb.GaplessBuffer()
}

func (e *EditLine) clearToEOL() {
	e.gb.DeleteAfterGap()
}

func (el *EditLine) move(pos int) bool {
	return el.moveCursor(pos)
}

// XXX This still lets the client pass the visual eol when the editbuffer is in normal mode...
func (e *EditLine) moveCursor(p int) bool {
	max := len(e.gb.GaplessBuffer())
	if e.nl {
		max -= 1
	}
	if p < 0 || p > max {
		return false
	}
	e.gb.moveGap(p)
	return true
}

func (e *EditLine) Cursor() int {
	return e.gb.gs
}

func (el *EditLine) cursor() int {
	return el.gb.gs
}

func (e *EditLine) afterCur() []byte {
	return e.gb.AfterGap()
}

func (e *EditLine) BeforeCursor() []byte {
	return e.gb.BeforeGap()
}
