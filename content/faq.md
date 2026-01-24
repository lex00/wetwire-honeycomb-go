---
title: "FAQ"
---

Frequently asked questions about wetwire-honeycomb-go.

---

## General

<details>
<summary>What is wetwire-honeycomb-go?</summary>

wetwire-honeycomb-go is a synthesis library that generates Honeycomb Query JSON from type-safe Go declarations. It does not execute queries or manage state - it only generates JSON output that you use with the Honeycomb API.
</details>

<details>
<summary>How does it relate to the Honeycomb API?</summary>

wetwire-honeycomb-go generates Query JSON that conforms to the Honeycomb Query Specification. You use the generated JSON with the Honeycomb API to execute queries. The library handles:
- Type-safe query construction in Go
- Validation of query structure
- JSON generation

You handle:
- Calling the Honeycomb API
- Authentication
- Response processing
</details>

<details>
<summary>Why use wetwire-honeycomb instead of the Honeycomb SDK?</summary>

| Feature | wetwire-honeycomb | Honeycomb SDK |
|---------|-------------------|---------------|
| Type safety | Compile-time validation | Runtime validation |
| Query organization | Code-based, version controlled | Often inline strings/maps |
| Reusability | Standard Go imports | Varies |
| Linting | Built-in best practices | Manual |
| Learning curve | Steeper (new patterns) | Gentler (familiar SDK patterns) |

Use wetwire-honeycomb when you want infrastructure-as-code patterns for queries. Use the SDK for ad-hoc queries or simpler use cases.
</details>

<details>
<summary>How do I manage Honeycomb API keys securely?</summary>

Never commit API keys to your repository. Use these approaches:

1. **Environment variables** - Set `HONEYCOMB_API_KEY` in your shell or CI/CD environment
2. **Secrets managers** - Use HashiCorp Vault, AWS Secrets Manager, or similar tools
3. **CI/CD secrets** - Store keys in GitHub Secrets, GitLab CI variables, etc.

Example usage in scripts:

```bash
curl -X POST https://api.honeycomb.io/1/queries/YOUR_DATASET \
  -H "X-Honeycomb-Team: $HONEYCOMB_API_KEY" \
  -H "Content-Type: application/json" \
  -d @queries.json
```

For local development, consider using `.envrc` with direnv (add `.envrc` to `.gitignore`).
</details>

<details>
<summary>Can I import existing board configurations?</summary>

Currently, there is no automated import tool for converting existing Honeycomb board JSON to Go code. To migrate existing boards:

1. Export your board configuration from Honeycomb (JSON format)
2. Manually translate the JSON structure to Go `board.Board{}` declarations
3. Use the linter to validate your translations: `wetwire-honeycomb lint ./boards/...`

An automated import feature is on the roadmap. For now, the manual process ensures proper type-safe declarations.
</details>

---

## Queries

<details>
<summary>How do I define a query?</summary>

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
</details>

<details>
<summary>Can I define queries inside functions?</summary>

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
</details>

<details>
<summary>What calculations are supported?</summary>

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
</details>

<details>
<summary>What filters are supported?</summary>

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
</details>

<details>
<summary>How do I combine multiple filters?</summary>

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
</details>

<details>
<summary>What time ranges are available?</summary>

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
</details>

---

## CLI

<details>
<summary>What commands are available?</summary>

| Command | Purpose |
|---------|---------|
| `build` | Generate Query JSON from Go queries |
| `lint` | Check queries for issues |
| `list` | List all discovered queries |

See [CLI](../cli/) for complete reference.
</details>

<details>
<summary>How do I build queries?</summary>

```bash
# Build all queries in current package
wetwire-honeycomb build

# Build queries in specific package
wetwire-honeycomb build ./queries/...

# Save to file
wetwire-honeycomb build -o queries.json ./queries/...
```
</details>

<details>
<summary>How do I use the generated JSON?</summary>

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
</details>

---

## Linting

<details>
<summary>What lint rules are available?</summary>

See [Lint Rules](../lint-rules/) for the complete list. Common rules include:

- **WHC001**: Use typed calculation functions
- **WHC002**: Use direct filter functions
- **WHC003**: Validate dataset references
- **WHC004**: Check time range validity
- **WHC005**: Ensure calculation names are unique
</details>

<details>
<summary>How does the linter help catch query errors?</summary>

The linter validates queries at build time, catching errors before they reach the Honeycomb API:

1. **Type validation** - Ensures calculations and filters use proper typed functions
2. **Field validation** - Checks that referenced fields exist in your schema (when configured)
3. **Best practices** - Warns about inefficient query patterns
4. **Consistency** - Enforces naming conventions and organization

Run the linter in CI to prevent broken queries from being deployed:

```yaml
- name: Lint queries
  run: wetwire-honeycomb lint --format json ./queries/...
```
</details>

<details>
<summary>How do I disable a lint rule?</summary>

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
</details>

<details>
<summary>Can lint auto-fix issues?</summary>

Yes, some rules support auto-fix with the `--fix` flag:

```bash
wetwire-honeycomb lint --fix ./queries/...
```

Not all rules can be auto-fixed. See [Lint Rules](../lint-rules/) for which rules support auto-fix.
</details>

---

## File Organization

<details>
<summary>What's the recommended project structure?</summary>

Organize queries by domain, team, or service:

