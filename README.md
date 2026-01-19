# wetwire-honeycomb-go

[![CI](https://github.com/lex00/wetwire-honeycomb-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-honeycomb-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-honeycomb-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-honeycomb-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-honeycomb-go)](https://goreportcard.com/report/github.com/lex00/wetwire-honeycomb-go)
[![Coverage](https://codecov.io/gh/lex00/wetwire-honeycomb-go/branch/main/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-honeycomb-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

Or add as a dependency:

```bash
go get github.com/lex00/wetwire-honeycomb-go
```

## Quick example

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
    Dataset:      "production",
    TimeRange:    query.Hours(2),
    Breakdowns:   query.Breakdown("endpoint", "service"),
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

Build to JSON:

```bash
wetwire-honeycomb build ./queries -o output.json
```

## Commands

| Command | Description |
|---------|-------------|
| `build` | Synthesize queries to Honeycomb Query JSON |
| `lint` | Check for issues and validate queries |
| `import` | Convert existing JSON to Go code |
| `validate` | Validate query structure |
| `list` | List all discovered queries |
| `init` | Scaffold a new project |
| `graph` | Visualize query dependencies |

## AI-Assisted Design

Use the `design` command for interactive, AI-assisted query creation:

```bash
# No API key required - uses Claude CLI
wetwire-honeycomb design "Create a query to find slow API requests"
```

Uses [Claude CLI](https://claude.ai/download) by default (no API key required). Falls back to Anthropic API if Claude CLI is not installed.

## Documentation

**Getting Started:**
- [Quick Start](docs/QUICK_START.md) - 5-minute tutorial
- [FAQ](docs/FAQ.md) - Common questions

**Reference:**
- [CLI Reference](docs/CLI.md) - All commands
- [Lint Rules](docs/LINT_RULES.md) - WHC rule reference

## Part of wetwire

This package follows the [wetwire specification](https://github.com/lex00/wetwire) for declarative infrastructure-as-code.

## License

MIT License - see [LICENSE](LICENSE)
