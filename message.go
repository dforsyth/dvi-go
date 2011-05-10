package main

import (
	"fmt"
	"os"
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
	fid  uint64
	lno  uint64
	text string
}

func NewLineRespMessage(text string, fid, lno uint64) *LineRespMessage {
	return &LineRespMessage{
		fid,
		lno,
		text,
	}
}

func (m *LineRespMessage) message() string {
	return fmt.Sprintf("FID: %d: LINE: %d: %s", m.fid, m.lno, m.text)
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

// commands from client

type InsertMessage struct {
	text     string
	line     uint64
	position uint64
}

func (m *InsertMessage) message() string {
	return fmt.Sprintf("INSERT %s @ LINE %d POSITION %d", m.text, m.line, m.position)
}

type OpenMessage struct {
	pathname string
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
	fid uint64
	lno uint64
}

func (m *LineMessage) message() string {
	return fmt.Sprintf("LINE: FID: %d LNO: %d", m.fid, m.lno)
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
