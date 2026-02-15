VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/cboone/gh-problemas/cmd.version=$(VERSION)
BINARY  := gh-problemas

.PHONY: build test lint vet install clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./...

lint: vet
	@echo "lint: ok (add golangci-lint when ready)"

vet:
	go vet ./...

install: build
	gh extension install .

clean:
	rm -f $(BINARY)
