# Quick Start Guide

Get started with wetwire-honeycomb-go in minutes.

## Prerequisites

- Go 1.21 or later
- A Honeycomb account (for validation)

## Installation

See [README.md](../README.md#installation) for installation instructions.

## Create your first query

### 1. Initialize a new project

```bash
mkdir my-queries
cd my-queries
wetwire-honeycomb init
```

This creates a basic project structure with example queries.

### 2. Define a query

Create `queries/latency.go`:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowRequestsQuery identifies requests with high latency
var SlowRequestsQuery = query.Query{
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
    Limit: 100,
}
```

### 3. Lint your queries

```bash
wetwire-honeycomb lint ./queries
```

The linter checks for common issues:
- Missing dataset
- Missing time range
- Invalid filter operators
- High cardinality breakdowns

### 4. Build JSON output

```bash
wetwire-honeycomb build ./queries -o output.json
```

This generates Honeycomb Query JSON:

```json
{
  "time_range": 7200,
  "breakdowns": ["endpoint", "service"],
  "calculations": [
    {"op": "P99", "column": "duration_ms"},
    {"op": "COUNT"}
  ],
  "filters": [
    {"column": "duration_ms", "op": ">", "value": 500}
  ],
  "limit": 100
}
```

### 5. List discovered queries

```bash
wetwire-honeycomb list ./queries
```

Output:
```
QUERY              PACKAGE   DATASET      FILE
SlowRequestsQuery  queries   production   queries/latency.go:8
```

## Common patterns

### Time ranges

```go
// Relative time
TimeRange: query.Hours(24),    // Last 24 hours
TimeRange: query.Days(7),      // Last 7 days
TimeRange: query.Minutes(30),  // Last 30 minutes

// Absolute time
TimeRange: query.Absolute(startTime, endTime),
```

### Calculations

```go
Calculations: []query.Calculation{
    query.Count(),              // Count events
    query.P99("duration_ms"),   // 99th percentile
    query.Avg("response_size"), // Average
    query.Sum("bytes"),         // Sum
    query.Max("memory_mb"),     // Maximum
    query.CountDistinct("user_id"), // Unique count
}
```

### Filters

```go
Filters: []query.Filter{
    query.Equals("status", "error"),      // Exact match
    query.GT("duration_ms", 1000),        // Greater than
    query.Contains("path", "/api/"),      // String contains
    query.Exists("user_id"),              // Field exists
    query.NotEquals("env", "test"),       // Not equal
}
```

### Ordering

```go
Orders: []query.Order{
    {Op: "COUNT", Order: "descending"},   // Sort by count
    {Column: "endpoint", Order: "ascending"}, // Sort by column
},
```

## Next steps

- Read the [CLI Reference](CLI.md) for all commands
- Check [LINT_RULES.md](LINT_RULES.md) for validation rules
- See [FAQ.md](FAQ.md) for common questions

## Troubleshooting

### "no queries found"

Ensure your query variables are:
- Exported (start with uppercase letter)
- Using `query.Query` type

### Lint errors

Run with `--fix` to auto-fix some issues:

```bash
wetwire-honeycomb lint --fix ./queries
```

### Build errors

Check that all referenced datasets and columns exist in Honeycomb.
