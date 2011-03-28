package main

import (
	"container/list"
	"syscall"
)

const (
	NORMAL  = 0
	INSERT  = 1
	COMMAND = 2
)

type GlobalState struct {
	Window        *Window
	Command       *Command
	CurrentMapper *Mapper
	Modeline      *Modeliner
	Buffers       *list.List
	CurrentBuffer *list.Element
	InputCh       chan int
	UpdateCh      chan int
	Mode          int
	Wd		string
}

func NewGlobalState() *GlobalState {
	gs := new(GlobalState)
	gs.Window = NewWindow(gs)
	gs.Command = NewCommand(gs)
	gs.CurrentMapper = nil
	gs.Buffers = list.New()
	gs.CurrentBuffer = nil
	gs.InputCh = make(chan int)
	gs.UpdateCh = make(chan int)
	return gs
}

func (gs *GlobalState) AddBuffer(buffer Interacter) {
	gs.CurrentBuffer = gs.Buffers.PushBack(buffer)
}

func (gs *GlobalState) RemoveBuffer(buffer Interacter) {
	for b := gs.Buffers.Front(); b != nil; b = b.Next() {
		if b.Value == buffer {
			gs.Buffers.Remove(b)
		}
	}
}

func (gs *GlobalState) SetMapper(mapper Mapper) {
	gs.CurrentMapper = &mapper
}

func (gs *GlobalState) SetModeline(modeliner Modeliner) {
	gs.Modeline = &modeliner
}

func Done(r int) {
	EndScreen()
	syscall.Exit(r)
}
