include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   d
GOFILES=    d.go \
            gapbuffer.go \
            editline.go \
            editbuffer.go \
            disk.go \
			screen.go \
		    normal.go \
            insert.go \
			ex.go

include $(GOROOT)/src/Make.cmd
