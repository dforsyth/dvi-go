package main

// Move the cursor in an editbuffer to the right by one and then enter insert mode.
func appendInsertMode(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		ln := eb.lines[eb.lno]
		l := ln.getLength()
		c := ln.cursor()
		if c+1 == l {
			ln.moveCursor(l)
		} else {
			c++
			ln.moveCursor(c)
		}
		insertMode(gs)
	}
}

func insertMode(gs *GlobalState) {
	InsertMode(gs)
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
