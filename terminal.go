package main

import (
	"container/list"
	"curses"
	"os"
	"strings"
)

const (
	ESC    = 27
	NOLINE = "~"
)

type Terminal struct {
	fid      uint64
	lno, col int
	cache    map[uint64]string
	client   *Client
	x, y     int
	cwin     *curses.Window
	k        int
	ex       bool
	exbuff   string
	q        *list.List
}

func NewTerminal(client *Client) *Terminal {
	t := new(Terminal)
	t.client = client
	t.cache = make(map[uint64]string)
	t.lno, t.col = 0, 0
	t.fid = 0
	t.q = list.New()
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
	t.ex = false
	t.q.Init()
}

func (t *Terminal) run() {
	t.basicNm()
}

func (t *Terminal) basicNm() {
	for {
		t.display()
		t.k = t.cwin.Getch()
		switch t.k {
		case ':':
			t.basicEx()
		case 'h':
			if t.col-1 >= 0 {
				t.col--
			}
		case 'j':
			if _, e := t.fetch(uint64(t.lno + 1)); e == nil {
				t.lno++
			}
		case 'k':
			if t.lno-1 >= 0 {
				t.lno--
				if s, _ := t.fetch(uint64(t.lno)); t.col > len(s)-1 {
					t.col = len(s) - 1
				}
			}
		case 'l':
			if s, _ := t.fetch(uint64(t.lno)); t.col+1 < len(s)-1 {
				t.col++
			}
		default:
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
		} else if k == curses.KEY_BACKSPACE || k == 127 {
			if len(t.exbuff) > 0 {
				t.exbuff = t.exbuff[:len(t.exbuff)-1]
			}
		} else if k == 0xa || k == 0xd {
			t.parseAndExecEx(t.exbuff)
			break
		} else {
			t.exbuff += string(k)
		}
	}
	t.ex = false
}

func (t *Terminal) parseAndExecEx(exbuff string) {
	exploded := strings.Split(exbuff, " ", -1)
	if len(exploded[0]) == 0 {
		return
	} else if len(exploded) == 1 && exploded[0] == "q" {
		if _, e := t.client.close(t.fid); e == nil {
			t.fid = 0
			for lno, _ := range t.cache {
				t.cache[lno] = "", false
			}
		} else {
			panic(e.String())
		}
	} else if len(exploded) >= 2 && exploded[0] == "e" {
		if o, e := t.client.open(exploded[1]); e == nil {
			t.fid = o.fid
		} else {
			t.qmsg(e.String())
		}
	}
}

func (t *Terminal) qmsg(msg string) {
	t.pushMsg(msg)
}

func (t *Terminal) pushMsg(msg string) {
	t.q.PushBack(msg)
}

func (t *Terminal) popMsg() string {
	if msg := t.q.Front(); msg != nil {
		t.q.Remove(msg)
		return msg.Value.(string)
	}
	return ""
}

func (t *Terminal) modeline() string {
	if msg := t.popMsg(); msg != "" {
		return msg
	}
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
		t.draw(0, t.y-1, ":"+t.exbuff)
	} else {
		t.draw(0, t.y-1, t.modeline())
		t.cwin.Move(t.lno, t.col)
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
		if lr, e := t.client.line(t.fid, lno, lno); e == nil {
			t.cache[lno] = lr.lnmap[lno]
			return t.cache[lno], nil
		}
	}
	return "", &DviError{"noline"}
}
