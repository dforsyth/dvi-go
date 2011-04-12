package main

import (
	"fmt"
	"os"
	"path"
)

func NewTempEditBuffer(gs *GlobalState, prefix string) *EditBuffer {
	// TODO: this.
	e := NewEditBuffer(gs, prefix)
	e.temp = true
	return e
}

func OpenBuffer(gs *GlobalState, pathname string) (Buffer, os.Error) {
	wd, e := os.Getwd()
	if e != nil {
		return nil, e
	}
	if !path.IsAbs(pathname) {
		pathname = path.Join(wd, pathname)
	}

	f, e := os.Open(pathname)
	if st, e := f.Stat(); e == nil {
		if st.IsRegular() {
			eb := NewEditBuffer(gs, pathname)
			if _, e = eb.readFile(f, 0); e == nil {
				return eb, nil
			}
		} else if st.IsDirectory() {
			eb := NewDirBuffer(gs, pathname)
			return eb, nil
		} else {
			e = &DviError{fmt.Sprintf("%s: can't deal with this filetype", pathname)}
		}
	}
	return nil, e
}
