package main

import (
	"curses"
	"os"
)

type Terminal struct {
	fid     uint64
	ln, col uint64
	cache   map[uint64]string
	client  *Client
	x, y    int
	cwin    *curses.Window
	input chan int
}

func NewTerminal(client *Client) *Terminal {
	t := new(Terminal)
	t.client = client
	t.cache = make(map[uint64]string)
	t.input = make(chan int)
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
	t.x = *curses.Rows
	t.y = *curses.Cols
	t.ln = 0
	t.col = 0
}

func (t *Terminal) run() {
	go func() {
		for {
			t.input <-t.cwin.Getch()
		}
	}()

	for {
		t.display()
		k := <-t.input
		switch k {
		case 'o':
			if o, e := t.client.open("Makefile"); e == nil {
				t.fid = o.fid
			} else {
				panic(e.String())
			}
		case 'q':
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

func (t *Terminal) display() {
	t.clear()
	if t.fid == 0 {
		t.draw(0, 0, "No file...")
		for lno, text := range t.cache {
			t.draw(0, int(lno+1), text)
		}
		goto refresh
	}
	for y := 0; y < t.y; y++ {
		if ln, e := t.fetch(uint64(y)); e == nil {
			t.draw(0, y, ln)
		} else if e.String() == "noline" {
			t.draw(0, y, "~")
		}
	}
refresh:
	// refresh since we can't rely on getch
	t.cwin.Refresh()
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
