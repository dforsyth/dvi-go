package main

/*
import (
	"curses"
)
*/

func AppendInsertMode() {
}

// insert mode
func InsertMode(gs *GlobalState) {

	gs.Mode = INSERT

	buffer := gs.CurrentBuffer.Value.(Interacter)

	if buffer == nil {
		panic("GlobalState has no CurrentBuffer in InsertMode")
	}

	m := NewInsertModeline()
	gs.SetModeline(m)
	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		switch k {
		case ESC:
			return
		default:
			buffer.SendInput(k)
		}
		m.Key = k
	}
}
