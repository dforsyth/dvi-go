package main

// Move the cursor in an editbuffer to the right by one and then enter insert mode (and force a
// remap)
func appendInputMode(gs *GlobalState) {
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
		eb.dirty = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

// Input a line below the current line in an editbuffer, move down to the new line, then enter
// insert mode (and force a remap)
func openInputMode(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		eb.AppendEmptyLine()
		eb.moveDown(1) // move down to the new line...
		eb.dirty = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

// Input a line above the current line in an editbuffer, move up to the new line, then enter insert
// mode (and force a remap)
func aboveOpenInputMode(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		eb.insertEmptyLine(eb.lno)
		eb.dirty = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

func replaceOne(gs *GlobalState) {
	gs.Mode = MODEREPLACE
	input(gs)
}

func replaceMany(gs *GlobalState) {
	gs.Mode = MODEREPLACE
	input(gs)
}

// Input mode
func input(gs *GlobalState) {
	buffer := gs.curbuf.Value.(Buffer)

	if buffer == nil {
		Die("GlobalState has no curbuf in input")
	}

	m := NewInputModeline()
	gs.SetModeline(m)
	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		buffer.SendInput(k)
		if k == ESC {
			return
		}
		m.Key = k
		m.LineNumber = buffer.(*EditBuffer).lno
		m.ColumnNumber = cap(buffer.(*EditBuffer).lines)
	}
}
