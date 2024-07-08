# VN Workday - Account

[//]: # (This project is a template for creating a new service in the VN Workday system. It provides a starting point for)
[//]: # (creating a new service that follows the best practices and conventions used widely across VN Workday microservices.)

[//]: # (- **Service Structure**: The template provides a basic structure for organizing the service code, including directories)
[//]: # (  for controllers, services, models, and routes.)
[//]: # (- **Configuration**: The template includes a configuration loader that loads configuration from environment variables)
[//]: # (  and a configuration file.)
[//]: # (- **Logging**: The template includes a logging package that provides structured logging with context and log levels.)
[//]: # (- **Instrumentation**: The template includes a request tracing middleware that adds trace and span IDs to the request)
[//]: # (  context.)
[//]: # (- **Testing**: The template includes a testing package that provides utilities for testing controllers and services.)
[//]: # (- **Dockerfile**: The template includes a Dockerfile that builds a Docker image for the service.)
[//]: # (- **Makefile**: The template includes a Makefile that provides commands for building, testing, and running the service.)
[//]: # (- **CI/CD**: The template includes GitHub Actions workflows for building, testing, and deploying the service.)
[//]: # (- **Documentation**: The template includes a README template that provides a starting point for documenting the service.)
[//]: # (- **License**: The template includes a license file that specifies the license under which the service is distributed.)
[//]: # (- **Contributing Guidelines**: The template includes a CONTRIBUTING file that specifies the guidelines for contributing)
[//]: # (  to the service.)

## Project Structure

The project structure follows the standard layout:

```
.
├── cmd
│   └── app
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── controller
│   │   └── controller.go
│   ├── model
│   │   └── model.go
│   ├── repository
│   │   └── repository.go
│   ├── router
│   │   └── router.go
│   ├── service
│   │   └── service.go
│   └── util
│       └── util.go
├── pkg
│   └── pkg.go
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Pre-requisites Installation

- [x] Install [Go 1.22+](https://golang.org/doc/install)
- [x] Install [Docker](https://docs.docker.com/get-docker/)
- [x] Install [Node.js](https://nodejs.org/en/download/) from the official website or using a package manager like `nvm`
- [x] Install [golangci-lint](https://golangci-lint.run/welcome/install/)

## Getting Started

1. Run the following commands to install the dependencies and generate the required files:

   ```bash
   npm install
   npm prepare
   ```

2. Start the service:

   ```bash
   go run main.go
   ```

## Coding Conventions

- **MUST** - Make sure you have successfully run `make lint` before committing your code. This will ensure that your code follows the
  coding conventions and best practices.
