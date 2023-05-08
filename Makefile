# Makefile for building and packaging binaries

APP_NAME = cyber
VERSION = v1.0.0

.PHONY: build-windows build-darwin build-linux package clean

all: build-windows build-darwin build-linux package

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/$(APP_NAME)-windows-amd64.exe

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o build/$(APP_NAME)-darwin-amd64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/$(APP_NAME)-linux-amd64

package:
	cd build && tar czvf $(APP_NAME)-$(VERSION)-windows-amd64.tar.gz $(APP_NAME)-windows-amd64.exe
	cd build && tar czvf $(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz $(APP_NAME)-darwin-amd64
	cd build && tar czvf $(APP_NAME)-$(VERSION)-linux-amd64.tar.gz $(APP_NAME)-linux-amd64

clean:
	rm -rf build/*
