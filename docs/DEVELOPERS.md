# Developer Guide

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
go test -v ./internal/discovery/...

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
│   ├── main.go                # Command definitions
│   ├── main_test.go           # E2E tests
│   ├── mcp.go                 # MCP server integration
│   ├── design.go              # AI-assisted design command
│   └── test.go                # Test command
│
├── internal/                  # Private implementation packages
│   ├── discovery/             # AST-based resource discovery
│   │   ├── discovery.go       # Query discovery
│   │   ├── board.go           # Board discovery
│   │   ├── slo.go             # SLO discovery
│   │   ├── trigger.go         # Trigger discovery
│   │   ├── ast.go             # AST utilities
│   │   └── *_test.go          # Unit tests
│   │
│   ├── serialize/             # JSON serialization
│   │   ├── serialize.go       # Query serialization
│   │   ├── board.go           # Board serialization
│   │   ├── slo.go             # SLO serialization
│   │   ├── trigger.go         # Trigger serialization
│   │   └── *_test.go          # Unit tests
│   │
│   ├── lint/                  # Lint rules and engine
│   │   ├── lint.go            # Core lint engine
│   │   ├── rules.go           # Query lint rules (WHC001-WHC023)
│   │   ├── board_rules.go     # Board lint rules (WHC030+)
│   │   ├── slo_rules.go       # SLO lint rules (WHC040+)
│   │   ├── trigger_rules.go   # Trigger lint rules (WHC050+)
│   │   └── *_test.go          # Unit tests
│   │
│   ├── builder/               # Build orchestration
│   │   ├── builder.go         # Build pipeline
│   │   ├── registry.go        # Resource registry
│   │   └── *_test.go          # Unit tests
│   │
│   ├── roundtrip/             # Round-trip tests
│   │   └── roundtrip_test.go  # JSON -> Go -> JSON validation
│   │
│   └── agent/                 # AI agent domain types
│       └── domain.go
│
├── query/                     # Public query types
│   ├── query.go               # Query struct
│   ├── calculation.go         # Calculation builders
│   ├── filter.go              # Filter builders
│   ├── time.go                # Time range utilities
│   ├── breakdown.go           # Breakdown utilities
│   └── *_test.go              # Unit tests
│
├── board/                     # Public board types
│   ├── board.go               # Board struct
│   ├── panel.go               # Panel types
│   └── *_test.go              # Unit tests
│
├── slo/                       # Public SLO types
│   ├── slo.go                 # SLO struct
│   ├── burn.go                # Burn rate utilities
│   └── *_test.go              # Unit tests
│
├── trigger/                   # Public trigger types
│   ├── trigger.go             # Trigger struct
│   ├── recipient.go           # Recipient types
│   └── *_test.go              # Unit tests
│
├── examples/                  # Example declarations
│   ├── latency/               # Latency queries
│   ├── errors/                # Error tracking
│   ├── slos/                  # SLO examples
│   ├── triggers/              # Trigger examples
│   ├── boards/                # Board examples
│   └── full_stack/            # Complete example
│
├── testdata/                  # Test fixtures
│   ├── queries/               # Query test files
│   └── roundtrip/             # Round-trip fixtures
│
└── docs/                      # Documentation
    ├── CLI.md                 # Command reference
    ├── LINT_RULES.md          # Lint rules
    ├── FAQ.md                 # Common questions
    └── DEVELOPERS.md          # This file
```

---

## Adding New Features

### Adding a New Resource Type

When adding support for a new Honeycomb resource type (e.g., derived columns), follow this pattern:

#### 1. Create Public Types

Create a new package in the repository root for the public API:

```go
// derivedcolumn/derivedcolumn.go
package derivedcolumn

type DerivedColumn struct {
    Name       string
    Dataset    string
    Expression string
    // ...
}
```

#### 2. Add Discovery Support

Create discovery in `internal/discovery/`:

```go
// internal/discovery/derivedcolumn.go
package discovery

