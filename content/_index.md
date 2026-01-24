---
title: "Wetwire Honeycomb"
---

[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-honeycomb-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-honeycomb-go)
[![CI](https://github.com/lex00/wetwire-honeycomb-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-honeycomb-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-honeycomb-go/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-honeycomb-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-honeycomb-go)](https://goreportcard.com/report/github.com/lex00/wetwire-honeycomb-go)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate Honeycomb queries, boards, and SLOs from Go structs with AI-assisted design.

## Philosophy

Wetwire uses typed constraints to reduce the model capability required for accurate code generation.

**Core hypothesis:** Typed input + smaller model ≈ Semantic input + larger model

The type system and lint rules act as a force multiplier — cheaper models produce quality output when guided by schema-generated types and iterative lint feedback.

## Documentation

| Document | Description |
|----------|-------------|
| [CLI Reference]({{< relref "/cli" >}}) | Command-line interface |
| [Quick Start]({{< relref "/quick-start" >}}) | Get started in 5 minutes |
| [Examples]({{< relref "/examples" >}}) | Sample observability projects |
| [FAQ]({{< relref "/faq" >}}) | Frequently asked questions |

## Installation

```bash
go install github.com/lex00/wetwire-honeycomb-go@latest
```

## Quick Example

```go
var LatencyQuery = query.Query{
    Dataset:     "production",
    Calculation: query.P99("duration_ms"),
    GroupBy:     []string{"service.name"},
}
```
