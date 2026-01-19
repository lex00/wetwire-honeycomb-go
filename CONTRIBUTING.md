# Contributing to wetwire-honeycomb-go

Thank you for your interest in contributing to wetwire-honeycomb-go!

## Getting Started

See the [Developer Guide](docs/DEVELOPERS.md) for:
- Development environment setup
- Project structure
- Running tests

## Code Style

- **Formatting**: Use `gofmt` (automatic in most editors)
- **Linting**: Use `go vet` and `golangci-lint`
- **Imports**: Use `goimports` for automatic import management

```bash
# Format code
gofmt -w .

# Lint
go vet ./...
golangci-lint run ./...

# Check for common issues
go build ./...
```

## Commit Messages

Follow conventional commits:

```
feat(query): Add support for HEATMAP calculation
fix(lint): Correct WHC005 false positive
docs: Update calculation examples
test: Add tests for board discovery
chore: Update dependencies
```

## Pull Request Process

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes with tests
3. Run tests: `go test ./...`
4. Run linter: `golangci-lint run ./...`
5. Commit with clear messages
6. Push and open Pull Request
7. Address review comments

## Adding a New Lint Rule

1. Add rule to appropriate file in `internal/lint/`
   - `rules.go` for query rules (WHC001-WHC019)
   - `board_rules.go` for board rules (WHC030+)
   - `slo_rules.go` for SLO rules (WHC040+)
   - `trigger_rules.go` for trigger rules (WHC050+)
2. Implement the check function
3. Add test case in corresponding `*_test.go` file
4. Update docs/LINT_RULES.md with the new rule

## Adding a New Resource Type

See the [Developer Guide](docs/DEVELOPERS.md#adding-new-features) for detailed instructions on:
- Creating public types
- Adding discovery support
- Adding serialization
- Adding lint rules
- Updating CLI commands

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include reproduction steps for bugs
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
