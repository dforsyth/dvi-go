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
		Eb.line.UpdateCursor()
	}

	UpdateDisplay()
	for {
		k := Vw.win.Getch()

		if fn, ok := NCmdMap[k]; ok {
			fn()
			UpdateDisplay()
		}
	}
}

func NCursorLeft() {
	Eb.MoveCursorLeft()
}

func NCursorDown() {
	Eb.MoveCursorDown()
}

func NCursorUp() {
	Eb.MoveCursorUp()
}

func NCursorRight() {
	Eb.MoveCursorRight()
}
