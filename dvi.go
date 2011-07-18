/* 
 * Copyright (c) 2011 David Forsythe.
 * See LICENSE file for license details.
 */

package main

import (
	"curses"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
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
	bufs    *Buffer
	w       *curses.Window
	lastcmd int
	showmsg bool
	status  DviStatus
	statusq chan DviStatus
	lastkey int
	currx   int
	curry   int
	buffers map[byte]*Buffer
	cmdDisp string
}

type DviMessage struct {
	message string
	color   int
	beep    bool
}

func (m *DviMessage) Display() string {
	return m.message
}

func (m *DviMessage) Color() int {
	return m.color
}

func (m *DviMessage) Beep() bool {
	return m.beep
}

type DviStatus interface {
	Display() string
	Color() int
	Beep() bool
}

func (d *Dvi) statusDisplay() (string, int, bool) {
	// not thrilled with the name of this fn
	if d.status == nil {
		return fmt.Sprintf("lastkey: %c(%d) | pos: x: %d/%d y: %d", d.lastkey, d.lastkey,
			d.currx, d.b.pos.off, d.curry), 3, false
	}
	return d.status.Display(), d.status.Color(), d.status.Beep()
}

func (d *Dvi) queueMsg(msg string, colors int, beep bool) {
	d.statusq <- &DviMessage{
		msg,
		colors,
		beep,
	}
}

func (d *Dvi) setStatus(s DviStatus) {
	d.status = s
}

func (d *Dvi) unsetStatus() {
	d.status = nil
}

var breakchars map[int]interface{} = map[int]interface{}{
	' ':  nil,
	'\n': nil,
	'\t': nil,
	'\r': nil,
	'(':  nil,
	')':  nil,
	'[':  nil,
	']':  nil,
	'/':  nil,
	'\\': nil,
	'.':  nil,
}

var blankchars map[int]interface{} = map[int]interface{}{
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
	startpos := &Position{d.b.pos.line, d.b.pos.off}
	for {
		draw(d)
		k := getCh(d)
		switch k {
		case 27:
			if p := prevChar(*d.b.pos); p.line == d.b.pos.line {
				d.b.pos = p
			}
			return
		case curses.KEY_LEFT, curses.KEY_RIGHT, curses.KEY_UP, curses.KEY_DOWN:
			curses.Beep()
		case ctrl('H'), 127, curses.KEY_BACKSPACE:
			// TODO: Don't let backspace travel past starting point of input session
			pp := prevChar2(*d.b.pos)
			if posEq(d.b.pos, startpos) {
				d.queueMsg("", 1, true)
				continue
			}
			d.b.remove(*prevChar2(*d.b.pos), *d.b.pos, false)
			d.b.pos = pp
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
func getCh(d *Dvi) int {
	// expose simple input
	return d.w.Getch()
}

func commandMode(d *Dvi) {
	ca := &CmdArgs{}
	count := 0
	for {
		k := getCh(d)
		if (k >= '1' && k <= '9') || (count != 0 && k == '0') {
			count *= 10
			count += k - '0'
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

			ca.start = &Position{d.b.pos.line, d.b.pos.off}
			if cmd.motion {
				// This is a motion command
				// XXX so much dup here.  really should pull this out.
				// XXX d0 doesnt work.
				ma := &CmdArgs{}
				mcount := 0
				for {
					mk := getCh(d)
					if (mk >= '1' && mk <= '9') || (mcount != 0 && mk == '0') {
						mcount *= 10
						mcount += mk - '0'
						continue
					}

					resetCmdArgs(ma)
					// if the initial command and the motion command are both given counts, then the two
					// counts are multiplied to form the final count
					// c1 should never be 0, i think...
					if mcount > 0 {
						ma.c1 = mcount
					} else {
						ma.c1 = 1
					}
					if count != 0 {
						ma.c1 *= count
					}

					if mk == k {
						ca.start.off = 0
						p := &Position{ca.start.line, 0}
						for i := 0; i < ma.c1-1; i++ {
							p = nextLine(*p)
						}
						ca.end = p
						// ca.end = &Position{ca.start.line, ca.start.line.length()}
						ca.line = true
					} else if mcmd, ok := vicmds[mk]; ok && mcmd.isMotion {
						ma.motion = true
						//if mcount == 0 && !mcmd.zerocount {
						//	ma.c1 = 1
						//}

						ma.start = ca.start
						ma.d = d
						if p, e := mcmd.fn(ma); e == nil {
							ca.end = p
							ca.line = mcmd.line
							if ca.line {
								ca.start.off = 0
								ca.end.off = ca.end.line.length()
							}
						} else {
							// error reporting should be set up in the cmd fn
							goto end
						}
					} else {
						d.queueMsg(fmt.Sprintf("%c is not a valid motion", mk), 2, true)
						goto end
					}
					break
				}
			}

			ca.d = d
			if ca.start != nil && ca.end != nil {
				ca.start, ca.end = orderPos(ca.start, ca.end)
			}
			if p, e := cmd.fn(ca); e == nil {
				d.b.pos = p
			} else {
				d.queueMsg(e.String(), 2, true)
			}
		} else {
			d.queueMsg(fmt.Sprintf("%c is not a dvi command", k), 2, true)
		}

	end:
		count = 0
		d.lastkey = k
		draw(d)
	}

}

func openFile(path string) (*Buffer, os.Error) {
	file, e := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	stat, e := file.Stat()
	if e != nil {
		return nil, e
	}
	if !stat.IsRegular() {
		return nil, &DviError{"Not a regular file", 0}
	}
	// TODO create lock file
	buf := newBuffer()
	if e := buf.loadFile(file); e != nil {
		return nil, e
	}
	buf.name = path
	buf.temp = false
	return buf, nil
}

func openTempFile() (*Buffer, os.Error) {
	tfile, e := ioutil.TempFile(config["tempdir"].(string), config["temppfx"].(string))
	if e != nil {
		return nil, e
	}
	defer tfile.Close()

	buf := newBuffer()
	buf.name = tfile.Name()
	buf.temp = true
	return buf, nil
}

func (d *Dvi) addBuf(buf *Buffer) {
	if d.bufs == nil {
		d.bufs = buf
	} else {
		var b *Buffer
		for b = d.bufs; b.next != nil; b = b.next {
		}
		b.next = buf
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

	// status (nil is default)
	d.status = nil
	d.statusq = make(chan DviStatus, 4)

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		for _, path := range args {
			if b, e := openFile(path); e == nil {
				d.addBuf(b)
				b.resetPos()
				b.disp = b.first
			}
		}
	}
	if d.bufs == nil {
		if b, e := openTempFile(); e == nil {
			d.addBuf(b)
			b.resetPos()
			b.disp = b.first
		}
	}

	// use the first buffer
	d.b = d.bufs

	go sighandlers()

	d.lastkey = 0
	initscreen(d)
	draw(d)

	commandMode(d)
}
