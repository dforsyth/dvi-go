package main

import (
	"fmt"
)

var cmdMap map[int]func() = map[int]func(){
	0: nil,
	'>': ExCmd,
	'i': InsertMode,
	'j': nil,
	'k': nil,
	'l': nil,
	';': nil,
}

// normal mode
func NormalMode() {

	if d.buf != nil && d.buf.Line() != nil {
		d.buf.Line().UpdateCursor()
	}

	for {
		Debug = ""
		k := d.view.win.Getch()

		if fn, ok := cmdMap[k]; ok {
			fn()
			Debug += fmt.Sprintf("(%s) normal: %x", string(k), k)
			UpdateDisplay()
		}
/*
		switch k {
		case 'i':
			Debug = "insert"
			d.InsertMode()
			if d.Buffer().Line() != nil {
				d.Buffer().Line().UpdateCursor()
			}
			Debug = "normal"
		case 'j':
			d.Buffer().MoveCursorLeft()
		case 'k':
			d.Buffer().MoveCursorDown()
		case 'l':
			d.Buffer().MoveCursorUp()
		case ';':
			d.Buffer().MoveCursorRight()
		case ':':
			Debug = "ex"
			d.ExCmd()
			Debug = "normal"
		}
*/
	}
}

func ExCmd() {
	ex := EXPROMPT
	UpdateLine(d.view.rows - 2, ex)
	cmdBuff := NewGapBuffer([]byte(""))
	for {
		k := d.view.win.Getch()

		switch k {
		case 27:
			return
		case 0x7f:
			if len(cmdBuff.String()) == 0 {
				/* vim behavior is to kill ex.  we beep. */
				Beep()
				continue
			} else {
				cmdBuff.DeleteSpan(cmdBuff.gs - 1, 1)
			}
		case 0xd:
			return
		default:
			cmdBuff.InsertChar(byte(k))
		}
		UpdateLine(d.view.rows - 2, ex + cmdBuff.String())
	}
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

