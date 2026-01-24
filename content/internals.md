---
title: "Internals"
---

This document describes the internal architecture and design of wetwire-honeycomb-go.

---

## Architecture Overview

wetwire-honeycomb-go is a synthesis library that converts Go struct declarations into Honeycomb Query JSON. It uses AST (Abstract Syntax Tree) analysis to discover resource definitions at compile time, validates them against Honeycomb constraints, and serializes them to the appropriate JSON format.

```
                         wetwire-honeycomb-go Architecture
 +============================================================================+
 |                                                                            |
 |   Go Source Files                                                          |
 |   +----------------+  +----------------+  +----------------+               |
 |   | queries/       |  | boards/        |  | alerts/        |               |
 |   | api.go         |  | dashboard.go   |  | triggers.go    |               |
 |   | auth.go        |  | overview.go    |  | slos.go        |               |
 |   +-------+--------+  +-------+--------+  +-------+--------+               |
 |           |                   |                   |                        |
 |           +-------------------+-------------------+                        |
 |                               |                                            |
 |                               v                                            |
 |   +------------------------------------------------------------------+     |
 |   |                     DISCOVERY PHASE                              |     |
 |   |  internal/discover/                                             |     |
 |   |  +--------------------+  +--------------------+                  |     |
 |   |  | AST Parser         |  | Type Detection     |                  |     |
 |   |  | - Parse Go files   |  | - query.Query      |                  |     |
 |   |  | - Walk syntax tree |  | - board.Board      |                  |     |
 |   |  | - Extract metadata |  | - slo.SLO          |                  |     |
 |   |  +--------------------+  | - trigger.Trigger  |                  |     |
 |   |                          +--------------------+                  |     |
 |   +------------------------------------------------------------------+     |
 |                               |                                            |
 |                               v                                            |
 |   +------------------------------------------------------------------+     |
 |   |                    VALIDATION PHASE                              |     |
 |   |  internal/lint/                                                  |     |
 |   |  +------------------+  +------------------+  +------------------+|     |
 |   |  | Query Rules      |  | Board Rules      |  | SLO/Trigger     ||     |
 |   |  | WHC001-023       |  | WHC030-034       |  | WHC040-056      ||     |
 |   |  +------------------+  +------------------+  +------------------+|     |
 |   +------------------------------------------------------------------+     |
 |                               |                                            |
 |                               v                                            |
 |   +------------------------------------------------------------------+     |
 |   |                  SERIALIZATION PHASE                             |     |
 |   |  internal/serialize/                                             |     |
 |   |  +------------------+  +------------------+  +------------------+|     |
 |   |  | Query -> JSON    |  | Board -> JSON    |  | SLO/Trigger     ||     |
 |   |  | serialize.go     |  | board.go         |  | slo.go          ||     |
 |   |  +------------------+  +------------------+  | trigger.go      ||     |
 |   |                                              +------------------+|     |
 |   +------------------------------------------------------------------+     |
 |                               |                                            |
 |                               v                                            |
 |   +------------------------------------------------------------------+     |
 |   |                      OUTPUT                                      |     |
 |   |  +--------------------+  +--------------------+                  |     |
 |   |  | queries.json       |  | boards.json        |                  |     |
 |   |  | slos.json          |  | triggers.json      |                  |     |
 |   |  +--------------------+  +--------------------+                  |     |
 |   +------------------------------------------------------------------+     |
 |                                                                            |
 +============================================================================+
```

---

## Core Components

### Discovery (`internal/discover/`)

The discovery package uses Go's `go/ast` and `go/parser` packages to find and extract resource definitions from Go source files.

#### Key Files

| File | Purpose |
|------|---------|
| `discovery.go` | Main entry point, orchestrates discovery for all resource types |
| `ast.go` | AST helper functions for extracting values from syntax nodes |
| `board.go` | Board-specific discovery logic |
| `slo.go` | SLO-specific discovery logic |
| `trigger.go` | Trigger-specific discovery logic |

#### How It Works

1. **File Scanning**: Walks the directory tree to find `.go` files (excluding `_test.go`)
2. **AST Parsing**: Uses `parser.ParseFile` to build an AST for each file
3. **Declaration Detection**: Inspects `*ast.GenDecl` nodes for `var` and `const` declarations
4. **Type Matching**: Checks if composite literals match known types (e.g., `query.Query`)
5. **Field Extraction**: Extracts field values from key-value expressions in the composite literal

---

### Serialization (`internal/serialize/`)

The serialize package converts Go structs into Honeycomb-compatible JSON format.

#### Key Files

