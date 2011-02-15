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
	prompt  string
	command string
}

func (e *Exline) String() string {
	return fmt.Sprintf("%s%s", e.prompt, e.command)
}

var curr *EditBuffer
var screen *Screen
var ex *Exline
var ml *Modeline

// options
var optLineNo = false

func SigHandler() {
	m := new(Modeline)
	for {
		s := <-signal.Incoming
		switch s.(signal.UnixSignal) {
		case syscall.SIGINT:
			Beep()
		case syscall.SIGTERM:
			Beep()
		case syscall.SIGWINCH:
			Beep()
		default:
			m.mode = s.String()
		}
	}
}

func Init(args []string) {
	// Setup view
	screen = NewScreen(curses.Stdwin)
	ml = new(Modeline)

	ex = new(Exline)
	ex.prompt = EXPROMPT

	var file *EditBuffer
	if len(args) == 0 {
		file = NewTempEditBuffer(TMPPREFIX)
		// XXX this is a workaround for my lazy design.  get rid
		// of this asap.
		file.InsertLine(NewLine([]byte("")))
		file.anchor = file.lines.Front()
		file.FirstLine()
	} else {
		for _, path := range args {
			if file, e := NewReadEditBuffer(path); e == nil {
				file.FirstLine()
			} else {
				file = NewTempEditBuffer(TMPPREFIX)
				file.FirstLine()
				// Ml.mode = "Error opening " + path + ": " + e.String()
			}
		}
	}
	SetCurrentFile(file)
}

func SetCurrentFile(eb *EditBuffer) {
	// call Map whenever a file becomes currfile
	curr = eb
	curr.Map()
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

func Run() {
	// UpdateDisplay()
	// enter normal mode
	NormalMode()
}

func main() {
	/* init */
	// Start in normal mode
	startScreen()
	defer endScreen()
	go SigHandler()
	// init has to happen after startscreen
	Init(os.Args[1:])
	Run()
}
