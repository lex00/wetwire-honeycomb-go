# Boards

Honeycomb boards are dashboards that combine multiple visualizations (panels) to provide a comprehensive view of your system's health and performance. The `board` package provides type-safe board declarations that compile to Honeycomb Board JSON.

## Overview

Boards in wetwire-honeycomb-go are declared as top-level variables and discovered automatically via AST analysis. A board consists of:

- Basic metadata (name, description)
- One or more panels (query visualizations, text notes, SLO displays)
- Optional preset filters applied to all query panels
- Optional tags for organization and categorization

## Quick Start

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "github.com/lex00/wetwire-honeycomb-go/query"
)

var ServiceOverview = board.Board{
    Name:        "Service Performance Overview",
    Description: "Real-time monitoring of API performance",
    Panels: []board.Panel{
        board.QueryPanel(RequestLatency,
            board.WithTitle("Request Latency P99"),
            board.WithPosition(0, 0, 6, 4),
        ),
        board.TextPanel("## Monitor during peak hours (9am-5pm EST)",
            board.WithTitle("Notes"),
            board.WithPosition(6, 0, 6, 2),
        ),
    },
    PresetFilters: []board.Filter{
        {Column: "service.name", Operation: "=", Value: "api"},
        {Column: "environment", Operation: "=", Value: "production"},
    },
    Tags: []board.Tag{
        {Key: "team", Value: "platform"},
        {Key: "service", Value: "api"},
    },
}
```

## Board Structure

### Board Type

```go
type Board struct {
    Name          string
    Description   string
    Panels        []Panel
    PresetFilters []Filter
    Tags          []Tag
}
```

#### Fields

**Name** (string, required)
- The display name of the board
- Shows in Honeycomb UI and board listings
- Should be descriptive and unique

**Description** (string, optional)
- Additional context about the board's purpose
- Displayed in the Honeycomb UI
- Use to document what the board monitors or when to use it

**Panels** ([]Panel, optional)
- The visual components that make up the board
- Can be query panels, text panels, or SLO panels
- Positioned using a grid layout system

**PresetFilters** ([]Filter, optional)
- Board-level filters applied to all query panels
- Useful for filtering entire board to specific environment, service, etc.
- Does not affect text or SLO panels

**Tags** ([]Tag, optional)
- Key-value metadata for organizing boards
- Used for categorization, ownership, and discovery
- Searchable in Honeycomb UI

## Panel Types

Boards support three types of panels:

### 1. Query Panel

Displays the results of a Honeycomb query visualization.

```go
board.QueryPanel(query.Query, ...PanelOption)
```

**Example:**

```go
// Define a query
var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Calculations: []query.Calculation{
        query.Count(),
        query.P99("duration_ms"),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 1000),
    },
    Breakdowns: []string{"endpoint"},
}

// Use in a board
Panels: []board.Panel{
    board.QueryPanel(SlowRequests,
        board.WithTitle("Slow Requests (>1s)"),
        board.WithPosition(0, 0, 12, 4),
    ),
}
```

**Notes:**
- Query must be a complete `query.Query` struct
- Query can be defined inline or referenced from another variable
- Preset filters from the board are automatically applied to the query

### 2. Text Panel

Displays markdown-formatted text for documentation and notes.

```go
board.TextPanel(content string, ...PanelOption)
```

**Example:**

```go
Panels: []board.Panel{
    board.TextPanel(`
## Dashboard Notes

This dashboard monitors API performance metrics.

**Alert thresholds:**
- P99 latency > 1000ms
- Error rate > 1%
- Success rate < 99%

**On-call:** Check PagerDuty rotation
    `,
        board.WithTitle("Documentation"),
        board.WithPosition(0, 8, 12, 3),
    ),
}
```

**Notes:**
- Content supports full markdown syntax
- Useful for documenting runbooks, alert thresholds, and context
- Does not respond to board preset filters

### 3. SLO Panel

Displays a Service Level Objective panel by referencing an SLO ID.

```go
board.SLOPanelByID(id string, ...PanelOption)
```

**Example:**

```go
Panels: []board.Panel{
    board.SLOPanelByID("api-availability-99-9",
        board.WithTitle("API Availability SLO"),
        board.WithPosition(6, 4, 6, 4),
    ),
}
```

**Notes:**
- References an existing SLO by its ID
- SLO must exist in Honeycomb (managed outside wetwire)
- Useful for displaying SLO burn rates and budget remaining
- Does not respond to board preset filters

## Panel Options

Panel options are functions that configure panel display properties.

### WithTitle

Sets the title displayed above the panel.

```go
board.WithTitle(title string) PanelOption
```

**Example:**

```go
board.QueryPanel(query,
    board.WithTitle("Request Rate"),
)
```

### WithPosition

Sets the position and size of the panel on the board grid.

```go
board.WithPosition(x, y, width, height int) PanelOption
```

**Parameters:**
- `x`: Horizontal position (0 = left edge)
- `y`: Vertical position (0 = top)
- `width`: Panel width in grid units (typical max: 12)
- `height`: Panel height in grid units

**Example:**

```go
board.QueryPanel(query,
    board.WithPosition(0, 0, 6, 4),  // Left half, 4 units tall
)

