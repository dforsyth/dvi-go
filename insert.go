package main

func AppendInsertMode() {
}

// insert mode
func InsertMode(gs *GlobalState) {

	gs.Mode = INSERT

	buffer := gs.curbuf.Value.(Buffer)

	if buffer == nil {
		panic("GlobalState has no curbuf in InsertMode")
	}

	m := NewInsertModeline()
	gs.SetModeline(m)
	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		buffer.SendInput(k)
		switch k {
		case ESC:
			return
		default:
		}
		m.Key = k
		m.LineNumber = buffer.(*EditBuffer).lno
		m.ColumnNumber = cap(buffer.(*EditBuffer).lines)
	}
}
