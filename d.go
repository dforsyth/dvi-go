/* A super minimal but competant programming editor.
 *
 * - modal
 * - auto indent
 * - tabs (routine per tab)
 * - jkl; instead of hjkl
 * - w, q
 * - temp files
 * - search (forward and backward)
 * - command history along the bottom of the screen in normal mode
 * - wrap or max line length
 * - tabs/spaces
 * - tabstop
 * - syntax highlighting
 * - word, line delete, copy, etc
 */

package main

import (
	"container/list"
	"curses"
	"fmt"
	// "os"
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

const (
	NORMAL  = 0
	INSERT  = 1
	COMMAND = 2
)


type Mapper interface {
	GetMap() []string
	GetCursor() (int, int)
	SetWindow(*Window)
	SetDimensions(int, int)
}

type Interacter interface {
	GetWindow() *Window
	SetWindow(*Window)
	SendInput(int)
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
	return fmt.Sprintf("INSERT -- Key: %c -- Line: %d -- Column: %d", m.Key, m.LineNumber, m.ColumnNumber)
}

func (m *InsertModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
}

type NormalModeline struct {
	Key int
}

func NewNormalModeline() *NormalModeline {
	m := new(NormalModeline)
	m.Key = ' '
	return m
}

func (m *NormalModeline) String() string {
	return fmt.Sprintf("NORMAL -- Key: %c", m.Key)
}

func (m *NormalModeline) GetCursor() int {
	// We never want the cursor for this modeline
	return -1
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
			case *EditBuffer: // Writable
				WriteFile(buffer.Pathname, buffer)
			}
		}
	}
	if quit {
		EndScreen()
		syscall.Exit(0)
	}
}

func (c *Command) Clear() {
	c.CommandBuffer = ""
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

func done() {
	EndScreen()
	syscall.Exit(0)
}

func main() {
	StartScreen()
	defer EndScreen()

	SignalsRoutine()

	gs := NewGlobalState()

	gs.Window.InputRoutine(gs.InputCh)
	gs.Window.UpdateRoutine(gs.UpdateCh)

	eb := NewTempEditBuffer(gs, TMPPREFIX)
	eb.insertLine(newEditLine([]byte("")))

	gs.AddBuffer(eb)
	gs.SetMapper(eb)
	eb.mapToScreen()

	NormalMode(gs)
	done()
}
