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

	//if curr != nil && curr.l != nil {
	//	curr.l.Value.(*editLine).UpdateCursor()
	//}

	screen.msg = ml // switch to modeline
	for {
		screen.RedrawAfter(0)
		screen.RedrawMessage()
		screen.RedrawCursor(curr.x, curr.y)
		screen.update <- 1
		k := <-inputch // screen.Window.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			// ml.lno = int(curr.line.Value.(*EditLine).lno)
			ml.col = int(curr.l.Value.(*editLine).b.gs)
		}
	}
}

func NCursorLeft() {
	curr.moveLeft()
}

func NCursorDown() {
	curr.moveDown()
}

func NCursorUp() {
	curr.moveUp()
}

func NCursorRight() {
	curr.moveRight()
}
