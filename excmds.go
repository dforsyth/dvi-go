package main

// XXX Everything in here should be considered a place holder

import (
	"curses"
	"os"
	"strconv"
	"syscall"
)

type ExStatus struct {
	inp string
}

func (x *ExStatus) Display() string {
	return ":" + x.inp
}

func (x *ExStatus) Color() int {
	return 0
}

func (x *ExStatus) Beep() bool {
	return false
}

type excmd struct {
	fn func(*ExArgs) (interface{}, os.Error)
}

type ExArgs struct {
	d     *Dvi
	args  []string
	force bool
}

func parseEx(x string) (*excmd, *ExArgs, os.Error) {
	// TODO: this
	a := &ExArgs{}
	if x == "w" {
		return excmds["w"], a, nil
	} else if x == "q" || x == "q!" {
		if x == "q!" {
			a.force = true
		}
		return excmds["q"], a, nil
		endscreen()
		syscall.Exit(0)
	} else if i, e := strconv.Atoi(x); e == nil && i >= 0 {
		a.args = make([]string, 1)
		a.args[0] = x
		return excmds["0"], a, nil
	} else if x == "db" {
		return excmds["db"], a, nil
	}
	/*
		else if msg.inp == "emacs" {
			emacs(d)
		} else if msg.inp == "file" {
			d.queueMsg(d.b.information(), 1, false)
			return
		} else if msg.inp == "showmsg" {
			d.showmsg = !d.showmsg
		} else {
			curses.Beep()
		}
	*/
	return nil, nil, &DviError{"no command: " + x, 0}
}

func exWriteFile(a *ExArgs) (interface{}, os.Error) {
	// only write if dirty
	if a.d.b.temp {
		return nil, &DviError{"Cannot write temporary file", 0}
	}
	if a.d.b.dirty {
		a.d.b.writeFile()
		a.d.b.dirty = false
	}
	return nil, nil
}

func exQuit(a *ExArgs) (interface{}, os.Error) {
	if !a.force {
		for b := a.d.bufs; b != nil; b = b.next {
			if b.dirty {
				return nil, &DviError{b.name + " has modifications", 0}
			}
		}
	}
	endscreen()
	syscall.Exit(0)
	return nil, nil
}

func exGoToLine(a *ExArgs) (interface{}, os.Error) {
	if i, e := strconv.Atoi(a.args[0]); e == nil {
		return i, nil
	}
	return nil, &DviError{"Cannot get to line " + a.args[0], 0}
}

func exDirBrowser(a *ExArgs) (interface{}, os.Error) {
	// XXX poc for extending ex
	directoryBrowser(a.d, ".")
	return nil, nil
}

func exmode(d *Dvi) {
	ex := &ExStatus{}
	d.setStatus(ex)
	defer d.unsetStatus()
	for {
		draw(d)
		k := getCh(d)
		switch k {
		case 27 /* ESC */ :
			return
		case 0xd, 0xa, curses.KEY_ENTER:
			if cmd, args, e := parseEx(ex.inp); e == nil {
				args.d = d
				if _, e := cmd.fn(args); e != nil {
					d.queueMsg(e.String(), 1, true)
				}
			}
			return
		default:
			ex.inp += string([]byte{byte(k)})
		}
	}
}
