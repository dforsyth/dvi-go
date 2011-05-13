package main

import (
	"bufio"
	"os"
)

type File struct {
	name  string // the identifying name on both the client and the host
	fid   uint64
	buf   [][]byte
	dirty bool
}

func NewFile(name string, f *os.File) (*File, os.Error) {
	nf := new(File)
	nf.name = name
	nf.fid = uint64(f.Fd())
	if e := nf.read(f); e != nil {
		return nil, e
	}
	return nf, nil
}

func (f *File) close() {
}

func (f *File) fileInfo() (*os.FileInfo, os.Error) {
	return os.Stat(f.name)
}

func (f *File) line(ln int) ([]byte, os.Error) {
	if ln > len(f.buf)-1 {
		return nil, &DviError{"Line number out of range"}
	}

	return f.buf[ln], nil
}

func (f *File) insert(ln, p int, b []byte) os.Error {
	return nil
}

func (f *File) update(ln uint64, b []byte) os.Error {
	return nil
}

func (f *File) newline(ln int) os.Error {
	return nil
}

func (f *File) delete(p, n int) os.Error {
	return nil
}

func (f *File) read(fi *os.File) os.Error {
	f.buf = make([][]byte, 0)
	r := bufio.NewReader(fi)
	for {
		if ln, e := r.ReadBytes('\n'); e == nil {
			f.buf = append(f.buf, ln)
		} else {
			if e != os.EOF {
				return e
			} else {
				if len(ln) > 0 {
					f.buf = append(f.buf, ln)
				}
				return nil
			}
		}
	}
	// NOT REACHED
	return nil
}

func (f *File) sync(fi *os.File) (uint64, os.Error) {
	tw := uint64(0)
	for _, ln := range f.buf {
		if len(ln) > 0 && ln[len(ln)-1] != '\n' {
			ln = append(ln, '\n')
		}
		if w, e := fi.Write(ln); e != nil {
			return tw, e
		} else {
			tw += uint64(w)
		}
	}

	return tw, nil
}
