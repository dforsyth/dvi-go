package main

import (
	"bufio"
	"os"
	"utf8"
)

type File struct {
	path  string
	dirty bool
	first *Line
	last  *Line
	disp  *Line
	pos   *Position
	lpos  *Position // last position
	// undo *Undo
	tags map[int]*Position
	clip []byte
}

type Line struct {
	text    []byte // string
	displen uint
	next    *Line
	prev    *Line
	dirty   bool
}

func (l *Line) length() int {
	return utf8.RuneCount(l.text)
}

func NewLine(text []byte) *Line {
	dst := make([]byte, len(text))
	copy(dst, text)
	return &Line{
		text: dst,
		next: nil,
		prev: nil,
	}
}

func (f *File) appendLine(l *Line) {
	if f.first == nil {
		f.first = l
		f.last = l
	} else {
		l.prev = f.last
		f.last.next = l
		f.last = l
	}
	l.dirty = true
}

func addLine(l, nl *Line) {
	nl.prev = l
	nl.next = l.next
	if nl.next != nil {
		nl.next.prev = nl
	}
	l.next = nl
}

func NewFile(path string) *File {
	nf := &File{
		path: path,
	}
	nf.pos = &Position{}
	nf.appendLine(NewLine([]byte{}))
	nf.pos.line = nf.first
	nf.pos.off = 0
	nf.dirty = false
	return nf
}

func readFile(path string) (*File, os.Error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	nf := NewFile(path)
	rdr := bufio.NewReader(f)
	buf := make([]byte, 4096)
	for {
		if n, e := rdr.Read(buf); n > 0 {
			nf.pos = add(*nf.pos, buf[:n])
		} else if n == 0 && e == os.EOF {
			break
		} else {
			return nil, e
		}
	}
	return nf, nil
}

func (f *File) writeFile() os.Error {
	wf, e := os.Create(f.path)
	if e != nil {
		return e
	}
	defer wf.Close()
	for l := f.first; l != nil; l = l.next {
		if l.next != nil {
			if _, e := wf.Write(append(l.text, '\n')); e != nil {
				return e
			}
		} else {
			if _, e := wf.Write(l.text); e != nil {
				return e
			}
		}
	}
	return nil
}

func (f *File) bof() {
	f.pos.line = f.first
	f.pos.off = 0
}

// insert text at f.pos (undo-able)
func (f *File) insert(text []byte) {
	f.pos = add(*f.pos, text)
}

// delete text between a and b (undo-able)
func (f *File) delete(a, b *Position) {
	// XXX either this needs to be called differently of there needs to be a check on whether
	// these position are actually in this file or not.
	fp, sp := orderPos(a, b)
	f.pos = remove(*fp, *sp)
}

func (f *File) undo() {
}

func (f *File) redo() {
}
