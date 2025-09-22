.PHONY: check
check: build test lint

.PHONY: build
build: generate migrate

.PHONY: generate
generate:
	go build -o bin/generate ./cmd/generate

.PHONY: migrate
migrate:
	go build -o bin/migrate ./cmd/migrate

.PHONY: install
install: build
	cp bin/generate /usr/local/bin/
	cp bin/migrate /usr/local/bin/

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf internal/generated/

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: test-clean
test-clean: clean test

.PHONY: demo
demo: build
	./bin/generate test-api.yaml
	./bin/migrate -status test-api.yaml
