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
	cd cmd && go build -o /usr/local/bin/kruise && cd -

.PHONY: test
