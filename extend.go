/* 
 * Copyright (c) 2011 David Forsythe.
 * See LICENSE file for license details.
 */

package main

import (
	"curses"
	"io/ioutil"
	"os"
)

func dirmode(d *Dvi) *Buffer {
	for {
		draw(d)
		switch k := getCh(d); k {
		case 'j':
			d.b.pos = nextLine(*d.b.pos)
		case 'k':
			d.b.pos = prevLine(*d.b.pos)
		case 0xd, 0xa, curses.KEY_ENTER:
			if b, e := openFile(string(d.b.pos.line.text)); e == nil {
				b.resetPos()
				b.disp = b.first
				d.addBuf(b)
				return b
			}
			return nil
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
	} else {
		d.b = o
	}
}

func emacs(d *Dvi) {
	for {
		draw(d)
		switch k := getCh(d); k {
		case 27:
			return
		case ctrl('N'):
			d.b.pos = nextLine(*d.b.pos)
		case ctrl('P'):
			d.b.pos = prevLine(*d.b.pos)
		case ctrl('B'):
			d.b.pos = prevChar2(*d.b.pos)
		case ctrl('F'):
			d.b.pos = nextChar2(*d.b.pos)
		case ctrl('H'), 127, curses.KEY_BACKSPACE:
			pp := prevChar2(*d.b.pos)
			d.b.remove(*prevChar2(*d.b.pos), *d.b.pos, false)
			d.b.pos = pp
		default:
			d.b.pos = d.b.add(*d.b.pos, []byte{byte(k)})
		}
	}
}

func nextBuffer(a *CmdArgs) (*Position, os.Error) {
	if a.d.b.next != nil {
		a.d.b = a.d.b.next
		return a.d.b.pos, nil
	}

	if a.d.b == a.d.bufs {
		return nil, &DviError{"single buffer", 0}
	}

	a.d.b = a.d.bufs
	return a.d.b.pos, nil
}
