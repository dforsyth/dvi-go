package main

import (
	"curses"
	"os"
)

func ExCmd() {

	/* Create an ex message and set it as the current displayed message. */
	if ex == nil {
		ex := new(Exline)
		ex.prompt = EXPROMPT
	}

	oldMsg := screen.msg

	buff := newGapBuffer([]byte(""))
	ex.buff = buff
	screen.msg = ex

	for {
		screen.RedrawAfter(0)
		screen.RedrawMessage()
		screen.update <- 1
		k := <-inputch

		switch k {
		case ESC:
			/* Help the GC out. */
			buff = nil
			/* Put the old message line back in place. */
			screen.msg = oldMsg
			return
		case curses.KEY_BACKSPACE:
			if len(ex.buff.String()) == 0 {
				/* vim behavior is to kill ex.  we beep. */
				Beep()
			} else {
				ex.buff.deleteSpan(ex.buff.gs-1, 1)
			}
		case 0xd, 0xa:
			handleCmd(ex.buff.String())
			return
		default:
			ex.buff.insertChar(byte(k))
		}
	}
}

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
		/*
			if currfile.dirty {
				Ml.mode = "Unsaved changes in " + currfile.title
				return
			}
		*/
		endScreen()
		os.Exit(0)
	}

	ml.mode = "Did not recognize: " + cmd
	screen.msg = ml
}
