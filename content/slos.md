---
title: "SLOs"
---

This document describes how to define type-safe Service Level Objectives (SLOs) that compile to Honeycomb SLO JSON.

## Overview

SLO synthesis provides a type-safe way to declare Service Level Objectives in Go that compile to Honeycomb SLO JSON format.

### Key Features

- Type-safe SLO declarations using Go structs
- Direct query references (no string IDs)
- Builder functions for targets and time periods
- Pre-configured burn alert helpers (FastBurn, SlowBurn)
- Automatic JSON serialization
- AST-based discovery (no manual registration)

### Basic Example

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

var APIAvailability = slo.SLO{
    Name:        "API Availability",
    Description: "99.9% of requests succeed",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  SuccessRequests,
        TotalEvents: AllRequests,
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
}
```

## SLO Structure

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | `string` | Yes | Display name of the SLO |
| `Description` | `string` | No | Additional context about the SLO |
| `Dataset` | `string` | Yes | Honeycomb dataset this SLO measures |
| `SLI` | `slo.SLI` | Yes | Service Level Indicator definition |
| `Target` | `slo.Target` | Yes | SLO target percentage |
| `TimePeriod` | `slo.TimePeriod` | Yes | Rolling window for SLO calculation |
| `BurnAlerts` | `[]slo.BurnAlert` | No | Error budget burn alerts |

## Service Level Indicators (SLI)

An SLI defines the ratio of good events to total events that determines SLO compliance.

### SLI Structure

```go
type SLI struct {
    GoodEvents  query.Query  // Query counting successful events
    TotalEvents query.Query  // Query counting all events
}
```

### SLI Ratio

The SLI ratio is calculated as:

```
SLI = (GoodEvents / TotalEvents) * 100
```

### Example SLIs

#### Availability SLI

```go
var GoodEvents = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("http.status_code", 500),
    },
}

var TotalEvents = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

#### Latency SLI

```go
var FastRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("duration_ms", 500),
    },
}
```

## Target and TimePeriod

### Target

Use `slo.Percentage()` to specify the target:

```go
Target: slo.Percentage(99.9)   // 99.9% (three nines)
Target: slo.Percentage(99.99)  // 99.99% (four nines)
Target: slo.Percentage(95.0)   // 95%
```

#### Common Targets

| Target | Description | Downtime per 30 days |
|--------|-------------|---------------------|
| 90.0% | Basic availability | 72 hours |
| 95.0% | Good availability | 36 hours |
| 99.0% | High availability | 7.2 hours |
| 99.9% | Three nines | 43 minutes |
| 99.99% | Four nines | 4.3 minutes |

### TimePeriod

Use `slo.Days()` to specify the rolling window:

```go
TimePeriod: slo.Days(7)   // 7-day rolling window
TimePeriod: slo.Days(30)  // 30-day rolling window
TimePeriod: slo.Days(90)  // 90-day rolling window
```

## Burn Alert Configuration

Burn alerts notify when error budget is being consumed too quickly.

### Builder Functions

#### FastBurn

Detects rapid error budget consumption using a 1-hour window:

```go
slo.FastBurn(2.0)  // Alert if 2% of budget burned in 1 hour
```

#### SlowBurn

Detects gradual error budget consumption using a 24-hour window:

```go
slo.SlowBurn(5.0)  // Alert if 5% of budget burned in 24 hours
```

### Recipients

Configure notification targets for alerts:

```go
Recipients: []slo.Recipient{
    {Type: "pagerduty", Target: "critical-service"},
    {Type: "slack", Target: "#incidents"},
    {Type: "email", Target: "team@example.com"},
}
```

### Recommended Configurations

For 99.9% SLO:
```go
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(2.0),   // 2% in 1 hour (urgent)
    slo.SlowBurn(5.0),   // 5% in 24 hours (warning)
}
```

For 99.99% SLO:
```go
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(1.0),   // 1% in 1 hour (critical)
    slo.SlowBurn(2.0),   // 2% in 24 hours (urgent)
}
```

## Complete Example

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

var GoodRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("http.status_code", 500),
    },
}

var TotalRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

var APIAvailability = slo.SLO{
    Name:        "API Availability",
    Description: "99.9% of API requests succeed",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  GoodRequests,
        TotalEvents: TotalRequests,
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(2.0),
        slo.SlowBurn(5.0),
    },
}
```

## Best Practices

1. **Organize SLOs by Service** - Group related SLOs in dedicated packages
2. **Use Descriptive Names** - Clear names like `AuthServiceAvailability`
3. **Include Meaningful Descriptions** - Explain what and why
4. **Reuse Query Definitions** - Define queries once, reuse in multiple SLOs
5. **Match Time Periods to Use Case** - 7 days for fast feedback, 30 days for monthly reporting
6. **Configure Appropriate Burn Alerts** - More sensitive for critical services
7. **Route Alerts Appropriately** - FastBurn to PagerDuty, SlowBurn to Slack

## See Also

- [CLI Documentation](../cli/) - Command reference
- [Lint Rules](../lint-rules/) - SLO lint rules (WHC040+)
- [FAQ](../faq/) - Common questions
