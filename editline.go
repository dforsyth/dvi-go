package main

type EditLine struct {
	lno        uint
	gb         *GapBuffer
	hasNewLine bool
	size       int
	cursor     int
	mark       int
}

func (l *EditLine) Screen() *Screen {
	return screen
}

func (l *EditLine) ScreenLines() int {
	return len(l.raw()) / l.Screen().Cols
}

func (l *EditLine) DisplayLength() int {
	if l.hasNewLine {
		return l.size - 1
	}
	return l.size
}

func NewLine(s []byte) *EditLine {
	l := new(EditLine)
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
func (l *EditLine) insertCharacter(c byte) {
	l.gb.InsertChar(c)
	l.size++
	if c == '\n' {
		l.hasNewLine = true
	}
	l.UpdateCursor()
}

// Get the bytes in this line
func (l *EditLine) raw() []byte {
	return []byte(l.gb.GaplessBuffer())
}

// Backspace
func (l *EditLine) backspace() {
	l.gb.DeleteSpan(l.gb.gs-1, 1)
	l.size--
	l.UpdateCursor()
}

// Move the cursor to p
func (l *EditLine) moveCursor(p int) int {
	if p < 0 || p > l.cursorMax() {
		return -1
	}

	l.cursor = p
	return l.cursor
}

func (l *EditLine) cursorMax() int {
	// dont allow the cursor to pass the newline char
	if l.hasNewLine {
		return l.size - 1
	}
	return l.size
}

// Mark at the cursor
func (l *EditLine) Mark() {
	l.mark = l.cursor
}

func (l *EditLine) delete(pos, length int) {
}

func (l *EditLine) UpdateCursor() {
	l.cursor = l.gb.gs
}

func (l *EditLine) UpdateGap() {
	l.gb.MoveGap(l.cursor)
}

func (l *EditLine) ClearAfterCursor() {
	l.gb.DeleteAfterGap()
}
