.PHONY: build clean install dev dev-verbose

OS := $(shell go env GOOS)
ifeq ($(OS),windows)
	OUTPUT := bin/skooma.exe
	NULL_DEVICE := nul
else
	OUTPUT := bin/skooma
	NULL_DEVICE := /dev/null
endif

# Build the binary to the bin/ directory
build:
	go build -o $(OUTPUT)

# Install the binary to GOPATH/bin (globally available)
install:
	go install .

# Watch .go files and rebuild + install on changes
dev:
	@watchexec --version > $(NULL_DEVICE) 2>&1 || (echo Please install watchexec first: https://github.com/watchexec/watchexec && exit 1)
	@echo Watching .go files for changes. Press Ctrl+C to stop.
	watchexec -e go -r "make --no-print-directory build && go install ."

# Same as dev but with verbose output showing which files changed
dev-verbose:
	@watchexec --version > $(NULL_DEVICE) 2>&1 || (echo Please install watchexec first: https://github.com/watchexec/watchexec && exit 1)
	@echo "Watching .go files for changes (verbose). Press Ctrl+C to stop."
	watchexec -e go -r --print-events "make build && go install ."

# Remove build artifacts
clean:
ifeq ($(OS),windows)
	@if exist bin rmdir /s /q bin
else
	rm -rf bin
endif

.DEFAULT_GOAL := build
