package main

import (
	"io/ioutil"
	"os"
)

func dirmode(d *Dvi) *Buffer {
	for {
		draw(d)
		switch k := getC(d); k {
		case 'j':
			d.b.pos = nextLine(*d.b.pos)
		case 'k':
			d.b.pos = prevLine(*d.b.pos)
		case 13:
			b := newBuffer()
			// XXX gross
			if f, e := os.Open(string(d.b.pos.line.text)); e == nil {
				defer f.Close()
				b.loadFile(f)
			}
			return b
		case 27:
			return nil
		default:
		}
	}
	return nil
}

func directoryBrowser(d *Dvi, path string) {
	// remember the old buffer
	o := d.b
	// set up a new buffer
	b := newBuffer()
	if ls, e := ioutil.ReadDir(path); e == nil {
		for i, d := range ls {
			b.add(*b.pos, []byte(d.Name))
			if i < len(ls)-1 {
				b.add(*b.pos, []byte("\n"))
			}
		}
	}
	b.disp = b.first
	// set the current buffer to the directory listing
	d.b = b
	// enter "dirmode"
	if n := dirmode(d); n != nil {
		d.b = n
		d.b.disp = d.b.first
		d.b.pos.line = d.b.first
		d.b.pos.off = 0
	} else {
		d.b = o
	}
}

func emacs(d *Dvi) {
	for {
		draw(d)
		switch k := getC(d); k {
		case ctrl('N'):
			d.b.pos = nextLine(*d.b.pos)
		case ctrl('P'):
			d.b.pos = prevLine(*d.b.pos)
		case ctrl('B'):
			d.b.pos = prevChar2(*d.b.pos)
		case ctrl('F'):
			d.b.pos = nextChar2(*d.b.pos)
		default:
			d.b.pos = d.b.add(*d.b.pos, []byte{byte(k)})
		}
	}
}

func gdb(d *Dvi) {
	// not implemented
	for {
		draw(d)
		switch k := getC(d); k {
		case ctrl('B'):
			// un/set breakpoint
		default:
		}
	}
}