board.QueryPanel(anotherQuery,
    board.WithPosition(6, 0, 6, 4),  // Right half, same height
)
```

**Grid Layout:**
- Boards typically use a 12-column grid system
- Position (0, 0) is the top-left corner
- Width of 12 spans the full board width
- Width of 6 takes up half the board width
- Height units are flexible and determine vertical size

**Common Layouts:**

```go
// Full width panel
board.WithPosition(0, 0, 12, 4)

// Two panels side-by-side
board.WithPosition(0, 0, 6, 4)  // Left
board.WithPosition(6, 0, 6, 4)  // Right

// Three columns
board.WithPosition(0, 0, 4, 4)  // Left
board.WithPosition(4, 0, 4, 4)  // Center
board.WithPosition(8, 0, 4, 4)  // Right

// Stacked panels
board.WithPosition(0, 0, 12, 3)  // Top
board.WithPosition(0, 3, 12, 3)  // Middle
board.WithPosition(0, 6, 12, 3)  // Bottom
```

### Multiple Options

Options can be combined using the variadic pattern:

```go
board.QueryPanel(query,
    board.WithTitle("Error Rate"),
    board.WithPosition(0, 4, 6, 4),
)
```

## Preset Filters

Preset filters apply board-level filtering to all query panels. They're useful for creating environment-specific, service-specific, or team-specific views.

### Filter Structure

```go
type Filter struct {
    Column    string
    Operation string
    Value     any
}
```

### Supported Operations

- `=` - Equals
- `!=` - Not equals
- `>` - Greater than
- `>=` - Greater than or equal
- `<` - Less than
- `<=` - Less than or equal
- `contains` - Substring match
- `does-not-contain` - Negative substring match
- `exists` - Field exists
- `does-not-exist` - Field does not exist
- `starts-with` - Prefix match
- `in` - Value in list
- `not-in` - Value not in list

### Examples

**Single Environment Filter:**

```go
PresetFilters: []board.Filter{
    {Column: "environment", Operation: "=", Value: "production"},
}
```

**Multiple Filters:**

```go
PresetFilters: []board.Filter{
    {Column: "service.name", Operation: "=", Value: "api"},
    {Column: "status_code", Operation: ">=", Value: 400},
    {Column: "user_id", Operation: "exists", Value: nil},
}
```

**Team-specific Board:**

```go
PresetFilters: []board.Filter{
    {Column: "team", Operation: "=", Value: "platform"},
    {Column: "environment", Operation: "in", Value: []string{"staging", "production"}},
}
```

**High-value Traffic:**

```go
PresetFilters: []board.Filter{
    {Column: "plan_type", Operation: "=", Value: "enterprise"},
    {Column: "request_size", Operation: ">", Value: 1000000},
}
```

### Notes

- Preset filters are applied to query panels only
- Text and SLO panels ignore preset filters
- Filters use AND logic (all must match)
- Use `Value: nil` for exists/does-not-exist operations

## Tags

Tags are key-value pairs used for board organization and discovery.

### Tag Structure

```go
type Tag struct {
    Key   string
    Value string
}
```

### Examples

**Team Ownership:**

```go
Tags: []board.Tag{
    {Key: "team", Value: "platform"},
    {Key: "owner", Value: "platform-team@example.com"},
}
```

**Service Classification:**

```go
Tags: []board.Tag{
    {Key: "service", Value: "api"},
    {Key: "tier", Value: "critical"},
    {Key: "environment", Value: "production"},
}
```

**Functional Grouping:**

```go
Tags: []board.Tag{
    {Key: "category", Value: "performance"},
    {Key: "type", Value: "monitoring"},
    {Key: "severity", Value: "p1"},
}
```

### Tag Conventions

Common tag keys used for organization:

- `team` - Owning team name
- `service` - Service name
- `environment` - Environment (prod, staging, dev)
- `category` - Functional category (performance, errors, business)
- `owner` - Email or identifier of board maintainer
- `tier` - Service tier (critical, high, medium, low)
- `type` - Board type (monitoring, debugging, business-metrics)

## Complete Examples

### Example 1: API Performance Dashboard

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "github.com/lex00/wetwire-honeycomb-go/query"
)

// Query definitions
var RequestRate = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Breakdowns: []string{"endpoint"},
}

var LatencyPercentiles = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P50("duration_ms"),
        query.P95("duration_ms"),
        query.P99("duration_ms"),
    },
    Breakdowns: []string{"endpoint"},
}

var ErrorRate = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.GTE("status_code", 400),
    },
    Breakdowns: []string{"status_code"},
}

// Board definition
var APIPerformance = board.Board{
    Name:        "API Performance Dashboard",
    Description: "Real-time monitoring of API latency, throughput, and errors",
    Panels: []board.Panel{
        board.QueryPanel(RequestRate,
            board.WithTitle("Request Rate"),
            board.WithPosition(0, 0, 6, 4),
        ),
        board.QueryPanel(LatencyPercentiles,
            board.WithTitle("Latency Percentiles"),
            board.WithPosition(6, 0, 6, 4),
        ),
        board.QueryPanel(ErrorRate,
            board.WithTitle("Error Rate"),
            board.WithPosition(0, 4, 12, 4),
        ),
        board.TextPanel(`
