.PHONY: fmt
fmt:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 fmt
	GOEXPERIMENT=jsonv2 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 fmt

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run
	GOEXPERIMENT=jsonv2 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run

.PHONY: test
test:
	go test -v -count 1 ./...
	GOEXPERIMENT=jsonv2 go test -v -count 1 ./...
