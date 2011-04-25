package main

/* XXX eventually i'll need to modify these functions so they take a command struct or a variable
 * number of interfaces so that i can pass arguments.
 */

// control commands in normal mode
var ebCmdMap map[int]func(*GlobalState) = map[int]func(*GlobalState){
	'a':  appendInsert,
	'i':  insert,
	'o':  openInsert,
	'O':  aboveOpenInsert,
	'n':  nextBuffer,
	'p':  prevBuffer,
	':':  ex,
	0x04: test, // ^D
	ESC:  cmdClear,
}

// commands on the editbuffer in normal mode
var normalFns map[int]func(*EditBuffer) = map[int]func(*EditBuffer){
	'j': left,
	'k': down,
	'l': up,
	';': right,
	'p': paste,
	'P': paste, // will fix later
	'G': maxLine,
	'u': undo,
}

func left(b *EditBuffer) {
	ln := b.line()
	if !ln.move(ln.cursor() - 1) {
		Beep()
	}
}

func down(b *EditBuffer) {
	if b.lno < len(b.lines)-1 {
		b.lno++
	} else {
		Beep()
	}
}

func up(b *EditBuffer) {
	if b.lno > 0 {
		b.lno--
	} else {
		Beep()
	}
}

func right(b *EditBuffer) {
	ln := b.line()
	if !ln.move(ln.cursor() + 1) {
		Beep()
	}
}

func paste(b *EditBuffer) {
	gs := b.gs
	gs.queueMessage(&Message{
		"paste.",
		false,
	})
}

func maxLine(b *EditBuffer) {
	b.lno = len(b.lines) - 1
}

func undo(b *EditBuffer) {
	// rewind
}

func test(gs *GlobalState) {
	gs.queueMessage(&Message{
		"this is a test",
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

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh // screen.Window.Getch()

		switch b := gs.curBuf(); t := b.(type) {
		case *EditBuffer:
			if fn, ok := ebCmdMap[k]; ok {
				fn(gs)
			} else {
				if fn, ok := normalFns[k]; ok {
					fn(gs.curbuf.Value.(*EditBuffer))
				}
			}
		}

		if gs.Mode != MODENORMAL {
			gs.Mode = MODENORMAL
			gs.SetModeline(m)
		}
		m.Key = k
		gs.curbuf.Value.(*EditBuffer).redraw = true
	}
}
