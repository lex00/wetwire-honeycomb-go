# wetwire-honeycomb-go

Honeycomb query synthesis for Go - type-safe declarations that compile to Honeycomb Query JSON.

## Overview

wetwire-honeycomb is a **synthesis library** - it generates Honeycomb Query JSON from typed Go declarations. It does not execute queries or manage state.

```
Go Structs → wetwire-honeycomb build → Query JSON → Honeycomb API
                                            ↓
                                   (user's responsibility)
```

## Installation

```bash
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
```

## Quick Example

```go
package main

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
)

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

## Commands

```bash
wetwire-honeycomb build ./queries/...  # Synthesize to Query JSON
wetwire-honeycomb lint ./queries/...   # Check for issues
wetwire-honeycomb list                 # List all queries
```

## Status

🚧 Under development - see [ROADMAP](https://github.com/lex00/wetwire-honeycomb-go/issues/18)
