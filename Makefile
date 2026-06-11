.PHONY: fmt
fmt:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2 fmt

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2 run

.PHONY: test
test:
	go test -v -count 1 ./...
	GOEXPERIMENT=jsonv2 go test -v -count 1 ./...

.PHONY: test-v1
test-v1:
	go test -v ./...

.PHONY: test-v2
test-v2:
	GOEXPERIMENT=jsonv2 go test -v ./...
