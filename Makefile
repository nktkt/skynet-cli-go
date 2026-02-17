.PHONY: build test fmt run

build:
	go build -o skynet .

test:
	go test ./...

fmt:
	gofmt -w main.go internal/skynet/*.go

run:
	go run . $(ARGS)
