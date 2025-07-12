.PHONY: fmt
fmt:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2 fmt

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2 run

.PHONY: test
test:
	go test -v ./...
