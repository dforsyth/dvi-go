package main

import (
	"fmt"
)

var NCmdMap map[int]func() = map[int]func(){
	// 0: nil,
	int(EXPROMPT[0]): ExCmd,
	'i': InsertMode,
	'a': AppendInsertMode,
	'j': NCursorLeft,
	'k': NCursorDown,
	'l': NCursorUp,
	';': NCursorRight,
}

// normal mode
func NormalMode() {

	if eb != nil && eb.line != nil {
		eb.line.UpdateCursor()
	}

	UpdateDisplay()
	for {
		k := vw.win.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			Debug = fmt.Sprintf("(%s) normal: %x", string(k), k)
			UpdateDisplay()
		}
	}
}

func NCursorLeft() {
	eb.MoveCursorLeft()
}

func NCursorDown() {
	eb.MoveCursorDown()
}

func NCursorUp() {
	eb.MoveCursorUp()
}

func NCursorRight() {
	eb.MoveCursorRight()
}

