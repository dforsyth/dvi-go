include $(GOROOT)/src/Make.inc
.DEFAULT_GOAL=	all
TARG=	dvi
GOFILES=	dvi.go draw.go vicmds.go position.go buffer.go config.go extend.go excmds.go

include $(GOROOT)/src/Make.cmd
