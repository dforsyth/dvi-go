package main

import ()

type vicmd struct {
	fn        func(*CmdArgs)
	motion    bool // need motion
	rw        bool // test writable
	zerocount bool
}

type CmdArgs struct {
	s      *State
	c1, c2 int
}

func cmdBackwards(a *CmdArgs) {
	f := a.s.f
	for i := a.c1; i > 0; i-- {
		f.pos = prevChar(*f.pos)
	}
}

func cmdForwards(a *CmdArgs) {
	f := a.s.f
	for i := a.c1; i > 0; i-- {
		f.pos = nextChar(*f.pos)
	}
}

func cmdUp(a *CmdArgs) {
	f := a.s.f
	for i := a.c1; i > 0; i-- {
		f.pos = prevLine(*f.pos)
	}
}

func cmdDown(a *CmdArgs) {
	f := a.s.f
	for i := a.c1; i > 0; i-- {
		f.pos = nextLine(*f.pos)
	}
}

func cmdInsert(a *CmdArgs) {
	insertmode(a.s)
}

func cmdAppend(a *CmdArgs) {
	f := a.s.f
	f.pos = nextChar(*f.pos)
	insertmode(a.s)
	for i := a.c1; i > 0; i-- {
		// append what happened a.c1 times...
	}
}

func cmdAppendEOL(a *CmdArgs) {
	eol(a.s.f)
	insertmode(a.s)
}

func cmdEOL(a *CmdArgs) {
	eol(a.s.f)
}

func cmdPrevWord(a *CmdArgs) {
}

func cmdPrevBigWord(a *CmdArgs) {
}

func cmdDelete(a *CmdArgs) {
}

func cmdDeleteEOL(a *CmdArgs) {
}

func cmdEndOfWord(a *CmdArgs) {
}

func cmdEndOfBigWord(a *CmdArgs) {
}

func cmdEx(a *CmdArgs) {
	exmode(a.s)
}

var vicmds map[int]*vicmd = map[int]*vicmd{
	'$': &vicmd{
		fn: cmdEOL,
	},
	':': &vicmd{
		fn: cmdEx,
	},
	'a': &vicmd{
		fn:        cmdAppend,
		motion:    false,
		zerocount: true,
	},
	'A': &vicmd{
		fn:        cmdAppendEOL,
		motion:    false,
		zerocount: true,
	},
	'b': &vicmd{
		fn: cmdPrevWord,
	},
	'B': &vicmd{
		fn: cmdPrevBigWord,
	},
	'd': &vicmd{
		fn: cmdDelete,
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
	'h': &vicmd{
		fn:     cmdBackwards,
		motion: false,
	},
	'i': &vicmd{
		fn:     cmdInsert,
		motion: false,
	},
	'j': &vicmd{
		fn:     cmdDown,
		motion: false,
	},
	'k': &vicmd{
		fn:     cmdUp,
		motion: false,
	},
	'l': &vicmd{
		fn:     cmdForwards,
		motion: false,
	},
}
