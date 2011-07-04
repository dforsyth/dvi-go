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
	cmdDisp string
}

type DviMessage struct {
	message string
	color   int32
	beep    bool
}

func message(s *Dvi) string {
	return fmt.Sprintf("lastkey: %c(%d) | pos: x: %d/%d y: %d", s.lastkey, s.lastkey, s.currx,
		s.b.pos.off, s.curry)
}

func (d *Dvi) queueMsg(msg string, colors int32, beep bool) {
	d.msg = &DviMessage{
		msg,
		colors,
		beep,
	}
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
		k := getC(d)
		switch k {
		case 27:
			if p := prevChar(*d.b.pos); p.line == d.b.pos.line {
				d.b.pos = p
			}
			return
		case curses.KEY_LEFT, curses.KEY_RIGHT, curses.KEY_UP, curses.KEY_DOWN:
			curses.Beep()
		case ctrl('H'), 127, curses.KEY_BACKSPACE:
			d.b.pos = d.b.remove(*prevChar(*d.b.pos), *d.b.pos, false)
		default:
			d.b.pos = d.b.add(*d.b.pos, []byte{byte(k)})
			d.b.dirty = true
		}
	}
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

// Read input off of the dvi input window
func getC(d *Dvi) int {
	// expose simple input
	return d.w.Getch()
}

func commandmode(d *Dvi) {
	ca := &CmdArgs{}
	count := 0
	for {
		k := getC(d)
		if (k >= '1' && k <= '9') || (count != 0 && k == '0') {
			count *= 10
			count += k-'0'
			continue
		}

		if cmd, ok := vicmds[k]; ok {
			// got a command
			resetCmdArgs(ca)
			if count != 0 {
				ca.c1 = count
			} else if !cmd.zerocount {
					ca.c1 = 1
			}

			ca.start = &(*d.b.pos)
			if cmd.motion {
				// This is a motion command
				// XXX so much dup here.  really should pull this out.
				ma := &CmdArgs{}
				mcount := 0
				mk := getC(d)
				if (mk >= '1' && mk <= '9') || (mcount != 0 && mk == '0') {
					mcount *= 10
					mcount += mk-'0'
					continue
				}
				resetCmdArgs(ma)
				if mk == k {
					ca.start.off = 0
					ca.end = &(*ca.start)
					ca.end.off = ca.end.line.length()
					ca.line = true
				} else if mcmd, ok := vicmds[k]; ok && mcmd.isMotion {
					ma.motion = true
					if mcount != 0 {
						ma.c1 = mcount
					} else if !mcmd.zerocount {
						ma.c1 = 1
					}

					// if the initial command and the motion command are both given counts, then the two
					// counts are multiplied to form the final count
					if count != 0 {
						ma.c1 *= count
					}

					ma.start = ca.start
					ma.d = d
					if p, e := mcmd.fn(ma); e != nil {
						ca.end = p
						ca.line = mcmd.line
						if ca.line {
							ca.start.off = 0
							ca.end.off = ca.end.line.length()
						}
					} else {
						// error reporting should be set up in the cmd fn
					}
				} else {
					d.queueMsg(fmt.Sprintf("%c is not a valid motion", mk), 1, true)
				}
			}

			ca.d = d
			if ca.start != nil && ca.end != nil {
				ca.start, ca.end = orderPos(ca.start, ca.end)
			}
			if p, e := cmd.fn(ca); e == nil {
				d.b.pos = p
			}
		} else {
			d.queueMsg(fmt.Sprintf("%c is not a dvi command", k), 1, true)
		}

		count = 0
		d.lastkey = k
		draw(d)
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

	commandmode(d)
}

func exmode(d *Dvi) {
	old := d.msg
	msg := &DviMessage{}
	for {
		d.msg = msg
		draw(d)
		k := getC(d)
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
					d.queueMsg("modifications made to buffer!", 1, true)
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