type DiscoveredDerivedColumn struct {
    Name    string
    Package string
    File    string
    Line    int
    // Extracted fields...
}

func DiscoverDerivedColumns(dir string) ([]DiscoveredDerivedColumn, error) {
    // Follow pattern from slo.go or trigger.go
}

func findDerivedColumnComposites(expr ast.Expr) []*ast.CompositeLit {
    // Check for derivedcolumn.DerivedColumn type
}
```

#### 3. Add Serialization

Create serialization in `internal/serialize/`:

```go
// internal/serialize/derivedcolumn.go
package serialize

func DerivedColumnToJSON(dc derivedcolumn.DerivedColumn) (json.RawMessage, error) {
    // Convert to Honeycomb JSON format
}
```

#### 4. Add Lint Rules

Create lint rules in `internal/lint/`:

```go
// internal/lint/derivedcolumn_rules.go
package lint

type DerivedColumnRule struct {
    Code     string
    Severity string
    Message  string
    Check    func(dc discovery.DiscoveredDerivedColumn) []LintResult
}

func AllDerivedColumnRules() []DerivedColumnRule {
    return []DerivedColumnRule{
        WHC060DerivedColumnMissingName(),
        WHC061InvalidExpression(),
        // ...
    }
}
```

#### 5. Update Discovery All

Add to `DiscoveredResources` in `internal/discovery/discovery.go`:

```go
type DiscoveredResources struct {
    Queries         []DiscoveredQuery
    SLOs            []DiscoveredSLO
    Triggers        []DiscoveredTrigger
    Boards          []DiscoveredBoard
    DerivedColumns  []DiscoveredDerivedColumn  // Add new type
}
```

#### 6. Update CLI Commands

Update `cmd/wetwire-honeycomb/main.go` to handle the new resource type in `build`, `lint`, and `list` commands.

#### 7. Add Tests

- Unit tests for discovery: `internal/discovery/derivedcolumn_test.go`
- Unit tests for serialization: `internal/serialize/derivedcolumn_test.go`
- Unit tests for lint rules: `internal/lint/derivedcolumn_rules_test.go`
- E2E tests in `cmd/wetwire-honeycomb/main_test.go`

### Adding New Lint Rules

Lint rules follow a consistent pattern. Here's how to add a new rule:

#### 1. Choose a Rule Code

Rule codes follow the format `WHC<NNN>`:
- WHC001-WHC019: Query rules
- WHC020-WHC029: Style enforcement rules
- WHC030-WHC039: Board rules
- WHC040-WHC049: SLO rules
- WHC050-WHC059: Trigger rules

#### 2. Implement the Rule

Add to the appropriate `*_rules.go` file:

```go
// internal/lint/rules.go

// WHC015NewRule checks for [specific condition].
func WHC015NewRule() Rule {
    return Rule{
        Code:     "WHC015",
        Severity: "warning",  // or "error"
        Message:  "Brief description of the issue",
        Check: func(query discovery.DiscoveredQuery) []LintResult {
            // Implement check logic
            if conditionViolated {
                return []LintResult{
                    {
                        Rule:     "WHC015",
                        Severity: "warning",
                        Message:  "Detailed message with context",
                        File:     query.File,
                        Line:     query.Line,
                        Query:    query.Name,
                    },
                }
            }
            return nil
        },
    }
}
```

#### 3. Register the Rule

Add to `AllRules()` function:

```go
func AllRules() []Rule {
    return []Rule{
        WHC001MissingDataset(),
        // ...existing rules...
        WHC015NewRule(),  // Add new rule
    }
}
```

#### 4. Add Tests

Add comprehensive tests covering:
- Positive cases (rule triggers)
- Negative cases (rule does not trigger)
- Edge cases

```go
// internal/lint/lint_test.go

