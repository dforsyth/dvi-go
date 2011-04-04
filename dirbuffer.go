package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type DirBuffer struct {
	Name    string
	Listing []*os.FileInfo
	Item    int
	Anchor  int
	X, Y    int
	CurY    int
	Window  *Window
	gs      *GlobalState
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
	db.Window = gs.Window
	db.CurY = 0
	db.Item = 0
	db.Anchor = 0
	db.X, db.Y = db.Window.Cols, db.Window.Rows-1
	db.gs = gs

	return db
}

func (db *DirBuffer) GetMap() *[]string {
	db.MapToScreen()
	return db.Window.ScreenMap
}

func (db *DirBuffer) GetCursor() (int, int) {
	return 0, db.CurY
}

func (db *DirBuffer) SetWindow(w *Window) {
}

func (db *DirBuffer) SetDimensions(x, y int) {
	db.X, db.Y = x, y
}

func (db *DirBuffer) GetWindow() *Window {
	return db.Window
}

func (db *DirBuffer) SendInput(k int) {
	gs := db.Window.gs
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

func (db *DirBuffer) RunRoutine(fn func(Interacter)) {
	go fn(db)
}

func (db *DirBuffer) Forward() {
	fi := db.Listing[db.Item]
	path := ""
	if db.Item == 0 && db.Name != "/" {
		path = filepath.Join(db.Name, "..")
	} else {
		path = filepath.Join(db.Name, fi.Name)
	}

	if fi.IsDirectory() {
		ndb := NewDirBuffer(db.gs, path)
		db.gs.AddBuffer(ndb)
		db.gs.SetMapper(ndb)
	} else if fi.IsRegular() {
		eb := NewEditBuffer(db.gs, path)
		if _, e := eb.readFile(path, 0); e == nil {
			db.gs.AddBuffer(eb)
			db.gs.SetMapper(eb)
			eb.GoToLine(1)
			// Now, remove this buffer
			db.gs.RemoveBuffer(db)
		} else {
			Beep()
		}
	}
}

func (db *DirBuffer) MoveUp() {
	if db.Item > 0 {
		db.Item -= 1
	} else {
		Beep()
	}
}

func (db *DirBuffer) MoveDown() {
	if db.Item < len(db.Listing)-1 {
		db.Item += 1
	} else {
		Beep()
	}
}

func (db *DirBuffer) MapToScreen() {
	smap := *db.Window.ScreenMap
	for i, fi := range db.Listing[db.Anchor:] {
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
		if fi == db.Listing[db.Item] {
			db.CurY = i
		}
	}
}
