/* 
 * Copyright (c) 2011 David Forsythe.
 * See LICENSE file for license details.
 */

package main

import (
	"os"
)

var config map[string]interface{} = map[string]interface{}{
	// edit options, set to the default specified by the spec
	// TODO: create a mapping from long to short option names
	"autoindent": false, // ai
	"autoprint":  true,  // ap
	"autowrite":  false, // aw
	"wrapscan":   true,  // ws
	"tempdir":    os.TempDir(),
	"temppfx":    "dvi.",
}

var vicmds map[int]*vicmd = map[int]*vicmd{
	' ': &vicmd{
		fn:       cmdForwards,
		motion:   false,
		isMotion: true,
	},
	'_': &vicmd{
		fn:       cmdCurrLineAndAbove,
		line:     true,
		isMotion: true,
	},
	'|': &vicmd{
		fn:       cmdMoveToColumn,
		isMotion: true,
	},
	'$': &vicmd{
		fn: cmdEOL,
	},
	'^': &vicmd{
		fn:       cmdFirstNonBlank,
		isMotion: true,
		line:     true,
	},
	'/': &vicmd{
		fn:       cmdFindRegex,
		isMotion: true,
		// can be both character and line mode
	},
	':': &vicmd{
		fn: cmdEx,
	},
	'~': &vicmd{
		fn: cmdReverseCase,
	},
	'<': &vicmd{
		fn:     cmdShiftLeft,
		motion: true,
	},
	'>': &vicmd{
		fn:     cmdShiftRight,
		motion: true,
	},
	'0': &vicmd{
		fn:       cmdBOL,
		isMotion: true,
	},
	'a': &vicmd{
		fn:     cmdAppend,
		motion: false,
	},
	'A': &vicmd{
		fn:     cmdAppendEOL,
		motion: false,
	},
	'b': &vicmd{
		fn:       cmdPrevWord,
		isMotion: true,
	},
	'B': &vicmd{
		fn: cmdPrevBigWord,
	},
	'c': &vicmd{
		fn:     cmdChange,
		motion: true,
	},
	'd': &vicmd{
		fn:     cmdDelete,
		motion: true,
	},
	'D': &vicmd{
		fn: cmdDeleteEOL,
	},
	'e': &vicmd{
		fn: cmdEndOfWord,
	},
	'E': &vicmd{
		fn: cmdEndOfBigWord,
	},
	'G': &vicmd{
		fn:        cmdToLine,
		zerocount: true,
	},
	'h': &vicmd{
		fn:       cmdBackwards,
		motion:   false,
		isMotion: true,
	},
	'i': &vicmd{
		fn:     cmdInsert,
		motion: false,
	},
	'j': &vicmd{
		fn:       cmdDown,
		motion:   false,
		isMotion: true,
		line:     true,
	},
	'k': &vicmd{
		fn:       cmdUp,
		motion:   false,
		isMotion: true,
		line:     true,
	},
	'l': &vicmd{
		fn:       cmdForwards,
		motion:   false,
		isMotion: true,
	},
	'o': &vicmd{
		fn: cmdInsertLineBelow,
	},
	'O': &vicmd{
		fn: cmdInsertLineAbove,
	},
	'p': &vicmd{
		fn:     cmdPut,
		motion: false,
	},
	'w': &vicmd{
		fn:       cmdNextWord,
		isMotion: true,
	},
	'x': &vicmd{
		fn: cmdDeleteAtCursor,
	},
	'X': &vicmd{
		fn: cmdDeleteBeforeCursor,
	},
	'y': &vicmd{
		fn:        cmdYank,
		motion:    true,
		zerocount: false,
	},
	ctrl('G'): &vicmd{
		fn: cmdDisplayInfo,
	},
	// XXX This is not a real/final command.
	ctrl('V'): &vicmd{
		fn: nextBuffer,
	},
}

var excmds map[string]*excmd = map[string]*excmd{
	"0": &excmd{
		fn: exGoToLine,
	},
	"w": &excmd{
		fn: exWriteFile,
	},
	"q": &excmd{
		fn: exQuit,
	},
	// XXX This is not a real/final command.
	"db": &excmd{
		fn: exDirBrowser,
	},
}
