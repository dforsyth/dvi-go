package main

import (
	"log"
	"netchan"
	"os"
)

type Host struct {
	openrecv chan *OpenMessage
	opensend chan *OpenRespMessage

	linerecv chan *LineMessage
	linesend chan *LineRespMessage

	updaterecv chan *UpdateMessage
	updatesend chan *UpdateRespMessage

	log   *log.Logger
	files map[uint64]*File

	laddr string
}

func ErrorResponse(e os.Error) *ErrorRespMessage {
	return &ErrorRespMessage{
		e.String(),
	}
}

func NewHost(laddr string) *Host {
	h := &Host {
		openrecv: make(chan *OpenMessage),
		opensend: make(chan *OpenRespMessage),
		linerecv: make(chan *LineMessage),
		linesend: make(chan *LineRespMessage),
		updaterecv: make(chan *UpdateMessage),
		updatesend: make(chan *UpdateRespMessage),
		files: make(map[uint64]*File),
		log: log.New(os.Stderr, "host(" + laddr + ")", 0),
		laddr: laddr,
	}

	return h
}

func (h *Host) serve() {
	exp := netchan.NewExporter()

	exp.Export("open", h.openrecv, netchan.Recv)
	exp.Export("openresp", h.opensend, netchan.Send)

	exp.Export("line", h.linerecv, netchan.Recv)
	exp.Export("lineresp", h.linesend, netchan.Send)

	exp.Export("update", h.updaterecv, netchan.Recv)
	exp.Export("updateresp", h.updatesend, netchan.Send)

	exp.ListenAndServe("tcp", h.laddr)

	var m Message
	for {
		select {
		case m = <-h.openrecv:
		case m = <-h.linerecv:
		// case m = <-h.statrecv:
		case m = <-h.updaterecv:
		// case m = <-h.syncrecv:
		}

		switch t := m.(type) {
		case *OpenMessage:
			r, _ := h.open(t)
			h.opensend <- r
		case *LineMessage:
			r, _ := h.line(t)
			h.linesend <- r
		case *UpdateMessage:
			h.log.Println("update message")
			r, _ := h.update(t)
			h.updatesend <- r
			h.log.Println("update finished")
		default:
		}

		h.log.Println(m.message())
	}

	return
}
