package main

import (
//	"fmt"
//	"time"
)

func main() {
	h := NewHost()
	c := NewClient(h)
	go h.serve()

	/*
		o, e := c.open("Makefile")
		if e != nil {
			panic(e.String())
		}
		fid := o.fid

		s, e := c.stat(fid)
		if e != nil {
			panic(e.String())
		}

		lcnt := s.lines
		if lcnt == 0 {
			panic("this file has no lines")
		}
		println(s.message())

		var at uint64 = 0
		for {
			l, e := c.line(fid, at)
			if e != nil {
				panic(e.String())
			}

			fmt.Printf("%d: %s", l.lno, l.text)

			if at++; at > lcnt-1 {
				at = 0
			}
			time.Sleep(1000000000)
		}
	*/
	t := NewTerminal(c)
	t.init()
	t.run()
}
