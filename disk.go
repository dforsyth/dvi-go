package main

import (
	"bufio"
	//"fmt"
	"io/ioutil"
	"os"
)

func NewTempEditBuffer(gs *GlobalState, prefix string) *EditBuffer {
	// TODO: this.
	return NewEditBuffer(gs, prefix)
}

func NewReadEditBuffer(gs *GlobalState, pathname string) (*EditBuffer, os.Error) {
	st, e := os.Stat(pathname)
	if e != nil {
		return nil, e
	}

	f, e := os.Open(pathname, os.O_RDONLY, 0444)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	b := NewEditBuffer(gs, st.Name)
	r := bufio.NewReader(f)
	for {
		l, e := r.ReadBytes(byte('\n'))
		if e != nil {
			// XXX gross.
			if e != os.EOF {
				return nil, e
			} else {
				b.InsertLine(NewEditLine(l))
				break
			}
		}
		b.InsertLine(NewEditLine(l))
	}
	b.fi = st

	return b, nil
}

// Do a naive write of the entire buffer to a temp file, then rename into place.
// XXX not adding newlines?
func WriteFile(pathname string, b *EditBuffer) (*os.FileInfo, os.Error) {

	f, e := ioutil.TempFile(TMPDIR, TMPPREFIX)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	i := 0
	wr := 0
	for _, l := range b.Lines {
		n, e := f.Write(l.GetRaw())
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
