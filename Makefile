.PHONY: all install fmt

all:
	env GO111MODULE=on go build -v .

install:
	env GO111MODULE=on go install -v .

fmt:
	pandoc -o tmp.md -s README.md
	mv tmp.md README.md
	pandoc -o tmp.md -s doc/server.md
	mv tmp.md doc/server.md
	pandoc -o tmp.md -s doc/notes.md
	mv tmp.md doc/notes.md
