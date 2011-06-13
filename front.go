package main

import (
	"curses"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type State struct {
	f       *File
	w       *curses.Window
	msg     *[]byte
	lastkey int
	currx   int
	curry   int
}

func message(s *State) string {
	return fmt.Sprintf("lastkey: %c(%d) | pos: x: %d/%d y: %d", s.lastkey, s.lastkey, s.currx,
		s.f.pos.off, s.curry)
}

func (s *State) queueMsg(msg string, colors int32, beep bool) {
}

var worddelim map[byte]interface{} = map[byte]interface{} {
	' ': nil,
	'\n': nil,
	'\t': nil,
	'\r': nil,
}

func ctrl(k int) int {
	return k & 0x1F
}

func insertmode(s *State) {
	for {
		draw(s)
		k := s.w.Getch()
		switch k {
		case 27:
			if p := prevChar(*s.f.pos); p.line == s.f.pos.line {
				s.f.pos = p
			}
			return
		case ctrl('H'), 127, curses.KEY_BACKSPACE:
			s.f.pos = remove(*prevChar(*s.f.pos), *s.f.pos)
		default:
			s.f.insert([]byte{byte(k)})
		}
	}
}

func main() {
	defer func() {
		endscreen()
	}()

	s := &State{}
	s.msg = nil

	path := os.Args[1]

	f, e := readFile(path)
	if e != nil {
		panic(e.String())
	}
	f.bof()

	s.f = f
	f.disp = f.first

	go func() {
		for {
			s := <-signal.Incoming
			switch s.(signal.UnixSignal) {
			case syscall.SIGINT:
				endscreen()
				panic("sigint")
			case syscall.SIGTERM:
				endscreen()
				panic("sigterm")
			case syscall.SIGWINCH:
				endscreen()
				panic("sigwinch")
			}
		}
	}()

	s.lastkey = 0
	initscreen(s)
	draw(s)
	// command mode

	count := 0

	resetargs := func(c *CmdArgs) {
		c.s = nil
		c.c1 = 0
		c.c2 = 0
	}

	cmdargs := &CmdArgs{}
	for {
		k := s.w.Getch()

		if (k >= '1' && k <= '9') || (count != 0 && k == '0') {
			count *= 10
			count += k - '0'
			continue
		}

		if cmd, ok := vicmds[k]; ok {
			resetargs(cmdargs)
			if count != 0 {
				cmdargs.c1 = count
				count = 0
			} else {
				cmdargs.c1 = 1
			}
			if cmd.motion {
			}

			// execute the command
			cmdargs.s = s
			cmd.fn(cmdargs)
		} else {
			curses.Beep()
		}

		s.lastkey = k
		draw(s)
	}
}

func exmode(s *State) {
	buf := []byte{}
	old := s.msg
	s.msg = &buf
	for {
		draw(s)
		k := s.w.Getch()
		switch k {
		case 27 /* ESC */ :
			s.msg = old
			return
		case 0xd, 0xa, curses.KEY_ENTER:
			if string(buf) == "w" {
				s.f.writeFile()
			} else if string(buf) == "q" {
				endscreen()
				syscall.Exit(0)
			} else if i, e := strconv.Atoi(string(buf)); e == nil {
				l := s.f.first
				for i > 1 && l != nil {
					l = l.next
					i--
				}
				if i == 1 {
					s.f.pos.line = l
					s.f.pos.off = 0
				} else {
					// error
				}
			}
			s.msg = old
			return
		default:
			buf = append(buf, byte(k))
		}
	}
}