## API Performance SLOs

**Latency:** P99 < 1000ms
**Availability:** 99.9% success rate
**Throughput:** > 1000 req/sec

Contact: api-team@example.com
        `,
            board.WithTitle("SLO Targets"),
            board.WithPosition(0, 8, 12, 2),
        ),
    },
    PresetFilters: []board.Filter{
        {Column: "service.name", Operation: "=", Value: "api"},
        {Column: "environment", Operation: "=", Value: "production"},
    },
    Tags: []board.Tag{
        {Key: "team", Value: "api"},
        {Key: "service", Value: "api"},
        {Key: "tier", Value: "critical"},
    },
}
```

### Example 2: Multi-Service Overview

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "github.com/lex00/wetwire-honeycomb-go/query"
)

var ServiceHealthQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
        query.P99("duration_ms"),
    },
    Filters: []query.Filter{
        query.GTE("status_code", 500),
    },
    Breakdowns: []string{"service.name"},
}

var ServiceOverview = board.Board{
    Name:        "Multi-Service Health Overview",
    Description: "System-wide health monitoring across all services",
    Panels: []board.Panel{
        board.QueryPanel(ServiceHealthQuery,
            board.WithTitle("Service Errors (5xx)"),
            board.WithPosition(0, 0, 8, 5),
        ),
        board.SLOPanelByID("api-availability",
            board.WithTitle("API SLO"),
            board.WithPosition(8, 0, 4, 5),
        ),
        board.SLOPanelByID("database-latency",
            board.WithTitle("Database SLO"),
            board.WithPosition(8, 5, 4, 5),
        ),
        board.TextPanel(`
## Incident Response

1. Check error rate across services
2. Review SLO burn rates
3. Escalate if P1 threshold exceeded

**Escalation:** PagerDuty @platform-oncall
        `,
            board.WithTitle("Runbook"),
            board.WithPosition(0, 5, 8, 5),
        ),
    },
    PresetFilters: []board.Filter{
        {Column: "environment", Operation: "=", Value: "production"},
    },
    Tags: []board.Tag{
        {Key: "team", Value: "platform"},
        {Key: "type", Value: "monitoring"},
        {Key: "scope", Value: "system-wide"},
    },
}
```

### Example 3: User Experience Dashboard

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "github.com/lex00/wetwire-honeycomb-go/query"
)

var PageLoadTime = query.Query{
    Dataset:   "frontend",
    TimeRange: query.Hours(2),
    Calculations: []query.Calculation{
        query.P50("page_load_ms"),
        query.P95("page_load_ms"),
        query.P99("page_load_ms"),
    },
    Breakdowns: []string{"page_name"},
    Orders: []query.Order{
        {Op: "P99", Order: "descending"},
    },
    Limit: 10,
}

var JSErrors = query.Query{
    Dataset:   "frontend",
    TimeRange: query.Hours(2),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Exists("error.message"),
    },
    Breakdowns: []string{"page_name", "error.type"},
}

var UserSessions = query.Query{
    Dataset:   "frontend",
    TimeRange: query.Hours(2),
    Calculations: []query.Calculation{
        query.CountDistinct("session_id"),
    },
    Breakdowns: []string{"browser", "device_type"},
}

var UserExperience = board.Board{
    Name:        "User Experience Metrics",
    Description: "Frontend performance and error tracking",
    Panels: []board.Panel{
        board.QueryPanel(PageLoadTime,
            board.WithTitle("Page Load Times (Top 10 Slowest)"),
            board.WithPosition(0, 0, 12, 4),
        ),
        board.QueryPanel(JSErrors,
            board.WithTitle("JavaScript Errors"),
            board.WithPosition(0, 4, 6, 4),
        ),
        board.QueryPanel(UserSessions,
            board.WithTitle("Active Sessions by Browser/Device"),
            board.WithPosition(6, 4, 6, 4),
        ),
        board.TextPanel(`
