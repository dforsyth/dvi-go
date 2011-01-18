package main

import (
	"curses"
	"fmt"
)

var cmdMap map[int]func(*D) = map[int]func(*D){
	0: nil,
	'>': ExCmd,
	'i': InsertMode,
	'j': nil,
	'k': nil,
	'l': nil,
	';': nil,
}

// normal mode
func (d *D) NormalMode() {

	if d.Buffer() != nil && d.Buffer().Line() != nil {
		d.Buffer().Line().UpdateCursor()
	}

	for {
		Debug = ""
		k := curses.Stdwin.Getch()

		fn, ok := cmdMap[k]
		if !ok {
			continue
		}
		fn(d)
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
		Debug += fmt.Sprintf("(%s) normal: %x", string(k), k)

		d.UpdateDisplay()
	}
}

func (d *D) ExCmd() {
	ex := EXPROMPT
	UpdateLine(d.e, ex)
	cmdBuff := NewGapBuffer([]byte(""))
	for {
		k := d.win.Getch()

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
		UpdateLine(d.e, ex + cmdBuff.String())
	}
}

// not sure if i like this way better yet...
func ExCmd(d *D) {
	d.ExCmd()
}

func (d *D) NextEditBuffer() {

	n := d.buf.Next()
	if n != nil {
		d.buf = n
	}
}

func (d *D) PrevEditBuffer() {
	p := d.buf.Prev()
	if p != nil {
		d.buf = p
	}
}

