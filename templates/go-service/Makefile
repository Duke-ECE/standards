# Thin convenience wrappers only — the same commands CI runs.
# GOTOOLCHAIN=local guards against dev machines with newer Go toolchains
# silently bumping go.mod past what the golang:1.25 Docker image accepts.

export GOTOOLCHAIN=local

build:
	go build ./...

vet:
	go vet ./...

test:
	go test ./...

run:
	go run ./cmd/server

check:
	./scripts/check.sh

tidy:
	go mod tidy
	@grep '^go ' go.mod

.PHONY: build vet test run check tidy
