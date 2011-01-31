package main

type Line struct {
	gb *GapBuffer
	hasNewLine bool
	size int
	next, prev *Line
	cursor int
	mark int
}

func NewLine(s []byte) *Line {
	l := new(Line)
	l.gb = NewGapBuffer(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		l.hasNewLine = true
	} else {
		l.hasNewLine = false
	}
	l.size = len(s)
	l.cursor = 0
	return l
}

// Insert a character
func (l *Line) insertCharacter(c byte) {
	l.gb.InsertChar(c)
	l.size++
	if c == '\n' {
		l.hasNewLine = true
	}
	l.UpdateCursor()
}

// Get the bytes in this line
func (l *Line) bytes() []byte {
	return []byte(l.gb.GaplessBuffer())
}

// Backspace
func (l *Line) backspace() {
	l.gb.DeleteSpan(l.gb.gs - 1, 1)
	l.size--
	l.UpdateCursor()
}

// Move the cursor to p
func (l *Line) moveCursor(p int) int {
	if p < 0 || p > l.cursorMax() {
		return -1
	}

	l.cursor = p
	return l.cursor
}

func (l *Line) cursorMax() int {
	// dont allow the cursor to pass the newline char
	if l.hasNewLine {
		Message = "has new line"
		return l.size - 1
	}
	return l.size
}

// Mark at the cursor
func (l *Line) Mark() {
	l.mark = l.cursor
}

func (l *Line) delete(pos, len int) {
}

func (l *Line) UpdateCursor() {
	l.cursor = l.gb.gs
}

func (l *Line) UpdateGap() {
	l.gb.MoveGap(l.cursor)
}

func (l *Line) ClearAfterCursor() {
	l.gb.DeleteAfterGap()
}
