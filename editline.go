package main

type EditLine struct {
	eb     *EditBuffer // The buffer containing this line
	gb     *GapBuffer
	nl     bool
	raw    []byte
	dirty  bool
	indent int // index of first character
	length int
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
	e.raw = e.gb.GaplessBuffer()
	e.length = len(e.raw)
	e.dirty = false
	return e
}

func (e *EditLine) insertChar(c byte) {
	e.gb.insertChar(c)
	if c == '\n' {
		e.nl = true
	}
	e.dirty = true
}

func (e *EditLine) Delete(d int) {
	e.gb.DeleteSpan(e.gb.gs-1, d)
	e.dirty = true
}

func (e *EditLine) getRaw() []byte {
	if e.dirty {
		e.raw = e.gb.GaplessBuffer()
		e.dirty = false
	}
	return e.raw
}

func (e *EditLine) getLength() int {
	return len(e.getRaw())
}

func (e *EditLine) ClearToEOL() {
	e.gb.DeleteAfterGap()
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

func (e *EditLine) AfterCursor() []byte {
	return e.gb.AfterGap()
}

func (e *EditLine) BeforeCursor() []byte {
	return e.gb.BeforeGap()
}
