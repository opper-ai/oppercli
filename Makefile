# Build the project
build:
	go build -o bin/opper cmd/opper/main.go

# Install dependencies and the project
install:
	cd cmd/opper && go install

# Clean build files
clean:
	go clean
	rm -f bin/opper

.PHONY: build install clean