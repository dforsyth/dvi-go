package main

// EditBuffer command map
var ebCmdMap map[int]func(*GlobalState) = map[int]func(*GlobalState) {
	'a': appendInputMode,
	'i': insertMode,
	'o': openInputMode,
	'O': aboveOpenInputMode,
	'n': nextBuffer,
	'p': prevBuffer,
	':': exMode,
	0x04: test, // ^D
}

func test(gs *GlobalState) {
	gs.queueMessage(&Message{
		"this is a test",
		true,
	})
}

// DirBuffer command map
var dbCmdMap map[int]func(*GlobalState) = map[int]func(*GlobalState) {
	'n': nextBuffer,
	'p': prevBuffer,
}

func nextBuffer(gs *GlobalState) {
	r := gs.NextBuffer()
	gs.queueMessage(&Message{
		gs.curBuf().ident(),
		r == nil,
	})
}

func prevBuffer(gs *GlobalState) {
	r := gs.PrevBuffer()
	gs.queueMessage(&Message{
		gs.curBuf().ident(),
		r == nil,
	})
}

// normal mode
func NormalMode(gs *GlobalState) {
	gs.Mode = NORMAL

	m := NewNormalModeline()
	gs.SetModeline(m)

	for {
		window := gs.Window
		window.PaintMapper(0, window.Rows-1, true)
		gs.UpdateCh <- 1
		k := <-gs.InputCh // screen.Window.Getch()

		switch b := gs.curbuf.Value.(Buffer); t := b.(type) {
		case *EditBuffer:
			if fn, ok := ebCmdMap[k]; ok {
				fn(gs)
			} else {
				b.SendInput(k)
			}
		case *DirBuffer:
			if fn, ok := dbCmdMap[k]; ok {
				fn(gs)
			}
		}

		if gs.Mode != NORMAL {
			gs.Mode = NORMAL
			gs.SetModeline(m)
		}
		m.Key = k
	}
}
