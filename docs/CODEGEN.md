# Code Generation

This document describes how wetwire-honeycomb-go generates Honeycomb Query JSON from Go declarations.

---

## Overview

Unlike domains that generate types from external schemas (like AWS CloudFormation), wetwire-honeycomb uses hand-crafted Go types that mirror the Honeycomb Query API. The code generation pipeline converts these typed Go structs into Honeycomb-compatible JSON.

---

## Directory Structure

```
wetwire-honeycomb-go/
├── query/                    # Public query types
│   ├── query.go              # Query struct definition
│   ├── calculation.go        # Calculation builders (P99, Count, etc.)
│   ├── filter.go             # Filter builders (GT, Exists, etc.)
│   ├── time.go               # Time range utilities
│   └── breakdown.go          # Breakdown utilities
│
├── board/                    # Public board types
│   ├── board.go              # Board struct
│   └── panel.go              # Panel types
│
├── slo/                      # Public SLO types
│   └── slo.go                # SLO struct
│
├── trigger/                  # Public trigger types
│   └── trigger.go            # Trigger struct
│
└── internal/
    ├── serialize/            # JSON serialization
    │   ├── serialize.go      # Query serialization
    │   ├── board.go          # Board serialization
    │   ├── slo.go            # SLO serialization
    │   └── trigger.go        # Trigger serialization
    └── discovery/            # AST-based discovery
```

---

## Generation Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                   Code Generation Pipeline                   │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐  │
│  │   DISCOVER  │ ─▶ │  SERIALIZE  │ ─▶ │     OUTPUT      │  │
│  └─────────────┘    └─────────────┘    └─────────────────┘  │
│                                                              │
│  AST parsing        Convert structs    Write JSON files     │
│  finds queries      to Honeycomb       to output dir        │
│                     JSON format                              │
└─────────────────────────────────────────────────────────────┘
```

### Stage 1: Discover

The discovery phase uses Go's AST package to find top-level variable declarations:

```go
import "github.com/lex00/wetwire-honeycomb-go/internal/discovery"

resources, err := discovery.DiscoverAll("./queries/...")
// resources.Queries, resources.Boards, resources.SLOs, resources.Triggers
```

Discovery extracts:
- Variable name
- File location and line number
- Struct field values (Dataset, TimeRange, Calculations, etc.)

### Stage 2: Serialize

The serialization phase converts Go structs to Honeycomb JSON format:

```go
import "github.com/lex00/wetwire-honeycomb-go/internal/serialize"

jsonBytes, err := serialize.QueryToJSON(query)
```

Key transformations:
- Calculation builders → Honeycomb calculation objects
- Filter builders → Honeycomb filter objects
- Time range types → Honeycomb time_range values
- Direct references → Resolved values

### Stage 3: Output

The output phase writes JSON files to the specified directory:

```bash
wetwire-honeycomb build ./queries/ -o ./output/
# Creates: output/queries/*.json
```

---

## Type Mapping

### Calculations

| Go Builder | Honeycomb JSON |
|------------|----------------|
| `query.Count()` | `{"op": "COUNT"}` |
| `query.P99("duration_ms")` | `{"op": "P99", "column": "duration_ms"}` |
| `query.Avg("response_time")` | `{"op": "AVG", "column": "response_time"}` |
| `query.Sum("bytes")` | `{"op": "SUM", "column": "bytes"}` |
| `query.Rate("errors", "requests")` | `{"op": "RATE", "column": "errors", "divisor_column": "requests"}` |

### Filters

| Go Builder | Honeycomb JSON |
|------------|----------------|
| `query.GT("duration", 500)` | `{"column": "duration", "op": ">", "value": 500}` |
| `query.Exists("user_id")` | `{"column": "user_id", "op": "exists"}` |
| `query.Contains("path", "/api")` | `{"column": "path", "op": "contains", "value": "/api"}` |

### Time Ranges

| Go Function | Honeycomb Value |
|-------------|-----------------|
| `query.Hours(2)` | `7200` (seconds) |
| `query.Days(7)` | `604800` (seconds) |
| `query.Minutes(30)` | `1800` (seconds) |

---

## Generated Output Structure

Each query generates a JSON file:

```json
{
  "time_range": 7200,
  "granularity": 0,
  "breakdowns": ["endpoint", "service"],
  "calculations": [
    {"op": "P99", "column": "duration_ms"},
    {"op": "COUNT"}
  ],
  "filters": [
    {"column": "duration_ms", "op": ">", "value": 500}
  ],
  "filter_combination": "AND"
}
```

---

## Validation

After generation, verify the output:

```bash
# Check syntax
wetwire-honeycomb lint ./queries/...

# Build and review output
wetwire-honeycomb build ./queries/... -o ./output/

# List all discovered resources
wetwire-honeycomb list ./queries/...
```

---

## Adding New Calculation Types

1. Add the builder function in `query/calculation.go`:

```go
func Heatmap(column string) Calculation {
    return Calculation{Op: "HEATMAP", Column: column}
}
```

2. Serialization handles it automatically (field-based mapping)

3. Add tests in `query/calculation_test.go`

4. Update documentation

---

## See Also

- [Developer Guide](DEVELOPERS.md) - Development workflow
- [CLI Reference](CLI.md) - Build command options
- [FAQ](FAQ.md) - Common questions