| File | Purpose |
|------|---------|
| `serialize.go` | Query serialization with `ToJSON` and `ToJSONPretty` |
| `board.go` | Board serialization with panel type handling |
| `slo.go` | SLO serialization with SLI and burn alert conversion |
| `trigger.go` | Trigger serialization with threshold and recipient handling |

#### JSON Format Mapping

```go
// Internal JSON representation matches Honeycomb API format
type queryJSON struct {
    TimeRange         int               `json:"time_range,omitempty"`
    StartTime         int               `json:"start_time,omitempty"`
    EndTime           int               `json:"end_time,omitempty"`
    Breakdowns        []string          `json:"breakdowns,omitempty"`
    Calculations      []calculationJSON `json:"calculations,omitempty"`
    Filters           []filterJSON      `json:"filters,omitempty"`
    FilterCombination string            `json:"filter_combination,omitempty"`
    Orders            []orderJSON       `json:"orders,omitempty"`
    Limit             int               `json:"limit,omitempty"`
    Granularity       int               `json:"granularity,omitempty"`
}
```

---

### Linting (`internal/lint/`)

The lint package provides domain-specific validation rules for Honeycomb resources. Rules are prefixed with `WHC` (Wetwire Honeycomb).

#### Key Files

| File | Purpose |
|------|---------|
| `lint.go` | Core linting engine, result aggregation, filtering |
| `rules.go` | Query lint rules (WHC001-023) |
| `board_rules.go` | Board lint rules (WHC030-034) |
| `slo_rules.go` | SLO lint rules (WHC040-047) |
| `trigger_rules.go` | Trigger lint rules (WHC050-056) |

#### Rule Categories

| Code Range | Resource Type | Examples |
|------------|---------------|----------|
| WHC001-023 | Query | Missing dataset, invalid filters, time range limits |
| WHC030-034 | Board | Empty panels, panel count limits |
| WHC040-047 | SLO | Missing name, target percentage, burn alerts |
| WHC050-056 | Trigger | Missing name, no recipients, frequency warnings |

---

## Build Pipeline Stages

The build pipeline processes Go source files through several stages:

```
+----------+     +----------+     +---------+     +-------+     +-----------+     +------+
| Discover | --> | Validate | --> | Extract | --> | Order | --> | Serialize | --> | Emit |
+----------+     +----------+     +---------+     +-------+     +-----------+     +------+
```

### Stage 1: Discover

Find all resource declarations in Go source files.

### Stage 2: Validate

Check resources against Honeycomb constraints and best practices.

### Stage 3: Extract

Convert discovered metadata to concrete resource structs.

### Stage 4: Order

Determine output order and handle dependencies.

### Stage 5: Serialize

Convert Go structs to Honeycomb JSON format.

### Stage 6: Emit

Write output to destination.

---

## Type-Safe Reference Pattern

wetwire-honeycomb-go uses Go's type system to ensure references between resources are valid at compile time.

```go
// Define a query
var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
    },
}

// Reference it in a trigger - type-safe!
var LatencyAlert = trigger.Trigger{
    Name:      "High Latency Alert",
    Query:     SlowRequests,  // Direct reference, not a string
    Threshold: trigger.GreaterThan(500),
}
```

### Benefits

1. **Compile-Time Validation**: Invalid references cause compile errors
2. **IDE Support**: Autocomplete and go-to-definition work
3. **Refactoring Safety**: Renaming updates all references
4. **No String Typos**: Can't misspell a variable name

---

## Design Decisions

### Why AST Analysis?

1. **No Runtime Overhead**: Discovery happens at build time, not execution
2. **No Registration Required**: Resources are found automatically
3. **Full Context**: Access to file location, package name, and surrounding code
4. **Style Metadata**: Can detect patterns like inline definitions and nesting depth

### Why Separate Discovery and Serialization?

1. **Separation of Concerns**: Discovery extracts metadata; serialization formats output
2. **Different Input Types**: Discovery reads AST nodes; serialization reads Go structs
3. **Flexibility**: Can lint discovered resources without serializing
4. **Testing**: Each phase can be tested independently

### Why Type-Safe References?

1. **Compile-Time Safety**: Invalid references fail at build time
2. **IDE Integration**: Works with Go tooling (autocomplete, navigation)
3. **Refactoring**: Rename refactoring updates all references automatically
4. **No String Matching**: Eliminates a class of runtime errors

---

## See Also

- [CLI Reference](../cli/) - Complete command documentation
- [Lint Rules](../lint-rules/) - All WHC lint rules
- [FAQ](../faq/) - Common questions
