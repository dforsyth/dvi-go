package main

import (
	"bufio"
	"io/ioutil"
	"os"
)

func NewTempFileEditBuffer(prefix string) *EditBuffer {
	// TODO: this.
	return NewEditBuffer(prefix)
}

func NewReadFileEditBuffer(pathname string) *EditBuffer {
	st, e := os.Stat(pathname)
	if e != nil {
		Debug = e.String()
		return nil
	}

	f, e := os.Open(pathname, os.O_RDONLY, 0444)
	if e != nil {
		Debug = e.String()
		return nil
	}
	defer f.Close()

	b := NewEditBuffer(st.Name)
	r := bufio.NewReader(f)
	for {
		l, e := r.ReadBytes(byte('\n'))
		if e != nil {
			Debug = e.String()
			break
		}
		b.InsertLine(NewGapBuffer(l))
	}

	return b
}

// Do a naive write of the entire buffer to a temp file, then rename into place.
func WriteEditBuffer(pathname string, b *EditBuffer) (*os.FileInfo, os.Error) {

	f, e := ioutil.TempFile(TMPDIR, TMPPREFIX)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	for l := b.Lines().Front(); l != nil; l = l.Next() {
		_, e := f.Write(l.Value.(*GapBuffer).GaplessBuffer())
		if e != nil {
			return nil, e
		}
	}

	st, e := f.Stat()
	if e != nil {
		return nil, e
	}

	// rename (i dont think this is dir smart.  i dont remember stat())
	e = os.Rename(st.Name, b.st.Name)
	if e != nil {
		return nil, e
	}

	b.st = st

	return st, nil
}
