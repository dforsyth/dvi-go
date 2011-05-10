package main

import (
	"fmt"
	"os"
)

func (h *Host) open(m *OpenMessage) (*OpenRespMessage, os.Error) {
	path := m.pathname
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}

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
	file, ok := h.files[m.fid]
	if !ok {
		return nil, &DviError{fmt.Sprintf("Fid %d not in files map")}
	}
	if m.lno > uint64(len(file.buf)-1) {
		return nil, &DviError{fmt.Sprintf("Line %d not in Fid %d", m.lno, m.fid)}
	}
	text := string(file.buf[m.lno])
	l := NewLineRespMessage(text, m.fid, m.lno)
	return l, nil
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
