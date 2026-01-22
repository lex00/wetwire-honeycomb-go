<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This FAQ covers questions specific to Honeycomb query synthesis. For general wetwire questions, see the [central FAQ](https://github.com/lex00/wetwire/blob/main/docs/FAQ.md).

---

## General

### What is wetwire-honeycomb-go?

wetwire-honeycomb-go is a synthesis library that generates Honeycomb Query JSON from type-safe Go declarations. It does not execute queries or manage state - it only generates JSON output that you use with the Honeycomb API.

### How does it relate to the Honeycomb API?

wetwire-honeycomb-go generates Query JSON that conforms to the Honeycomb Query Specification. You use the generated JSON with the Honeycomb API to execute queries. The library handles:
- Type-safe query construction in Go
- Validation of query structure
- JSON generation

You handle:
- Calling the Honeycomb API
- Authentication
- Response processing

### Why use wetwire-honeycomb instead of the Honeycomb SDK?

| Feature | wetwire-honeycomb | Honeycomb SDK |
|---------|-------------------|---------------|
| Type safety | Compile-time validation | Runtime validation |
| Query organization | Code-based, version controlled | Often inline strings/maps |
| Reusability | Standard Go imports | Varies |
| Linting | Built-in best practices | Manual |
| Learning curve | Steeper (new patterns) | Gentler (familiar SDK patterns) |

Use wetwire-honeycomb when you want infrastructure-as-code patterns for queries. Use the SDK for ad-hoc queries or simpler use cases.

---

## Queries

### How do I define a query?

Create a top-level `var` with type `query.Query`:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

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

### Can I define queries inside functions?

No. Queries must be top-level `var` declarations to be discovered by the build command. Queries inside functions, methods, or closures will not be found.

```go
// GOOD: Top-level var
var MyQuery = query.Query{...}

// BAD: Won't be discovered
func GetQuery() query.Query {
    return query.Query{...}
}

// BAD: Won't be discovered
func init() {
    myQuery := query.Query{...}
}
```

### What calculations are supported?

Common Honeycomb calculations:

| Function | Description |
|----------|-------------|
| `query.Count()` | Count of events |
| `query.Sum(field)` | Sum of field values |
| `query.Avg(field)` | Average of field values |
| `query.Max(field)` | Maximum field value |
| `query.Min(field)` | Minimum field value |
| `query.P50(field)` | 50th percentile (median) |
| `query.P95(field)` | 95th percentile |
| `query.P99(field)` | 99th percentile |
| `query.Rate(numerator, denominator)` | Rate calculation |
| `query.CountDistinct(field)` | Unique count |

See Honeycomb documentation for the complete list.

### What filters are supported?

Common filter operations:

| Function | Description | Example |
|----------|-------------|---------|
| `query.Equals(field, value)` | Equality | `query.Equals("status", 200)` |
| `query.NotEquals(field, value)` | Inequality | `query.NotEquals("error", nil)` |
| `query.GT(field, value)` | Greater than | `query.GT("duration_ms", 500)` |
| `query.GTE(field, value)` | Greater than or equal | `query.GTE("retries", 3)` |
| `query.LT(field, value)` | Less than | `query.LT("score", 100)` |
| `query.LTE(field, value)` | Less than or equal | `query.LTE("age", 30)` |
| `query.Exists(field)` | Field exists | `query.Exists("user_id")` |
| `query.NotExists(field)` | Field does not exist | `query.NotExists("error")` |
| `query.Contains(field, value)` | String contains | `query.Contains("endpoint", "/api/")` |
| `query.NotContains(field, value)` | String does not contain | `query.NotContains("path", "admin")` |

### How do I combine multiple filters?

By default, multiple filters are combined with AND:

```go
Filters: []query.Filter{
    query.GT("duration_ms", 500),      // AND
    query.Contains("endpoint", "/api/"), // AND
    query.Exists("user_id"),            // AND
}
```

For OR logic, use `query.Or()`:

```go
Filters: []query.Filter{
    query.Or(
        query.Equals("status", 500),
        query.Equals("status", 502),
        query.Equals("status", 503),
    ),
}
```

### What time ranges are available?

| Function | Description |
|----------|-------------|
| `query.Minutes(n)` | Last n minutes |
| `query.Hours(n)` | Last n hours |
| `query.Days(n)` | Last n days |
| `query.Weeks(n)` | Last n weeks |

Example:

```go
TimeRange: query.Hours(24)  // Last 24 hours
TimeRange: query.Days(7)    // Last 7 days
```

---

## CLI

### What commands are available?

| Command | Purpose |
|---------|---------|
| `build` | Generate Query JSON from Go queries |
| `lint` | Check queries for issues |
| `list` | List all discovered queries |

See [CLI.md](CLI.md) for complete reference.

### How do I build queries?

```bash
# Build all queries in current package
wetwire-honeycomb build

# Build queries in specific package
wetwire-honeycomb build ./queries/...

# Save to file
wetwire-honeycomb build -o queries.json ./queries/...
```

### How do I use the generated JSON?

The generated JSON conforms to the Honeycomb Query Specification. Use it with the Honeycomb API:

```bash
curl -X POST https://api.honeycomb.io/1/queries/YOUR_DATASET \
  -H "X-Honeycomb-Team: $HONEYCOMB_API_KEY" \
  -H "Content-Type: application/json" \
  -d @queries.json
```

Or use with the Honeycomb Go SDK:

```go
import "github.com/honeycombio/libhoney-go"

// Load generated JSON
queryJSON := loadQueryJSON("queries.json")

// Use with Honeycomb client
client := libhoney.NewClient(...)
result, err := client.Query(ctx, queryJSON)
```

---

## Linting

### What lint rules are available?

See [LINT_RULES.md](LINT_RULES.md) for the complete list. Common rules include:

- **WHC001**: Use typed calculation functions
- **WHC002**: Use direct filter functions
- **WHC003**: Validate dataset references
- **WHC004**: Check time range validity
- **WHC005**: Ensure calculation names are unique

### How do I disable a lint rule?

Use the `--disable` flag:

```bash
wetwire-honeycomb lint --disable WHC003 ./queries/...
```

Or in `.wetwire-honeycomb.yaml`:

```yaml
lint:
  disabled_rules:
    - WHC003
    - WHC005
```

### Can lint auto-fix issues?

Yes, some rules support auto-fix with the `--fix` flag:

```bash
wetwire-honeycomb lint --fix ./queries/...
```

Not all rules can be auto-fixed. See [LINT_RULES.md](LINT_RULES.md) for which rules support auto-fix.

---

## File Organization

### Where should I put query files?

Organize queries by domain, team, or service:

```
queries/
├── auth/
│   ├── login.go
│   └── session.go
├── api/
│   ├── performance.go
│   └── errors.go
└── backend/
    └── database.go
```

Each file should use a consistent package name (e.g., `package queries`).

### Can I split a query across multiple files?

No. Each query must be a complete top-level `var` declaration in a single file.

### Should queries be in the same package?

It depends on your organization:

- **Single package** (`package queries`): Simple, easy to import, good for smaller projects
- **Multiple packages** (`package auth`, `package api`): Better organization, clearer ownership, good for larger projects

Both approaches work. Choose based on your team's needs.

---

## Integration

### How do I integrate with CI/CD?

Add lint and build checks to your CI pipeline:

```yaml
# .github/workflows/queries.yml
name: Query Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install wetwire-honeycomb
        run: go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest

      - name: Lint queries
        run: wetwire-honeycomb lint --format json ./queries/...

      - name: Build queries
        run: wetwire-honeycomb build -o queries.json ./queries/...
```

### Can I use wetwire-honeycomb with existing Honeycomb queries?

Not directly. wetwire-honeycomb generates JSON from Go code. To use existing queries:

1. Manually convert JSON to Go query definitions
2. Or use both approaches (wetwire for new queries, JSON for legacy)

There is no automated import tool yet (see roadmap).

### How do I version queries?

Query files are Go code, so use standard version control:

```bash
git add queries/
git commit -m "Add performance monitoring queries"
git tag v1.0.0
```

You can track generated JSON separately or regenerate it on deployment.

---

## Troubleshooting

### Build fails with "query not discovered"

Ensure your query is:
1. A top-level `var` (not inside a function)
2. Of type `query.Query`
3. In a `.go` file in the scanned path

### Lint reports "use typed calculation function"

Replace raw structs with typed functions:

```go
// BAD
Calculations: []query.Calculation{
    {Type: "P99", Field: "duration_ms"},
}

// GOOD
Calculations: []query.Calculation{
    query.P99("duration_ms"),
}
```

### Generated JSON doesn't work with Honeycomb API

Check:
1. Dataset name matches your Honeycomb dataset exactly
2. Field names match your Honeycomb schema
3. Calculation types are supported by Honeycomb
4. Time range is valid

Enable verbose output for debugging:

```bash
wetwire-honeycomb build -v ./queries/...
```

### "No queries found" error

Ensure:
1. Path includes `.go` files with query declarations
2. Using correct path pattern (e.g., `./queries/...` for recursive)
3. Files have valid Go syntax (run `go build` first)

---

## Performance

### How fast is query discovery?

Discovery uses Go's AST parser, which is very fast. Typical projects with hundreds of queries build in under 1 second.

### Is there a query limit?

No hard limit, but practical considerations:
- Keep individual query files focused (< 500 lines)
- Split large query sets into multiple packages
- Use meaningful names for easier management

### Can I parallelize builds?

The build command is already optimized for parallel package parsing. For very large projects, you can split builds by package:

```bash
wetwire-honeycomb build -o auth.json ./queries/auth/...
wetwire-honeycomb build -o api.json ./queries/api/...
```

---

## Advanced Usage

### Can I generate queries dynamically?

No. Queries must be static top-level declarations. For dynamic queries, use the Honeycomb SDK directly.

wetwire-honeycomb is designed for infrastructure-as-code patterns, not runtime query generation.

### Can I share queries across services?

Yes. Publish queries as a Go module:

```go
// In github.com/myorg/honeycomb-queries
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SharedQuery = query.Query{...}
```

Import in other services:

```go
import "github.com/myorg/honeycomb-queries"

// Reference in your code
var CustomQuery = queries.SharedQuery  // Copy or extend
```

### Can I extend query types?

The `query.Query` type is fixed to match Honeycomb's schema. For custom abstractions, create helper functions:

```go
// Helper for common SLO query pattern
func SLOQuery(dataset string, threshold int) query.Query {
    return query.Query{
        Dataset:   dataset,
        TimeRange: query.Hours(24),
        Calculations: []query.Calculation{
            query.P99("duration_ms"),
            query.P95("duration_ms"),
        },
        Filters: []query.Filter{
            query.LT("duration_ms", threshold),
        },
    }
}

// Use helper
var APISLOQuery = SLOQuery("production", 500)
```

---

## Resources

- [CLI Reference](CLI.md) - Complete command documentation
- [Lint Rules](LINT_RULES.md) - All WHC rules
- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md) - Core patterns
- [Honeycomb Query Specification](https://docs.honeycomb.io/api/query-specification/) - API reference
- [Honeycomb SDK](https://github.com/honeycombio/libhoney-go) - Query execution
