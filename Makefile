all: update-deps vet fmt

fmt:
	gofmt -w -s .

fmt-dry-run:
	gofmt -w -d .

vet:
	go vet ./...

update-deps:
	go mod tidy
