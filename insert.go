package main

import (
	"curses"
)

func AppendInsertMode() {
}

// insert mode
func InsertMode() {

	// we shouldn't hit these anymore, but if we do we should be ready to deal with them...
	if curr == nil {
		curr = NewTempEditBuffer(TMPPREFIX)
	}

	if curr.l == nil {
		curr.appendLine()
	}

	// curr.line.Value.(*EditLine).UpdateGap()

	ml.mode = "insert"
	ml.name = curr.name
	screen.msg = ml // switch to modeline

	for {
		screen.RedrawAfter(0)
		screen.RedrawMessage()
		screen.RedrawCursor(curr.y, curr.x)
		screen.update <- 1
		k := <-inputch

		switch k {
		case ESC:
			return
		case 127, curses.KEY_BACKSPACE:
			// improperly handles the newline at the end of the prev line
			curr.backspace()
		case 0xd, 0xa:
			curr.newLine(byte('\n'))
		case 0x9:
			// currfilefer().InsertTab()
		default:
			curr.insertChar(byte(k))
		}
		ml.char = k
		// ml.lno = int(curr.line.Value.(*EditLine).lno)
		// ml.col = int(curr.l.Value.(*editLine).cursor)
	}
}
