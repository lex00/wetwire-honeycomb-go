# Internals

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
 |   |  internal/discovery/                                             |     |
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

### Discovery (`internal/discovery/`)

The discovery package uses Go's `go/ast` and `go/parser` packages to find and extract resource definitions from Go source files. Discovery is the entry point for the build pipeline.

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

```go
// Discovery walks the AST looking for patterns like:
var SlowRequests = query.Query{
    Dataset:   "production",    // Extracted as string literal
    TimeRange: query.Hours(2),  // Extracted as function call with args
    //...
}
```

#### Type Detection

Each resource type has a dedicated type checker:

```go
// From ast.go
func isQueryType(expr ast.Expr) bool {
    // Matches: query.Query
    switch t := expr.(type) {
    case *ast.SelectorExpr:
        if ident, ok := t.X.(*ast.Ident); ok {
            return ident.Name == "query" && t.Sel.Name == "Query"
        }
    }
    return false
}

// Similar functions exist for:
// - isBoardType() -> board.Board
// - isSLOType() -> slo.SLO
// - isTriggerType() -> trigger.Trigger
```

#### Discovered Resource Types

Each resource type has a corresponding `Discovered*` struct that holds extracted metadata:

```go
type DiscoveredQuery struct {
    Name         string        // Variable name (e.g., "SlowRequests")
    Package      string        // Package name
    File         string        // Absolute file path
    Line         int           // Line number
    Dataset      string        // Honeycomb dataset
    TimeRange    TimeRange     // Time window
    Breakdowns   []string      // Group-by columns
    Calculations []Calculation // Aggregations
    Filters      []Filter      // Filter conditions
    // ... additional fields
}

type DiscoveredResources struct {
    Queries  []DiscoveredQuery
    SLOs     []DiscoveredSLO
    Triggers []DiscoveredTrigger
    Boards   []DiscoveredBoard
}
```

---

### Serialization (`internal/serialize/`)

The serialize package converts Go structs (from the public `query`, `board`, `slo`, and `trigger` packages) into Honeycomb-compatible JSON format.

#### Key Files

| File | Purpose |
|------|---------|
| `serialize.go` | Query serialization with `ToJSON` and `ToJSONPretty` |
| `board.go` | Board serialization with panel type handling |
| `slo.go` | SLO serialization with SLI and burn alert conversion |
| `trigger.go` | Trigger serialization with threshold and recipient handling |

#### JSON Format Mapping

The serializer maps Go struct fields to Honeycomb API JSON keys (snake_case):

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

#### Serialization Flow

```
Go Struct (query.Query)
         |
         v
Internal JSON Struct (queryJSON)  <- Field name conversion
         |
         v
json.Marshal() / json.MarshalIndent()
         |
         v
Honeycomb Query JSON ([]byte)
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

#### Rule Structure

Each rule is defined as a struct with a check function:

```go
type Rule struct {
    Code     string  // e.g., "WHC001"
    Severity string  // "error", "warning", or "info"
    Message  string  // Human-readable description
    Check    func(query DiscoveredQuery) []LintResult
}

type LintResult struct {
    Rule     string  // Rule code
    Severity string  // Severity level
    Message  string  // Specific message for this occurrence
    File     string  // File path
    Line     int     // Line number
    Query    string  // Resource name
}
```

#### Rule Categories

| Code Range | Resource Type | Examples |
|------------|---------------|----------|
| WHC001-023 | Query | Missing dataset, invalid filters, time range limits |
| WHC030-034 | Board | Empty panels, panel count limits |
| WHC040-047 | SLO | Missing name, target percentage, burn alerts |
| WHC050-056 | Trigger | Missing name, no recipients, frequency warnings |

#### Lint Execution

```go
// Run all rules on all resources
func LintAll(resources *DiscoveredResources) []LintResult {
    var results []LintResult
    results = append(results, LintQueries(resources.Queries)...)
    results = append(results, LintBoards(resources.Boards)...)
    results = append(results, LintSLOs(resources.SLOs)...)
    results = append(results, LintTriggers(resources.Triggers)...)
    return results
}
```

---

### Builder (`internal/builder/`)

The builder package orchestrates the discovery and registration pipeline, providing options for filtering and namespacing.

#### Key Files

| File | Purpose |
|------|---------|
| `builder.go` | `Builder` struct and fluent configuration API |
| `registry.go` | `Registry` for storing discovered resources, duplicate detection |

#### Builder Pattern

```go
b, _ := builder.NewBuilder("./queries")
result, _ := b.
    WithNamespacing(true).      // Enable package-based namespacing
    WithStrictMode(true).       // Fail on duplicates
    WithPackageFilter("api").   // Filter by package
    WithDatasetFilter("prod").  // Filter by dataset
    Build()

