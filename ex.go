package main

import (
	"container/list"
	"fmt"
	"os"
)

func exMode(gs *GlobalState) {

	gs.SetModeline(gs.ex)
	gs.ex.Reset()

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh

		switch k {
		case ESC:
			return
		case 0xd, 0xa:
			gs.ex.execute()
			return
		default:
			gs.ex.SendInput(k)
		}
	}
}

type exBuffer struct {
	buffer string
	gs     *GlobalState
}

func newExBuffer(gs *GlobalState) *exBuffer {
	c := new(exBuffer)
	c.buffer = ""
	c.gs = gs
	return c
}

func (c *exBuffer) String() string {
	return fmt.Sprintf(":%s", c.buffer)
}

func (c *exBuffer) GetCursor() int {
	return len(c.String()) - 1
}

func (c *exBuffer) msgOverride(m *Message) {
}

func (c *exBuffer) SendInput(k int) {
	c.buffer += string(k)
}

func (c *exBuffer) execute() {
	save := false
	quit := false
	all := false
	targets := list.New()
	targets.Init()

	for _, c := range c.buffer {
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
		targets.PushFront(gs.curbuf.Value)
	} else {
		targets.PushFrontList(gs.Buffers)
	}

	for t := targets.Front(); t != nil; t = t.Next() {
		if save {
			switch bt := t.Value.(type) {
			case *EditBuffer: // I should make these io.Writer s
				eb := t.Value.(*EditBuffer)
				if eb.temp == true {
					c.gs.queueMessage(&Message{
						"Buffer has no non-temp pathname",
						true,
					})
					continue
				}
				// XXX rewrite the entire file, like a boss.
				if f, e := os.Create(eb.pathname); e == nil {
					eb.writeFile(f)
				} else {
					EndScreen()
					panic(e.String())
				}
			}
		}
	}
	if quit {
		Done(0)
	}
}

func (c *exBuffer) Reset() {
	c.buffer = ""
}
