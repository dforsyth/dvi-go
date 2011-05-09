include $(GOROOT)/src/Make.inc

.DEFAULT_GOAL=  all

TARG=   dvi
GOFILES=    host.go client.go file.go main.go message.go hostcmd.go terminal.go

include $(GOROOT)/src/Make.cmd
