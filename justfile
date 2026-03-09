deps:
    go mod tidy

test: deps
    go test -v ./...
