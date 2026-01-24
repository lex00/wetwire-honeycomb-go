---
title: "Triggers"
---

This document describes Honeycomb trigger synthesis in wetwire-honeycomb-go.

---

## Overview

Triggers are Honeycomb alerts that monitor query results and send notifications when threshold conditions are met. wetwire-honeycomb provides type-safe trigger declarations that compile to Honeycomb Trigger JSON.

```
Go Structs -> wetwire-honeycomb build -> Trigger JSON -> Honeycomb API
                                              |
                                   (user's responsibility)
```

### Key Concepts

- **Type safety** - Go structs provide compile-time validation
- **Direct query references** - Triggers reference existing query.Query variables
- **Auto-discovery** - AST-based, no registration required
- **Synthesis only** - Generates JSON, does not manage trigger lifecycle

---

## Trigger Declaration

### Basic Structure

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

var HighLatencyAlert = trigger.Trigger{
    Name:        "High P99 Latency",
    Description: "Alert when P99 exceeds 500ms",
    Dataset:     "production",
    Query:       SlowRequestsQuery,
    Threshold:   trigger.GreaterThan(500),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
        trigger.PagerDutyService("api-team"),
    },
    Disabled: false,
}
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Display name for the trigger |
| `Dataset` | `string` | Honeycomb dataset to monitor |
| `Query` | `query.Query` | Query that defines the metric |
| `Threshold` | `Threshold` | Condition that fires the trigger |
| `Frequency` | `Frequency` | How often to evaluate |

### Optional Fields

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `Description` | `string` | Additional context | `""` |
| `Recipients` | `[]Recipient` | Notification targets | `[]` |
| `Disabled` | `bool` | Whether trigger is active | `false` |

---

## Direct Query Reference

Triggers reference existing `query.Query` variables directly, ensuring consistency and reusability.

```go
// 1. Define the query
var ErrorRateQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint"},
    Calculations: []query.Calculation{
        query.Rate("errors", "requests"),
    },
}

// 2. Reference it in the trigger
var ErrorRateAlert = trigger.Trigger{
    Name:      "High Error Rate",
    Dataset:   "production",
    Query:     ErrorRateQuery,  // Direct reference
    Threshold: trigger.GreaterThan(0.05),
    Frequency: trigger.Minutes(5),
}
```

### Benefits

1. **Reusability** - Same query can be used for multiple triggers
2. **Type safety** - Compiler validates query structure
3. **Consistency** - Query and trigger stay in sync
4. **Testability** - Query and trigger can be tested separately

---

## Threshold Configuration

Thresholds define the condition that must be met to fire the trigger.

### Threshold Operators

| Function | Operator | Description |
|----------|----------|-------------|
| `GreaterThan(value)` | `>` | Fire when metric exceeds value |
| `GreaterThanOrEqual(value)` | `>=` | Fire when metric is at least value |
| `LessThan(value)` | `<` | Fire when metric is below value |
| `LessThanOrEqual(value)` | `<=` | Fire when metric is at most value |

### Examples

```go
// Latency alert
Threshold: trigger.GreaterThan(500)  // > 500ms

// Error rate alert
Threshold: trigger.GreaterThan(0.05)  // > 5%

// Low success rate alert
Threshold: trigger.LessThan(0.99)  // < 99%

// Capacity warning
Threshold: trigger.LessThanOrEqual(10)  // <= 10 connections
```

---

## Frequency Configuration

Frequency determines how often the trigger evaluates its query.

### Frequency Functions

| Function | Description |
|----------|-------------|
| `Minutes(m)` | Evaluate every m minutes |
| `Seconds(s)` | Evaluate every s seconds |

### Choosing Frequency

| Alert Type | Frequency | Rationale |
|------------|-----------|-----------|
| Critical production issues | 1-2 minutes | Fast detection required |
| High-priority warnings | 5 minutes | Balance speed and noise |
| General monitoring | 10-15 minutes | Sufficient for most use cases |
| Low-priority checks | 30+ minutes | Reduce alert fatigue |

---

## Recipient Types

Recipients define where notifications are sent when a trigger fires.

### Recipient Functions

| Function | Type | Target Format |
|----------|------|---------------|
| `SlackChannel(channel)` | `slack` | Channel name (e.g., `#alerts`) |
| `PagerDutyService(serviceID)` | `pagerduty` | Service ID or key |
| `EmailAddress(email)` | `email` | Email address |
| `WebhookURL(url)` | `webhook` | HTTPS URL |

### Examples

```go
// Multiple recipients
Recipients: []trigger.Recipient{
    trigger.PagerDutyService("oncall"),
    trigger.SlackChannel("#incidents"),
    trigger.EmailAddress("leadership@example.com"),
    trigger.WebhookURL("https://status.example.com/api/incident"),
}
```

---

## Complete Examples

### Multi-Tier Alerting

```go
var MemoryUsageQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Minutes(10),
    Breakdowns: []string{"service"},
    Calculations: []query.Calculation{
        query.Max("memory_mb"),
    },
}

// Warning at 75% capacity
var MemoryWarning = trigger.Trigger{
    Name:      "Memory Usage Warning",
    Dataset:   "production",
    Query:     MemoryUsageQuery,
    Threshold: trigger.GreaterThan(6000),
    Frequency: trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#warnings"),
    },
}

// Critical at 90% capacity
var MemoryCritical = trigger.Trigger{
    Name:      "Memory Usage Critical",
    Dataset:   "production",
    Query:     MemoryUsageQuery,
    Threshold: trigger.GreaterThan(7200),
    Frequency: trigger.Minutes(2),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("platform-oncall"),
        trigger.SlackChannel("#incidents"),
    },
}
```

---

## Best Practices

1. **Naming Conventions** - Use clear names indicating severity and metric
2. **Multi-Tier Alerting** - Create warning and critical triggers for same metric
3. **Appropriate Frequency** - Balance detection speed with alert fatigue
4. **Query Reuse** - Share queries between multiple triggers
5. **Meaningful Descriptions** - Add context to help on-call engineers
6. **Test Disabled First** - Create new triggers in disabled state
7. **Organize by Domain** - Structure trigger files by domain or team
8. **Document Thresholds** - Comment why specific thresholds were chosen

---

## Common Patterns

### SLA Monitoring

```go
var SLAViolation = trigger.Trigger{
    Name:        "SLA Violation - P99 > 500ms",
    Description: "Alert when P99 exceeds SLA target of 500ms",
    Dataset:     "production",
    Query:       LatencyQuery,
    Threshold:   trigger.GreaterThan(500),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#sla-alerts"),
    },
}
```

### Health Check

```go
var ServiceDown = trigger.Trigger{
    Name:        "Service Health Check Failed",
    Description: "Alert when health check success rate falls below 100%",
    Dataset:     "health-checks",
    Query:       HealthQuery,
    Threshold:   trigger.LessThan(1.0),
    Frequency:   trigger.Minutes(1),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),
        trigger.SlackChannel("#incidents"),
    },
}
```

---

## See Also

- [CLI Reference](../cli/) - Complete command documentation
- [Lint Rules](../lint-rules/) - Trigger validation rules (WHC050+)
- [FAQ](../faq/) - Common questions
