include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   dvi
GOFILES=    d.go \
            dirbuffer.go \
            disk.go \
            gapbuffer.go \
            global.go \
            editbuffer.go \
            editline.go \
            ex.go \
            insert.go \
            normal.go \
            window.go

include $(GOROOT)/src/Make.cmd
