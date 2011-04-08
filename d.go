package main

import (
	"curses"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	// strings
	NaL       string = "~" // char that shows that a line does not exist
	EXPROMPT  string = ":"
	TMPDIR    = "."
	TMPPREFIX = "d." // temp file prefix
	ESC       = 27
)

// This could/should probably be called Buffer(er) or something...
type Mapper interface {
	GetMap() *[]string
	GetCursor() (int, int)
	SetWindow(*Window)
	SetDimensions(int, int)
}

type Interacter interface {
	GetWindow() *Window
	SetWindow(*Window)
	SendInput(int)
	RunRoutine(func(Interacter))
}

type Modeliner interface {
	String() string
	GetCursor() int
}

type InsertModeline struct {
	Key          int
	LineNumber   int
	ColumnNumber int
}

func NewInsertModeline() *InsertModeline {
	m := new(InsertModeline)
	m.Key = ' '
	m.LineNumber = -1
	m.ColumnNumber = -1
	return m
}

func (m *InsertModeline) String() string {
	return fmt.Sprintf("INSERT -- Key: %c (%d)-- Line: %d -- Column: %d", m.Key, m.Key, m.LineNumber, m.ColumnNumber)
}

func (m *InsertModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
}

type NormalModeline struct {
	Key      int
	Message  string
	Row, Col int
}

func NewNormalModeline() *NormalModeline {
	m := new(NormalModeline)
	m.Key = ' '
	m.Message = ""
	m.Row, m.Col = 0, 0
	return m
}

func (m *NormalModeline) String() string {
	return fmt.Sprintf("NORMAL -- Key: %c", m.Key)
}

func (m *NormalModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
}


// options
var optLineNo = true

func SignalsRoutine() {
	go func() {
		for {
			s := <-signal.Incoming
			switch s.(signal.UnixSignal) {
			case syscall.SIGINT:
				EndScreen()
				panic("sigint")
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

func StartScreen() {
	curses.Initscr()
	curses.Cbreak()
	curses.Noecho()
	curses.Nonl()
	curses.Stdwin.Keypad(true)
}

func EndScreen() {
	curses.Endwin()
}

func main() {
	StartScreen()
	defer EndScreen()

	SignalsRoutine()

	gs := NewGlobalState()
	wd, e := os.Getwd()
	if e != nil {
		panic(e.String())
	}
	gs.Wd = wd
	gs.Window.InputRoutine(gs.InputCh)
	gs.Window.UpdateRoutine(gs.UpdateCh)

	if len(os.Args) > 1 {
		for _, path := range os.Args[1:] {
			if fi, e := os.Stat(path); e == nil {
				if fi.IsDirectory() {
					if fi.Name == "" {
						fi.Name = "/"
					}
					db := NewDirBuffer(gs, path)
					gs.AddBuffer(db)
					gs.SetMapper(db)
				} else if fi.IsRegular() {
					eb := NewEditBuffer(gs, path)
					f, e := os.Open(path)
					if e != nil {
						panic(e.String())
					}

					if _, e := eb.readFile(f, 0); e == nil {
						gs.AddBuffer(eb)
						gs.SetMapper(eb)
						eb.gotoLine(1)
					} else {
						panic(e.String())
					}
					f.Close()
				}
			} else {
				panic(e.String())
			}
		}
	} else {
		eb := NewTempEditBuffer(gs, TMPPREFIX)
		eb.insert(NewEditLine([]byte("")), 0) // Insert the initial line per vi
		gs.AddBuffer(eb)
		gs.SetMapper(eb)
	}

	NormalMode(gs)
	Done(0)
}
