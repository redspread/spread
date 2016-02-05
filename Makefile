.PHONY: build test get

TESTOPTS = -timeout 3m
GETOPTS = -v -d -t

get:
	go get $(GETOPTS) ./cli
	go get $(GETOPTS) ./pkg/...
	go get $(GETOPTS) ./cmd/spread

test:
	go test ./cli
	go test ./pkg/...

build:
	go build -v -o ./bin/spread ./cmd/spread
