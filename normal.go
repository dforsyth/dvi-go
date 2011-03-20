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
			ExCmd(gs)
			gs.Mode = NORMAL
			gs.SetModeline(m)
		} else if k == int('i') {
			InsertMode(gs)
			gs.Mode = NORMAL
			gs.SetModeline(m)
		} else {
			buffer := gs.CurrentBuffer.Value.(Interacter)
			buffer.SendInput(k)
		}
		m.Key = k
	}
}
