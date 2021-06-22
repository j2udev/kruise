all: update-deps vet fmt install

fmt:
	gofmt -w -s .

fmt-dry-run:
	gofmt -w -d .

vet:
	go vet ./...

update-deps:
	go mod tidy

install:
	go install ./cmd/kruise/main.go
