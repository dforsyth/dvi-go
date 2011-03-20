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

type Window struct {
	Curses     *curses.Window
	Cols, Rows int
	gs         *GlobalState
}

func NewWindow(gs *GlobalState) *Window {
	w := new(Window)
	w.Curses = curses.Stdwin
	w.Cols = *curses.Cols
	w.Rows = *curses.Rows
	w.gs = gs
	return w
}

func (w *Window) InputRoutine(ch chan int) {
	go func() {
		for {
			ch <- w.Curses.Getch()
		}
	}()
}

func (w *Window) UpdateRoutine(ch chan int) {
	go func() {
		for {
			<-ch
			w.PaintMapper(0, w.Rows-1, true)
			w.Curses.Refresh()
		}
	}()
}

func (w *Window) PaintMapper(start, end int, paintCursor bool) {
	cols, rows := w.Cols, w.Rows-1

	gs := w.gs
	mapper := *gs.CurrentMapper

	if start < 0 || start > rows || end > rows {
		EndScreen()
		panic(fmt.Sprintf("Window.Paint: Bad range (%d, %d) [%d, %d]", start, end, cols, rows))
	}

	smap := mapper.GetMap()
	for i := start; i < end; i++ {
		w.Curses.Move(i, 0)
		w.Curses.Clrtoeol()
		w.Curses.Mvwaddnstr(i, 0, smap[i], cols)
	}

	cX, cY := mapper.GetCursor()
	if paintCursor {
		if cX < 0 || cY < 0 || cX > cols || cY > rows {
			EndScreen()
			panic(fmt.Sprintf("Window.Paint: Bad cursor (%d, %d) [%d, %d]", start, end, cols, rows))
		}
	}
}

const (
	NORMAL  = 0
	INSERT  = 1
	COMMAND = 2
)

type GlobalState struct {
	Window        *Window
	Command       *Command
	CurrentMapper *Mapper
	Buffers       *list.List
	CurrentBuffer *list.Element
	InputCh       chan int
	UpdateCh      chan int
	Mode          int
}

func NewGlobalState() *GlobalState {
	gs := new(GlobalState)
	gs.Window = NewWindow(gs)
	gs.Command = NewCommand()
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

func (gs *GlobalState) SetMapper(mapper Mapper) {
	gs.CurrentMapper = &mapper
}

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

type ModeLiner interface {
	String() string
}

type Command struct {
	// implements ModeLiner
}

func NewCommand() *Command {
	command := new(Command)
	return command
}

func (c *Command) String() string {
	return "command"
}

func (c *Command) SendInput(k int) {
}

func (c *Command) Execute() {
}

// Modeline
type Modeline struct {
	mode          string
	char          int
	lno, lco, col int
	name          string
}

func (m *Modeline) String() string {
	return fmt.Sprintf("%s %c/%d %d/%d-%d %s", m.mode, m.char, m.char, m.lno, m.lco, m.col, m.name)
}

// ex line
type Exline struct {
	prompt string
	buff   *gapBuffer
}

func (e *Exline) String() string {
	return fmt.Sprintf("%s%s", e.prompt, e.buff.String())
}

// options
var optLineNo = true

func SignalsRoutine() {
	m := new(Modeline)
	go func() {
		for {
			s := <-signal.Incoming
			switch s.(signal.UnixSignal) {
			case syscall.SIGINT:
				EndScreen()
				panic("sigterm")
				Beep()
			case syscall.SIGTERM:
				EndScreen()
				panic("sigterm")
				// Beep()
			case syscall.SIGWINCH:
				Beep()
			default:
				m.mode = s.String()
			}
		}
	}()
}

/*
func initialize(args []string) {
	// Setup view
	screen = NewScreen(curses.Stdwin)
	ml = new(Modeline)

	// Don't allocate the cmd buffer here
	ex = new(Exline)
	ex.prompt = EXPROMPT

	var file *EditBuffer
	if len(args) == 0 {
		file = NewTempEditBuffer(TMPPREFIX)
		// XXX this is a workaround for my lazy design.  get rid
		// of this asap.
		file.insertLine(newEditLine([]byte("")))
		// file.anchor = file.lines.Front()
		// file.FirstLine()
	} else {
		for _, path := range args {
			if file, e := NewReadEditBuffer(path); e == nil {
				file.top()
			} else {
				file = NewTempEditBuffer(TMPPREFIX)
				file.top()
				// Ml.mode = "Error opening " + path + ": " + e.String()
			}
		}
	}
	setCurrentBuffer(file)
}
*/

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

	eb := NewTempEditBuffer(gs, "dtemp")
	eb.insertLine(newEditLine([]byte("")))

	gs.AddBuffer(eb)
	gs.SetMapper(eb)
	eb.mapToScreen()

	NormalMode(gs)
	done()
}
