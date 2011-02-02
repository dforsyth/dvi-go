package main

import (
	"os"
)

func ExCmd() {
	ex := new(Exline)
	ex.prompt = EXPROMPT
	ex.command = ""
	cmdBuff := NewGapBuffer([]byte(""))
	UpdateModeLine(ex)
	for {
		k := vw.win.Getch()

		switch k {
		case 27:
			return
		case 0x7f:
			if len(cmdBuff.String()) == 0 {
				/* vim behavior is to kill ex.  we beep. */
				Beep()
			} else {
				cmdBuff.DeleteSpan(cmdBuff.gs-1, 1)
			}
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

	if cmd == "w" {
		go WriteEditBuffer(eb.title, eb)
	}
	if cmd == "q" {
		// XXX make a real exit fn
		endScreen()
		os.Exit(0)
	}
}
