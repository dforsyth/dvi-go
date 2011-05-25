package main

import (
	"fmt"
	"os"
	"strings"
)

type Message interface {
	message() string
}

// responses from host

type ErrorRespMessage struct {
	m string
}

func (m *ErrorRespMessage) message() string {
	return m.m
}

type OpenRespMessage struct {
	Fid      uint64
	Pathname string
}

func NewOpenRespMessage(f *File) *OpenRespMessage {
	m := new(OpenRespMessage)
	m.Fid = f.fid
	m.Pathname = f.name
	return m
}

func (m *OpenRespMessage) message() string {
	return fmt.Sprintf("OPEN %s : %d", m.Pathname, m.Fid)
}

type LineRespMessage struct {
	Fid   uint64
	Lnmap map[uint64]string
}

func NewLineRespMessage(lnmap map[uint64]string, fid uint64) *LineRespMessage {
	return &LineRespMessage{
		fid,
		lnmap,
	}
}

func (m *LineRespMessage) message() string {
	rval := fmt.Sprintf("FID: %d", m.Fid)
	for k, v := range m.Lnmap {
		// have to do this out here because the compiler thinks k and v aren't used if we do
		// this inside of the join call
		arr := []string{rval, fmt.Sprintf("LINE: %d: %s", k, v)}
		rval = strings.Join(arr, "\n")
	}
	return rval
}

type StatRespMessage struct {
	name  string
	lines uint64
	dirty bool
}

func NewStatRespMessage(file *File, fi *os.FileInfo) *StatRespMessage {
	s := new(StatRespMessage)
	s.name = file.name
	s.lines = uint64(len(file.buf))
	s.dirty = file.dirty
	return s
}

func (m *StatRespMessage) message() string {
	return fmt.Sprintf("NAME: %s: LINES: %d: DIRTY: %t", m.name, m.lines, m.dirty)
}

type CloseRespMessage struct {
	fid uint64
}

func (m *CloseRespMessage) message() string {
	return fmt.Sprintf("CLOSED: FID: %d", m.fid)
}

func NewCloseRespMessage(fid uint64) *CloseRespMessage {
	c := new(CloseRespMessage)
	c.fid = fid
	return c
}

type ListRespMessage struct {
	files map[uint64]string
}

func (m *ListRespMessage) message() string {
	r := "LIST: "
	for fid, name := range m.files {
		r += fmt.Sprintf("(%d:%s) ", fid, name)
	}
	return r
}

func NewListRespMessage(files map[uint64]string) *ListRespMessage {
	l := new(ListRespMessage)
	l.files = files
	return l
}

type UpdateRespMessage struct {
	Success bool
}

func (m *UpdateRespMessage) message() string {
	return fmt.Sprintf("UPDATE: SUCCESS: %T", m.Success)
}

func NewUpdateRespMessage(success bool) *UpdateRespMessage {
	u := &UpdateRespMessage {
		Success: success,
	}

	return u
}

type NewlineRespMessage struct {

}

type SyncRespMessage struct {
	w uint64
}

func (m *SyncRespMessage) message() string {
	return ""
}

func NewSyncRespMessage(w uint64) *SyncRespMessage {
	return &SyncRespMessage{w}
}

// commands from client

type UpdateMessage struct {
	Fid uint64
	Upd map[uint64]string
}

func (m *UpdateMessage) message() string {
	r := fmt.Sprintf("UPDATE: FID: %d:", m.Fid)
	for k, v := range m.Upd {
		arr := []string{r, fmt.Sprintf("LNO: %d: %s", k, v)}
		r = strings.Join(arr, "\n")
	}
	return r
}

type OpenMessage struct {
	Pathname string
	Force    bool
}

func (m *OpenMessage) message() string {
	return fmt.Sprintf("OPEN: %s", m.Pathname)
}

type StatMessage struct {
	fid uint64
}

func (m *StatMessage) message() string {
	return fmt.Sprintf("STAT: %d", m.fid)
}

type LineMessage struct {
	Fid         uint64
	First, Last uint64
}

func (m *LineMessage) message() string {
	return fmt.Sprintf("LINE: FID: %d START: %d: FINISH: %d", m.Fid, m.First, m.Last)
}

type CloseMessage struct {
	fid  uint64
	sync bool
}

func (m *CloseMessage) message() string {
	return fmt.Sprintf("CLOSE: FID: %d", m.fid)
}

type ListMessage struct {

}

func (m *ListMessage) message() string {
	return "LIST"
}

type NewlineMessage struct {

}

type SyncMessage struct {
	fid  uint64
	path string
}

func (m *SyncMessage) message() string {
	return ""
}
