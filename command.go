package main

import (
	"container/list"
	"fmt"
)

func CommandMode(gs *GlobalState) {

	gs.SetModeline(gs.Command)
	gs.Command.Reset()

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		switch k {
		case ESC:
			return
		case 0xd, 0xa:
			gs.Command.Execute()
			return
		default:
			gs.Command.SendInput(k)
		}
	}
}

type Command struct {
	CommandBuffer string
	gs            *GlobalState
}

func NewCommand(gs *GlobalState) *Command {
	c := new(Command)
	c.CommandBuffer = ""
	c.gs = gs
	return c
}

func (c *Command) String() string {
	return fmt.Sprintf(":%s", c.CommandBuffer)
}

func (c *Command) GetCursor() int {
	return len(c.String()) - 1
}

func (c *Command) SendInput(k int) {
	c.CommandBuffer += string(k)
}

func (c *Command) Execute() {
	save := false
	quit := false
	all := false
	targets := list.New()
	targets.Init()

	for _, c := range c.CommandBuffer {
		switch c {
		case 'w':
			save = true
		case 'q':
			quit = true
		case 'a':
			all = true
		}
	}

	gs := c.gs

	if !all {
		targets.PushFront(gs.CurrentBuffer.Value)
	} else {
		targets.PushFrontList(gs.Buffers)
	}

	for t := targets.Front(); t != nil; t = t.Next() {
		if save {
			switch buffer := t.Value.(type) {
			case *EditBuffer: // I should make these io.Writer s
				WriteFile(buffer.Pathname, buffer)
			}
		}
	}
	if quit {
		Done(0)
	}
}

func (c *Command) Reset() {
	c.CommandBuffer = ""
}
