package main

import (
	"fmt"
)


// insert mode
func InsertMode() {

	// we shouldn't hit these anymore, but if we do we should be ready to deal with them...
	if d.buf == nil {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
	}

	if d.buf.line == nil {
		d.buf.AppendLine()
	}

	d.buf.line.UpdateGap()

	UpdateDisplay()
	for {
		Debug = ""
		k := d.view.win.Getch()
		switch k {
		case 27:
			return
		case 0x7f:
			// improperly handles the newline at the end of the prev line
			d.buf.BackSpace()
		case 0xd, 0xa:
			d.buf.NewLine(byte('\n'))
		case 0x9:
			// d.Buffer().InsertTab()
		default:
			Debug = "adding char "
			d.buf.InsertChar(byte(k))
		}
		Debug += fmt.Sprintf("insert: %x", k)

		UpdateDisplay()
	}
}
