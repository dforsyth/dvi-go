package main

// XXX Some position pointers in these commands are NOT COPIED.  FIX!!!

import (
	"os"
)

type vicmd struct {
	// TODO: swap out all of these bools for masks
	fn        func(*CmdArgs) (*Position, os.Error)
	motion    bool // need motion
	isMotion  bool // is a motion command
	line      bool // in motion, this command used start and end line positions (must be isMotion)
	rw        bool // test writable
	zerocount bool // count default is zero instead of 1
}

type CmdArgs struct {
	d          *Dvi
	b          *Buffer
	k          int
	c1, c2     int
	start, end *Position
	line       bool // linemode
	motion     bool // using cmd as a motion
	buffer     byte
}

func cmdBackwards(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	for i := a.c1; i > 0 && p.off > 0; i-- {
		p = prevChar(*p)
	}
	return p, nil
}

func cmdForwards(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	for i := a.c1; i > 0 && p.off < len(p.line.text)-1; i-- {
		p = nextChar(*p)
	}
	return p, nil
}

func cmdUp(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	for i := a.c1; i > 0; i-- {
		n := prevLine(*p)
		if n.line == p.line {
			break
		}
		p = n
	}
	// XXX this needs to be visually oriented
	if p.off > len(p.line.text)-1 {
		p.off = len(p.line.text) - 1
		if p.off < 0 {
			p.off = 0
		}
	}
	return p, nil
}

func cmdDown(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	for i := a.c1; i > 0; i-- {
		n := nextLine(*p)
		if n.line == p.line {
			break
		}
		p = n
	}
	// XXX this needs to be visually oriented
	if p.off > len(p.line.text)-1 {
		p.off = len(p.line.text) - 1
		if p.off < 0 {
			p.off = 0
		}
	}
	return p, nil
}

func cmdInsert(a *CmdArgs) (*Position, os.Error) {
	s := *a.d.b.pos
	insertmode(a.d)
	if a.c1 > 1 {
		text := get(&s, a.d.b.pos)
		for i := a.c1; i > 1; i-- {
			a.d.b.pos = a.d.b.add(*a.d.b.pos, text)
		}
	}
	return a.d.b.pos, nil
}

func cmdAppend(a *CmdArgs) (*Position, os.Error) {
	if p := nextChar(*a.d.b.pos); p.line == a.d.b.pos.line {
		a.d.b.pos = p
	}
	return cmdInsert(a)
}

func cmdAppendEOL(a *CmdArgs) (*Position, os.Error) {
	eol(*a.d.b.pos)
	return cmdInsert(a)
}

func cmdEOL(a *CmdArgs) (*Position, os.Error) {
	eol(*a.d.b.pos)
	return a.d.b.pos, nil
}

func cmdBOL(a *CmdArgs) (*Position, os.Error) {
	p := *a.start
	p.off = 0
	return &p, nil
}

func cmdPrevWord(a *CmdArgs) (*Position, os.Error) {
	p := prevChar2(*a.start)
	if posEq(p, a.start) {
		// no where to go
		return p, nil
	}
	for i := a.c1; i > 0; i-- {
		for {
			// find a non-breakchar
			if c, e := p.getChar(); e == nil {
				if _, ok := breakchars[c]; !ok {
					break
				}
			}
			pp := prevChar2(*p)
			if posEq(p, pp) {
				return p, nil
			}
			p = pp
		}
		// prevChar to the first character of the word
		pp := prevChar2(*p)
		if posEq(p, pp) {
			return p, nil
		}
		for {
			if c, e := pp.getChar(); e == nil {
				if _, ok := breakchars[c]; ok {
					return p, nil
				}
			}
			p = pp
			pp = prevChar2(*pp)
			if posEq(p, pp) || pp.line != p.line {
				return p, nil
			}
		}
	}
	return nil, &DviError{}
}

func cmdPrevBigWord(a *CmdArgs) (*Position, os.Error) {
	return nil, nil
}

