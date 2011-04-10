package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type DirBuffer struct {
	Name    string
	Listing []*os.FileInfo
	item    int
	head  int
	X, Y    int
	CurY    int
	gs      *GlobalState
	dirty bool
}

func NewDirBuffer(gs *GlobalState, name string) *DirBuffer {
	db := new(DirBuffer)
	db.Name = name
	if listing, err := ioutil.ReadDir(db.Name); err == nil {
		if db.Name != "/" {
			db.Listing = make([]*os.FileInfo, len(listing)+1)
			if par, e := os.Lstat(filepath.Join(db.Name, "..")); e == nil {
				db.Listing[0] = par
			} else {
				panic(e.String())
			}
			copy(db.Listing[1:], listing)
		} else {
			db.Listing = listing
		}
	} else {
		// For now, panic.  Really should be sending a message to the modeline and returning nil
		panic(err.String())
	}
	db.gs = gs
	db.CurY = 0
	db.item = 0
	db.head = 0
	db.X, db.Y = db.gs.Window.Cols, db.gs.Window.Rows-1

	return db
}

func (db *DirBuffer) mapScreen() {
	db.MapToScreen()
}

func (db *DirBuffer) getCursor() (int, int) {
	return 0, db.CurY
}

func (db *DirBuffer) getWindow() *Window {
	return db.gs.Window
}

func (db *DirBuffer) SetWindow(w *Window) {
}

func (db *DirBuffer) SetDimensions(x, y int) {
	db.X, db.Y = x, y
}

func (db *DirBuffer) GetWindow() *Window {
	return db.gs.Window
}

func (db *DirBuffer) SendInput(k int) {
	gs := db.gs
	switch gs.Mode {
	case INSERT, NORMAL:
		switch k {
		case 0xd, 0xa:
			db.Forward()
		case 'k':
			db.MoveDown()
		case 'l':
			db.MoveUp()
		}
	case COMMAND: // How did you get here?
	}
}

func (db *DirBuffer) RunRoutine(fn func(Buffer)) {
	go fn(db)
}

func (db *DirBuffer) Forward() {
	fi := db.Listing[db.item]
	path := ""
	if db.item == 0 && db.Name != "/" {
		path = filepath.Join(db.Name, "..")
	} else {
		path = filepath.Join(db.Name, fi.Name)
	}

	if fi.IsDirectory() {
		ndb := NewDirBuffer(db.gs, path)
		db.gs.AddBuffer(ndb)
	} else if fi.IsRegular() {
		eb := NewEditBuffer(db.gs, path)
		f, e := os.Open(path)
		if e != nil {
			panic(e.String())
		}

		if _, e := eb.readFile(f, 0); e == nil {
			db.gs.AddBuffer(eb)
			eb.gotoLine(1)
			// Now, remove this buffer
			db.gs.RemoveBuffer(db)
		} else {
			panic(e)
		}
	}
}

func (db *DirBuffer) MoveUp() {
	if db.item > 0 {
		db.item -= 1
		if db.item < db.head {
			db.head = db.item
			db.dirty = true
		}
	} else {
		Beep()
	}
}

func (db *DirBuffer) MoveDown() {
	if db.item < len(db.Listing)-1 {
		db.item += 1
		if db.item > db.head + db.Y-1 {
			db.head = db.item - db.Y + 1
			db.dirty = true
		}
	} else {
		Beep()
	}
}

func (db *DirBuffer) MapToScreen() {
	smap := db.gs.Window.screenMap
	for i, _ := range smap {
		smap[i] = ""
	}
	for i, fi := range db.Listing[db.head:] {
		if i > db.Y-1 {
			break
		}
		smap[i] = fi.Name
		if i == 0 && db.Name != "/" {
			smap[i] = ".."
		}
		if fi.IsDirectory() {
			smap[i] += "/"
		}
		if fi.IsSymlink() {
			smap[i] = "@" + smap[i]
		}
		if fi == db.Listing[db.item] {
			db.CurY = i
		}
	}
}
