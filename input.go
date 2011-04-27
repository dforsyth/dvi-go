package main

import (
	"curses"
)

/*
type incmd struct {
	fn func(*GlobalState)
	usage string
}

var inputFns map[int]*incmd = map[int]*incmd {
}
*/

// Move the cursor in an editbuffer to the right by one and then enter insert mode (and force a
// remap)
func appendInsert(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		ln := eb.lines[eb.lno]
		l := len(ln.raw())
		c := ln.cursor()
		if c+1 == l {
			ln.moveCursor(l)
		} else {
			c++
			ln.moveCursor(c)
		}
		eb.redraw = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

// Input a line below the current line in an editbuffer, move down to the new line, then enter
// insert mode (and force a remap)
func openInsert(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		eb.AppendEmptyLine()
		eb.moveDown(1) // move down to the new line...
		eb.redraw = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

// Input a line above the current line in an editbuffer, move up to the new line, then enter insert
// mode (and force a remap)
func aboveOpenInsert(gs *GlobalState) {
	if eb, ok := gs.curbuf.Value.(*EditBuffer); ok {
		eb.insertEmptyLine(eb.lno)
		eb.redraw = true
		gs.Mode = MODEINSERT
		input(gs)
	}
}

func insert(gs *GlobalState) {
	if _, ok := gs.curbuf.Value.(*EditBuffer); ok {
		// dont really need to mark redraw on this one
		gs.Mode = MODEINSERT
		input(gs)
	}
}

func replaceOne(gs *GlobalState) {
	gs.Mode = MODEREPLACE
	input(gs)
}

func replaceMany(gs *GlobalState) {
	gs.Mode = MODEREPLACE
	input(gs)
}

var inputFns map[int]func(*EditBuffer) = map[int]func(*EditBuffer){
	ESC:                  inputEscape,
	curses.KEY_BACKSPACE: inputBackspace,
	127:                  inputBackspace,
	0xd:                  inputNewline,
	0xa:                  inputNewline,
}

func inputEscape(b *EditBuffer) {
	ln := b.line()
	ln.move(ln.cursor() - 1)
}

func inputBackspace(b *EditBuffer) {
	// TODO: this
	b.backspace()
}

func inputNewline(b *EditBuffer) {
	ln := b.line()
	if nl := ln.splitLn(ln.cursor()); nl != nil {
		ln.insert(byte('\n'))
		if len(nl.raw()) > 0 {
			b.insertLn(nl, b.lno+1)
		} else {
			b.insertLn(newEditLine([]byte("")), b.lno+1)
		}
		b.lno++
		b.dirty = true
		b.redraw = true
	} else {
		Beep()
	}
}


// Input mode
func input(gs *GlobalState) {
	buf := gs.curbuf.Value.(*EditBuffer)

	if buf == nil {
		Die("GlobalState has no curbuf in input")
	}

	m := NewInputModeline()
	gs.SetModeline(m)
	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		if fn, ok := inputFns[k]; ok {
			fn(buf)
		} else {
			ln := buf.line()
			if gs.Mode == MODEINSERT {
				ln.insert(byte(k))
				buf.dirty = true
			} else if gs.Mode == MODEREPLACE {
				ln.replace(byte(k))
				buf.dirty = true
			} else {
				gs.queueMessage(&Message{
					"no mode",
					true,
				})
			}
		}
		buf.redraw = true

		if k == ESC {
			return
		}

		m.Key = k
		m.LineNumber = buf.lno
		m.ColumnNumber = cap(buf.lines)
	}
}
