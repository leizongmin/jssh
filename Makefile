OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

.PHONY: all
all: jssh

.PHONY: jssh
jssh:
	@./build.js
	@ls -alh release/*

.PHONY: jssh_m
jssh_m: jssh
	@cp release/$(OS)/jssh release/$(OS)/jssh_m
	@upx --best --lzma release/$(OS)/jssh_m

.PHONY: go-nm
go-nm:
	@go build -o release/jssh .
	@go tool nm -size -sort size release/jssh 2>/dev/null | head -50
