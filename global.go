package main

import (
	"container/list"
	"os"
	"os/signal"
	"syscall"
)

const (
	NORMAL  = 0
	INSERT  = 1
	COMMAND = 2
)

type DviError struct {
	msg string
}

func (e *DviError) String() string {
	return e.msg
}

type Message struct {
	text string
	beep bool
}

const (
	MODEINSERT = iota
	MODEREPLACE
	MODENORMAL
	MODEEX
)

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
	yb       []string            // yank buffer
	ub       map[int][]string    // unnamed buffers
	nb       map[string][]string // named buffers

	cmd string
	x   *Ex
	n   *Nm

	version   string
	buildDate string
	author    string
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

	// lol
	gs.version = "0.0"
	gs.buildDate = "0/0/20XX"
	gs.author = "David Forsythe"
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

func (gs *GlobalState) NextBuffer() Buffer {
	if n := gs.curbuf.Next(); n != nil {
		gs.curbuf = n
		gs.Window.buf = gs.curbuf.Value.(Buffer)
		return gs.Window.buf
	}
	return nil
}

func (gs *GlobalState) PrevBuffer() Buffer {
	if p := gs.curbuf.Prev(); p != nil {
		gs.curbuf = p
		gs.Window.buf = gs.curbuf.Value.(Buffer)
		return gs.Window.buf
	}
	return nil
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
					"Interrupted",
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

func (gs *GlobalState) parseConfig(pathname string) os.Error {
	f, e := os.Open(pathname)
	if e != nil {
		return e
	}
	defer f.Close()

	return nil
}
