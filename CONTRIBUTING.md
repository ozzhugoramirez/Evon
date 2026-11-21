# Contributing to EVON

Thank you for your interest in contributing to EVON.

EVON is an experimental programming language project written in Go. Its architecture and core concepts are still evolving.

## Development Guidelines

- Keep the code simple and easy to understand.
- Maintain clear separation of responsibilities between packages.
- Follow standard Go formatting conventions using `gofmt`.
- Handle errors explicitly.
- Avoid unnecessary external dependencies.

## Project Structure

```text
cmd/
  Evon/          # Application entry point

internal/
  ai/            # AI-related functionality
  cli/           # Command-line interface
  compiler/      # Language processing and evaluation
  memory/        # Memory and data storage
```

The project structure may change as EVON develops.

## Contributing Changes

1. Fork the repository.
2. Create a descriptive branch.
3. Make your changes.
4. Run the tests:

   ```bash
   go test ./...
   ```

5. Format your code:

   ```bash
   gofmt -w .
   ```

6. Submit a pull request describing your changes.

Contributions should focus on clarity, consistency, and the long-term direction of EVON.