func cmdDelete(a *CmdArgs) (*Position, os.Error) {
	if buf, ok := a.d.buffers[a.buffer]; ok {
		y := get(a.start, a.end)
		buf.clear()
		buf.add(*(&Position{buf.first, 0}), y)
		buf.line = a.line
	} else {
		return nil, &DviError{"Buffer does not exist", 1}
	}
	p := a.d.b.remove(a.start, a.end, a.line)
	return p, nil
}

func cmdDeleteEOL(a *CmdArgs) (*Position, os.Error) {
	return nil, nil
}

func cmdEndOfWord(a *CmdArgs) (*Position, os.Error) {
	return nil, nil
}

func cmdEndOfBigWord(a *CmdArgs) (*Position, os.Error) {
	return nil, nil
}

func cmdToLine(a *CmdArgs) (*Position, os.Error) {
	l := a.d.b.first
	c := a.c1
	if c != 0 {
		for c > 1 && l != nil {
			l = l.next
			c--
		}
		if c == 1 {
			a.d.b.pos.line = l
			a.d.b.pos.off = 0
		}
	} else {
		for l.next != nil {
			l = l.next
		}
		a.d.b.pos.line = l
		a.d.b.pos.off = 0
	}
	// center line on screen
	return a.d.b.pos, nil
}

func cmdInsertLineBelow(a *CmdArgs) (*Position, os.Error) {
	l := NewLine([]byte{})
	a.d.b.insertLineBelow(a.d.b.pos.line, l)
	a.d.b.pos.line = l
	a.d.b.pos.off = 0

	s := a.d.b.pos
	insertmode(a.d)
	if a.c1 > 1 {
		text := get(s, a.d.b.pos)
		for i := a.c1; i > 1; i-- {
			l := NewLine([]byte{})
			a.d.b.insertLineBelow(a.d.b.pos.line, l)
			a.d.b.pos.line = l
			a.d.b.pos.off = 0
			a.d.b.pos = a.d.b.add(*a.d.b.pos, text)
		}
	}
	return a.d.b.pos, nil
}

func cmdInsertLineAbove(a *CmdArgs) (*Position, os.Error) {
	l := NewLine([]byte{})
	a.d.b.insertLineAbove(a.d.b.pos.line, l)
	a.d.b.pos.line = l
	a.d.b.pos.off = 0

	s := a.d.b.pos
	insertmode(a.d)
	if a.c1 > 1 {
		text := get(s, a.d.b.pos)
		for i := a.c1; i > 1; i-- {
			l := NewLine([]byte{})
			a.d.b.insertLineAbove(a.d.b.pos.line, l)
			a.d.b.pos.line = l
			a.d.b.pos.off = 0
			a.d.b.pos = a.d.b.add(*a.d.b.pos, text)
		}
	}
	return a.d.b.pos, nil
}

func cmdDeleteAtCursor(a *CmdArgs) (*Position, os.Error) {
	for i := 0; i < a.c1; i++ {
		a.d.b.pos = a.d.b.remove(a.d.b.pos, nextChar(*a.d.b.pos), a.line)
	}
	return a.d.b.pos, nil
}

func cmdDeleteBeforeCursor(a *CmdArgs) (*Position, os.Error) {
	for i := 0; i < a.c1; i++ {
		a.d.b.pos = a.d.b.remove(prevChar(*a.d.b.pos), a.d.b.pos, a.line)
	}
	return a.d.b.pos, nil
}

func cmdPut(a *CmdArgs) (*Position, os.Error) {
	buf := a.d.buffers[a.buffer]
	s := *a.start
	// TODO: support a.c1
	if buf.line {
		// add the new line position s at the beginning of it
		n := NewLine([]byte{})
		a.d.b.insertLineBelow(s.line, n)
		s.line = n
		s.off = 0
	} else {
		// for a non line buffer, move down one char before doing the insert
		s.off = nextChar(s).off
	}
	p := a.d.b.add(s, buf.getAll())
	// this command should put the cursor at the beginning of the put (one after a.start or at
	// the start of the first new line)
	if buf.line {
		p.off = 0
	} else {
		p = nextChar(*a.start)
	}
	return p, nil
}

