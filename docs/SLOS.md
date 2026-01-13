# SLO Synthesis

This document describes how to define type-safe Service Level Objectives (SLOs) that compile to Honeycomb SLO JSON.

## Table of Contents

1. [Overview](#overview)
2. [SLO Structure](#slo-structure)
3. [Service Level Indicators (SLI)](#service-level-indicators-sli)
4. [Direct Query References](#direct-query-references)
5. [Target and TimePeriod](#target-and-timeperiod)
6. [Burn Alert Configuration](#burn-alert-configuration)
7. [Complete Examples](#complete-examples)
8. [JSON Output Format](#json-output-format)
9. [Best Practices](#best-practices)

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

The `slo.SLO` type represents a complete SLO specification.

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

### Example

```go
var LatencySLO = slo.SLO{
    Name:        "API Latency",
    Description: "95% of requests complete within 500ms",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  FastRequests,
        TotalEvents: AllRequests,
    },
    Target:     slo.Percentage(95.0),
    TimePeriod: slo.Days(7),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(2.0),
        slo.SlowBurn(5.0),
    },
}
```

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

For a 99.9% availability SLO:
- If 999 out of 1000 requests succeed, SLI = 99.9% (meets target)
- If 990 out of 1000 requests succeed, SLI = 99.0% (below target)

### Defining SLIs

#### Availability SLI

Measures the percentage of successful requests:

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

var APIAvailability = slo.SLO{
    Name:    "API Availability",
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  GoodEvents,
        TotalEvents: TotalEvents,
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
}
```

#### Latency SLI

Measures the percentage of requests completing within a time threshold:

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

var AllRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

var LatencySLO = slo.SLO{
    Name:    "API Latency",
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  FastRequests,
        TotalEvents: AllRequests,
    },
    Target:     slo.Percentage(95.0),
    TimePeriod: slo.Days(7),
}
```

#### Error Rate SLI

Measures the percentage of requests without errors:

```go
var SuccessfulRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.NotExists("error"),
    },
}

var AllRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

var ErrorRateSLO = slo.SLO{
    Name:    "Low Error Rate",
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  SuccessfulRequests,
        TotalEvents: AllRequests,
    },
    Target:     slo.Percentage(99.5),
    TimePeriod: slo.Days(30),
}
```

## Direct Query References

SLOs use direct Go variable references for queries, not string IDs.

### Pattern

```go
// Define queries as top-level vars
var GoodEvents = query.Query{...}
var TotalEvents = query.Query{...}

// Reference queries directly
var MySLO = slo.SLO{
    SLI: slo.SLI{
        GoodEvents:  GoodEvents,   // Direct reference
        TotalEvents: TotalEvents,  // Direct reference
    },
}
```

### Benefits

- **Type safety**: Compiler ensures queries exist
- **Refactoring support**: Rename refactoring works across references
- **No string IDs**: No risk of typos or stale references
- **Clear dependencies**: Easy to see what queries an SLO uses

### Example

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

// Define queries
var SuccessRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("http.status_code", 500),
    },
}

var AllRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

// Reference queries directly in SLO
var APIAvailability = slo.SLO{
    Name:    "API Availability",
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  SuccessRequests,  // Direct reference
        TotalEvents: AllRequests,       // Direct reference
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
}
```

## Target and TimePeriod

SLOs require a target percentage and a rolling time window.

### Target

Use `slo.Percentage()` to specify the target:

```go
type Target struct {
    Percentage float64  // Target value (e.g., 99.9 for 99.9%)
}

func Percentage(p float64) Target
```

#### Examples

```go
Target: slo.Percentage(99.9)   // 99.9% (three nines)
Target: slo.Percentage(99.99)  // 99.99% (four nines)
Target: slo.Percentage(95.0)   // 95%
Target: slo.Percentage(99.5)   // 99.5%
```

#### Common Targets

| Target | Description | Downtime per 30 days |
|--------|-------------|---------------------|
| 90.0% | Basic availability | 72 hours |
| 95.0% | Good availability | 36 hours |
| 99.0% | High availability | 7.2 hours |
| 99.5% | Very high availability | 3.6 hours |
| 99.9% | Three nines | 43 minutes |
| 99.95% | Three and a half nines | 21 minutes |
| 99.99% | Four nines | 4.3 minutes |

### TimePeriod

Use `slo.Days()` to specify the rolling window:

```go
type TimePeriod struct {
    Days  int
    Hours int
}

func Days(d int) TimePeriod
```

#### Examples

```go
TimePeriod: slo.Days(7)   // 7-day rolling window
TimePeriod: slo.Days(14)  // 14-day rolling window
TimePeriod: slo.Days(30)  // 30-day rolling window
TimePeriod: slo.Days(90)  // 90-day rolling window
```

#### Common Time Periods

| Period | Use Case |
|--------|----------|
| 7 days | Weekly sprint alignment, fast feedback |
| 14 days | Bi-weekly cycles |
| 30 days | Monthly reporting, standard practice |
| 90 days | Quarterly objectives |

### Complete Example

```go
var APIAvailability = slo.SLO{
    Name:       "API Availability",
    Dataset:    "production",
    SLI:        slo.SLI{...},
    Target:     slo.Percentage(99.9),  // 99.9% target
    TimePeriod: slo.Days(30),          // 30-day rolling window
}
```

## Burn Alert Configuration

Burn alerts notify when error budget is being consumed too quickly.

### BurnAlert Structure

```go
type BurnAlert struct {
    Name       string        // Display name
    AlertType  AlertType     // "exhaustion_time" or "budget_rate"
    Threshold  float64       // Trigger threshold
    Window     TimePeriod    // Time window for calculation
    Recipients []Recipient   // Notification targets
}
```

### Alert Types

```go
const (
    ExhaustionTime AlertType = "exhaustion_time"  // Alert when budget will be exhausted soon
    BudgetRate     AlertType = "budget_rate"      // Alert when burn rate is too high
)
```

### Builder Functions

Pre-configured helpers for common burn alert patterns:

#### FastBurn

Detects rapid error budget consumption using a 1-hour window:

```go
func FastBurn(budgetPercent float64) BurnAlert

// Example
slo.FastBurn(2.0)  // Alert if 2% of budget burned in 1 hour
```

Use FastBurn for:
- Immediate incidents requiring urgent response
- Short detection window (1 hour)
- High burn rate threshold (2-5%)

#### SlowBurn

Detects gradual error budget consumption using a 24-hour window:

```go
func SlowBurn(budgetPercent float64) BurnAlert

// Example
slo.SlowBurn(5.0)  // Alert if 5% of budget burned in 24 hours
```

Use SlowBurn for:
- Sustained degradation over time
- Longer detection window (24 hours)
- Lower burn rate threshold (5-10%)

### Recipients

Configure notification targets for alerts:

```go
type Recipient struct {
    Type   string  // "slack", "pagerduty", "email", "webhook"
    Target string  // Channel, service ID, email, or URL
}
```

#### Recipient Types

| Type | Target Format | Example |
|------|---------------|---------|
| `slack` | Channel name | `#alerts` |
| `pagerduty` | Service ID or key | `api-team` or `service-123` |
| `email` | Email address | `team@example.com` |
| `webhook` | Full URL | `https://example.com/webhook` |

### Examples

#### Basic Burn Alerts

```go
var APIAvailability = slo.SLO{
    Name:    "API Availability",
    Dataset: "production",
    SLI:     slo.SLI{...},
    Target:  slo.Percentage(99.9),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(2.0),   // Alert if 2% burned in 1 hour
        slo.SlowBurn(5.0),   // Alert if 5% burned in 24 hours
    },
}
```

#### Custom Burn Alert with Recipients

```go
var CriticalSLO = slo.SLO{
    Name:    "Critical Service",
    Dataset: "production",
    SLI:     slo.SLI{...},
    Target:  slo.Percentage(99.99),
    BurnAlerts: []slo.BurnAlert{
        {
            Name:      "Critical Fast Burn",
            AlertType: slo.BudgetRate,
            Threshold: 1.0,
            Window:    slo.TimePeriod{Hours: 1},
            Recipients: []slo.Recipient{
                {Type: "pagerduty", Target: "critical-service"},
                {Type: "slack", Target: "#incidents"},
            },
        },
        {
            Name:      "Warning Slow Burn",
            AlertType: slo.BudgetRate,
            Threshold: 10.0,
            Window:    slo.TimePeriod{Hours: 24},
            Recipients: []slo.Recipient{
                {Type: "slack", Target: "#monitoring"},
                {Type: "email", Target: "team@example.com"},
            },
        },
    },
}
```

#### Exhaustion Time Alert

```go
var APIAvailability = slo.SLO{
    Name:    "API Availability",
    Dataset: "production",
    SLI:     slo.SLI{...},
    Target:  slo.Percentage(99.9),
    BurnAlerts: []slo.BurnAlert{
        {
            Name:      "Budget Exhaustion Warning",
            AlertType: slo.ExhaustionTime,
            Threshold: 2.0,  // Alert if budget will be exhausted in < 2 days
            Window:    slo.Days(1),
            Recipients: []slo.Recipient{
                {Type: "slack", Target: "#alerts"},
            },
        },
    },
}
```

### Burn Alert Best Practices

1. **Use both fast and slow burn** - Catches immediate spikes and gradual degradation
2. **Set appropriate thresholds** - Too sensitive causes alert fatigue, too loose misses issues
3. **Match severity to recipients** - FastBurn to PagerDuty, SlowBurn to Slack
4. **Consider SLO target** - Higher targets (99.99%) need more sensitive alerts

#### Recommended Configurations

For 99.9% SLO (43 minutes downtime per month):
```go
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(2.0),   // 2% in 1 hour = ~12 minutes (urgent)
    slo.SlowBurn(5.0),   // 5% in 24 hours = ~2 hours (warning)
}
```

For 99.99% SLO (4.3 minutes downtime per month):
```go
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(1.0),   // 1% in 1 hour = ~2.5 minutes (critical)
    slo.SlowBurn(2.0),   // 2% in 24 hours = ~5 minutes (urgent)
}
```

## Complete Examples

### Basic Availability SLO

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

### Latency SLO

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

var FastRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("duration_ms", 500),
        query.Equals("endpoint", "/api/users"),
    },
}

var AllUserRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("endpoint", "/api/users"),
    },
}

var UserEndpointLatency = slo.SLO{
    Name:        "User Endpoint Latency",
    Description: "95% of /api/users requests complete within 500ms",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  FastRequests,
        TotalEvents: AllUserRequests,
    },
    Target:     slo.Percentage(95.0),
    TimePeriod: slo.Days(7),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(5.0),
        slo.SlowBurn(10.0),
    },
}
```

### Multi-Service SLO

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

var AuthSuccessful = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("service", "auth"),
        query.LT("http.status_code", 500),
        query.NotExists("error"),
    },
}

var AuthTotal = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("service", "auth"),
    },
}

var AuthServiceAvailability = slo.SLO{
    Name:        "Auth Service Availability",
    Description: "99.95% of authentication requests succeed",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  AuthSuccessful,
        TotalEvents: AuthTotal,
    },
    Target:     slo.Percentage(99.95),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        {
            Name:      "Auth Fast Burn",
            AlertType: slo.BudgetRate,
            Threshold: 1.0,
            Window:    slo.TimePeriod{Hours: 1},
            Recipients: []slo.Recipient{
                {Type: "pagerduty", Target: "auth-team"},
                {Type: "slack", Target: "#auth-alerts"},
            },
        },
        {
            Name:      "Auth Slow Burn",
            AlertType: slo.BudgetRate,
            Threshold: 3.0,
            Window:    slo.TimePeriod{Hours: 24},
            Recipients: []slo.Recipient{
                {Type: "slack", Target: "#auth-monitoring"},
            },
        },
    },
}
```

### Composite SLO with Multiple Conditions

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
)

