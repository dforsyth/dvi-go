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
			/*
				case 127, curses.KEY_BACKSPACE:
					// improperly handles the newline at the end of the prev line
					curr.backspace()
				case 0xd, 0xa:
					curr.newLine(byte('\n'))
				case 0x9:
					// currfilefer().InsertTab()
				default:
					curr.insertChar(byte(k))
			*/
		}
		m.Key = k
	}
}