## Performance Budget

- **Page Load:** P95 < 3000ms
- **Time to Interactive:** P95 < 5000ms
- **JS Error Rate:** < 0.1% of page views

Monitor during peak traffic: 12pm-2pm, 6pm-9pm EST
        `,
            board.WithTitle("Performance Targets"),
            board.WithPosition(0, 8, 12, 2),
        ),
    },
    PresetFilters: []board.Filter{
        {Column: "environment", Operation: "=", Value: "production"},
        {Column: "bot", Operation: "!=", Value: true},
    },
    Tags: []board.Tag{
        {Key: "team", Value: "frontend"},
        {Key: "category", Value: "user-experience"},
        {Key: "service", Value: "web-app"},
    },
}
```

## Using Template Functions

Boards can reference queries defined elsewhere using Go variables and template functions from the `query` package.

### Time Range Functions

```go
// Relative time ranges
query.Hours(2)        // Last 2 hours
query.Minutes(30)     // Last 30 minutes
query.Days(7)         // Last 7 days
query.Seconds(300)    // Last 300 seconds

// Convenience functions
query.Last24Hours()   // Last 24 hours
query.Last7Days()     // Last 7 days

// Absolute time ranges
query.Absolute(startTime, endTime)  // time.Time values
```

### Calculation Functions

```go
// Basic aggregations
query.Count()
query.Sum("column")
query.Avg("column")
query.Max("column")
query.Min("column")

// Percentiles
query.P50("column")   // Median
query.P75("column")
query.P90("column")
query.P95("column")
query.P99("column")
query.P999("column")

// Rates
query.Rate("column")
query.RateSum("column")
query.RateAvg("column")
query.RateMax("column")

// Special
query.Heatmap("column")
query.Concurrency()
query.CountDistinct("column")
```

### Filter Functions

```go
// Comparison
query.Equals("column", value)      // or query.Eq()
query.NotEquals("column", value)   // or query.Ne()
query.GreaterThan("column", value) // or query.GT()
query.GreaterThanOrEqual("column", value) // or query.GTE()
query.LessThan("column", value)    // or query.LT()
query.LessThanOrEqual("column", value)    // or query.LTE()

// String operations
query.Contains("column", value)
query.DoesNotContain("column", value)
query.StartsWith("column", value)

// Existence
query.Exists("column")
query.DoesNotExist("column")

// List operations
query.In("column", []any{"value1", "value2"})
query.NotIn("column", []any{"value1", "value2"})
```

### Composing Queries

Queries can be composed and reused across multiple boards:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// Shared query definitions
var BaseLatencyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P50("duration_ms"),
        query.P95("duration_ms"),
        query.P99("duration_ms"),
    },
}

var SlowRequestsFilter = []query.Filter{
    query.GT("duration_ms", 1000),
}
```

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "yourorg/queries"
)

var PerformanceBoard = board.Board{
    Name: "Performance",
    Panels: []board.Panel{
        board.QueryPanel(queries.BaseLatencyQuery,
            board.WithTitle("Latency"),
        ),
    },
}
```

## JSON Output Format

Boards compile to Honeycomb Board JSON format. Here's the structure:

### Basic Board JSON

```json
{
  "name": "Service Performance",
  "description": "Real-time monitoring of API performance",
  "panels": [...],
  "preset_filters": [...],
}
```

### Query Panel JSON

```json
{
  "type": "query",
  "title": "Request Rate",
  "position": {
    "x": 0,
    "y": 0,
    "width": 6,
    "height": 4
  },
  "query": {
    "dataset": "production",
    "time_range": 3600,
    "calculations": [
      {"op": "COUNT"}
    ],
    "breakdowns": ["endpoint"]
  }
}
```

### Text Panel JSON

