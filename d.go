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
)

const (
	// strings
	NaL string = "+" // char that shows that a line does not exist

	TMPDIR    = "."
	TMPPREFIX = "d.tmp." // temp file prefix

	// modes
	MODENORMAL int = 1
	MODEINSERT int = 0
)

var Debug string = ""

// global state for the editor
// TODO kill this struct
type D struct {
	mode int // current mode

	bufs *list.List
	buf  *list.Element

	yank *list.List // list of lines in the current yank buff

	s int // the first row we give to the edit buffer
	e int // number of rows we give to the editor

	err string // error string to display

	win *curses.Window
}

var win *curses.Window

func Log(msg string) {
	// send msg to dbg.txt
}

func (d *D) init(args []string) {
	d.mode = MODENORMAL

	d.bufs = list.New()
	if len(args) == 0 {
		d.InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
		d.Buffer().FirstLine()
	} else {
		for _, path := range args {
			d.InsertBuffer(NewReadFileEditBuffer(path))
			d.Buffer().FirstLine()
		}
	}

	d.s = 1
	d.e = *curses.Rows - 1
	d.win = curses.Stdwin

	// newer init
	win = curses.Stdwin
}

func (d *D) Buffer() *EditBuffer {
	if d.buf == nil {
		return nil
	}
	return d.buf.Value.(*EditBuffer)
}

func (d *D) InsertBuffer(b *EditBuffer) {
	if d.buf == nil {
		d.buf = d.bufs.PushFront(b)
	} else {
		d.buf = d.bufs.InsertAfter(b, d.buf)
	}
}

func (d *D) NextBuffer() {
	if d.buf == nil {
		return
	}

	d.buf = d.buf.Next()
}

func (d *D) Mode() int {
	return d.mode
}

func (d *D) SetError(err string) {
	d.err = err
}

func (d *D) Error() string {
	return d.err
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

func (d *D) run() {

	Debug = "normal"
	d.UpdateDisplay()
	// enter normal mode
	d.NormalMode()

}


func (d *D) StatusLine() string {
	return Debug
}

func main() {
	/* init */
	// Start in normal mode
	d := new(D)
	startScreen()
	defer endScreen()

	// init has to happen after startscreen
	d.init(os.Args[1:])
	d.run()
}
