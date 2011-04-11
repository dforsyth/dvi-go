package main

import (
	"container/list"
	"os/signal"
	"syscall"
)

const (
	NORMAL  = 0
	INSERT  = 1
	COMMAND = 2
)

type Message struct {
	text string
	beep bool
}

type GlobalState struct {
	Window   *Window
	ex       *exBuffer
	Modeline *Modeliner
	Buffers  *list.List
	curbuf   *list.Element
	InputCh  chan int
	UpdateCh chan int
	Mode     int
	Wd       string
	config   map[string]interface{}
	msgQueue *list.List
	yb       []string // yank buffer
}

func NewGlobalState() *GlobalState {
	gs := new(GlobalState)
	gs.Window = NewWindow(gs)
	gs.ex = newExBuffer(gs)
	gs.Buffers = list.New()
	gs.curbuf = nil
	gs.InputCh = make(chan int)
	gs.UpdateCh = make(chan int)
	gs.msgQueue = list.New()
	return gs
}

func (gs *GlobalState) AddBuffer(buf Buffer) {
	gs.curbuf = gs.Buffers.PushBack(buf)
	gs.Window.buf = gs.curbuf.Value.(Buffer)
}

func (gs *GlobalState) RemoveBuffer(buf Buffer) {
	for b := gs.Buffers.Front(); b != nil; b = b.Next() {
		if b.Value == buf {
			gs.Buffers.Remove(b)
			if b == gs.curbuf {
				EndScreen()
				panic("removing curbuf is not supported yet")
			}
		}
	}
}

func (gs *GlobalState) NextBuffer() {
	if gs.curbuf.Next() != nil {
		gs.curbuf = gs.curbuf.Next()
		gs.Window.buf = gs.curbuf.Value.(Buffer)
	}
}

func (gs *GlobalState) PrevBuffer() {
	if gs.curbuf.Prev() != nil {
		gs.curbuf = gs.curbuf.Prev()
		gs.Window.buf = gs.curbuf.Value.(Buffer)
	}
}

func (gs *GlobalState) curBuf() Buffer {
	return gs.curbuf.Value.(Buffer)
}

func (gs *GlobalState) SetModeline(modeliner Modeliner) {
	gs.Modeline = &modeliner
}

func Done(r int) {
	EndScreen()
	syscall.Exit(r)
}

func (gs *GlobalState) queueMessage(msg *Message) {
	gs.msgQueue.PushBack(msg)
}

func (gs *GlobalState) getMessage() *Message {
	if f := gs.msgQueue.Front(); f != nil {
		return gs.msgQueue.Remove(f).(*Message)
	}
	return nil
}

func (gs *GlobalState) SignalsRoutine() {
	go func() {
		for {
			s := <-signal.Incoming
			switch s.(signal.UnixSignal) {
			case syscall.SIGINT:
				gs.queueMessage(&Message{
					"quit with :q",
					true,
				})
				gs.UpdateCh <- 1
				// EndScreen()
				// panic("sigint")
				// Beep()
			case syscall.SIGTERM:
				EndScreen()
				panic("sigterm")
				// Beep()
			case syscall.SIGWINCH:
				Beep()
			}
		}
	}()
}