```json
{
  "type": "text",
  "title": "Notes",
  "content": "## Dashboard Notes\nMonitor during peak hours",
  "position": {
    "x": 6,
    "y": 0,
    "width": 6,
    "height": 2
  }
}
```

### SLO Panel JSON

```json
{
  "type": "slo",
  "title": "API SLO",
  "slo_id": "api-availability-99-9",
  "position": {
    "x": 0,
    "y": 4,
    "width": 12,
    "height": 4
  }
}
```

### Preset Filters JSON

```json
{
  "preset_filters": [
    {
      "column": "service.name",
      "op": "=",
      "value": "api"
    },
    {
      "column": "status_code",
      "op": ">=",
      "value": 400
    }
  ]
}
```

### Complete Example JSON

```json
{
  "name": "API Performance Dashboard",
  "description": "Real-time monitoring of API latency, throughput, and errors",
  "panels": [
    {
      "type": "query",
      "title": "Request Rate",
      "position": {"x": 0, "y": 0, "width": 6, "height": 4},
      "query": {
        "dataset": "production",
        "time_range": 3600,
        "calculations": [{"op": "COUNT"}],
        "breakdowns": ["endpoint"]
      }
    },
    {
      "type": "text",
      "title": "SLO Targets",
      "content": "## API Performance SLOs\n\n**Latency:** P99 < 1000ms",
      "position": {"x": 0, "y": 4, "width": 12, "height": 2}
    },
    {
      "type": "slo",
      "title": "API SLO",
      "slo_id": "api-availability",
      "position": {"x": 6, "y": 0, "width": 6, "height": 4}
    }
  ],
  "preset_filters": [
    {
      "column": "environment",
      "op": "=",
      "value": "production"
    }
  ]
}
```

## Best Practices

### Board Organization

1. **One board per concern** - Create focused boards for specific services or use cases
2. **Logical panel layout** - Place related panels near each other
3. **Use text panels** - Document runbooks, SLO targets, and escalation procedures
4. **Consistent naming** - Use clear, descriptive board names

### Panel Layout

1. **Top-to-bottom flow** - Most important metrics at the top
2. **Grid alignment** - Use consistent widths (6 or 12 units common)
3. **Reasonable heights** - 4 units typical for charts, 2-3 for text
4. **Related panels grouped** - Keep similar metrics visually close

### Preset Filters

1. **Environment isolation** - Filter to single environment in production boards
2. **Service focus** - Use service filters for service-specific boards
3. **Minimal filters** - Only filter what's necessary, allow panel queries to be flexible
4. **Document filters** - Use text panels to explain what preset filters are applied

### Tags

1. **Consistent keys** - Establish tag key conventions across organization
2. **Ownership clarity** - Always tag with team/owner
3. **Service mapping** - Tag boards with associated services
4. **Searchability** - Use tags that make boards easy to find

### Query Reuse

1. **Define once** - Create query variables and reuse across boards
2. **Shared library** - Keep common queries in a shared package
3. **Naming convention** - Use clear query variable names
4. **Documentation** - Document complex queries with comments

## CLI Commands

### Build Boards to JSON

```bash
wetwire-honeycomb build --boards ./boards/...
```

### List All Boards

```bash
wetwire-honeycomb list --boards
```

### Lint Boards

```bash
wetwire-honeycomb lint --boards ./boards/...
```

## Troubleshooting

### Board not discovered

**Problem:** Board declared but not found by CLI

**Solution:**
- Ensure board is a top-level `var` declaration
- Check that package is in the path specified to CLI
- Verify board type is `board.Board`

### Panel not rendering in JSON

**Problem:** Panel appears in board but missing from JSON output

**Solution:**
- Verify panel is added to `Panels` slice
- Check that query referenced in QueryPanel is valid
- Ensure panel options are applied correctly

### Preset filters not working

**Problem:** Filters don't seem to apply to queries

**Solution:**
- Verify filter column names match your dataset schema
- Check operation syntax (use `=` not `==`, `contains` not `CONTAINS`)
- Remember filters only apply to query panels, not text/SLO panels

### Position/layout issues

**Problem:** Panels overlap or don't appear where expected

**Solution:**
- Check X coordinates don't exceed 12 (standard grid width)
- Ensure Y coordinates don't conflict (panels at same Y will overlap)
- Verify width + x doesn't exceed 12 for panels on same row
- Use height consistently (4 is typical for most visualizations)

## Next Steps

- See [Query Documentation](../README.md) for query construction
- See [CLI Documentation](CLI.md) for command reference
- See [Lint Rules](LINT_RULES.md) for board validation rules
- Explore example boards in `examples/boards/`
