deps:
    go mod tidy

test: deps
    go test -v ./...

install: deps
    go install ./cmd/health
