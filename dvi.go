package main

import (
	"curses"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode"
)

const (
	ERRRANGE = iota
)

type DviError struct {
	message string
	eno     int
}

func (e *DviError) String() string {
	return e.message
}

type Dvi struct {
	b       *Buffer
	w       *curses.Window
	lastcmd int
	msg     *DviMessage
	lastkey int
	currx   int
	curry   int
	buffers map[byte]*Buffer
}

type DviMessage struct {
	message string
	color   string
	beep    bool
}

func message(s *Dvi) string {
	return fmt.Sprintf("lastkey: %c(%d) | pos: x: %d/%d y: %d", s.lastkey, s.lastkey, s.currx,
		s.b.pos.off, s.curry)
}

func (s *Dvi) queueMsg(msg string, colors int32, beep bool) {
}

var breakchars map[int]interface{} = map[int]interface{}{
	' ':  nil,
	'\n': nil,
	'\t': nil,
	'\r': nil,
}

func validBufName(n byte) bool {
	// 0 is unnamed
	if n >= 'a' && n <= 'z' {
		return true
	}
	return false
}

func ctrl(k int) int {
	// just watching out for you :)
	return unicode.ToUpper(k) & 0x1F
}

func insertmode(d *Dvi) {
	for {
		draw(d)
		k := d.w.Getch()
		switch k {
		case 27:
			if p := prevChar(*d.b.pos); p.line == d.b.pos.line {
				d.b.pos = p
			}
			return
		case curses.KEY_LEFT, curses.KEY_RIGHT, curses.KEY_UP, curses.KEY_DOWN:
			curses.Beep()
		case ctrl('H'), 127, curses.KEY_BACKSPACE:
			d.b.pos = d.b.remove(prevChar(*d.b.pos), d.b.pos, false)
		default:
			d.b.pos = d.b.add(*d.b.pos, []byte{byte(k)})
			d.b.dirty = true
		}
	}
}

func motionInput() *CmdArgs {
	return nil
}

func getBuffer() byte {
	return 0
}

func sighandlers() {
	for {
		s := <-signal.Incoming
		switch s.(os.UnixSignal) {
		case syscall.SIGINT:
			endscreen()
			panic("sigint")
		case syscall.SIGTERM:
			endscreen()
			panic("sigterm")
			//case syscall.SIGWINCH:
			//	endscreen()
			//	panic("sigwinch")
		}
	}
}

func main() {
	defer func() {
		endscreen()
	}()

	d := &Dvi{}

	// buffers
	d.buffers = make(map[byte]*Buffer)
	d.buffers[0] = newBuffer()

	// msg (nil is default)
	d.msg = nil

	flag.Parse()
	args := flag.Args()

	b := newBuffer()

	if len(args) > 0 {
		path := args[0]
		f, e := os.Open(path)
		if e != nil {
			panic(e.String())
		}
		b.name = path
		if e := b.loadFile(f); e != nil {
			panic(e.String())
		}
		f.Close()
	}

	d.b = b

	// reset position and jump display the first line
	b.resetPos()
	b.disp = b.first

	go sighandlers()

	d.lastkey = 0
	initscreen(d)
	draw(d)

	// command mode
	count := 0

	resetargs := func(c *CmdArgs) {
		c.d = nil
		c.c1 = 0
		c.c2 = 0
		c.buffer = 0
		c.line = false
	}

	cmdargs := &CmdArgs{}
	for {
		k := d.w.Getch()

		if (k >= '1' && k <= '9') || (count != 0 && k == '0') {
			count *= 10
			count += k - '0'
			continue
		}

		if cmd, ok := vicmds[k]; ok {
			resetargs(cmdargs)

			if count != 0 {
				cmdargs.c1 = count
			} else {
				if !cmd.zerocount {
					cmdargs.c1 = 1
				} else {
					cmdargs.c1 = 0
				}
			}

			cmdargs.start = d.b.pos
			if cmd.motion {
				// TODO: move all of this into motionInput
				mcount := 0
				mcmdargs := &CmdArgs{}
				resetargs(mcmdargs)
				for {
					mk := d.w.Getch()

					if (mk >= '1' && mk <= '9') || (mcount != 0 && mk == '0') {
						mcount *= 10
						mcount += mk - '0'
						continue
					}

					if mk == k {
						l := d.b.pos.line
						cmdargs.start = &Position{l, 0}
						cmdargs.end = &Position{l, l.length()}
						cmdargs.line = true
						curses.Beep()
						break
					}

					if mcmd, ok := vicmds[mk]; ok {
						if !mcmd.isMotion {
							curses.Beep()
							break
						}

						if mcount != 0 {
							mcmdargs.c1 = mcount
						} else {
							if !mcmd.zerocount {
								mcmdargs.c1 = 1
							} else {
								mcmdargs.c1 = 0
							}
						}

						if count != 0 {
							mcmdargs.c1 *= count
						}

						// XXX Need to sort start/end position
						mcmdargs.start = d.b.pos
						mcmdargs.motion = true
						if end, e := mcmd.fn(mcmdargs); e == nil {
							cmdargs.end = end
							if mcmd.line {
								cmdargs.start.off = 0
								if end.line.length() > 0 {
									cmdargs.end.off = cmdargs.end.line.length()
								} else {
									cmdargs.end.off = 0
								}
							}
							cmdargs.line = mcmd.line
							break
						}
						panic("got a bad motion command")
					} else {
						curses.Beep()
					}
					mcount = 0
				}
			}

			// execute the command
			cmdargs.d = d
			if cmdargs.start != nil && cmdargs.end != nil {
				cmdargs.start, cmdargs.end = orderPos(cmdargs.start, cmdargs.end)
			}
			if p, e := cmd.fn(cmdargs); e == nil {
				d.b.pos = p
			}
			count = 0
			// XXX The cursor actually needs to be corrected since we're in command
			// mode.  If p.off == p.line.length() { p = prevChar(*p) }
		} else {
			curses.Beep()
		}

		d.lastkey = k
		draw(d)
	}
}

func exmode(d *Dvi) {
	old := d.msg
	msg := &DviMessage{}
	for {
		d.msg = msg
		draw(d)
		k := d.w.Getch()
		switch k {
		case 27 /* ESC */ :
			d.msg = nil
			return
		case 0xd, 0xa, curses.KEY_ENTER:
			if msg.message == "w" {
				d.b.writeFile()
				d.b.dirty = false
			} else if msg.message == "q" || msg.message == "q!" {
				if d.b.dirty && msg.message != "q!" {
					curses.Beep()
					return
				}
				endscreen()
				syscall.Exit(0)
			} else if i, e := strconv.Atoi(msg.message); e == nil {
				l := d.b.first
				for i > 1 && l != nil {
					l = l.next
					i--
				}
				if i == 1 {
					d.b.pos.line = l
					d.b.pos.off = 0
				} else {
					// error
				}
			} else if msg.message == "db" {
				directoryBrowser(d, ".")
			} else if msg.message == "emacs" {
				emacs(d)
			} else {
				curses.Beep()
			}
			d.msg = old
			return
		default:
			msg.message += string([]byte{byte(k)})
		}
	}
}
