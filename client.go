package main

import (
	"os"
)

type Client struct {
	out chan Message
	in  chan Message
}

type DviError struct {
	message string
}

func (e *DviError) String() string {
	return e.message
}

func NewClient(host *Host) *Client {
	c := new(Client)
	c.out = host.in
	c.in = host.out
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
