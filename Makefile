.PHONY: build install clean test

# Binary name
BINARY_NAME=gh-problemas

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Install the extension locally
install: build
	mkdir -p ~/.local/share/gh/extensions/gh-problemas
	cp $(BINARY_NAME) ~/.local/share/gh/extensions/gh-problemas/

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Run tests (if any)
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...