queries := result.Queries()
```

#### Registry

The registry stores discovered resources and tracks duplicates:

```go
type Registry struct {
    queries     map[string]DiscoveredQuery
    duplicates  []DuplicateEntry
    namespacing bool  // If true, uses "pkg.Name" as key
    strictMode  bool  // If true, returns error on duplicates
}
```

---

## Build Pipeline Stages

The build pipeline processes Go source files through several stages:

```
+----------+     +----------+     +---------+     +-------+     +-----------+     +------+
| Discover | --> | Validate | --> | Extract | --> | Order | --> | Serialize | --> | Emit |
+----------+     +----------+     +---------+     +-------+     +-----------+     +------+
```

### Stage 1: Discover

**Purpose**: Find all resource declarations in Go source files

**Input**: Directory path (e.g., `./queries/...`)

**Process**:
1. Walk directory tree for `.go` files
2. Parse each file into AST
3. Inspect declarations for known types
4. Extract metadata from composite literals

**Output**: `DiscoveredResources` containing all found resources

### Stage 2: Validate

**Purpose**: Check resources against Honeycomb constraints and best practices

**Input**: `DiscoveredResources`

**Process**:
1. Apply all applicable lint rules
2. Collect errors, warnings, and info messages
3. Sort results by file and line number

**Output**: `[]LintResult` with any issues found

### Stage 3: Extract

**Purpose**: Convert discovered metadata to concrete resource structs

**Input**: `DiscoveredQuery`, `DiscoveredBoard`, etc.

**Process**:
1. Map discovered fields to public struct fields
2. Handle type conversions (e.g., time range seconds to duration)
3. Resolve references where possible

**Output**: `query.Query`, `board.Board`, `slo.SLO`, `trigger.Trigger`

### Stage 4: Order

**Purpose**: Determine output order and handle dependencies

**Input**: Concrete resource structs

**Process**:
1. Group resources by type
2. Sort by name or other criteria
3. Resolve cross-resource references (e.g., Board -> Query)

**Output**: Ordered resource collections

### Stage 5: Serialize

**Purpose**: Convert Go structs to Honeycomb JSON format

**Input**: Concrete resource structs

**Process**:
1. Map struct fields to JSON fields
2. Apply snake_case naming
3. Handle nested structures (panels, SLI, thresholds)
4. Omit empty fields

**Output**: JSON bytes for each resource

### Stage 6: Emit

**Purpose**: Write output to destination

**Input**: JSON bytes

**Process**:
1. Write to stdout or file
2. Apply formatting (compact or pretty)
3. Group by resource type if multiple

**Output**: JSON files or stdout

---

## Resource Types and Relationships

### Resource Type Hierarchy

```
                    +------------------+
                    |  Honeycomb API   |
                    +------------------+
                            ^
                            |
    +----------+------------+------------+-----------+
    |          |            |            |           |
+-------+  +-------+    +-------+    +---------+    |
| Query |  | Board |    |  SLO  |    | Trigger |    |
+-------+  +-------+    +-------+    +---------+    |
    ^          |            |            |          |
    |          v            v            v          |
    |     +--------+   +--------+   +--------+      |
    +<----| Panels |   |  SLI   |   |  Query |------+
          +--------+   +--------+   +--------+
               |           |
               v           v
          QueryRef    GoodEvents
                      TotalEvents
```

### Query

The foundational resource type - defines aggregation queries against Honeycomb data.

```go
type Query struct {
    Dataset      string
    TimeRange    TimeRange
    Breakdowns   []string
    Calculations []Calculation
    Filters      []Filter
    Orders       []Order
    Limit        int
    Granularity  int
}
```

### Board

Visual dashboards containing multiple panels, each referencing queries.

```go
type Board struct {
    Name          string
    Description   string
    Panels        []Panel        // Can reference queries
    PresetFilters []BoardFilter
}

