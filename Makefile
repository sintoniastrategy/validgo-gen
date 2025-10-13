.PHONY: check
check: generate test lint

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test ./...
