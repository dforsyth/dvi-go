package main

var NCmdMap map[int]func() = map[int]func(){
	// 0: nil,
	int(EXPROMPT[0]): ExCmd,
	'i':              InsertMode,
	'a':              AppendInsertMode,
	'j':              NCursorLeft,
	'k':              NCursorDown,
	'l':              NCursorUp,
	';':              NCursorRight,
}

// normal mode
func NormalMode() {

	if curr != nil && curr.line != nil {
		curr.line.Value.(*EditLine).UpdateCursor()
	}

	screen.SetMessage(ml) // switch to modeline
	screen.RedrawAfter(0)
	screen.RedrawCursor(curr.CursorCoord())
	screen.Window.Refresh()
	for {
		k := <-input // screen.Window.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			ml.lno = int(curr.line.Value.(*EditLine).lno)
			ml.col = int(curr.line.Value.(*EditLine).cursor)
			screen.RedrawAfter(0)
			screen.RedrawMessage()
		}
		screen.RedrawCursor(curr.CursorCoord())
		screen.Window.Refresh()
	}
}

func NCursorLeft() {
	curr.MoveLeft()
}

func NCursorDown() {
	curr.MoveDown()
}

func NCursorUp() {
	curr.MoveUp()
}

func NCursorRight() {
	curr.MoveRight()
}
