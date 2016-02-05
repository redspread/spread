.PHONY: build test

TESTOPTS = "-timeout 3m"

test:
	go test ./cli
	go test ./pkg/...

build:
	go build -v -o ./bin/spread ./cmd/spread
