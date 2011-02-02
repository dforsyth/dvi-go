package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func NewTempFileEditBuffer(prefix string) *EditBuffer {
	// TODO: this.
	b := NewEditBuffer(prefix)
	b.AppendLine()
	return b
}

func NewReadFileEditBuffer(pathname string) (*EditBuffer, os.Error) {
	st, e := os.Stat(pathname)
	if e != nil {
		return nil, e
	}

	f, e := os.Open(pathname, os.O_RDONLY, 0444)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	b := NewEditBuffer(st.Name)
	r := bufio.NewReader(f)
	for {
		l, e := r.ReadBytes(byte('\n'))
		if e != nil {
			return nil, e
		}
		b.InsertLine(NewLine(l))
	}
	b.st = st

	return b, nil
}

// Do a naive write of the entire buffer to a temp file, then rename into place.
func WriteEditBuffer(pathname string, b *EditBuffer) (*os.FileInfo, os.Error) {

	f, e := ioutil.TempFile(TMPDIR, TMPPREFIX)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	i := 0
	wr := 0
	for l := b.lines; l != nil; l = l.next {
		n, e := f.Write(l.raw())
		if e != nil {
			return nil, e
		}
		i++
		wr += n
	}

	Ml.mode = fmt.Sprintf("\"%s\", %d bytes", pathname, wr)

	st, e := f.Stat()
	if e != nil {
		return nil, e
	}

	if b.st != nil {
		pathname = b.st.Name
	}
	e = os.Rename(st.Name, pathname)
	if e != nil {
		return nil, e
	}

	b.st = st
	return st, nil
}
