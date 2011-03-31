package main

type EditLine struct {
	b       *GapBuffer
	nl      bool
	raw     []byte
	dirty   bool
	indent  int // index of first character
	MapInfo []MapInfo
}

type MapInfo struct {
	ls, le int // position in line start and end
	ss, se int // position on screen start and end
}

func NewEditLine(s []byte) *EditLine {
	e := new(EditLine)
	e.b = NewGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		e.nl = true
	} else {
		e.nl = false
	}
	e.indent = 0
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

func (e *EditLine) ClearToEOL() {
	e.b.DeleteAfterGap()
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

func (e *EditLine) AfterCursor() []byte {
	return e.b.AfterGap()
}

func (e *EditLine) BeforeCursor() []byte {
	return e.b.BeforeGap()
}