var FastAndSuccessful = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.LT("duration_ms", 1000),
        query.LT("http.status_code", 500),
    },
    FilterCombination: "AND",
}

var AllAPIRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

var APIQuality = slo.SLO{
    Name:        "API Quality",
    Description: "99% of requests are both fast (<1s) and successful (<500)",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  FastAndSuccessful,
        TotalEvents: AllAPIRequests,
    },
    Target:     slo.Percentage(99.0),
    TimePeriod: slo.Days(7),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(3.0),
        slo.SlowBurn(7.0),
    },
}
```

## JSON Output Format

SLOs serialize to Honeycomb-compatible JSON format.

### Structure

```json
{
  "name": "string",
  "description": "string",
  "dataset": "string",
  "sli": {
    "good_events": {
      "time_range": 2592000,
      "calculations": [...],
      "filters": [...]
    },
    "total_events": {
      "time_range": 2592000,
      "calculations": [...]
    }
  },
  "target_per_million": 999000,
  "time_period_days": 30,
  "burn_alerts": [
    {
      "alert_type": "budget_rate",
      "threshold": 2.0,
      "window_hours": 1
    }
  ]
}
```

### Field Mappings

| Go Field | JSON Field | Conversion |
|----------|------------|------------|
| `Name` | `name` | Direct |
| `Description` | `description` | Direct |
| `Dataset` | `dataset` | Direct |
| `Target.Percentage` | `target_per_million` | `Percentage * 10000` |
| `TimePeriod.Days` | `time_period_days` | Direct |
| `SLI.GoodEvents` | `sli.good_events` | Query JSON |
| `SLI.TotalEvents` | `sli.total_events` | Query JSON |
| `BurnAlerts` | `burn_alerts` | Array |

### Example Output

For this SLO:

```go
var APIAvailability = slo.SLO{
    Name:        "API Availability",
    Description: "99.9% of requests succeed",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents: query.Query{
            Dataset:   "production",
            TimeRange: query.Days(30),
            Calculations: []query.Calculation{
                query.Count(),
            },
            Filters: []query.Filter{
                query.LT("http.status_code", 500),
            },
        },
        TotalEvents: query.Query{
            Dataset:   "production",
            TimeRange: query.Days(30),
            Calculations: []query.Calculation{
                query.Count(),
            },
        },
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(2.0),
        slo.SlowBurn(5.0),
    },
}
```

The JSON output is:

```json
{
  "name": "API Availability",
  "description": "99.9% of requests succeed",
  "dataset": "production",
  "sli": {
    "good_events": {
      "time_range": 2592000,
      "calculations": [
        {
          "op": "COUNT"
        }
      ],
      "filters": [
        {
          "column": "http.status_code",
          "op": "<",
          "value": 500
        }
      ]
    },
    "total_events": {
      "time_range": 2592000,
      "calculations": [
        {
          "op": "COUNT"
        }
      ]
    }
  },
  "target_per_million": 999000,
  "time_period_days": 30,
  "burn_alerts": [
    {
      "alert_type": "budget_rate",
      "threshold": 2.0,
      "window_hours": 1
    },
    {
      "alert_type": "budget_rate",
      "threshold": 5.0,
      "window_hours": 24
    }
  ]
}
```

### Generating JSON

Use the CLI to build SLO JSON:

```bash
# Build all SLOs
wetwire-honeycomb build ./slos/...