```
queries/
├── auth/
│   ├── login.go
│   └── session.go
├── api/
│   ├── performance.go
│   └── errors.go
├── backend/
│   └── database.go
├── boards/
│   ├── overview.go
│   └── sre_dashboard.go
└── slos/
    ├── availability.go
    └── latency.go
```

Key recommendations:
- Group related queries in subdirectories
- Use descriptive file names that indicate the query domain
- Keep board and SLO definitions in separate directories
- Each file should use a consistent package name
</details>

<details>
<summary>Can I split a query across multiple files?</summary>

No. Each query must be a complete top-level `var` declaration in a single file.
</details>

<details>
<summary>Should queries be in the same package?</summary>

It depends on your organization:

- **Single package** (`package queries`): Simple, easy to import, good for smaller projects
- **Multiple packages** (`package auth`, `package api`): Better organization, clearer ownership, good for larger projects

Both approaches work. Choose based on your team's needs.
</details>

---

## SLOs

<details>
<summary>How do I define SLOs programmatically?</summary>

Define SLOs using the `slo.SLO` type:

```go
package slos

import "github.com/lex00/wetwire-honeycomb-go/slo"

var APIAvailability = slo.SLO{
    Name:        "API Availability",
    Description: "99.9% of requests succeed",
    Dataset:     "production",
    SLI: slo.SLI{
        Alias: "availability",
        Calculation: slo.Ratio{
            Numerator:   slo.Count(slo.LT("status_code", 500)),
            Denominator: slo.Count(nil),
        },
    },
    TargetPerMillion: 999000, // 99.9%
    TimePeriodDays:   30,
}

var P99Latency = slo.SLO{
    Name:        "P99 Latency",
    Description: "99th percentile under 500ms",
    Dataset:     "production",
    SLI: slo.SLI{
        Alias: "latency",
        Calculation: slo.Ratio{
            Numerator:   slo.Count(slo.LT("duration_ms", 500)),
            Denominator: slo.Count(nil),
        },
    },
    TargetPerMillion: 990000, // 99.0%
    TimePeriodDays:   7,
}
```

Build SLOs the same way as queries:

```bash
wetwire-honeycomb build -t slo ./slos/...
```
</details>

---

## Integration

<details>
<summary>How do I integrate with CI/CD?</summary>

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
</details>

<details>
<summary>Can I use wetwire-honeycomb with existing Honeycomb queries?</summary>

Not directly. wetwire-honeycomb generates JSON from Go code. To use existing queries:

1. Manually convert JSON to Go query definitions
2. Or use both approaches (wetwire for new queries, JSON for legacy)

There is no automated import tool yet (see roadmap).
</details>

<details>
<summary>How do I version queries?</summary>

Query files are Go code, so use standard version control:

```bash
git add queries/
git commit -m "Add performance monitoring queries"
git tag v1.0.0
```

You can track generated JSON separately or regenerate it on deployment.
</details>

---

## Troubleshooting

<details>
<summary>Build fails with "query not discovered"</summary>

Ensure your query is:
1. A top-level `var` (not inside a function)
2. Of type `query.Query`
3. In a `.go` file in the scanned path
</details>

<details>
<summary>Lint reports "use typed calculation function"</summary>

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
</details>

<details>
<summary>Generated JSON doesn't work with Honeycomb API</summary>

Check:
1. Dataset name matches your Honeycomb dataset exactly
2. Field names match your Honeycomb schema
3. Calculation types are supported by Honeycomb
4. Time range is valid

Enable verbose output for debugging:

```bash
wetwire-honeycomb build -v ./queries/...
```
</details>

<details>
<summary>"No queries found" error</summary>

Ensure:
1. Path includes `.go` files with query declarations
2. Using correct path pattern (e.g., `./queries/...` for recursive)
3. Files have valid Go syntax (run `go build` first)
</details>

---

## Performance

<details>
<summary>How fast is query discovery?</summary>

Discovery uses Go's AST parser, which is very fast. Typical projects with hundreds of queries build in under 1 second.
</details>

<details>
<summary>Is there a query limit?</summary>

No hard limit, but practical considerations:
- Keep individual query files focused (< 500 lines)
- Split large query sets into multiple packages
- Use meaningful names for easier management
</details>

<details>
<summary>Can I parallelize builds?</summary>

The build command is already optimized for parallel package parsing. For very large projects, you can split builds by package:

```bash
wetwire-honeycomb build -o auth.json ./queries/auth/...
wetwire-honeycomb build -o api.json ./queries/api/...
```
</details>

---

## Advanced Usage

<details>
<summary>Can I generate queries dynamically?</summary>

No. Queries must be static top-level declarations. For dynamic queries, use the Honeycomb SDK directly.

wetwire-honeycomb is designed for infrastructure-as-code patterns, not runtime query generation.
</details>

<details>
<summary>Can I share queries across services?</summary>

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
</details>

<details>
<summary>Can I extend query types?</summary>

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
</details>

---

## Resources

- [CLI Reference](../cli/) - Complete command documentation
- [Lint Rules](../lint-rules/) - All WHC rules
- [Honeycomb Query Specification](https://docs.honeycomb.io/api/query-specification/) - API reference
- [Honeycomb SDK](https://github.com/honeycombio/libhoney-go) - Query execution
