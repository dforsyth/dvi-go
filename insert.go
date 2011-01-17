package main

import (
	"curses"
	"fmt"
)

// insert mode
func (d *D) InsertMode() {

	if d.Buffer() == nil {
		d.InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
	}

	if d.Buffer().Line() == nil {
		d.Buffer().AppendLine()
	}

	d.Buffer().Line().MoveGapToCursor()

	d.UpdateDisplay()
	for {
		Debug = ""
		k := curses.Stdwin.Getch()
		switch k {
		case 27 :
			return
		case 0x7f:
			// improperly handles the newline at the end of the prev line
			d.Buffer().BackSpace()
		case 0xd, 0xa:
			d.Buffer().InsertChar(byte('\n'))
			d.Buffer().InsertLine(NewGapBuffer([]byte("")))
		case 0x9:
			// d.Buffer().InsertTab()
		default:
			Debug = "adding char "
			d.Buffer().InsertChar(byte(k))
		}
		Debug += fmt.Sprintf("insert: %x", k)

		if d.Buffer().Line() != nil {
			d.Buffer().Line().UpdateCursor()
		}
		d.UpdateDisplay()
	}
}

