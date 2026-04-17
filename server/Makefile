# Build the application
all: build

build:
	@echo "Building..."
	@go build -v -o bbs-go main.go

build-tkdev:
	@echo "Building tkdev..."
	@go build -v -o tkdev main.go

buildlinux:
	@echo "Building..."
	@GOOS=linux GOARCH=amd64 go build

build-tkdev-linux:
	@echo "Building tkdev for Linux..."
	@GOOS=linux GOARCH=amd64 go build -v -o tkdev main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bbs-go

generator:
	@go run cmd/generator/generator.go

.PHONY: all build run test clean
