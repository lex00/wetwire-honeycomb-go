# wetwire-honeycomb-go

Honeycomb query synthesis for Go - type-safe declarations that compile to Honeycomb Query JSON.

## Package Structure

```
wetwire-honeycomb-go/
├── cmd/wetwire-honeycomb/  # CLI entry point
│   └── main.go
├── internal/
│   ├── discovery/          # AST-based query discovery
│   ├── codegen/            # Query JSON generation
│   ├── lint/               # Honeycomb-specific lint rules
│   └── query/              # Query builder types
├── examples/               # Example queries
├── testdata/               # Test fixtures
└── docs/                   # Documentation
    ├── CLI.md              # Command reference
    ├── FAQ.md              # Common questions
    └── LINT_RULES.md       # Lint rules (WHC001, etc.)
```

## Core Components

### Query Discovery

AST-based discovery finds top-level query declarations:

```go
var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Breakdowns: []string{"endpoint", "service"},
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
        query.Count(),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
    },
}
```

Discovery extracts:
- Query name (`SlowRequests`)
- Query type (`query.Query`)
- File location and line number
- Dependencies (if any)

### Code Generation

Converts Go structs to Honeycomb Query JSON format:

```
Go Query struct → Intermediate map → Query JSON
```

Handles:
- Field name conversion (Go → JSON)
- Calculation serialization
- Filter serialization
- Time range formatting

### Linting

Uses the `WHC` prefix (Wetwire Honeycomb). See [docs/LINT_RULES.md](docs/LINT_RULES.md) for the complete rule reference.

## Key Principles

1. **Synthesis library** - Generates Query JSON, does not execute queries
2. **Type safety** - Go structs provide compile-time validation
3. **Flat declarations** - Top-level vars, no nested definitions
4. **Direct references** - Use Go variables, not string references
5. **Auto-discovery** - AST-based, no registration required

## Syntax Rules

### Query Declaration

```go
// GOOD: Top-level var with typed fields
var ErrorRate = query.Query{
    Dataset:   "backend",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Rate("errors", "requests"),
    },
}

// BAD: Inline definition
func GetQuery() query.Query {
    return query.Query{...}  // Won't be discovered
}

// BAD: String-based calculations
var BadQuery = query.Query{
    Calculations: []query.Calculation{
        {Type: "RATE", Field: "errors"},  // Use query.Rate()
    },
}
```

### Time Ranges

```go
// GOOD: Typed time range functions
TimeRange: query.Hours(2)
TimeRange: query.Days(7)
TimeRange: query.Minutes(30)

// BAD: String or raw values
TimeRange: "2h"  // Use query.Hours(2)
```

### Filters

```go
// GOOD: Typed filter functions
Filters: []query.Filter{
    query.GT("duration_ms", 500),
    query.Exists("user_id"),
    query.Contains("endpoint", "/api/"),
}

// BAD: Raw maps
Filters: []query.Filter{
    {"column": "duration_ms", "op": ">", "value": 500},
}
```

### Calculations

```go
// GOOD: Typed calculation functions
Calculations: []query.Calculation{
    query.P99("duration_ms"),
    query.Count(),
    query.Avg("response_size"),
    query.Sum("bytes_sent"),
}

// BAD: String-based
Calculations: []query.Calculation{
    {Type: "P99", Field: "duration_ms"},
}
```

## Common Patterns

### Multi-dimensional breakdown

```go
var RequestsByEndpointAndUser = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(24),
    Breakdowns: []string{"endpoint", "user_id", "region"},
    Calculations: []query.Calculation{
        query.Count(),
        query.P95("duration_ms"),
    },
}
```

### Complex filtering

```go
var SlowAuthRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Filters: []query.Filter{
        query.GT("duration_ms", 1000),
        query.Contains("endpoint", "/auth"),
        query.Exists("user_id"),
    },
    Calculations: []query.Calculation{
        query.Count(),
        query.P99("duration_ms"),
    },
}
```

### Multiple calculations

