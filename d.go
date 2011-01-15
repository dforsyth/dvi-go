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
	"fmt"
)

const (
	// strings
	NaL string = "+" // char that shows that a line does not exist

	TMPDIR    = "."
	TMPPREFIX = "d.tmp." // temp file prefix

	// modes
	MNORMAL string = "NORMAL"
	MINSERT string = "INSERT"
)

var Debug string = ""

// global state for the editor
type D struct {
	m string // current mode

	bufs *list.List
	buf  *list.Element

	yank *list.List // list of lines in the current yank buff

	s int // the first row we give to the edit buffer
	e int // number of rows we give to the editor

	err string // error string to display
}

func Log(msg string) {
	// send msg to dbg.txt
}

func (d *D) init(args []string) {
	d.m = MNORMAL

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

func (d *D) Mode() string {
	return d.m
}

func (d *D) ModeNormal() {
	d.m = MNORMAL

	if d.Buffer() != nil && d.Buffer().Line() != nil {
		d.Buffer().Line().UpdateCursor()
	}
}

func (d *D) ModeInsert() {
	d.m = MINSERT

	if d.Buffer() == nil {
		d.InsertBuffer(NewTempFileEditBuffer(TMPPREFIX))
	}

	if d.Buffer().Line() == nil {
		d.Buffer().AppendLine()
	}

	d.Buffer().Line().MoveGapToCursor()
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

	d.UpdateDisplay()

	for {

		i := curses.Stdwin.Getch()

		if d.Mode() == MNORMAL {
			if i == 'j' {
				Debug = "j"
				d.Buffer().MoveCursorLeft()
			} else if i == 'k' {
				Debug = "k"
				d.Buffer().MoveCursorDown()
			} else if i == 'l' {
				Debug = "l"
				d.Buffer().MoveCursorUp()
			} else if i == ';' {
				Debug = ";"
				d.Buffer().MoveCursorRight()
			} else if i == 'i' {
				Debug = "insert"
				d.ModeInsert()
			} else if i == 'w' {
				fi, e := WriteEditBuffer(d.Buffer().Title(), d.Buffer())
				if e != nil {
					Debug = "write failed " + e.String()
				} else {
					Debug = fmt.Sprintf("wrote %d bytes", fi.Size)
				}
			}
		} else if d.Mode() == MINSERT {
			if i == 27 {
				d.ModeNormal()
			} else if i == 0x7f {
				// improperly handles the newline at the end of the prev line
				d.Buffer().BackSpace()
			} else if i == 0xd {
				if d.Buffer().Line() != nil {
					d.Buffer().InsertChar(byte('\n'))
				}
				d.Buffer().InsertLine(NewGapBuffer([]byte("")))
			} else {
				if d.Buffer().Line() == nil {
					d.Buffer().InsertLine(NewGapBuffer([]byte("")))
				}
				d.Buffer().InsertChar(byte(i))
			}
		}

		d.UpdateDisplay()
	}
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
