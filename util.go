package main

import (
	"fmt"
)

func editBufferInfo(b *EditBuffer) string {

	info := b.pathname + ": "
	if b.nameChanged() {
		info += "name changed: "
	}
	if !b.isDirty() {
		info += "un"
	}
	info += "modified: "
	if lns := len(b.lines); lns > 1 || len(b.line().raw()) > 0 {
		lno := b.lno + 1
		per := int((float32(lno) / float32(lns)) * 100)
		info += fmt.Sprintf("line %d of %d [%d%]", lno, lns, per)
	} else {
		info += "empty file"
	}

	return info
}
