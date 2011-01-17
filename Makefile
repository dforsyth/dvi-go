include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   d
GOFILES=    d.go \
            buffer.go \
            gapbuffer.go \
            file.go \
			view.go \
		    normal.go \
            insert.go

include $(GOROOT)/src/Make.cmd
