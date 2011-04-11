package main

import (
//"fmt"
)

func NewTempEditBuffer(gs *GlobalState, prefix string) *EditBuffer {
	// TODO: this.
	e := NewEditBuffer(gs, prefix)
	e.temp = true
	return e
}
