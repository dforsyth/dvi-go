package main

import (
	"curses"
	"fmt"
	"os"
)

const (
	// strings
	NaL       string = "~" // char that shows that a line does not exist
	EXPROMPT  string = ":"
	TMPDIR    = "."
	TMPPREFIX = "dvi." // temp file prefix
	ESC       = 27
)

// This could/should probably be called Buffer(er) or something...
type Buffer interface {
	mapScreen()
	getWindow() *Window
	getCursor() (int, int)
	SetDimensions(int, int)
	SendInput(int)
	RunRoutine(func(Buffer))
}

type Modeliner interface {
	String() string
	GetCursor() int
	msgOverride(*Message)
}

type InsertModeline struct {
	Key          int
	LineNumber   int
	ColumnNumber int
	msg          *Message
}

func NewInsertModeline() *InsertModeline {
	m := new(InsertModeline)
	m.Key = ' '
	m.LineNumber = -1
	m.ColumnNumber = -1
	m.msg = nil
	return m
}

func (m *InsertModeline) String() string {
	show := "INSERT"
	if m.msg != nil {
		show = m.msg.text
	}

	ml := fmt.Sprintf("%s -- Key: %c (%d)-- Line: %d -- Column: %d", show, m.Key, m.Key, m.LineNumber, m.ColumnNumber)

	m.msg = nil
	return ml
}

func (m *InsertModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
}

func (m *InsertModeline) msgOverride(msg *Message) {
	m.msg = msg
}

type NormalModeline struct {
	Key      int
	info     string
	Row, Col int
	msg      *Message
}

func NewNormalModeline() *NormalModeline {
	m := new(NormalModeline)
	m.Key = ' '
	m.info = ""
	m.Row, m.Col = 0, 0
	m.msg = nil
	return m
}

func (m *NormalModeline) String() string {
	show := "NORMAL"
	if m.msg != nil {
		show = m.msg.text
	}

	ml := fmt.Sprintf("%s -- Key: %c", show, m.Key)
	m.msg = nil
	return ml
}

func (m *NormalModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
}

func (m *NormalModeline) msgOverride(msg *Message) {
	m.msg = msg
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

	gs := NewGlobalState()
	gs.SignalsRoutine()
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
				} else if fi.IsRegular() {
					eb := NewEditBuffer(gs, path)
					f, e := os.Open(path)
					if e != nil {
						panic(e.String())
					}

					if _, e := eb.readFile(f, 0); e == nil {
						gs.AddBuffer(eb)
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
	}

	NormalMode(gs)
	Done(0)
}
