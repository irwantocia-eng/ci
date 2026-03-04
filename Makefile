.PHONY: run build test lint lint-install clean coverage coverage-html sonar-scan

run:
	go run main.go

build:
	go build -o server main.go

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run

test: lint
	go test ./...

coverage:
	go test ./handlers -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

sonar-scan: coverage
	golangci-lint run --out-format=sarif:golangci-lint.sarif
	bash scripts/sonar-scanner.sh

clean:
	rm -f server koban.db coverage.out coverage.html golangci-lint.sarif
