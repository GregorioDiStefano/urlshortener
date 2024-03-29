.PHONY: lint

export PATH := bin:$(PATH):/home/greg/go/bin/

build:
	CGO_ENABLED=0 go build .

lint:
	goimports -w .
	revive .
	golangci-lint run 

test:
	GORACE="halt_on_error=1" go test -race -count 1 -v ./...

coverage:
	go test $(go list ./... | grep -v mocks ) -v -coverprofile cover.out
	go tool cover -html cover.out -o cover.html

mock:
	mockery --all

