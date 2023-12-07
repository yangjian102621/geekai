SHELL=/usr/bin/env bash
NAME := chatgpt-plus
all: amd64 arm64

amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux main.go
.PHONY: amd64

arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=7 go build -o bin/$(NAME)-linux main.go
.PHONY: arm64

clean:
	rm -rf bin/$(NAME)-*
.PHONY: clean
