package main

// XXX Some position pointers in these commands are NOT COPIED.  FIX!!!

import (
	"curses" // for regex
	"os"
	"regexp"
	"unicode"
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

/* XXX not sure if this is a good idea
type ViCmd interface {
	Function() func(*CmdArgs) (*Position, os.Error)
	Motion() bool
	IsMotion() bool
	LineMode() bool
	ZeroCount() bool
}
*/

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

func resetCmdArgs(a *CmdArgs) {
	a.d = nil
	a.c1 = 0
	a.c2 = 0
	a.start = nil
	a.end = nil
	a.buffer = 0
	a.line = false
}

/* XXX wrapper to force cursor fixes
func doViCmd(cmd *vicmd, a *CmdArgs) (*Position, os.Error) {
	p, e := cmd.fn(a)
	if p != nil {
		p = fixCursor(p)
	}
	return p, e
}
*/

func fixCursor(pos *Position) *Position {
	for pos.off > pos.line.length()-1 {
		pos = prevChar(*pos)
	}
	return pos
}

func cmdCurrLineAndAbove(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	// TODO: If there are less than a.c1-1 line after p.line in the buffer, it's an error.
	for i := 0; i < a.c1-1; i++ {
		p = nextLine(*p)
		// TODO: go to first non-blank
	}
	return p, nil
}

func cmdMoveToColumn(a *CmdArgs) (*Position, os.Error) {
	p := a.start
	c := a.c1
	if c > p.line.length() {
		c = p.line.length()
	}
	// c-1 because positon.off is 0 based
	c--
	if !a.motion {
		return &Position{p.line, c}, nil
	} else {
		if p.line.length() == 0 || p.off+1 == a.c1 {
			// TODO: if the line is empty or the cursor is at the countth position in the
			// current line, it shall be an error.
			return nil, &DviError{}
		} else {
			return &Position{p.line, c}, nil
		}
	}
	return nil, &DviError{}
}

func cmdFirstNonBlank(a *CmdArgs) (*Position, os.Error) {
	p := &Position{a.start.line, 0}

	for n := nextChar(*p); ; n = nextChar(*p) {
		if c, e := p.getChar(); e == nil {
			if _, ok := blankchars[c]; !ok {
				break
			}
		}
		if posEq(n, p) {
			if a.motion {
				// XXX This doesn't seem to be nvi behavior, but the spec says this 
				// should be an error...
				return nil, &DviError{}
			}
			break
		}
		p = n
	}
	return fixCursor(p), nil
}

type REMessage struct {
	re *string
}

func (m *REMessage) Message() string {
	return "/" + *m.re
}

func (m *REMessage) Color() int {
	return 0
}

func (m *REMessage) Beep() bool {
	return false
}

func cmdFindRegex(a *CmdArgs) (*Position, os.Error) {
	// TODO: implement this properly
	// TODO: pull code out into shared find op, probably on buffer
	// TODO: scroll

	p := nextChar(*a.start)
	msg := &REMessage{}
	re := ""
	msg.re = &re
	for {
		// XXX this is so ghetto.  really need to fix up status/messaging
		a.d.msg = msg
		draw(a.d)
		switch k := getCh(a.d); k {
		case 0xd, 0xa, curses.KEY_ENTER:
			if r, e := regexp.Compile(re); e == nil {
				wrap := false
				var i []int = nil
				for i = r.FindIndex(p.line.text[p.off:]); i == nil; i = r.FindIndex(p.line.text[p.off:]) {
					if p.line == a.start.line && wrap {
						return nil, &DviError{"Pattern not found", 0}
					}
					if p.line == a.d.b.last {
						wrap = true
						p.line = a.d.b.first
						p.off = 0
					} else {
						p = nextLine(*p)
						p.off = 0
					}
				}
				p.off = i[0] + p.off
				if wrap {
					a.d.queueMsg("Search wrapped", 0, false)
				}
				// XXX if motion, we need to determine whether or not this is in line or char mode
			} else {
				return nil, &DviError{"regex compile failed: " + e.String(), 0}
			}
			return p, nil
		default:
			re += string([]byte{byte(k)})
		}
	}

	return p, &DviError{"Not reached", 0}
}

func cmdReverseCase(a *CmdArgs) (*Position, os.Error) {
	p := &Position{a.start.line, a.start.off}
	for i := 0; i < a.c1; p, i = nextChar2(*p), i+1 {
		if p.off == p.line.length() {
			// work around: end of line position doesn't count.
			i--
			continue
		}
		if c, e := p.getChar(); e == nil {
			if unicode.IsLower(c) {
				p.setChar(unicode.ToUpper(c))
			} else {
				p.setChar(unicode.ToLower(c))
			}
			// really shouldn't mark dirty if it's a symbol, but meh.
			a.d.b.dirty = true
		}
	}
	return fixCursor(p), nil
}

func cmdShiftLeft(a *CmdArgs) (*Position, os.Error) {
	return nil, nil
}

func cmdShiftRight(a *CmdArgs) (*Position, os.Error) {
	if a.d.b.lineNumber(a.start.line)+a.c1-1 > a.d.b.lineCount() {
		return nil, &DviError{}
	}
	p := &Position{a.start.line, 0}
	for ; p.line != a.end.line.next; p = nextLine(*p) {
		a.d.b.add(*p, []byte{'\t'})
	}
	return &Position{a.start.line, a.start.off + 1}, nil
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
	if a.start.off == 0 {
		return nil, &DviError{}
	}
	return &Position{a.start.line, 0}, nil
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
	return nil, &DviError{"Not yet implemented", 0}
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
	lno := a.d.b.lineNumber(a.start.line)
	a.d.b.remove(*a.start, *a.end, a.line)
	ln := a.d.b.getLine(lno)
	off := a.start.off
	if a.line {
		off = 0
	}
	return &Position{ln, off}, nil
}

func cmdDeleteEOL(a *CmdArgs) (*Position, os.Error) {
	return nil, &DviError{"Not yet implemented", 0}
}

func cmdEndOfWord(a *CmdArgs) (*Position, os.Error) {
	return nil, &DviError{"Not yet implemented", 0}
}

func cmdEndOfBigWord(a *CmdArgs) (*Position, os.Error) {
	return nil, &DviError{"Not yet implemented", 0}
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
	p := a.d.b.pos
	for i := 0; i < a.c1; i++ {
		pp := prevChar(*p)
		a.d.b.remove(*p, *nextChar(*p), a.line)
		if p.off > p.line.length()-1 {
			p = pp
		}
	}

	return p, nil
}

func cmdDeleteBeforeCursor(a *CmdArgs) (*Position, os.Error) {
	p := a.d.b.pos
	for i := 0; i < a.c1; i++ {
		pp := prevChar(*p)
		a.d.b.remove(*prevChar(*p), *p, a.line)
		p = pp
	}
	return p, nil
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
	// TODO: fix implementation
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
	a.d.queueMsg(a.d.b.information(), 1, false)
	return a.d.b.pos, nil
}
