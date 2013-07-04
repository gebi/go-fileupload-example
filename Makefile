all: build

build:
	go build

test:
	go test

format:
	gofmt -s -w=true *.go
