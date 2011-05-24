package main

import (
	"flag"
	//	"fmt"
	//	"time"
)

func main() {
	noterm := flag.Bool("noterm", false, "start a terminal")
	flag.Parse()

	h := NewHost()
	h.serve()
	c := NewClient()

	if !*noterm {
		t := NewTerminal(c)
		t.init()
		t.run()
	} else {
		c.open("test")
		c.open("test2")
	}
}
