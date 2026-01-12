# CLI Reference

Complete command reference for wetwire-honeycomb.

---

## Commands

### build

Generate Honeycomb Query JSON from Go query declarations.

```bash
wetwire-honeycomb build [OPTIONS] [PATH]
```

**Description:**

Discovers query declarations in Go source files and generates Honeycomb Query JSON. The output can be used with the Honeycomb Query API.

**Arguments:**

| Argument | Description | Default |
|----------|-------------|---------|
| `PATH` | Path to Go package(s) containing queries | `.` |

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output FILE` | Write output to FILE instead of stdout | stdout |
| `-f, --format FORMAT` | Output format: `json`, `yaml` | `json` |
| `--pretty` | Pretty-print JSON output | `false` |
| `-v, --verbose` | Verbose output (show discovery details) | `false` |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Build failed (invalid queries, references, etc.) |
| 2 | Invalid arguments or options |

**Examples:**

```bash
# Build all queries in current directory
wetwire-honeycomb build

# Build queries in specific package
wetwire-honeycomb build ./queries/...

# Build and save to file
wetwire-honeycomb build -o queries.json ./queries/...

# Build with pretty-printed JSON
wetwire-honeycomb build --pretty ./queries/...

# Build with verbose output
wetwire-honeycomb build -v ./queries/...

# Build YAML format
wetwire-honeycomb build -f yaml ./queries/...
```

**Output Format:**

The build command generates an array of Honeycomb Query JSON objects:

```json
[
  {
    "time_range": 7200,
    "granularity": 0,
    "breakdowns": ["endpoint", "service"],
    "calculations": [
      {
        "op": "P99",
        "column": "duration_ms"
      },
      {
        "op": "COUNT"
      }
    ],
    "filters": [
      {
        "column": "duration_ms",
        "op": ">",
        "value": 500
      }
    ],
    "filter_combination": "AND",
    "orders": [],
    "havings": [],
    "limit": 1000
  }
]
```

---

### lint

Check Go query declarations for issues and anti-patterns.

```bash
wetwire-honeycomb lint [OPTIONS] [PATH]
```

**Description:**

Analyzes query declarations and reports issues based on Honeycomb best practices and wetwire conventions.

**Arguments:**

| Argument | Description | Default |
|----------|-------------|---------|
| `PATH` | Path to Go package(s) to lint | `.` |

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--fix` | Automatically fix issues where possible | `false` |
| `--severity LEVEL` | Minimum severity: `error`, `warning`, `info` | `warning` |
| `--rules RULES` | Comma-separated list of rules to check | all |
| `--disable RULES` | Comma-separated list of rules to skip | none |
| `-v, --verbose` | Show rule explanations | `false` |
| `--format FORMAT` | Output format: `text`, `json` | `text` |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | No issues found (or only info-level) |
| 1 | Issues found (warnings or errors) |
| 2 | Invalid arguments or options |

**Examples:**

```bash
# Lint all queries
wetwire-honeycomb lint ./queries/...

# Lint with auto-fix
wetwire-honeycomb lint --fix ./queries/...

# Show only errors
wetwire-honeycomb lint --severity error ./queries/...

# Check specific rules
wetwire-honeycomb lint --rules WHC001,WHC002 ./queries/...

# Disable specific rules
wetwire-honeycomb lint --disable WHC003 ./queries/...

# Verbose output with explanations
wetwire-honeycomb lint -v ./queries/...

# JSON output for CI/CD
wetwire-honeycomb lint --format json ./queries/...
```

**Output Format (text):**

```
queries/performance.go:15:5: WHC001 (error): Use typed calculation function
    Found: {Type: "P99", Field: "duration_ms"}
    Expected: query.P99("duration_ms")

queries/api.go:23:5: WHC003 (warning): Dataset 'prod' not validated
    Consider using a constant for dataset names

2 issues found (1 error, 1 warning)
```

**Output Format (json):**

```json
{
  "issues": [
    {
      "file": "queries/performance.go",
      "line": 15,
      "column": 5,
      "rule": "WHC001",
      "severity": "error",
      "message": "Use typed calculation function",
      "details": "Found: {Type: \"P99\", Field: \"duration_ms\"}\nExpected: query.P99(\"duration_ms\")"
    }
  ],
  "summary": {
    "total": 2,
    "errors": 1,
    "warnings": 1,
    "info": 0
  }
}
```

---

### list

List all discovered query declarations.

