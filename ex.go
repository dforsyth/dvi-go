package main

/*
import (
	"curses"
	"os"
)
*/

func ExCmd(gs *GlobalState) {

	gs.SetModeline(gs.Command)

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		switch k {
		case ESC:
			return
		case 0xd, 0xa:
			gs.Command.Execute()
			return
		default:
			gs.Command.SendInput(k)
			/*
				case curses.KEY_BACKSPACE:
					if len(ex.buff.String()) == 0 {
						// vim behavior is to kill ex.  we beep.
						Beep()
					} else {
						ex.buff.deleteSpan(ex.buff.gs-1, 1)
					}
				case 0xd, 0xa:
					handleCmd(ex.buff.String())
					return
				default:
					ex.buff.insertChar(byte(k))
			*/
		}
	}
}

/*
func handleCmd(cmd string) {
	if cmd == "" {
		return
	}

	if cmd == "w" {
		go WriteFile(curr.name, curr)
		return
	}
	if cmd == "q" {
		// XXX make a real exit fn
		endScreen()
		os.Exit(0)
	}

	ml.mode = "Did not recognize: " + cmd
	screen.msg = ml
}
*/