# Build specific SLO
wetwire-honeycomb build ./slos/api_availability.go

# Pretty-print JSON
wetwire-honeycomb build --format=pretty ./slos/...
```

## Best Practices

### 1. Organize SLOs by Service

```
slos/
├── auth/
│   ├── availability.go
│   └── latency.go
├── api/
│   ├── availability.go
│   ├── latency.go
│   └── quality.go
└── database/
    ├── query_performance.go
    └── connection_pool.go
```

### 2. Use Descriptive Names

```go
// GOOD: Clear, specific names
var AuthServiceAvailability = slo.SLO{...}
var UserEndpointLatency = slo.SLO{...}
var PaymentProcessingQuality = slo.SLO{...}

// BAD: Vague names
var SLO1 = slo.SLO{...}
var Test = slo.SLO{...}
var API = slo.SLO{...}
```

### 3. Include Meaningful Descriptions

```go
// GOOD: Explains what and why
Description: "99.9% of API requests succeed - critical for user experience"

// BAD: Restates the obvious
Description: "SLO for API"
```

### 4. Reuse Query Definitions

```go
// GOOD: Define queries once, reuse in multiple SLOs
var SuccessRequests = query.Query{...}
var TotalRequests = query.Query{...}

var AvailabilitySLO = slo.SLO{
    SLI: slo.SLI{
        GoodEvents:  SuccessRequests,
        TotalEvents: TotalRequests,
    },
}