func TestLintQueries_WHC015_NewRule_Triggers(t *testing.T) {
    queries := []discovery.DiscoveredQuery{
        {
            Name:    "TestQuery",
            Package: "test",
            File:    "/test/file.go",
            Line:    10,
            // Set up condition that should trigger the rule
        },
    }

    results := LintQueries(queries)

    if !hasResult(results, "WHC015") {
        t.Error("Expected WHC015 warning for [condition]")
    }

    result := findResult(results, "WHC015")
    if result.Severity != "warning" {
        t.Errorf("Expected warning severity, got %s", result.Severity)
    }
}

func TestLintQueries_WHC015_NewRule_NoTrigger(t *testing.T) {
    queries := []discovery.DiscoveredQuery{
        {
            Name:    "TestQuery",
            // Set up condition that should NOT trigger the rule
        },
    }

    results := LintQueries(queries)

    if hasResult(results, "WHC015") {
        t.Error("Did not expect WHC015 warning for valid query")
    }
}
```

#### 5. Document the Rule

Add documentation to `docs/LINT_RULES.md`:

```markdown
### WHC015: New Rule

**Severity:** warning

**Description:**

[Explain what the rule checks for]

**Why:**

- [Reason 1]
- [Reason 2]

**Bad:**

\`\`\`go
// Example that triggers the rule
\`\`\`

**Good:**

\`\`\`go
// Example that passes the rule
\`\`\`
```

### Adding New CLI Commands

CLI commands are implemented using Cobra in `cmd/wetwire-honeycomb/main.go`:

```go
func newMyCommand() *cobra.Command {
    var flagOption string

    cmd := &cobra.Command{
        Use:   "mycommand [args]",
        Short: "Brief description",
        Long:  `Longer description with examples...`,
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
            return nil
        },
    }

    cmd.Flags().StringVarP(&flagOption, "flag", "f", "default", "Flag description")

    return cmd
}
```

Register the command in `main()`:

```go
rootCmd.AddCommand(
    // ...existing commands...
    newMyCommand(),
)
```

---

## Testing Requirements

### Unit Tests

Every package should have comprehensive unit tests:

```go
// internal/lint/lint_test.go

func TestLintQueries_WHC001_MissingDataset(t *testing.T) {
    queries := []discovery.DiscoveredQuery{
        {
            Name:      "TestQuery",
            Package:   "test",
            File:      "/test/file.go",
            Line:      10,
            Dataset:   "", // Missing dataset
            TimeRange: discovery.TimeRange{TimeRange: 3600},
            Calculations: []discovery.Calculation{
                {Op: "COUNT"},
            },
        },
    }

    results := LintQueries(queries)

    if !hasResult(results, "WHC001") {
        t.Error("Expected WHC001 error for missing dataset")
    }
}
```

### Integration Tests

Test interactions between packages:

```go
// internal/roundtrip/roundtrip_test.go

func TestRoundTrip_BasicQuery(t *testing.T) {
    // JSON -> Go code -> Parse -> Build -> JSON
    // Verify semantic equivalence
}
```

### E2E Tests

End-to-end tests validate the full workflow using the `TestE2E_*` pattern:

```go
// cmd/wetwire-honeycomb/main_test.go

func TestE2E_InitImportList(t *testing.T) {
    // Full e2e test: init -> import -> list
    tmpDir := t.TempDir()
    projectPath := filepath.Join(tmpDir, "full-e2e")

    // Step 1: Create directory and initial files
    if err := os.MkdirAll(projectPath, 0755); err != nil {
        t.Fatalf("failed to create dir: %v", err)
    }

    // Step 2: Create Go files with declarations
    content := `package queries
import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    // ...
}
`
    if err := os.WriteFile(filepath.Join(projectPath, "queries.go"), []byte(content), 0644); err != nil {
        t.Fatalf("failed to write file: %v", err)
    }

    // Step 3: Discover and verify
    resources, err := discovery.DiscoverAll(projectPath)
    if err != nil {
        t.Fatalf("DiscoverAll failed: %v", err)
    }

    // Verify expected resources were found
    if len(resources.Queries) < 1 {
        t.Errorf("expected at least 1 query, got %d", len(resources.Queries))
    }
}

func TestE2E_AllResourceTypes(t *testing.T) {
    // Test discovering queries, SLOs, triggers, and boards
    // in a single codebase
}
```

