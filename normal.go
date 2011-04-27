package main

import (
	"strconv"
	"unicode"
)

/* XXX eventually i'll need to modify these functions so they take a command struct or a variable
 * number of interfaces so that i can pass arguments.
 */

type nmcmd struct {
	fn     func(*GlobalState)
	usage  string
	motion bool
}

// normal commands
var normalFns map[int]*nmcmd = map[int]*nmcmd{
	'h': &nmcmd{
		normalh,
		"[count]h",
		false,
	},
	'j': &nmcmd{
		normalj,
		"[count]j",
		false,
	},
	'k': &nmcmd{
		normalk,
		"[count]k",
		false,
	},
	'l': &nmcmd{
		normall,
		"[count]l",
		false,
	},
	';': &nmcmd{
		normalSemiColon,
		"[count];",
		false,
	},
	'p': &nmcmd{
		normalp,
		"",
		false,
	},
	'P': &nmcmd{
		normalP,
		"",
		false,
	},
	'G': &nmcmd{
		normalG,
		"[count]G",
		false,
	},
	'u': &nmcmd{
		normalu,
		"",
		false,
	},
	'a': &nmcmd{
		normala,
		"",
		false,
	},
	'i': &nmcmd{
		normali,
		"",
		false,
	},
	'o': &nmcmd{
		normalo,
		"",
		false,
	},
	'O': &nmcmd{
		normalO,
		"[count]O",
		false,
	},
	'^': &nmcmd{
		normalCaret,
		"^",
		false,
	},
	':': &nmcmd{
		normalColon,
		":",
		false,
	},
	'-': &nmcmd{
		normalMinus,
		"",
		false,
	},
	'+': &nmcmd{
		normalPlus,
		"",
		false,
	},
	'#': &nmcmd{
		normalHash,
		"",
		false,
	},
	' ': &nmcmd{
		normalSpace,
		"",
		false,
	},
	'!': &nmcmd{
		normalBang,
		"",
		false,
	},
	'<': &nmcmd{
		normalLShift,
		"",
		false,
	},
	'>': &nmcmd{
		normalRShift,
		"",
		false,
	},
	'$': &nmcmd{
		normalDollar,
		"",
		false,
	},
	'0': &nmcmd{
		normal0,
		"",
		false,
	},
	1: &nmcmd{
		normalCtlA,
		"",
		false,
	}, // ^A
	2: &nmcmd{
		normalCtlB,
		"",
		false,
	}, // ^B
	// 3: normalCtlC, // ^C
	4: &nmcmd{
		normalCtlD,
		"",
		false,
	}, // ^D
	5: &nmcmd{
		normalCtlE,
		"",
		false,
	}, // ^E
	6: &nmcmd{
		normalCtlF,
		"",
		false,
	}, // ^F
	7: &nmcmd{
		normalCtlG,
		"",
		false,
	}, // ^G
	8: &nmcmd{
		normalCtlH,
		"",
		false,
	}, // ^H
	9: &nmcmd{
		normalCtlI,
		"",
		false,
	}, // ^I
	10: &nmcmd{
		normalCtlJ,
		"",
		false,
	}, // ^J
	11: &nmcmd{
		normalCtlK,
		"",
		false,
	}, // ^K
	12: &nmcmd{
		normalCtlL,
		"",
		false,
	}, // ^L
	13: &nmcmd{
		normalCtlM,
		"",
		false,
	}, // ^M
	16: &nmcmd{
		normalCtlP,
		"",
		false,
	}, // ^P
	20: &nmcmd{
		normalCtlT,
		"",
		false,
	}, // ^T
	21: &nmcmd{
		normalCtlU,
		"",
		false,
	}, // ^U
	23: &nmcmd{
		normalCtlW,
		"",
		false,
	}, // ^W
	25: &nmcmd{
		normalCtlY,
		"",
		false,
	}, // ^Y
	26: &nmcmd{
		normalCtlZ,
		"",
		false,
	}, // ^Z
	29: &nmcmd{
		normalCtlRSB,
		"",
		false,
	}, // ^] (right square bracket)
	// x: normalCtlCaret
	ESC: &nmcmd{
		cmdClear,
		"",
		false,
	},
}

func normalh(gs *GlobalState) {
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

func normalj(gs *GlobalState) {
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

func normalk(gs *GlobalState) {
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

func normall(gs *GlobalState) {
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

func normalSemiColon(gs *GlobalState) {
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
		b.line().move(0)
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

func normalCaret(gs *GlobalState) {
	switch t := gs.curBuf(); b := t.(type) {
	case *EditBuffer:
		ln := b.line()
		raw := ln.raw()
		for i, c := range raw {
			if c != ' ' && c != '\t' {
				ln.move(i)
				return
			}
		}
		ln.move(len(raw) - 1)
	}
}

func normalColon(gs *GlobalState) {
	switch c := gs.curBuf(); b := c.(type) {
	case *EditBuffer:
		ex(gs)
	}
}

func normalMinus(gs *GlobalState) {
	normalk(gs)
}

func normalPlus(gs *GlobalState) {
	normalj(gs)
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
			Beep()
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
		gs.queueMessage(&Message{
			editBufferInfo(b),
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
	normalj(gs)
}

func normalCtlN(gs *GlobalState) {
	normalj(gs)
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
	cmd int
	cnt int
	mtn int
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

			if cmd, ok := normalFns[k]; ok {
				// XXX motion
				if cmd.motion {
					m := <-gs.InputCh
					gs.n.mtn = m
				}
				// " is going to have to steal input/ go into its own loop or
				// something
				cmd.fn(gs)
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
