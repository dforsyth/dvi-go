package main

import (
	"curses"
	"os"
)

func ExCmd() {
	ex := new(Exline)
	ex.prompt = EXPROMPT
	ex.command = ""
	cmdBuff := NewGapBuffer([]byte(""))
	UpdateModeLine(ex)
	for {
		k := Vw.win.Getch()

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
		// case curses.KEY_ENTER:
		case 0xd:
			handleCmd(cmdBuff.String())
			return
		default:
			cmdBuff.InsertChar(byte(k))
		}
		ex.command = cmdBuff.String()
		UpdateModeLine(ex)
	}
}

func handleCmd(cmd string) {
	if cmd == "" {
		return
	}

	if cmd == "w" {
		go WriteEditBuffer(Eb.title, Eb)
		return
	}
	if cmd == "q" {
		// XXX make a real exit fn
		/*
			if Eb.dirty {
				Ml.mode = "Unsaved changes in " + Eb.title
				return
			}
		*/
		endScreen()
		os.Exit(0)
	}

	Ml.mode = "Did not recognize: " + cmd
}