### Test Helpers

Common test helpers are defined at the bottom of test files:

```go
func hasResult(results []LintResult, rule string) bool {
    for _, r := range results {
        if r.Rule == rule {
            return true
        }
    }
    return false
}

func findResult(results []LintResult, rule string) *LintResult {
    for i, r := range results {
        if r.Rule == rule {
            return &results[i]
        }
    }
    return nil
}

func getRepoRoot(t *testing.T) string {
    _, currentFile, _, ok := runtime.Caller(0)
    if !ok {
        t.Fatal("Failed to get current file path")
    }
    return filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
}
```

---

## Code Style and Conventions

### Go Style

Follow standard Go conventions:
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names

### Package Design

- **Public packages** (`query/`, `board/`, `slo/`, `trigger/`): Contain user-facing types
- **Internal packages** (`internal/`): Implementation details, not for external use
- Keep packages focused on a single responsibility

### Error Handling

```go
// Good: Return wrapped errors with context
if err != nil {
    return nil, fmt.Errorf("failed to parse file: %w", err)
}

// Good: Use errors for recoverable conditions
if dataset == "" {
    return []LintResult{{Rule: "WHC001", Message: "missing dataset"}}
}
```

### Documentation

- Document all exported types and functions
- Use package-level doc comments
- Include examples in doc comments where helpful

```go
// Query represents a Honeycomb query configuration.
// It maps to the Honeycomb Query API JSON format.
type Query struct {
    // Dataset is the Honeycomb dataset to query.
    // This is a required field.
    Dataset string

    // TimeRange specifies the time window for the query.
    TimeRange TimeRange
    // ...
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

Use descriptive branch names:

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

```
fix(discovery): handle nested composite literals

Previously, deeply nested query definitions were not
discovered correctly. This fixes the AST traversal
to find queries at any depth.

Fixes #456
```

```
docs(cli): add examples for diff command

- Add semantic diff example
- Document exit codes
- Add CI/CD integration example
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

Current version is defined in `cmd/wetwire-honeycomb/main.go`:

```go
const version = "0.3.0"
```

### Creating a Release

1. **Update version**

   Edit `cmd/wetwire-honeycomb/main.go`:
   ```go
   const version = "0.4.0"
   ```

2. **Update CHANGELOG** (if exists)

3. **Create PR**

   ```bash
   git checkout -b release/v0.4.0
   git add .
   git commit -m "chore: release v0.4.0"
   git push origin release/v0.4.0
   ```

4. **Merge PR after review**

5. **Tag the release**

   ```bash
   git checkout main
   git pull
   git tag -a v0.4.0 -m "Release v0.4.0"
   git push origin v0.4.0
   ```

6. **GitHub Release** (automated via CI)

### Post-Release

After a release:
- Update any dependent projects
- Announce in relevant channels
- Monitor for issues

---

## Troubleshooting

### Common Development Issues

**Tests fail with "failed to access directory"**

The test is likely using a relative path. Ensure tests use `t.TempDir()` or absolute paths:

```go
tmpDir := t.TempDir()
projectPath := filepath.Join(tmpDir, "test-project")
```

**Discovery doesn't find queries**

Check that:
1. Queries are top-level `var` declarations
2. Query names are exported (start with uppercase)
3. The file is a `.go` file (not `_test.go`)
4. The composite literal uses `query.Query` type

**Lint rules don't trigger**

Verify:
1. The rule is registered in `AllRules()` (or `AllBoardRules()`, etc.)
2. The rule's Check function returns results
3. Test data matches the rule's conditions

---

## See Also

- [CLI Reference](CLI.md) - Command documentation
- [Lint Rules](LINT_RULES.md) - Complete rule reference
- [FAQ](FAQ.md) - Common questions
