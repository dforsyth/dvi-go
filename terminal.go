package main

import (
	"curses"
	"os"
)

const (
	ESC = 27
	NOLINE = "~"

)

type Terminal struct {
	fid     uint64
	ln, col uint64
	cache   map[uint64]string
	client  *Client
	x, y    int
	cwin    *curses.Window
	k int
	ex bool
	exbuff string
}

func NewTerminal(client *Client) *Terminal {
	t := new(Terminal)
	t.client = client
	t.cache = make(map[uint64]string)
	t.ln, t.col = 0, 0
	t.fid = 0
	return t
}

func (t *Terminal) init() {
	curses.Initscr()
	curses.Cbreak()
	curses.Noecho()
	curses.Nonl()
	curses.Stdwin.Keypad(true)

	t.cwin = curses.Stdwin
	t.x = *curses.Cols
	t.y = *curses.Rows
	t.ln = 0
	t.col = 0
	t.ex = false
}

func (t *Terminal) run() {
	for {
		t.display()
		t.k = t.cwin.Getch()
		switch t.k {
		case 'o':
			if o, e := t.client.open("Makefile"); e == nil {
				t.fid = o.fid
			} else {
				panic(e.String())
			}
		case ':':
			t.basicEx()
		case ESC:
			if _, e := t.client.close(t.fid); e == nil {
				t.fid = 0
				for lno, _ := range t.cache {
					t.cache[lno] = "", false
				}
			} else {
				panic(e.String())
			}
		}
	}
}

func (t *Terminal) basicEx() {
	t.ex = true
	t.exbuff = ""
	for {
		t.display()
		k := t.cwin.Getch()
		if k == ESC {
			break
		}
		t.exbuff += string(k)
	}
	t.ex = false
}

func (t *Terminal) modeline() string {
	return string(t.k)
}

func (t *Terminal) display() {
	t.clear()
	if t.fid == 0 {
		t.draw(0, 0, "No file...")
		goto lastline
	}
	for y := 0; y < t.y-1; y++ {
		if ln, e := t.fetch(uint64(y)); e == nil {
			t.draw(0, y, ln)
		} else if e.String() == "noline" {
			t.draw(0, y, NOLINE)
		}
	}
lastline:
	if t.ex {
		t.draw(0, t.y-1, ":" + t.exbuff)
	} else {
		t.draw(0, t.y-1, t.modeline())
	}
}

func (t *Terminal) clear() {
	t.cwin.Clear()
}

func (t *Terminal) draw(x, y int, ln string) {
	t.cwin.Mvwaddnstr(y, x, ln, t.x)
}

func (t *Terminal) fetch(lno uint64) (string, os.Error) {
	if ln, ok := t.cache[lno]; ok {
		return ln, nil
	} else {
		if lr, e := t.client.line(t.fid, lno); e == nil {
			t.cache[lno] = lr.text
			return t.cache[lno], nil
		}
	}
	return "", &DviError{"noline"}
}
