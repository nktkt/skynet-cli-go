.PHONY: build build-skynet build-history test run-skynet run-history fmt

build: build-skynet build-history

build-skynet:
	go build -o skynet .

build-history:
	go build -o codex-history-cli ./cmd/codex-history-cli

test:
	go test ./...

fmt:
	gofmt -w main.go cmd/codex-history-cli/*.go internal/skynet/*.go internal/history/*.go

run-skynet:
	go run . $(ARGS)

run-history:
	go run ./cmd/codex-history-cli $(ARGS)
