package main

import (
	"bufio"
	//"fmt"
	"io/ioutil"
	"os"
)

func NewTempEditBuffer(prefix string) *EditBuffer {
	// TODO: this.
	return NewEditBuffer(prefix, screen)
}

func NewReadEditBuffer(pathname string) (*EditBuffer, os.Error) {
	st, e := os.Stat(pathname)
	if e != nil {
		return nil, e
	}

	f, e := os.Open(pathname, os.O_RDONLY, 0444)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	b := NewEditBuffer(st.Name, screen)
	r := bufio.NewReader(f)
	for {
		l, e := r.ReadBytes(byte('\n'))
		if e != nil {
			// XXX gross.
			if e != os.EOF {
				return nil, e
			} else {
				b.InsertLine(NewLine(l))
				break
			}
		}
		b.InsertLine(NewLine(l))
	}
	b.fi = st

	// XXX as in d.go, this is a workaround for my lazy design.  fix asap.
	b.anchor = b.lines.Front()

	return b, nil
}

// Do a naive write of the entire buffer to a temp file, then rename into place.
func WriteFile(pathname string, b *EditBuffer) (*os.FileInfo, os.Error) {

	f, e := ioutil.TempFile(TMPDIR, TMPPREFIX)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	i := 0
	wr := 0
	for l := b.lines.Front(); l != nil; l = l.Next() {
		n, e := f.Write(l.Value.(*EditLine).raw())
		if e != nil {
			return nil, e
		}
		i++
		wr += n
	}

	// Ml.mode = fmt.Sprintf("\"%s\", %d bytes", pathname, wr)

	st, e := f.Stat()
	if e != nil {
		return nil, e
	}

	if b.fi != nil {
		pathname = b.fi.Name
	}
	e = os.Rename(st.Name, pathname)
	if e != nil {
		return nil, e
	}

	b.fi = st
	return st, nil
}
