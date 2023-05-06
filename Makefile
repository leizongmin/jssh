RELEASE_TARGET := $(shell go run . eval "__os + '-' + __arch" 2>/dev/null)

.PHONY: all
all: jssh jssh_m

.PHONY: info
info:
	@echo "RELEASE_TARGET: $(RELEASE_TARGET)"

.PHONY: clean
clean:
	@rm -rf release node_modules

.PHONY: jssh
jssh:
	@./build.js
	@ln -sf $(RELEASE_TARGET)/jssh release/jssh
	@ls -alh release/*

.PHONY: jssh_m
jssh_m: jssh
	@cp release/$(RELEASE_TARGET)/jssh release/$(RELEASE_TARGET)/jssh_m
	@upx --best --lzma release/$(RELEASE_TARGET)/jssh_m
	@ln -sf $(RELEASE_TARGET)/jssh_m release/jssh_m
	@ls -alh release/*

.PHONY: go-nm
go-nm:
	@go build -o release/jssh .
	@go tool nm -size -sort size -type release/jssh 2>/dev/null | head -500
