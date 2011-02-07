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
	NaL      string = "~" // char that shows that a line does not exist
	EXPROMPT string = ":"

	TMPDIR    = "."
	TMPPREFIX = "d." // temp file prefix
	ESC       = 27
)

// Message lines
type Message interface {
	String() string
}

// Modeline
type Modeline struct {
	mode          string
	char          int
	lno, lco, col int
}

func (m *Modeline) String() string {
	return fmt.Sprintf("%s %c/%d %d/%d-%d", m.mode, m.char, m.char, m.lno, m.lco, m.col)
}

// ex line
type Exline struct {
	prompt  string
	command string
}

func (e *Exline) String() string {
	return fmt.Sprintf("%s%s", e.prompt, e.command)
}

var OptLineNumbers = true

var Eb *EditBuffer
var Ml *Modeline
var Vw *View

func SigHandler() {
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
			Ml.mode = s.String()
		}
	}
}

func Init(args []string) {
	// Setup modeline
	Ml = new(Modeline)
	Ml.mode = ""
	Ml.char = '@'
	Ml.lno = 0
	Ml.lco = 0
	Ml.col = 0

	// Setup view
	Vw = NewView()

	if len(args) == 0 {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
		// XXX this is a workaround for my lazy design.  get rid
		// of this asap.
		Eb.InsertLine(NewLine([]byte("")))
		Eb.anchor = Eb.lines.Front()
		Eb.FirstLine()
	} else {
		for _, path := range args {
			if b, e := NewReadFileEditBuffer(path); e == nil {
				InsertBuffer(b)
				Eb.FirstLine()
			} else {
				InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
				Eb.FirstLine()
				Ml.mode = "Error opening " + path + ": " + e.String()
			}
		}
	}
}

func InsertBuffer(b *EditBuffer) {
	if Eb == nil {
		Eb = b
	} else {
		Eb.next = b
	}
}

func NextBuffer() *EditBuffer {
	if Eb == nil {
		return nil
	}

	Eb = Eb.next
	return Eb
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
	UpdateDisplay()
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
