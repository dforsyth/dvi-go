include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   d
GOFILES=    d.go \
            file.go \
            gapbuffer.go \
            line.go \
            disk.go \
			screen.go \
		    normal.go \
            insert.go \
			ex.go

include $(GOROOT)/src/Make.cmd
