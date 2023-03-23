SHELL=/usr/bin/env bash
NAME := wechatGPT
all: window_x86 window_amd64 linux_x86 linux_amd64 mac_x86 mac_64


window_x86:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/$(NAME).exe main.go
.PHONY: window_x86

window_amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-amd64.exe main.go
.PHONY: window_amd64

linux_x86:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o bin/$(NAME)-386-linux main.go
.PHONY: linux_x86

linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-amd64-linux main.go
.PHONY: linux_amd64

mac_x86:
	CGO_ENABLED=1 GOOS=darwin GOARCH=386 go build -o bin/$(NAME)-386-darwin main.go
.PHONY: mac_x86

mac_64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-amd64-darwin main.go
.PHONY: mac_64

clean:
	rm -rf bin/$(NAME)-*
.PHONY: clean
