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
	NaL string = "~" // char that shows that a line does not exist
	EXPROMPT string = ":"

	TMPDIR    = "."
	TMPPREFIX = "d." // temp file prefix
)

var Debug string = ""

type Message interface {
	String() string
}

// Modeline
type Modeline struct {
	mode     string
	char     byte
	lno, col int
}

type Exline struct {
	prompt  string
	command string
}

func (e *Exline) String() string {
	return fmt.Sprintf("%s%s", e.prompt, e.command)
}

func (m *Modeline) String() string {
	return fmt.Sprintf("%s %b %d-%d", m.mode, m.char, m.lno, m.col)
}

var eb *EditBuffer
var ml *Modeline
var vw *View

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
			ml.mode = s.String()
		}
	}
}

func Init(args []string) {
	if len(args) == 0 {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
		eb.FirstLine()
	} else {
		for _, path := range args {
			if b, e := NewReadFileEditBuffer(path); e == nil {
				InsertBuffer(b)
				eb.FirstLine()
			}
		}
	}

	// Setup view
	vw = new(View)
	vw.win = curses.Stdwin
	vw.rows = *curses.Rows
	vw.cols = *curses.Cols

	// Setup modeline
	ml = new(Modeline)
	ml.mode = ""
	ml.char = '@'
	ml.lno = 0
	ml.col = 0
}

func InsertBuffer(b *EditBuffer) {
	if eb == nil {
		eb = b
	} else {
		eb.next = b
	}
}

func NextBuffer() *EditBuffer {
	if eb == nil {
		return nil
	}

	eb = eb.next
	return eb
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

func statusLine() string {
	return Debug
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
