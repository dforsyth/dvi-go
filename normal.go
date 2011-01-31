package main

import (
	"fmt"
)

var NCmdMap map[int]func() = map[int]func(){
	// 0: nil,
	'>': ExCmd,
	'i': InsertMode,
	'j': NCursorLeft,
	'k': NCursorDown,
	'l': NCursorUp,
	';': NCursorRight,
}

// normal mode
func NormalMode() {

	if d.buf != nil && d.buf.line != nil {
		d.buf.line.UpdateCursor()
	}

	for {
		Debug = ""
		k := d.view.win.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			Debug += fmt.Sprintf("(%s) normal: %x", string(k), k)
			UpdateDisplay()
		}
	}
}

func NCursorLeft() {
	d.buf.MoveCursorLeft()
}

func NCursorDown() {
	d.buf.MoveCursorDown()
}

func NCursorUp() {
	d.buf.MoveCursorUp()
}

func NCursorRight() {
	d.buf.MoveCursorRight()
}

func (d *D) NextEditBuffer() {

	n := d.buf.next
	if n != nil {
		d.buf = n
	}
}

func (d *D) PrevEditBuffer() {
	p := d.buf.prev
	if p != nil {
		d.buf = p
	}
}