// Panels reference queries via type-safe functions
board.QueryPanel(MyQuery)  // References query.Query variable
```

### SLO (Service Level Objective)

Error budget tracking with good/total event queries.

```go
type SLO struct {
    Name        string
    Description string
    Dataset     string
    SLI         SLI          // References two queries
    Target      Target       // Percentage target
    TimePeriod  TimePeriod
    BurnAlerts  []BurnAlert
}

type SLI struct {
    GoodEvents  query.Query  // Query for successful events
    TotalEvents query.Query  // Query for all events
}
```

### Trigger

Alert definitions that fire when query results exceed thresholds.

```go
type Trigger struct {
    Name        string
    Description string
    Dataset     string
    Query       query.Query  // The query to evaluate
    Threshold   Threshold    // When to fire
    Frequency   Frequency    // How often to check
    Recipients  []Recipient  // Where to send alerts
    Disabled    bool
}
```

---

## Type-Safe Reference Pattern

wetwire-honeycomb-go uses Go's type system to ensure references between resources are valid at compile time.

### The Pattern

Instead of using string IDs to reference resources, use Go variables directly:

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

// Reference it in a board
var PerformanceBoard = board.Board{
    Name: "Performance Dashboard",
    Panels: []board.Panel{
        board.QueryPanel(SlowRequests),  // Direct reference
    },
}
```

### Benefits

1. **Compile-Time Validation**: Invalid references cause compile errors
2. **IDE Support**: Autocomplete and go-to-definition work
3. **Refactoring Safety**: Renaming updates all references
4. **No String Typos**: Can't misspell a variable name

### Discovery of References

During AST analysis, references are extracted by examining identifier nodes:

```go
// In trigger discovery
case "Query":
    if ident, ok := kv.Value.(*ast.Ident); ok {
        trigger.QueryRef = ident.Name  // Stores "SlowRequests"
    }

// In board discovery
case "QueryPanel":
    if len(call.Args) > 0 {
        if ident, ok := call.Args[0].(*ast.Ident); ok {
            queryRefs = append(queryRefs, ident.Name)
        }
    }
```

---

## Extension Points

### Adding a New Resource Type

To add support for a new Honeycomb resource type (e.g., Markers):

#### 1. Create Public Package

```go
// marker/marker.go
package marker

type Marker struct {
    Name      string
    Dataset   string
    Type      MarkerType
    StartTime int
    EndTime   int
    Message   string
    URL       string
}
```

#### 2. Add Discovery Logic

```go
// internal/discovery/marker.go
package discovery

type DiscoveredMarker struct {
    Name       string
    Package    string
    File       string
    Line       int
    MarkerName string
    Dataset    string
    // ... extracted fields
}

func DiscoverMarkers(dir string) ([]DiscoveredMarker, error) {
    // Walk files, parse AST, find marker.Marker composites
}

func isMarkerType(expr ast.Expr) bool {
    // Check for marker.Marker type
}

func extractMarkerFromComposite(comp *ast.CompositeLit, ...) DiscoveredMarker {
    // Extract field values
}
```

#### 3. Update DiscoveredResources

```go
// internal/discovery/discovery.go
type DiscoveredResources struct {
    Queries  []DiscoveredQuery
    SLOs     []DiscoveredSLO
    Triggers []DiscoveredTrigger
    Boards   []DiscoveredBoard
    Markers  []DiscoveredMarker  // Add new type
}

func DiscoverAll(dir string) (*DiscoveredResources, error) {
    // ... existing code ...

    markers, err := DiscoverMarkers(dir)
    if err != nil {
        return nil, fmt.Errorf("failed to discover markers: %w", err)
    }
    resources.Markers = markers

    return resources, nil
}
```

#### 4. Add Serialization

```go
// internal/serialize/marker.go
package serialize

type markerJSON struct {
    Name      string `json:"name"`
    Dataset   string `json:"dataset,omitempty"`
    Type      string `json:"type"`
    StartTime int    `json:"start_time,omitempty"`
    EndTime   int    `json:"end_time,omitempty"`
    Message   string `json:"message,omitempty"`
    URL       string `json:"url,omitempty"`
}

func MarkerToJSON(m marker.Marker) ([]byte, error) {
    jm := toMarkerJSON(m)
    return json.Marshal(jm)
}
```

#### 5. Add Lint Rules

