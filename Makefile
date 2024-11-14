VERSION ?= 0.1.0
LDFLAGS := -X main.Version=$(VERSION)

.PHONY: build install clean release test test-race test-cover test-all

# Build the project
build:
	go build -o bin/opper cmd/opper/main.go

# Install dependencies and the project
install:
	cd cmd/opper && go install

# Clean build files
clean:
	go clean
	rm -rf bin/ dist/
	mkdir -p bin/ dist/

# Create release builds
release: clean
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/opper-darwin-arm64 ./cmd/opper
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/opper-darwin-amd64 ./cmd/opper
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/opper-linux-amd64 ./cmd/opper
	cd dist && \
		sha256sum opper-darwin-arm64 > opper-darwin-arm64.sha256 && \
		sha256sum opper-darwin-amd64 > opper-darwin-amd64.sha256 && \
		sha256sum opper-linux-amd64 > opper-linux-amd64.sha256

# Add these to your existing Makefile

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race -v ./...

.PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-all
test-all: test test-race test-cover
