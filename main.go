package main

import (
	"flag"
	//	"fmt"
	"time"
)

func main() {
	noterm := flag.Bool("noterm", false, "start a terminal")
	hostonly := flag.Bool("hostonly", false, "this be a host")
	clientonly := flag.Bool("clientonly", false, "this be a client")
	host := flag.String("host", "localhost:4334", "this be a host")
	flag.Parse()


	if !*clientonly {
		h := NewHost(*host)
		h.serve()
		for {
			time.Sleep(100000)
		}
	}

	if !*noterm && !*hostonly {
		c := NewClient(*host)
		t := NewTerminal(c)
		t.init()
		t.run()
	}
}
