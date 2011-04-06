package main

import (
	//"fmt"
	"io/ioutil"
	"os"
)

func NewTempEditBuffer(gs *GlobalState, prefix string) *EditBuffer {
	// TODO: this.
	return NewEditBuffer(gs, prefix)
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
	for _, l := range b.lines {
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
