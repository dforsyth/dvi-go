include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   d
GOFILES=    d.go \
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
