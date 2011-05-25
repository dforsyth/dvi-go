package main

import (
	"fmt"
	"os"
)

func (h *Host) open(m *OpenMessage) (*OpenRespMessage, os.Error) {
	path := m.Pathname
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	file, e := NewFile(path, f)
	if e != nil {
		return nil, e
	}

	if _, ok := h.files[file.fid]; !ok {
		h.files[file.fid] = file
	} else {
		return nil, &DviError{"Fid already exists"}
	}
	o := NewOpenRespMessage(file)
	return o, nil
}

func (h *Host) stat(m *StatMessage) (*StatRespMessage, os.Error) {
	file, ok := h.files[m.fid]
	if !ok {
		return nil, &DviError{fmt.Sprintf("Fid %d not in files map")}
	}
	i, e := file.fileInfo()
	if e != nil {
		return nil, e
	}
	s := NewStatRespMessage(file, i)
	return s, nil
}

func (h *Host) line(m *LineMessage) (*LineRespMessage, os.Error) {
	file, ok := h.files[m.Fid]
	if !ok {
		return NewLineRespMessage(nil, m.Fid), &DviError{fmt.Sprintf("Fid %d not in files map")}
	}
	if m.Last > m.First {
		return NewLineRespMessage(nil, m.Fid), &DviError{fmt.Sprintf("Last and first out of order: %d > %d", m.Last,
			m.First)}
	}
	if m.First > uint64(len(file.buf)-1) {
		return NewLineRespMessage(nil, m.Fid), &DviError{fmt.Sprintf("First is out of range: %d > %d", m.First,
			len(file.buf)-1)}
	}
	if m.Last > uint64(len(file.buf)-1) {
		return NewLineRespMessage(nil, m.Fid), &DviError{fmt.Sprintf("First is out of range: %d > %d", m.Last,
			len(file.buf)-1)}
	}

	first, last := m.First, m.Last
	lnmap := make(map[uint64]string)
	for i := first; i < uint64(len(file.buf)) && i < last+1; i++ {
		lnmap[i] = string(file.buf[i])
	}
	l := NewLineRespMessage(lnmap, m.Fid)
	return l, nil
}

func (h *Host) update(m *UpdateMessage) (*UpdateRespMessage, os.Error) {
	f, ok := h.files[m.Fid]
	if !ok {
		return NewUpdateRespMessage(false), &DviError{fmt.Sprintf("Fid %d not in files map")}
	}

	rb := make(map[uint64][]byte)
	max := uint64(len(f.buf) - 1)
	for lno, text := range m.Upd {
		if lno > max {
			// rollback
			return nil, &DviError{fmt.Sprintf("Line out of range: %d > %d", lno, max)}
		}
		rb[lno] = f.buf[lno]
		f.buf[lno] = []byte(text)
	}
	return NewUpdateRespMessage(true), nil
}

func (h *Host) newline(m *NewlineMessage) (*NewlineRespMessage, os.Error) {
	return nil, nil
}

func (h *Host) sync(m *SyncMessage) (*SyncRespMessage, os.Error) {
	file, ok := h.files[m.fid]
	if !ok {
		return nil, &DviError{fmt.Sprintf("Fid %d not in files map")}
	}

	path := file.name
	if len(m.path) > 0 {
		path = m.path
	}

	f, e := os.Create(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	w, e := file.sync(f)
	if e != nil {
		return nil, e
	}
	file.name = path

	return NewSyncRespMessage(w), nil
}

func (h *Host) close(m *CloseMessage) (*CloseRespMessage, os.Error) {
	if f, ok := h.files[m.fid]; ok {
		f.close()
		h.files[m.fid] = nil, false
		c := NewCloseRespMessage(m.fid)
		return c, nil
	}
	return nil, &DviError{fmt.Sprintf("Fid %d not in files map")}
}

// deceptively named -- you actually get a map
func (h *Host) list(m *ListMessage) (*ListRespMessage, os.Error) {
	files := make(map[uint64]string)
	for fid, file := range h.files {
		files[fid] = file.name
	}
	return NewListRespMessage(files), nil
}
