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
	"os"
	"os/signal"
	"syscall"
)

const (
	// strings
	NaL string = "+" // char that shows that a line does not exist

	TMPDIR    = "."
	TMPPREFIX = "d." // temp file prefix

	// modes
	MODENORMAL int = 1
	MODEINSERT int = 0

	EXPROMPT string = ":"
)

var Debug string = ""
var Message string = ""

// global state for the editor
type D struct {
	mode int // current mode

	buf *EditBuffer

	yank *list.List // list of lines in the current yank buff

	view *View

}

type Status struct {
	mode     string
	char     byte
	row, col int
}

// editor state
var d D

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
			Message = s.String()
		}
	}
}

func dInit(args []string) {
	d.mode = -1
	if len(args) == 0 {
		InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
		d.buf.FirstLine()
	} else {
		for _, path := range args {
			InsertBuffer(NewReadFileEditBuffer(path))
			d.buf.FirstLine()
		}
	}

	d.yank = list.New()

	// ready view
	d.view = new(View)
	d.view.win = curses.Stdwin
	d.view.rows = *curses.Rows
	d.view.cols = *curses.Cols
}

func InsertBuffer(b *EditBuffer) {
	if d.buf == nil {
		d.buf = b
	} else {
		d.buf.next = b
	}
}

func NextBuffer() *EditBuffer {
	if d.buf == nil {
		return nil
	}

	d.buf = d.buf.next
	return d.buf
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

func dRun() {

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
	dInit(os.Args[1:])
	dRun()
}

