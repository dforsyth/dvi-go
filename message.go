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
	fid      uint64
	pathname string
}

func NewOpenRespMessage(f *File) *OpenRespMessage {
	m := new(OpenRespMessage)
	m.fid = f.fid
	m.pathname = f.name
	return m
}

func (m *OpenRespMessage) message() string {
	return fmt.Sprintf("OPEN %s : %d", m.pathname, m.fid)
}

type LineRespMessage struct {
	fid   uint64
	lnmap map[uint64]string
}

func NewLineRespMessage(lnmap map[uint64]string, fid uint64) *LineRespMessage {
	return &LineRespMessage{
		fid,
		lnmap,
	}
}

func (m *LineRespMessage) message() string {
	rval := fmt.Sprintf("FID: %d", m.fid)
	for k, v := range m.lnmap {
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

}

func (m *UpdateRespMessage) message() string {
	return ""
}

func NewUpdateRespMessage() *UpdateRespMessage {
	return new(UpdateRespMessage)
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
	fid uint64
	upd map[uint64]string
}

func (m *UpdateMessage) message() string {
	r := "UPDATE: FID: %d:"
	for k, v := range m.upd {
		arr := []string{r, fmt.Sprintf("LNO: %d: %s", k, v)}
		r = strings.Join(arr, "\n")
	}
	return r
}

type OpenMessage struct {
	pathname string
	force    bool
}

func (m *OpenMessage) message() string {
	return fmt.Sprintf("OPEN: %s", m.pathname)
}

type StatMessage struct {
	fid uint64
}

func (m *StatMessage) message() string {
	return fmt.Sprintf("STAT: %d", m.fid)
}

type LineMessage struct {
	fid         uint64
	first, last uint64
}

func (m *LineMessage) message() string {
	return fmt.Sprintf("LINE: FID: %d START: %d: FINISH: %d", m.fid, m.first, m.last)
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