```bash
wetwire-honeycomb list [OPTIONS] [PATH]
```

**Description:**

Discovers and lists all query declarations with metadata (name, file, line, dataset).

**Arguments:**

| Argument | Description | Default |
|----------|-------------|---------|
| `PATH` | Path to Go package(s) to scan | `.` |

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--format FORMAT` | Output format: `table`, `json`, `csv` | `table` |
| `--sort FIELD` | Sort by: `name`, `file`, `dataset` | `name` |
| `-v, --verbose` | Include additional details | `false` |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Invalid arguments or options |

**Examples:**

```bash
# List all queries
wetwire-honeycomb list

# List queries in specific package
wetwire-honeycomb list ./queries/...

# List as JSON
wetwire-honeycomb list --format json

# List as CSV
wetwire-honeycomb list --format csv > queries.csv

# Sort by file
wetwire-honeycomb list --sort file

# Verbose output with details
wetwire-honeycomb list -v
```

**Output Format (table):**

```
NAME                 DATASET      FILE                      LINE
SlowRequests         production   queries/performance.go    12
ErrorRate            backend      queries/api.go            25
LoginMetrics         auth         queries/auth.go           8

3 queries found
```

**Output Format (json):**

```json
{
  "queries": [
    {
      "name": "SlowRequests",
      "dataset": "production",
      "file": "queries/performance.go",
      "line": 12,
      "calculations": ["P99", "COUNT"],
      "breakdowns": ["endpoint", "service"],
      "filters": 1
    }
  ],
  "total": 3
}
```

---

## Global Options

These options work with all commands:

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help for command |
| `--version` | Show version information |
| `--no-color` | Disable colored output |

---

## Exit Codes

All commands use consistent exit codes:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Command-specific failure (build failed, lint issues, etc.) |
| 2 | Invalid usage (bad arguments, missing files, etc.) |

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WETWIRE_HONEYCOMB_CACHE` | Cache directory for query metadata | `~/.cache/wetwire-honeycomb` |
| `WETWIRE_HONEYCOMB_LOG` | Log level: `debug`, `info`, `warn`, `error` | `info` |
| `NO_COLOR` | Disable colored output (set to any value) | - |

**Examples:**

```bash
# Enable debug logging
export WETWIRE_HONEYCOMB_LOG=debug
wetwire-honeycomb build ./queries/...

# Disable colored output
export NO_COLOR=1
wetwire-honeycomb lint ./queries/...

# Use custom cache directory
export WETWIRE_HONEYCOMB_CACHE=/tmp/wetwire-cache
wetwire-honeycomb build ./queries/...
```

---

## Configuration File

wetwire-honeycomb can be configured via `.wetwire-honeycomb.yaml` in the project root:

```yaml
# Lint configuration
lint:
  severity: warning
  disabled_rules:
    - WHC003
  auto_fix: false

# Build configuration
build:
  format: json
  pretty: true
  output: queries.json

# List configuration
list:
  format: table
  sort: name
```

**Precedence:** CLI flags > environment variables > config file > defaults

---

## Examples

### Basic workflow

```bash
# 1. Create query file
cat > queries/performance.go <<EOF
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Breakdowns: []string{"endpoint"},
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
        query.Count(),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
    },
}
EOF

# 2. Lint queries
wetwire-honeycomb lint ./queries/...

# 3. Build Query JSON
wetwire-honeycomb build -o queries.json ./queries/...

# 4. Use JSON with Honeycomb API (user's responsibility)
curl -X POST https://api.honeycomb.io/1/queries/dataset \
  -H "X-Honeycomb-Team: $HONEYCOMB_API_KEY" \
  -d @queries.json
```

### CI/CD integration

```bash
#!/bin/bash
# ci-check.sh

set -e

echo "Linting queries..."
wetwire-honeycomb lint --format json ./queries/... > lint-results.json

# Fail if errors found
if [ $(jq '.summary.errors' lint-results.json) -gt 0 ]; then
  echo "Lint errors found!"
  jq '.issues[] | select(.severity == "error")' lint-results.json
  exit 1
fi

echo "Building queries..."
wetwire-honeycomb build -o queries.json ./queries/...

echo "All checks passed!"
```

---

## See Also

- [FAQ](FAQ.md) - Common questions
- [LINT_RULES.md](LINT_RULES.md) - Complete lint rule reference
- [Honeycomb Query Specification](https://docs.honeycomb.io/api/query-specification/)