var QualitySLO = slo.SLO{
    SLI: slo.SLI{
        GoodEvents:  HighQualityRequests,
        TotalEvents: TotalRequests,  // Reused
    },
}
```

### 5. Match Time Periods to Use Case

```go
// Short-term SLO for fast-changing services
TimePeriod: slo.Days(7)

// Standard SLO for monthly reporting
TimePeriod: slo.Days(30)

// Long-term SLO for strategic planning
TimePeriod: slo.Days(90)
```

### 6. Configure Appropriate Burn Alerts

```go
// For critical services (99.99%)
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(1.0),  // More sensitive
    slo.SlowBurn(2.0),
}

// For standard services (99.9%)
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(2.0),
    slo.SlowBurn(5.0),
}

// For lenient services (95%)
BurnAlerts: []slo.BurnAlert{
    slo.FastBurn(5.0),  // Less sensitive
    slo.SlowBurn(10.0),
}
```

### 7. Route Alerts Appropriately

```go
BurnAlerts: []slo.BurnAlert{
    {
        Name: "Critical Fast Burn",
        // ... config ...
        Recipients: []slo.Recipient{
            {Type: "pagerduty", Target: "oncall"},  // Urgent response
        },
    },
    {
        Name: "Warning Slow Burn",
        // ... config ...
        Recipients: []slo.Recipient{
            {Type: "slack", Target: "#monitoring"},  // Awareness
        },
    },
}
```

### 8. Top-Level Variables Only

```go
// GOOD: Top-level var (will be discovered)
var APIAvailability = slo.SLO{...}