```go
var PerformanceMetrics = query.Query{
    Dataset:   "backend",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"service"},
    Calculations: []query.Calculation{
        query.Count(),
        query.Avg("duration_ms"),
        query.P50("duration_ms"),
        query.P95("duration_ms"),
        query.P99("duration_ms"),
        query.Max("duration_ms"),
    },
}
```

## File Organization

### queries/

Organize queries by domain or team:

```
queries/
├── auth/
│   ├── login_performance.go
│   └── session_metrics.go
├── api/
│   ├── endpoint_latency.go
│   └── error_rates.go
└── backend/
    ├── database_queries.go
    └── cache_performance.go
```

Each file should:
- Use `package queries` (or domain-specific package)
- Import `"github.com/lex00/wetwire-honeycomb-go/query"`
- Declare one or more queries as top-level vars

## Gotchas

1. **Discovery requires top-level vars** - Queries inside functions or methods won't be discovered
2. **Package names matter** - Use consistent package naming for easier imports
3. **Dataset names** - String values, must match Honeycomb datasets exactly
4. **Field names** - Column names must match Honeycomb schema
5. **Time range validation** - Some combinations may not be supported by Honeycomb API

## Development Workflow

### 1. Define queries

```bash
# Create query file
vim queries/performance.go
```

### 2. Build to Query JSON

```bash
wetwire-honeycomb build ./queries/...
```

### 3. Lint for issues

```bash
wetwire-honeycomb lint ./queries/...
```

### 4. List all queries

```bash
wetwire-honeycomb list
```

### 5. Use JSON with Honeycomb API

Generated JSON can be used with Honeycomb Query API (user's responsibility).

## Running Tests

```bash
go test -v ./...
```

## Running Scenarios

**IMPORTANT:** Scenarios use the Claude CLI, NOT the Anthropic API. No `ANTHROPIC_API_KEY` needed.

Scenarios are run using `wetwire-core-go`'s scenario runner, which invokes `claude` CLI directly:

```bash
# Clone wetwire-core-go if needed
# Then from that directory:
go run ./cmd/run_scenario /path/to/wetwire-honeycomb-go/examples/tasks_api_scenario [persona] --verbose

# Examples:
go run ./cmd/run_scenario ./examples/tasks_api_scenario beginner --verbose
go run ./cmd/run_scenario ./examples/tasks_api_scenario intermediate --verbose
go run ./cmd/run_scenario ./examples/tasks_api_scenario expert --verbose
go run ./cmd/run_scenario ./examples/tasks_api_scenario --all --verbose
```

The scenario runner:
1. Reads `scenario.yaml` for config (model, timeout, validation rules)
2. Reads `system_prompt.md` for domain knowledge
3. Reads persona-specific prompts from `prompts/{persona}.md`
4. Invokes Claude CLI to generate output
5. Scores and validates results

Do NOT use `wetwire-honeycomb test` for scenarios - that command is for ad-hoc testing with the Anthropic API.

## Build Commands

```bash
# Build CLI tool
go build ./cmd/wetwire-honeycomb

# Install globally
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
```


## Diff

Compare Honeycomb query configurations semantically:

```bash
# Compare two files
wetwire-honeycomb diff file1 file2

# JSON output for CI/CD
wetwire-honeycomb diff file1 file2 -f json

# Ignore array ordering differences
wetwire-honeycomb diff file1 file2 --ignore-order
```

The diff command performs semantic comparison by resource name, detecting:
- Added resources
- Removed resources
- Modified resources (with property-level change details)

Exit code is 1 if differences are found, enabling CI pipeline validation.

## Key Files

| File | Purpose |
|------|---------|
| `cmd/wetwire-honeycomb/main.go` | CLI entry point |
| `internal/discover/` | AST-based query discovery |
| `internal/codegen/` | Query JSON generation |
| `internal/lint/` | Lint rule implementations |
| `internal/query/` | Query type definitions |
| `internal/agent/domain.go` | AI system prompt for scenarios |
| `examples/tasks_api_scenario/` | Example scenario (run via wetwire-core-go) |
| `README.md` | Quick start guide |
| `docs/CLI.md` | Complete command reference |
| `docs/FAQ.md` | Common questions |
| `docs/LINT_RULES.md` | All WHC rules |
