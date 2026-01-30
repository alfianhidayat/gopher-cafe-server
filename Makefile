GO_FILES_NO_MOCKS_NO_SQL_NO_PB := $(shell find . -type f -iname '*.go' -not -iname "mock_*.go" -not -iname "*.sql.go" -not -iname "*.pb.go")
TOOLS_BIN=./var/bin/
GO_VERSION=1.24.0
LINT_VERSION=v1.64.2
LINT_BIN=${TOOLS_BIN}/golangci-lint

run:
	go run cmd/grpc/main.go

genmock:
	go generate ./...

test:
	go test ./...

test_coverage:
	-go test ./... -covermode=count -coverprofile=coverage.out ; go tool cover -html=coverage.out
	-rm coverage.out

test_verbose:
	go test -v ./...

setup-dev: install_imports install-golangci-lint
	cp ./pre-commit .git/hooks/pre-commit && chmod u+x .git/hooks/pre-commit

install_imports:
	@echo "Installing Go Imports..."
	go install golang.org/x/tools/cmd/goimports@v0.29.0

install-golangci-lint:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${TOOLS_BIN} ${LINT_VERSION}

fix: fmt imports

fmt:
	@echo "Simplifying the code format..."
	@gofmt -s -w $(GO_FILES_NO_MOCKS_NO_SQL_NO_PB)
	@echo "Code format is simplified."

imports: install_imports
	@echo "Fixing imports..."
	@goimports -w $(GO_FILES_NO_MOCKS_NO_SQL_NO_PB)
	@echo "Imports are fixed."

lint:
	@$(LINT_BIN) run --fix=false ./...
