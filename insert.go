package main

func AppendInsertMode() {
}

// insert mode
func InsertMode() {

	// we shouldn't hit these anymore, but if we do we should be ready to deal with them...
	if eb == nil {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
	}

	if eb.line == nil {
		eb.AppendLine()
	}

	eb.line.UpdateGap()

	UpdateDisplay()
	for {
		Debug = ""
		k := vw.win.Getch()
		switch k {
		case 27:
			return
		case 0x7f:
			// improperly handles the newline at the end of the prev line
			eb.BackSpace()
		case 0xd, 0xa:
			eb.NewLine(byte('\n'))
		case 0x9:
			// ebfer().InsertTab()
		default:
			eb.InsertChar(byte(k))
		}
		UpdateDisplay()
	}
}
