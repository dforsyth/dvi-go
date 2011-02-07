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

	Eb.line.Value.(*Line).UpdateGap()

	oldMode := Ml.mode
	Ml.mode = "insert"

	UpdateDisplay()
	for {
		k := Vw.win.Getch()
		switch k {
		case 27:
			Ml.mode = oldMode
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
		Ml.char = k
		Ml.lno = int(Eb.line.Value.(*Line).lno)
		Ml.lco = int(Eb.lco)
		Ml.col = int(Eb.line.Value.(*Line).cursor)
		UpdateDisplay()
	}
}
