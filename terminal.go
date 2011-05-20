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
	fid        uint64
	lno, col   uint64
	upd, cache map[uint64]string
	client     *Client
	x, y       int
	cwin       *curses.Window
	k          int
	ex         bool
	exbuff     string
	q          *list.List
}

func NewTerminal(client *Client) *Terminal {
	t := new(Terminal)
	t.client = client
	t.cache = make(map[uint64]string)
	t.upd = make(map[uint64]string)
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
		case 'i':
			if t.fid > 0 {
				t.basicIn()
			}
		case 'h':
			if t.col-1 >= 0 {
				t.col--
			}
		case 'j':
			if _, e := t.fetch(t.lno + 1); e == nil {
				t.lno++
			}
		case 'k':
			if t.lno-1 >= 0 {
				t.lno--
				if s, _ := t.fetch(t.lno); t.col > uint64(len(s)-1) {
					t.col = uint64(len(s) - 1)
				}
			}
		case 'l':
			if s, _ := t.fetch(t.lno); t.col+1 < uint64(len(s)-1) {
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
		switch k := t.cwin.Getch(); k {
		case ESC:
			goto exit
		case curses.KEY_BACKSPACE, 127:
			if len(t.exbuff) > 0 {
				t.exbuff = t.exbuff[:len(t.exbuff)-1]
			}
		case 0xa, 0xd:
			t.parseAndExecEx(t.exbuff)
			goto exit
		default:
			t.exbuff += string(k)
		}
	}
exit:
	t.ex = false
}

func (t *Terminal) basicIn() {
	for {
		t.display()
		switch k := t.cwin.Getch(); k {
		case ESC:
			if _, e := t.client.update(t.fid, t.upd); e == nil {
				for k, _ := range t.upd {
					t.upd[k] = "", false
				}
			} else {
				t.pushMsg("update failed")
			}
			goto exit
		default:
			s, _ := t.fetch(t.lno)
			l := s[:t.col] + string(k) + s[t.col:]
			t.update(t.lno, l)
			t.col++
		}
	}
exit:
	return
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
	} else if len(exploded) == 1 && exploded[0] == "w" {
		if _, e := t.client.update(t.fid, t.upd); e != nil {
			panic(e.String())
		} else {
			for k, _ := range t.upd {
				t.upd[k] = "", false
			}
		}
		if s, e := t.client.sync(t.fid); e != nil {
			panic(e.String())
		} else {
			t.qmsg(s.message())
		}
	}
}

func (t *Terminal) qmsg(msg string) {
	t.pushMsg(msg)
}

func (t *Terminal) pushMsg(msg string) {
	t.q.PushBack(msg)
}

func (t *Terminal) nextMsg() string {
	if msg := t.q.Front(); msg != nil {
		t.q.Remove(msg)
		return msg.Value.(string)
	}
	return ""
}

func (t *Terminal) modeline() string {
	if msg := t.nextMsg(); msg != "" {
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
		// XXX needs to be relative
		t.cwin.Move(int(t.lno), int(t.col))
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

func (t *Terminal) update(lno uint64, text string) {
	// XXX add undo stack
	t.cache[lno] = text
	t.upd[lno] = t.cache[lno]
}
