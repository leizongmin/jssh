.PHONY: all
all: jssh

.PHONY: jssh
jssh:
	@./build.js
	@ls -alh release/*

.PHONY: go-nm
go-nm:
	@go build -o release/jssh .
	@go tool nm -size -sort size release/jssh 2>/dev/null | head -50
