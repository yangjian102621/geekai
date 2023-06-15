SHELL=/usr/bin/env bash
NAME := chatgpt-v3
all: window linux darwin


window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-amd64.exe main.go
.PHONY: window

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-amd64-linux main.go
.PHONY: linux

darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-amd64-darwin main.go
.PHONY: darwin

clean:
	rm -rf bin/$(NAME)-*
.PHONY: clean
