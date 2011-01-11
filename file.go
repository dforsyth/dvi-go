package main

import (
	"bufio"
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

func WriteEditBuffer(pathname string, b *EditBuffer) *os.FileInfo {
	return nil
}
