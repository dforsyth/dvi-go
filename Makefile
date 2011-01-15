include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   d
GOFILES=    d.go \
            buffer.go \
            gapbuffer.go \
            file.go \
			view.go

include $(GOROOT)/src/Make.cmd
