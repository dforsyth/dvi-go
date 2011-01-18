package main

func ExCmd() {
	ex := EXPROMPT
	UpdateLine(d.view.rows - 2, ex)
	cmdBuff := NewGapBuffer([]byte(""))
	for {
		k := d.view.win.Getch()

		switch k {
		case 27:
			return
		case 0x7f:
			if len(cmdBuff.String()) == 0 {
				/* vim behavior is to kill ex.  we beep. */
				Beep()
				continue
			} else {
				cmdBuff.DeleteSpan(cmdBuff.gs - 1, 1)
			}
		case 0xd:
			return
		default:
			cmdBuff.InsertChar(byte(k))
		}
		UpdateLine(d.view.rows - 2, ex + cmdBuff.String())
	}
}

