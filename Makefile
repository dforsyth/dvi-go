include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   dvi
GOFILES=    d.go \
            disk.go \
            gapbuffer.go \
            global.go \
            editbuffer.go \
            editline.go \
            ex.go \
            input.go \
            normal.go \
            window.go \
	    util.go

include $(GOROOT)/src/Make.cmd