```go
// internal/lint/marker_rules.go
package lint

type MarkerRule struct {
    Code     string
    Severity string
    Message  string
    Check    func(marker DiscoveredMarker) []LintResult
}

func AllMarkerRules() []MarkerRule {
    return []MarkerRule{
        WHC060MarkerMissingName(),
        WHC061MarkerInvalidTimeRange(),
    }
}

func WHC060MarkerMissingName() MarkerRule {
    return MarkerRule{
        Code:     "WHC060",
        Severity: "error",
        Message:  "Marker missing name",
        Check: func(marker DiscoveredMarker) []LintResult {
            // Validation logic
        },
    }
}
```

#### 6. Update CLI Commands

```go
// cmd/wetwire-honeycomb/main.go

// In newBuildCmd() - add marker serialization
if len(resources.Markers) > 0 {
    markerMap := make(map[string]json.RawMessage)
    for _, dm := range resources.Markers {
        m := discoveredToMarker(dm)
        data, _ := serialize.MarkerToJSON(m)
        markerMap[dm.Name] = data
    }
    data, _ := json.Marshal(markerMap)
    outputData["markers"] = data
}
```

### Adding a New Lint Rule

To add a new lint rule:

#### 1. Define the Rule

```go
// internal/lint/rules.go

// WHC024DuplicateBreakdown checks for duplicate breakdown columns.
func WHC024DuplicateBreakdown() Rule {
    return Rule{
        Code:     "WHC024",
        Severity: "warning",
        Message:  "Duplicate breakdown column",
        Check: func(query DiscoveredQuery) []LintResult {
            seen := make(map[string]bool)
            var results []LintResult

            for _, bd := range query.Breakdowns {
                if seen[bd] {
                    results = append(results, LintResult{
                        Rule:     "WHC024",
                        Severity: "warning",
                        Message:  fmt.Sprintf("Duplicate breakdown column: %s", bd),
                        File:     query.File,
                        Line:     query.Line,
                        Query:    query.Name,
                    })
                }
                seen[bd] = true
            }
            return results
        },
    }
}
```

#### 2. Register the Rule

```go
// internal/lint/rules.go
func AllRules() []Rule {
    return []Rule{
        // ... existing rules ...
        WHC024DuplicateBreakdown(),  // Add new rule
    }
}
```

#### 3. Add Tests

```go
// internal/lint/rules_test.go
func TestWHC024DuplicateBreakdown(t *testing.T) {
    query := DiscoveredQuery{
        Name:       "TestQuery",
        Breakdowns: []string{"service", "endpoint", "service"},  // Duplicate
    }

    rule := WHC024DuplicateBreakdown()
    results := rule.Check(query)

    if len(results) != 1 {
        t.Errorf("expected 1 result, got %d", len(results))
    }
}
```

#### 4. Document the Rule

Update `docs/LINT_RULES.md` with the new rule documentation.

---

## Key Files Reference

### Internal Packages

| Package | Key Files | Purpose |
|---------|-----------|---------|
| `internal/discovery` | `discovery.go`, `ast.go` | AST-based resource discovery |
| `internal/discovery` | `board.go`, `slo.go`, `trigger.go` | Type-specific discovery |
| `internal/serialize` | `serialize.go` | Query JSON serialization |
| `internal/serialize` | `board.go`, `slo.go`, `trigger.go` | Type-specific serialization |
| `internal/lint` | `lint.go` | Lint engine and result handling |
| `internal/lint` | `rules.go` | Query lint rules |
| `internal/lint` | `board_rules.go`, `slo_rules.go`, `trigger_rules.go` | Type-specific rules |
| `internal/builder` | `builder.go` | Build pipeline orchestration |
| `internal/builder` | `registry.go` | Resource storage and deduplication |

### Public Packages

| Package | Purpose |
|---------|---------|
| `query` | Query type definitions and builder functions |
| `board` | Board and panel type definitions |
| `slo` | SLO, SLI, and burn alert definitions |
| `trigger` | Trigger, threshold, and recipient definitions |

### CLI

| File | Purpose |
|------|---------|
| `cmd/wetwire-honeycomb/main.go` | CLI entry point, command definitions |

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

- [CLI Reference](CLI.md) - Complete command documentation
- [Lint Rules](LINT_RULES.md) - All WHC lint rules
- [FAQ](FAQ.md) - Common questions
