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
	"curses"
	"fmt"
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
	return fmt.Sprintf("INSERT -- Key: %c -- Line: %d -- Column: %d", m.Key, m.LineNumber, m.ColumnNumber)
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

	gs.Window.InputRoutine(gs.InputCh)
	gs.Window.UpdateRoutine(gs.UpdateCh)

	eb := NewTempEditBuffer(gs, TMPPREFIX)
	eb.insertLine(newEditLine([]byte("")))

	gs.AddBuffer(eb)
	gs.SetMapper(eb)
	eb.MapToScreen()

	NormalMode(gs)
	Done(0)
}
