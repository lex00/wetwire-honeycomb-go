---
title: "Developers"
---

This guide covers development setup, code organization, and contribution guidelines for wetwire-honeycomb-go.

---

## Development Setup

### Prerequisites

- **Go 1.23+** - Required for the module
- **Git** - Version control
- **Make** (optional) - For running common tasks

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/lex00/wetwire-honeycomb-go.git
cd wetwire-honeycomb-go

# Build the CLI
go build ./cmd/wetwire-honeycomb

# Install locally
go install ./cmd/wetwire-honeycomb
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/lint/...
go test -v ./internal/discover/...

# Run tests with race detection
go test -race ./...
```

### Verify Installation

```bash
# Check version
./wetwire-honeycomb version

# List queries in examples
./wetwire-honeycomb list ./examples/...

# Lint examples
./wetwire-honeycomb lint ./examples/...
```

---

## Code Organization

```
wetwire-honeycomb-go/
├── cmd/wetwire-honeycomb/     # CLI entry point and commands
│
├── internal/                  # Private implementation packages
│   ├── discovery/             # AST-based resource discovery
│   ├── serialize/             # JSON serialization
│   ├── lint/                  # Lint rules and engine
│   ├── builder/               # Build orchestration
│   └── agent/                 # AI agent domain types
│
├── query/                     # Public query types
├── board/                     # Public board types
├── slo/                       # Public SLO types
├── trigger/                   # Public trigger types
│
├── examples/                  # Example declarations
├── testdata/                  # Test fixtures
└── docs/                      # Documentation
```

---

## Adding New Features

### Adding a New Resource Type

When adding support for a new Honeycomb resource type:

1. **Create Public Types** - New package in repository root
2. **Add Discovery Support** - `internal/discover/`
3. **Add Serialization** - `internal/serialize/`
4. **Add Lint Rules** - `internal/lint/`
5. **Update Discovery All** - Add to `DiscoveredResources`
6. **Update CLI Commands** - Handle in `build`, `lint`, `list`
7. **Add Tests** - Unit, integration, and E2E tests

### Adding New Lint Rules

#### 1. Choose a Rule Code

Rule codes follow the format `WHC<NNN>`:
- WHC001-WHC019: Query rules
- WHC020-WHC029: Style enforcement rules
- WHC030-WHC039: Board rules
- WHC040-WHC049: SLO rules
- WHC050-WHC059: Trigger rules

#### 2. Implement the Rule

```go
// internal/lint/rules.go

func WHC015NewRule() Rule {
    return Rule{
        Code:     "WHC015",
        Severity: "warning",
        Message:  "Brief description of the issue",
        Check: func(query discovery.DiscoveredQuery) []LintResult {
            // Implement check logic
            if conditionViolated {
                return []LintResult{{
                    Rule:     "WHC015",
                    Severity: "warning",
                    Message:  "Detailed message with context",
                    File:     query.File,
                    Line:     query.Line,
                    Query:    query.Name,
                }}
            }
            return nil
        },
    }
}
```

#### 3. Register the Rule

Add to `AllRules()` function.

#### 4. Add Tests and Documentation

---

## Code Style and Conventions

### Go Style

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names

### Package Design

- **Public packages** (`query/`, `board/`, etc.): User-facing types
- **Internal packages** (`internal/`): Implementation details
- Keep packages focused on a single responsibility

### Error Handling

```go
// Return wrapped errors with context
if err != nil {
    return nil, fmt.Errorf("failed to parse file: %w", err)
}
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Rule codes | WHC + 3 digits | `WHC001`, `WHC040` |
| Rule functions | `WHC<NNN><Description>` | `WHC001MissingDataset()` |
| Test functions | `Test<Package>_<Function>_<Case>` | `TestLintQueries_WHC001_MissingDataset` |
| Discovery types | `Discovered<Resource>` | `DiscoveredQuery`, `DiscoveredSLO` |

---

## PR Guidelines

### Branch Naming

```
feat/add-derived-column-support
fix/whc001-false-positive
docs/update-lint-rules
refactor/discovery-performance
```

### Commit Message Format

Follow conventional commits:

```
feat(lint): add WHC015 rule for excessive breakdowns

- Check for queries with more than 10 breakdowns
- Add warning with suggestion to reduce cardinality
- Document rule in LINT_RULES.md

Closes #123
```

| Prefix | Description |
|--------|-------------|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation only |
| `refactor:` | Code change that neither fixes a bug nor adds a feature |
| `test:` | Adding or updating tests |
| `chore:` | Maintenance tasks |

### CI Requirements

All PRs must pass:

1. **Tests**: `go test -v ./...`
2. **Race detection**: `go test -race ./...`
3. **Lint**: `golangci-lint run`
4. **Build**: `go build ./...`

### PR Checklist

- [ ] Tests added/updated for new functionality
- [ ] Documentation updated if needed
- [ ] Lint rules documented in `LINT_RULES.md`
- [ ] No breaking changes to public API (or documented in PR)
- [ ] Commit messages follow convention

---

## Release Process

### Versioning

This project uses [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes to public API
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Creating a Release

1. Update version in `cmd/wetwire-honeycomb/main.go`
2. Update CHANGELOG (if exists)
3. Create PR and merge after review
4. Tag the release: `git tag -a v0.4.0 -m "Release v0.4.0"`
5. Push tags: `git push origin v0.4.0`

---

## Troubleshooting

### Tests fail with "failed to access directory"

Ensure tests use `t.TempDir()` or absolute paths.

### Discovery doesn't find queries

Check that:
1. Queries are top-level `var` declarations
2. Query names are exported (start with uppercase)
3. The file is a `.go` file (not `_test.go`)
4. The composite literal uses `query.Query` type

### Lint rules don't trigger

Verify:
1. The rule is registered in `AllRules()`
2. The rule's Check function returns results
3. Test data matches the rule's conditions

---

## See Also

- [CLI Reference](../cli/) - Command documentation
- [Lint Rules](../lint-rules/) - Complete rule reference
- [FAQ](../faq/) - Common questions
