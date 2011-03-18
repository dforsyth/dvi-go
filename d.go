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

var curr *editBuffer
var screen *Screen
var ex *Exline
var ml *Modeline

var inputch chan int

// options
var optLineNo = true

func sigHandlerRoutine() {
	m := new(Modeline)
	for {
		s := <-signal.Incoming
		switch s.(signal.UnixSignal) {
		case syscall.SIGINT:
			panic("sigterm")
			Beep()
		case syscall.SIGTERM:
			panic("sigterm")
			// Beep()
		case syscall.SIGWINCH:
			Beep()
		default:
			m.mode = s.String()
		}
	}
}

func inputRoutine() {
	inputch = make(chan int)
	for {
		inputch <- screen.Window.Getch()
	}
}

func initialize(args []string) {
	// Setup view
	screen = NewScreen(curses.Stdwin)
	ml = new(Modeline)

	// Don't allocate the cmd buffer here
	ex = new(Exline)
	ex.prompt = EXPROMPT

	var file *editBuffer
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

func setCurrentBuffer(eb *editBuffer) {
	// call Map whenever a file becomes currfile
	curr = eb
	curr.mapToScreen()
}

func startScreen() {
	curses.Initscr()
	curses.Cbreak()
	curses.Noecho()
	curses.Nonl()
	curses.Stdwin.Keypad(true)
}

func endScreen() {
	curses.Endwin()
}

func run() {
	go screen.ScreenRoutine()
	// enter normal mode
	NormalMode()
}

func done() {
	endScreen()
	syscall.Exit(0)
}

func main() {
	startScreen()
	defer endScreen()
	go sigHandlerRoutine()
	initialize(os.Args[1:])
	go inputRoutine()
	run()
	done()
}
