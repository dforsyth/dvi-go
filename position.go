package main

type Position struct {
	line *Line
	off  int
}

func orderPos(a, b *Position) (*Position, *Position) {
	if a.line != b.line {
		for l := b.line; l != nil; l = l.next {
			if l == a.line {
				return b, a
			}
		}
		return a, b
	}
	if b.off < a.off {
		return b, a
	}
	return a, b
}

// XXX These should really be renamed nextPos and prevPos
func prevChar(p Position) *Position {
	if p.off > 0 {
		p.off--
	} else if p.line.prev != nil {
		p.line = p.line.prev
		p.off = p.line.length()
	}
	return &p
}

func nextChar(p Position) *Position {
	if p.off < p.line.length() {
		p.off++
	} else if p.line.next != nil {
		p.line = p.line.next
		p.off = 0
	}
	return &p
}

func prevWord(p Position) *Position {
	return &p
}

func nextWord(p Position) *Position {
	return &p
}

func prevLine(p Position) *Position {
	if p.line.prev == nil {
		p.off = 0
	} else {
		p.line = p.line.prev
		// TODO utf8-itize this
		if p.off > p.line.length() {
			p.off = p.line.length()
		}
	}
	return &p
}

func nextLine(p Position) *Position {
	if p.line.next == nil {
		p.off = p.line.length()
	} else {
		p.line = p.line.next
		if p.off > p.line.length() {
			p.off = p.line.length()
		}
	}
	return &p
}

func eol(f *File) {
	f.pos.off = len(f.pos.line.text)
}

func bol(f *File) {
	f.pos.off = 0
}
