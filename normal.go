package main

// normal mode
func NormalMode(gs *GlobalState) {
	gs.Mode = NORMAL

	m := NewNormalModeline()
	gs.SetModeline(m)

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh // screen.Window.Getch()

		if k == int(EXPROMPT[0]) {
			exMode(gs)
			gs.Mode = NORMAL
			gs.SetModeline(m)
		} else {
			buffer := gs.curbuf.Value.(Buffer)
			buffer.SendInput(k)
			switch k {
			case 'i', 'a', 'o':
				InsertMode(gs)
				gs.Mode = NORMAL
				gs.SetModeline(m)
			case 'n':
				r := gs.NextBuffer()
				gs.queueMessage(&Message{
					"buffer: " + gs.curBuf().getIdent(),
					!r,
				})
			case 'p':
				r := gs.PrevBuffer()
				gs.queueMessage(&Message{
					"buffer: " + gs.curBuf().getIdent(),
					!r,
				})
			default:
			}
		}
		m.Key = k
	}
}
