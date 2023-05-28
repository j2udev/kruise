path := /usr/local/bin

all: tidy fmt lint vet test install

lint:
	golangci-lint run ./...

fmt:
	gofmt -w -s .

fmtd:
	gofmt -w -d .

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./... -v

install:
	go build -o $(path)/kruise

.PHONY: test
