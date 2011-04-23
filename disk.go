package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func NewTempEditBuffer(gs *GlobalState, prefix string) (*EditBuffer, os.Error) {
	// TODO: this.
	f, e := ioutil.TempFile(os.TempDir(), prefix)
	if e != nil {
		return nil, e
	}
	defer f.Close() // close this now

	eb := NewEditBuffer(gs, f.Name())
	eb.temp = true
	return eb, nil
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
	if e != nil {
		return nil, e
	}
	defer f.Close()

	if st, e := f.Stat(); e == nil {
		if st.IsRegular() {
			eb := NewEditBuffer(gs, pathname)
			if _, e = eb.readFile(f, 0); e == nil {
				return eb, nil
			}
		} else {
			e = &DviError{fmt.Sprintf("%s: can't deal with this filetype", pathname)}
		}
	}
	return nil, e
}
