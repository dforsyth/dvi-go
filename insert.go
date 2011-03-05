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

	if curr.line == nil {
		curr.AppendLine()
	}

	curr.line.Value.(*EditLine).UpdateGap()

	screen.SetMessage(ml) // switch to modeline
	ml.mode = "insert"
	ml.name = curr.name

	screen.RedrawCursor(curr.CursorCoord())
	for {
		k := <-input // screen.Window.Getch()
		switch k {
		case ESC:
			return
		case 127, curses.KEY_BACKSPACE:
			// improperly handles the newline at the end of the prev line
			curr.BackSpace()
		case 0xd, 0xa:
			curr.NewLine(byte('\n'))
		case 0x9:
			// currfilefer().InsertTab()
		default:
			curr.Insert(byte(k))
		}
		ml.char = k
		ml.lno = int(curr.line.Value.(*EditLine).lno)
		ml.col = int(curr.line.Value.(*EditLine).cursor)
		screen.RedrawAfter(0)
		screen.RedrawMessage()
		screen.RedrawCursor(curr.CursorCoord())
		// screen.update <-1
		screen.Window.Refresh()
	}
}
