package main

import (
	"fmt"
	"strconv"
	"unicode"
)

var aliases map[string]string = map[string]string{
	"write":   "w",
	"quit":    "q",
	"version": "ve",
}

var exFns map[string]func(*GlobalState) = map[string]func(*GlobalState){
	"w":  write,
	"wq": writeQuit,
	"q":  quit,
	"nb": nextBuffer,
	"pb": prevBuffer,
	"ve": version,
}

func writeEditBuffer(b *EditBuffer, path string, force bool) (int, *Message) {
	if b.isTemp() && path == b.ident() && !force {
		return -1, &Message{
			fmt.Sprintf("%s is a temporary file", b.ident()),
			true,
		}
	}

	if !b.isDirty() {
		return -1, nil
	}

	b.dirty = false
	return -1, &Message{
		"writeEditBuffer not implemented",
		true,
	}
}

func write(gs *GlobalState) {
	if t, ok := gs.curBuf().(*EditBuffer); ok {
		if _, msg := writeEditBuffer(t, t.ident(), gs.x.frc); msg != nil {
			gs.queueMessage(msg)
		}
		return
	}
}

func quit(gs *GlobalState) {
	if t, ok := gs.curBuf().(*EditBuffer); ok {
		if t.isDirty() && !gs.x.frc {
			gs.queueMessage(&Message{
				fmt.Sprintf("%s has unsaved changes", t.ident()),
				false,
			})
			return
		}
	}
	Done(0)
}

func writeQuit(gs *GlobalState) {
	if t, ok := gs.curBuf().(*EditBuffer); ok {
		if lw, msg := writeEditBuffer(t, t.ident(), gs.x.frc); lw >= 0 {
			Done(0)
		} else {
			gs.queueMessage(msg)
		}
	}
}

func version(gs *GlobalState) {
	gs.queueMessage(&Message{
		fmt.Sprintf("Version %s (%s) %s", gs.version, gs.buildDate, gs.author),
		false,
	})
}

func ex(gs *GlobalState) {

	gs.SetModeline(gs.ex)
	gs.ex.Reset()

	gs.Mode = MODEEX

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		switch k {
		case ESC:
			return
		case 0xd, 0xa:
			x := new(Ex)
			gs.x = x
			x.cmd = gs.ex.buffer
			f := x.parse()
			if f == nil {
				gs.queueMessage(&Message{
					fmt.Sprintf("cmd: %s, st: %d, end: %d, cnt: %d", x.cmd,
						x.st, x.end, x.cnt),
					true,
				})
			} else {
				f(gs)
			}
			// gs.ex.execute()
			return
		default:
			gs.ex.buffer += string(k)
		}
	}
}

type exBuffer struct {
	buffer string
	gs     *GlobalState
}

func newExBuffer(gs *GlobalState) *exBuffer {
	c := new(exBuffer)
	c.buffer = ""
	c.gs = gs
	return c
}

func (c *exBuffer) String() string {
	return fmt.Sprintf(":%s", c.buffer)
}

func (c *exBuffer) GetCursor() int {
	return len(c.String()) - 1
}

func (c *exBuffer) msgOverride(m *Message) {
}

type Ex struct {
	cmd  string
	st   int
	end  int
	cnt  int
	frc  bool
	args []string
	gs   *GlobalState
}

func (x *Ex) clear() {
	x.st = 0
	x.end = 0
	x.cnt = 0
	x.frc = false
	x.args = make([]string, 1)
}

// parse a single ex cmd
func (x *Ex) parse() func(*GlobalState) {
	x.clear()

	cmd := ""
	// get rid of extra colons and spaces
	for i, c := range x.cmd {
		if c != ':' || c != ' ' {
			cmd = x.cmd[i:]
			break
		}
	}
	if len(cmd) == 0 {
		Die("0 len: " + cmd + " vs " + x.cmd)
		return nil
	}

	// if the line is a comment, leave
	if cmd[0] == '"' {
		return nil
	}

	r := false
	comma := false
	a := false
	p := ""
	for i, c := range x.cmd {
		if c == ' ' {
			continue
		}
		if c == ';' {
			// cmd split, not supported yet.
			goto lookup
		}

		if unicode.IsDigit(c) {
			if a {
				return nil
			}
			r = true
			p += string(c)
		} else {
			if r {
				if c == ',' {
					if comma {
						return nil
					}
					if st, err := strconv.Atoi(p); err == nil {
						x.st = st
						p += string(c)
						p = ""
						comma = true
					} else {
						return nil
					}
				} else {
					if comma {
						if end, err := strconv.Atoi(p); err == nil {
							x.end = end
							p = ""
						}
					} else {
						if cnt, err := strconv.Atoi(p); err == nil {
							x.cnt = cnt
							p = ""
						}
					}
					// start building the command
					p += string(c)
				}
			} else {
				a = true
				if c == '!' && i == len(x.cmd)-1 {
					x.frc = true
				} else {
					p += string(c)
				}
			}
		}
	}
lookup:
	if alias, ok := aliases[p]; ok {
		p = alias
	}

	if fn, ok := exFns[p]; ok {
		return fn
	}

	return nil
}

func (c *exBuffer) Reset() {
	c.buffer = ""
}
