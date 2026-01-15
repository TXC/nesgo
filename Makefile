# This Makefile exists to help make a Mac OS X app bundle.  If you
# want to build the regular nesgo binary, just use 'go build'.

all: nesgo Nesgo.app

Nesgo.app:
	cp nesgo Nesgo.app/Contents/MacOS

nesgo:
	go generate ./...
	go build

.PHONY: nesgo Nesgo.app
