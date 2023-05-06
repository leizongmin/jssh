RELEASE_TARGET := $(shell go run . eval "__os + '-' + __arch" 2>/dev/null)
JSBUILTIN_FILES := $(shell cd internal/jsbuiltin/src && find . -type f -name '*.js' -print0 | xargs -0)

.PHONY: all
all: jssh jssh_m

.PHONY: info
info:
	@echo "RELEASE_TARGET: $(RELEASE_TARGET)"
	@echo "JSBUILTIN_FILES: $(JSBUILTIN_FILES)"

.PHONY: clean
clean:
	@rm -rf release node_modules

.PHONY: jssh
jssh: jsbuiltin
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

.PHONY: jsbuiltin
jsbuiltin:
	@echo "Building JS builtin files..."
	@mkdir -p internal/jsbuiltin/dist
	@for f in $(JSBUILTIN_FILES); do \
		echo "  $$f"; \
		npx uglifyjs -c -- internal/jsbuiltin/src/$$f > internal/jsbuiltin/dist/$$f; \
	done