// BAD: Inside function (won't be discovered)
func GetSLO() slo.SLO {
    return slo.SLO{...}
}

// BAD: Inside struct (won't be discovered)
type Config struct {
    SLO slo.SLO
}
```

### 9. Keep Queries and SLOs in Sync

Ensure GoodEvents and TotalEvents queries have matching:
- Dataset
- Time range
- Base filters (except the condition that defines "good")

```go
// GOOD: Queries are in sync
var GoodEvents = query.Query{
    Dataset:   "production",           // Same
    TimeRange: query.Days(30),         // Same
    Calculations: []query.Calculation{
        query.Count(),                 // Same
    },
    Filters: []query.Filter{
        query.LT("http.status_code", 500),  // Defines "good"
    },
}

var TotalEvents = query.Query{
    Dataset:   "production",           // Same
    TimeRange: query.Days(30),         // Same
    Calculations: []query.Calculation{
        query.Count(),                 // Same
    },
    // No filters = all events
}
```

### 10. Use Lint Rules

Run linting to catch common issues:

```bash
wetwire-honeycomb lint ./slos/...
```

Common lint rules:
- **WHC040**: SLO missing name
- **WHC044**: Target percentage out of range (0-100)
- **WHC047**: SLO has no burn alerts configured

## See Also

- [CLI Documentation](CLI.md) - Command reference
- [Query Documentation](../README.md) - Query syntax
- [Lint Rules](LINT_RULES.md) - All lint rules
- [FAQ](FAQ.md) - Common questions
