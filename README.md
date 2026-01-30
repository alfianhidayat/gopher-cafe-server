# ☕ gopher-cafe

## Features

* **Clean Architecture Implementation**: Strict separation between Transport (gRPC), Usecase (Logic), and Entity (Domain) layers.
* **Concurrency-Ready**: Designed for Barista worker pools and Equipment semaphores.
* **Built-in Reflection**: Self-documenting gRPC server compatible with Postman and Evans CLI.
* **Production Tooling**: Includes automated linting, pre-commit hooks, and code formatting.

## Requirements

* **Go**: `1.24.0`
* **Docker**
* **Make**

## Setup

### **1. Initial Development Setup**

This will install `goimports`, `golangci-lint` to the local bin, and configure your git pre-commit hooks:

```sh
make setup-dev

```

### **2. Install Dependencies**

```sh
go mod tidy

```

### **3. Setup Environment Variables**

```sh
cp config/.env.example config/.env

```

## Build & Run

### **Generate Mocks**

Ensure the mocks are updated before running tests or the server:

```sh
make genmock

```

### **Run the Server**

Starts the server located at `cmd/grpc/main.go`:

```sh
make run

```

## Testing & Quality Control

### **Run Tests**

Standard test execution:

```sh
make test

```

### **Test with Coverage**

Runs tests and opens a browser window with the HTML coverage report:

```sh
make test_coverage

```

### **Linting & Formatting**

Check for code smells and fix formatting/imports:

```sh
make lint
make fix

```

## Notes

* **End-to-End Design**: This boilerplate is functional out-of-the-box. It includes gRPC reflection and table-driven tests.
* **Architectural Freedom**: This is a starter kit. **Modify it however you see fit.** If you find a better structure for the simulation logic, refactor as needed.
* **Strict Linting**: The `Makefile` ignores generated files (`.pb.go`, `mock_*.go`) during formatting and imports to ensure your custom business logic remains clean without touching auto-generated code.
* **Tooling Isolation**: Tools like `golangci-lint` are installed into `./var/bin/` to avoid polluting your global environment.