func cmdNextWord(a *CmdArgs) (*Position, os.Error) {
	// TODO: Add next/prevword to position api
	// TODO: Account for error conditions
	p := a.start
	for i := a.c1; i > 0; i-- {
		for {
			// scan forward until we find a breakchar
			if c, e := p.getChar(); e == nil {
				if _, ok := breakchars[c]; ok {
					break
				}
			}
			n := nextChar2(*p)
			// newline is a break
			if n.line != p.line {
				break
			}
			if posEq(n, p) {
				// eob
				// beep?
				return p, nil
			}
			p = n
		}
		p = nextChar2(*p)
		for {
			// scan forward until we find the beginning of the next word
			if c, e := p.getChar(); e == nil {
				if _, ok := breakchars[c]; !ok {
					if i == 1 {
						return p, nil
					}
					break
				}
			}
			n := nextChar2(*p)
			if posEq(n, p) {
				if a.motion {
					return p, nil
				}
				// XXX place holder: eof, bail
				return nil, &DviError{}
			}
			p = n
		}
	}
	return nil, &DviError{}
}

func cmdYank(a *CmdArgs) (*Position, os.Error) {
	// XXX This needs to work correctly for [count]yy
	if buf, ok := a.d.buffers[a.buffer]; ok {
		y := get(a.start, a.end)
		buf.clear()
		buf.add(*(&Position{buf.first, 0}), y)
		buf.line = a.line
	}
	return a.start, nil
}

func cmdEx(a *CmdArgs) (*Position, os.Error) {
	exmode(a.d)
	return a.d.b.pos, nil
}

func cmdDisplayInfo(a *CmdArgs) (*Position, os.Error) {
	a.d.msg = &DviMessage{message: "fileinformation"}
	return a.d.b.pos, nil
}

var vicmds map[int]*vicmd = map[int]*vicmd{
	'$': &vicmd{
		fn: cmdEOL,
	},
	':': &vicmd{
		fn: cmdEx,
	},
	'0': &vicmd{
		fn: cmdBOL,
	},
	'a': &vicmd{
		fn:     cmdAppend,
		motion: false,
	},
	'A': &vicmd{
		fn:     cmdAppendEOL,
		motion: false,
	},
	'b': &vicmd{
		fn:       cmdPrevWord,
		isMotion: true,
	},
	'B': &vicmd{
		fn: cmdPrevBigWord,
	},
	'd': &vicmd{
		fn:     cmdDelete,
		motion: true,
	},
	'D': &vicmd{
		fn: cmdDeleteEOL,
	},
	'e': &vicmd{
		fn: cmdEndOfWord,
	},
	'E': &vicmd{
		fn: cmdEndOfBigWord,
	},
	'G': &vicmd{
		fn:        cmdToLine,
		zerocount: true,
	},
	'h': &vicmd{
		fn:       cmdBackwards,
		motion:   false,
		isMotion: true,
	},
	'i': &vicmd{
		fn:     cmdInsert,
		motion: false,
	},
	'j': &vicmd{
		fn:       cmdDown,
		motion:   false,
		isMotion: true,
		line:     true,
	},
	'k': &vicmd{
		fn:       cmdUp,
		motion:   false,
		isMotion: true,
		line:     true,
	},
	'l': &vicmd{
		fn:       cmdForwards,
		motion:   false,
		isMotion: true,
	},
	'o': &vicmd{
		fn: cmdInsertLineBelow,
	},
	'O': &vicmd{
		fn: cmdInsertLineAbove,
	},
	'p': &vicmd{
		fn:     cmdPut,
		motion: false,
	},
	'w': &vicmd{
		fn:       cmdNextWord,
		isMotion: true,
	},
	'x': &vicmd{
		fn: cmdDeleteAtCursor,
	},
	'X': &vicmd{
		fn: cmdDeleteBeforeCursor,
	},
	'y': &vicmd{
		fn:        cmdYank,
		motion:    true,
		zerocount: false,
	},
	ctrl('G'): &vicmd{
		fn: cmdDisplayInfo,
	},
}
