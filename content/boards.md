---
title: "Boards"
---

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

## Panel Options

### WithTitle

Sets the title displayed above the panel.

```go
board.WithTitle(title string) PanelOption
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
```

## Preset Filters

Preset filters apply board-level filtering to all query panels.

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

## Tags

Tags are key-value pairs used for board organization and discovery.

### Common Tag Keys

- `team` - Owning team name
- `service` - Service name
- `environment` - Environment (prod, staging, dev)
- `category` - Functional category (performance, errors, business)
- `owner` - Email or identifier of board maintainer
- `tier` - Service tier (critical, high, medium, low)
- `type` - Board type (monitoring, debugging, business-metrics)

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

### Tags

1. **Consistent keys** - Establish tag key conventions across organization
2. **Ownership clarity** - Always tag with team/owner
3. **Service mapping** - Tag boards with associated services
4. **Searchability** - Use tags that make boards easy to find

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

## See Also

- [Query Documentation](https://github.com/lex00/wetwire-honeycomb-go#readme) - Query construction
- [CLI Documentation](../cli/) - Command reference
- [Lint Rules](../lint-rules/) - Board validation rules
