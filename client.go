package main

import (
	"netchan"
	"os"
)

type Client struct {
	in, out chan Message

	opensend chan *OpenMessage
	openrecv chan *OpenRespMessage

	linesend chan *LineMessage
	linerecv chan *LineRespMessage

	updatesend chan *UpdateMessage
	updaterecv chan *UpdateRespMessage

	imp *netchan.Importer

	laddr string
}

type DviError struct {
	message string
}

func (e *DviError) String() string {
	return e.message
}

func NewClient(laddr string) *Client {
	c := new(Client)

	c.in, c.out = make(chan Message), make(chan Message)

	c.opensend = make(chan *OpenMessage)
	c.openrecv = make(chan *OpenRespMessage)

	c.linesend = make(chan *LineMessage)
	c.linerecv = make(chan *LineRespMessage)

	c.updatesend = make(chan *UpdateMessage)
	c.updaterecv = make(chan *UpdateRespMessage)

	imp, e := netchan.Import("tcp", laddr)
	if e != nil {
		panic(e.String())
	}
	c.laddr = laddr
	c.imp = imp

	c.imp.Import("open", c.opensend, netchan.Send, 1)
	c.imp.Import("openresp", c.openrecv, netchan.Recv, 1)

	c.imp.Import("line", c.linesend, netchan.Send, 1)
	c.imp.Import("lineresp", c.linerecv, netchan.Recv, 1)

	c.imp.Import("update", c.updatesend, netchan.Send, 1)
	c.imp.Import("updateresp", c.updaterecv, netchan.Recv, 1)

	return c
}

func (c *Client) send(cmd Message) {
	c.out <- cmd
}

func (c *Client) receive() Message {
	return <-c.in
}

// Send a message to clients host telling host to open up the file at pathname
func (c *Client) open(pathname string) (*OpenRespMessage, os.Error) {
	o := &OpenMessage{pathname, false}
	c.opensend <- o
	r := <-c.openrecv
	if r == nil {
		return nil, &DviError{"nil recieved"}
	}
	return r, nil
}

// Send a message to clients host asking for lno in file fid
func (c *Client) line(fid, first, last uint64) (*LineRespMessage, os.Error) {
	l := &LineMessage{fid, first, last}
	c.linesend <- l
	r := <-c.linerecv
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	return r, nil
}

func (c *Client) update(fid uint64, upd map[uint64]string) (*UpdateRespMessage, os.Error) {
	u := &UpdateMessage{fid, upd}
	c.updatesend <- u
	r := <-c.updaterecv
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	return r, nil
}

func (c *Client) sync(fid uint64) (*SyncRespMessage, os.Error) {
	u := &SyncMessage{fid, ""}
	c.send(u)
	r := c.receive()
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	switch m := r.(type) {
	case *SyncRespMessage:
		return m, nil
	case *ErrorRespMessage:
		return nil, &DviError{m.message()}
	default:
		return nil, &DviError{"Recieved unexpected Message type"}
	}
	return nil, nil
}

func (c *Client) stat(fid uint64) (*StatRespMessage, os.Error) {
	s := &StatMessage{fid}
	c.send(s)
	r := c.receive()
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	switch m := r.(type) {
	case *StatRespMessage:
		return m, nil
	case *ErrorRespMessage:
		return nil, &DviError{m.message()}
	default:
		return nil, &DviError{"Recieved unexpected Message type"}
	}
	return nil, nil
}

func (c *Client) close(fid uint64) (*CloseRespMessage, os.Error) {
	return nil, nil
}
