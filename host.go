package main

import (
	"log"
	"netchan"
	"os"
)

type Host struct {
	in    chan Message
	out   chan Message
	log   *log.Logger
	files map[uint64]*File
}

func ErrorResponse(e os.Error) *ErrorRespMessage {
	return &ErrorRespMessage{
		e.String(),
	}
}

func NewHost() *Host {
	h := new(Host)
	h.log = log.New(os.Stderr, "host", 0)
	h.in = make(chan Message)
	h.out = make(chan Message)
	h.files = make(map[uint64]*File)
	return h
}

func (h *Host) serve() {
	exp := netchan.NewExporter()

	exp.Export("dviToHost", h.in, netchan.Recv)
	exp.Export("dviToClient", h.out, netchan.Send)

	exp.ListenAndServe("tcp", "localhost:4334")

	go func() {
		for {
			// wait for a command
			switch c := <-h.in; m := c.(type) {
			case *OpenMessage:
				// fmt.Println(c.message())
				r, e := h.open(m)
				if e != nil {
					// log.Panicln(e.String())
					h.out <- ErrorResponse(e)
					break
				}
				h.out <- r
			case *StatMessage:
				// fmt.Println(c.message())
				r, e := h.stat(m)
				if e != nil {
					log.Panicln(e.String())
					h.out <- nil
				}
				h.out <- r
			case *LineMessage:
				// fmt.Println(c.message())
				r, e := h.line(m)
				if e != nil {
					// log.Panicln(e.String())
					h.out <- ErrorResponse(e)
					break
				}
				h.out <- r
			case *UpdateMessage:
				r, e := h.update(m)
				if e != nil {
					h.out <- ErrorResponse(e)
					break
				}
				h.out <- r
			case *SyncMessage:
				r, e := h.sync(m)
				if e != nil {
					h.out <- ErrorResponse(e)
					break
				}
				h.out <- r
			default:
				h.out <- ErrorResponse(&DviError{"unknown message"})
			}
		}
	}()
	return
}
