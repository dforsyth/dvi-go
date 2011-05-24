package main

import (
	"netchan"
	"os"
)

type Client struct {
	out chan Message
	in  chan Message
	imp *netchan.Importer
}

type DviError struct {
	message string
}

func (e *DviError) String() string {
	return e.message
}

func NewClient() *Client {
	c := new(Client)
	c.out = make(chan Message)
	c.in = make(chan Message)

	imp, e := netchan.Import("tcp", "localhost:4334")
	if e != nil {
		panic(e.String())
	}
	c.imp = imp
	c.imp.Import("dviToHost", c.out, netchan.Send, 1)
	c.imp.Import("dviToClient", c.in, netchan.Recv, 1)

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
	c.send(o)
	r := c.receive()
	if r == nil {
		return nil, &DviError{"nil recieved"}
	}

	switch m := r.(type) {
	case *OpenRespMessage:
		return m, nil
	case *ErrorRespMessage:
		return nil, &DviError{m.message()}
	default:
		return nil, &DviError{"Recieved unexpected Message type"}
	}
	// NOT REACHED
	return nil, nil
}

// Send a message to clients host asking for lno in file fid
func (c *Client) line(fid, first, last uint64) (*LineRespMessage, os.Error) {
	l := &LineMessage{fid, first, last}
	c.send(l)
	r := c.receive()
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	switch m := r.(type) {
	case *LineRespMessage:
		return m, nil
	case *ErrorRespMessage:
		return nil, &DviError{m.message()}
	default:
		return nil, &DviError{"Recieved unexpected Message type"}
	}
	return nil, nil
}

func (c *Client) update(fid uint64, upd map[uint64]string) (*UpdateRespMessage, os.Error) {
	u := &UpdateMessage{fid, upd}
	c.send(u)
	r := c.receive()
	if r == nil {
		return nil, &DviError{"nil received"}
	}
	switch m := r.(type) {
	case *UpdateRespMessage:
		return m, nil
	case *ErrorRespMessage:
		return nil, &DviError{m.message()}
	default:
		return nil, &DviError{"Recieved unexpected Message type"}
	}
	return nil, nil
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
