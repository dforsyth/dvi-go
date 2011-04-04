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
	ex       *exBuffer
	CurrentMapper *Mapper
	Modeline      *Modeliner
	Buffers       *list.List
	CurrentBuffer *list.Element
	InputCh       chan int
	UpdateCh      chan int
	Mode          int
	Wd            string
	config        map[string]interface{}
}

func NewGlobalState() *GlobalState {
	gs := new(GlobalState)
	gs.Window = NewWindow(gs)
	gs.ex = newExBuffer(gs)
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

func (gs *GlobalState) NextBuffer() {
	if gs.CurrentBuffer.Next() != nil {
		gs.CurrentBuffer = gs.CurrentBuffer.Next()
		gs.SetMapper(gs.CurrentBuffer.Value.(Mapper))
	}
}

func (gs *GlobalState) PrevBuffer() {
	if gs.CurrentBuffer.Prev() != nil {
		gs.CurrentBuffer = gs.CurrentBuffer.Prev()
		gs.SetMapper(gs.CurrentBuffer.Value.(Mapper))
	}
}

func (gs *GlobalState) SetMapper(mapper Mapper) {
	newMapper := &mapper
	if newMapper != gs.CurrentMapper {
		gs.Window.ClearMap()
	}
	gs.CurrentMapper = newMapper
}

func (gs *GlobalState) SetModeline(modeliner Modeliner) {
	gs.Modeline = &modeliner
}

func Done(r int) {
	EndScreen()
	syscall.Exit(r)
}
