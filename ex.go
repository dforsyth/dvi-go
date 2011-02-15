package main

import (
	"curses"
	"os"
)

func ExCmd() {
	if ex == nil {
		ex := new(Exline)
		ex.prompt = EXPROMPT
	}

	screen.SetMessage(ex) // switch to ex
	cmdBuff := NewGapBuffer([]byte(""))
	ex.command = cmdBuff.String()
	screen.RedrawMessage()
	for {
		k := screen.Window.Getch()

		switch k {
		case ESC:
			return
		case curses.KEY_BACKSPACE:
			if len(cmdBuff.String()) == 0 {
				/* vim behavior is to kill ex.  we beep. */
				Beep()
			} else {
				cmdBuff.DeleteSpan(cmdBuff.gs-1, 1)
			}
		case 0xd, 0xa:
			handleCmd(cmdBuff.String())
			return
		default:
			cmdBuff.InsertChar(byte(k))
		}
		ex.command = cmdBuff.String()
		screen.RedrawMessage()
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
	screen.SetMessage(ml)
}
