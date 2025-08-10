# Project: logwild

## Overview

`logwild` is a program that runs an http server exposing an endpoint
`/loggen` which responds to a `GET` request by emitting a burst of log lines to
a specified file. The primary use case is for testing the behavior of log forwarders,
`logrotate` and other programs that might need to respond to a text file
that is constantly growing.

## Project Structure

- `cmd/`: Contains the main application entry point
- `pkg/`: Contains all functionality and subcommands used by the main entrypoint.

## Coding Standards

- **Formatting**: All code must be formatted using `gofmt`.
- **Linting**: Use `golint` and `go vet` to ensure code quality.
- **Logging**: The project uses `structlog` for logging and prefers rich log messages.
- **Error Handling**:
  - Always check for errors returned by functions.
  - Provide context when returning errors using `fmt.Errorf("context: %w", err)`.
- **Naming Conventions**:
  - Use `PascalCase` for exported identifiers.
  - Use `camelCase` for unexported identifiers.
  - Interface names should be descriptive (e.g., `Reader`, not `IReader`)

## Dependencies

- [cobra](https://github.com/spf13/cobra): Used for building the CLI interface

## Build Instructions

- see the `justfile` for building type stuff
