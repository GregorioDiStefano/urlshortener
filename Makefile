.PHONY: lint

export PATH := bin:$(PATH):/home/greg/go/bin/

build:
	CGO_ENABLED=0 go build .

check:
	govulncheck ./...

lint:
	goimports -w .
	revive .
	golangci-lint run 

test:
	GORACE="halt_on_error=1" go test -race -count 1 -v ./...
	./smoke-test.sh

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

mock:
	mockery --all

