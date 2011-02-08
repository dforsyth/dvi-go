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

	if Eb != nil && Eb.line != nil {
		Eb.line.Value.(*EditLine).UpdateCursor()
	}

	UpdateDisplay()
	for {
		k := Vw.win.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			Ml.lno = int(Eb.line.Value.(*EditLine).lno)
			Ml.lco = int(Eb.lco)
			Ml.col = int(Eb.line.Value.(*EditLine).cursor)
			UpdateDisplay()
		}
	}
}

func NCursorLeft() {
	Eb.MoveLeft()
}

func NCursorDown() {
	Eb.MoveDown()
}

func NCursorUp() {
	Eb.MoveUp()
}

func NCursorRight() {
	Eb.MoveRight()
}
