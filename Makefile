all: tidy fmt lint vet install

lint:
	golangci-lint run

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
	cd cmd/kruise && go install && cd -

.PHONY: test
