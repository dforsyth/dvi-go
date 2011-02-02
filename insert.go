package main

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

	Eb.line.UpdateGap()

	UpdateDisplay()
	for {
		DEbug = ""
		k := Vw.win.Getch()
		switch k {
		case 27:
			return
		case 0x7f:
			// improperly handles the newline at the end of the prev line
			Eb.BackSpace()
		case 0xd, 0xa:
			Eb.NewLine(byte('\n'))
		case 0x9:
			// Ebfer().InsertTab()
		default:
			Eb.InsertChar(byte(k))
		}
		UpdateDisplay()
	}
}
