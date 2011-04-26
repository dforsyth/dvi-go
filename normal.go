package main

import (
	"fmt"
	"strconv"
	"unicode"
)

/* XXX eventually i'll need to modify these functions so they take a command struct or a variable
 * number of interfaces so that i can pass arguments.
 */

// normal commands
var normalFns map[int]func(*GlobalState) = map[int]func(*GlobalState){
	'j': normalj,
	'k': normalk,
	'l': normall,
	';': normalSemiColon,
	'p': normalp,
	'P': normalP, // will fix later
	'G': normalG,
	'u': normalu,
	'a': normala,
	'i': normali,
	'o': normalo,
	'O': normalO,
	// 'n':  nextBuffer,
	// 'p':  prevBuffer,
	':': normalColon,
	'-': normalMinus,
	'+': normalPlus,
	'#': normalHash,
	' ': normalSpace,
	'!': normalBang,
	'<': normalLShift,
	'>': normalRShift,
	'$': normalDollar,
	'0': normal0,
	1:   normalCtlA, // ^A
	2:   normalCtlB, // ^B
	// 3: normalCtlC, // ^C
	4:  normalCtlD,   // ^D
	5:  normalCtlE,   // ^E
	6:  normalCtlF,   // ^F
	7:  normalCtlG,   // ^G
	8:  normalCtlH,   // ^H
	9:  normalCtlI,   // ^I
	10: normalCtlJ,   // ^J
	11: normalCtlK,   // ^K
	12: normalCtlL,   // ^L
	13: normalCtlM,   // ^M
	16: normalCtlP,   // ^P
	20: normalCtlT,   // ^T
	21: normalCtlU,   // ^U
	23: normalCtlW,   // ^W
	25: normalCtlY,   // ^Y
	26: normalCtlZ,   // ^Z
	29: normalCtlRSB, // ^] (right square bracket)
	// x: normalCtlCaret
	ESC: cmdClear,
}

func normalj(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// left
		ln := b.line()
		if ln.cursor()-gs.n.cnt < 0 {
			ln.move(0)
			Beep()
		} else {
			ln.move(ln.cursor() - gs.n.cnt)
		}
	}
}

func normalk(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// down
		if b.lno+gs.n.cnt < len(b.lines)-1 {
			b.lno += gs.n.cnt
		} else {
			b.lno = len(b.lines) - 1
			Beep()
		}

		// TODO column needs to be maintained for the down/up commands (even if the line you
		// move to is not long enough).
	}
}

func normall(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// up
		if b.lno-gs.n.cnt > 0 {
			b.lno -= gs.n.cnt
		} else {
			b.lno = 0
			Beep()
		}
	}
}

func normalSemiColon(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// right
		ln := b.line()
		if ln.cursor()+gs.n.cnt < len(ln.raw()) {
			ln.move(ln.cursor() + gs.n.cnt)
		} else {
			ln.move(len(ln.raw()) - 1)
			Beep()
		}
	}
}

func normalp(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// paste
		gs.queueMessage(&Message{
			"paste.",
			false,
		})
	}
}

func normalP(gs *GlobalState) {
}

func normalG(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		b.lno = len(b.lines) - 1
	}
}

func normalu(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// rewind
	}
}

func normala(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// appendInput
		appendInsert(gs)
	}
}

func normali(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		insert(gs)
	}
}

func normalo(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		openInsert(gs)
	}
}

func normalO(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		aboveOpenInsert(gs)
	}
}

func normalColon(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		ex(gs)
	}
}

func normalMinus(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalPlus(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalHash(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalSpace(gs *GlobalState) {
	normalSemiColon(gs)
}

func normalBang(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalLShift(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalRShift(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalDollar(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		cnt := gs.n.cnt - 1
		if b.lno+cnt > len(b.lines)-1 {
			return
		}

		b.lno += cnt
		ln := b.line()
		ln.move(len(ln.raw()) - 1)
	}
}

func normal0(gs *GlobalState) {
}

func normalCtlA(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlB(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlD(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		// scroll down
	}
}

func normalCtlE(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlF(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlG(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		mod := "modified"
		if !b.isDirty() {
			mod = "un" + mod
		}
		// XXX This is actual not correct.  When the file is empty, we want to show "empty
		// file" rather than file position information.
		gs.queueMessage(&Message{
			fmt.Sprintf("%s: %s: line %d of %d [%d%]", b.ident(), mod, b.lno+1,
				len(b.lines), b.lno/len(b.lines)),
			false,
		})
	}
}

func normalCtlH(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlI(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlJ(gs *GlobalState) {
	normalj(gs)
}

func normalCtlK(gs *GlobalState) {
	normalk(gs)
}

func normalCtlL(gs *GlobalState) {
	// repaint
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlM(gs *GlobalState) {
	normalPlus(gs)
}

func normalCtlP(gs *GlobalState) {
	normalk(gs)
}

func normalCtlT(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlU(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlW(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

// XXX ^y and ^z are fucked
func normalCtlY(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlZ(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func normalCtlRSB(gs *GlobalState) {
	gs.queueMessage(&Message{
		"not implemented",
		true,
	})
}

func nextBuffer(gs *GlobalState) {
	r := gs.NextBuffer()
	gs.queueMessage(&Message{
		gs.curBuf().ident(),
		r == nil,
	})
}

func prevBuffer(gs *GlobalState) {
	r := gs.PrevBuffer()
	gs.queueMessage(&Message{
		gs.curBuf().ident(),
		r == nil,
	})
}

func cmdClear(gs *GlobalState) {
	gs.cmd = ""
	gs.n.cnt = 1
}

type Nm struct {
	buf string
	cnt int
}

// normal mode
func NormalMode(gs *GlobalState) {
	gs.Mode = MODENORMAL

	m := NewNormalModeline()
	gs.SetModeline(m)

	// advertise the current buffer
	gs.queueMessage(&Message{
		gs.curBuf().ident(),
		false,
	})

	gs.n = new(Nm)

	buf := ""
	gs.n.cnt = 1
	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh // screen.Window.Getch()

		if !unicode.IsDigit(k) {
			if len(buf) == 0 {
				gs.n.cnt = 1
				buf = string(k)
			} else {
				if cnt, e := strconv.Atoi(buf); e == nil {
					gs.n.cnt = cnt
					buf = string(k)
				}
			}

			if fn, ok := normalFns[k]; ok {
				fn(gs)
			}
			buf = ""
		} else {
			buf += string(k)
		}

		if gs.Mode != MODENORMAL {
			gs.Mode = MODENORMAL
			gs.SetModeline(m)
		}
		m.Key = k
		gs.curbuf.Value.(*EditBuffer).redraw = true
	}
}
