package main

import (
	"curses"
)

func AppendInsertMode() {
}

// insert mode
func InsertMode() {

	// we shouldn't hit these anymore, but if we do we should be ready to deal with them...
	if Eb == nil {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
	}

	if Eb.line == nil {
		Eb.AppendLine()
	}

	Eb.line.Value.(*EditLine).UpdateGap()

	oldMode := Ml.mode
	Ml.mode = "insert"

	UpdateDisplay()
	for {
		k := Vw.win.Getch()
		switch k {
		case ESC:
			Ml.mode = oldMode
			return
		case curses.KEY_BACKSPACE:
			// improperly handles the newline at the end of the prev line
			Eb.BackSpace()
		case 0xd, 0xa:
			Eb.NewLine(byte('\n'))
		case 0x9:
			// Ebfer().InsertTab()
		default:
			Eb.InsertChar(byte(k))
		}
		Ml.char = k
		Ml.lno = int(Eb.line.Value.(*EditLine).lno)
		Ml.lco = int(Eb.lco)
		Ml.col = int(Eb.line.Value.(*EditLine).cursor)
		UpdateDisplay()
	}
}
