# Triggers

This document describes Honeycomb trigger synthesis in wetwire-honeycomb-go.

---

## Overview

Triggers are Honeycomb alerts that monitor query results and send notifications when threshold conditions are met. wetwire-honeycomb provides type-safe trigger declarations that compile to Honeycomb Trigger JSON.

```
Go Structs → wetwire-honeycomb build → Trigger JSON → Honeycomb API
                                            ↓
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

### Top-Level Declaration Required

Triggers must be declared as top-level package variables to be discovered:

```go
// GOOD: Top-level var
var HighLatencyAlert = trigger.Trigger{...}

// BAD: Function return (won't be discovered)
func GetAlert() trigger.Trigger {
    return trigger.Trigger{...}
}

// BAD: Nested in struct (won't be discovered)
type Alerts struct {
    HighLatency trigger.Trigger
}
```

---

## Trigger Fields

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

### Field Details

#### Name

Human-readable name displayed in Honeycomb UI.

```go
Name: "High P99 Latency"
Name: "Error Rate Spike"
Name: "Database Connection Pool Exhausted"
```

#### Description

Optional context about what the trigger monitors and why it matters.

```go
Description: "Alert when P99 latency exceeds 500ms for more than 5 minutes"
Description: "Monitors error rate across all API endpoints"
```

#### Dataset

The Honeycomb dataset this trigger monitors. Must match an existing dataset.

```go
Dataset: "production"
Dataset: "backend-api"
Dataset: "frontend-events"
```

#### Query

Reference to a `query.Query` variable that defines the metric to monitor. See [Direct Query Reference](#direct-query-reference) for details.

```go
var LatencyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
    },
}

var HighLatencyAlert = trigger.Trigger{
    Query: LatencyQuery,  // Direct reference
    // ...
}
```

#### Threshold

Defines the condition that fires the trigger. See [Threshold Configuration](#threshold-configuration).

```go
Threshold: trigger.GreaterThan(500)
Threshold: trigger.LessThan(0.99)
Threshold: trigger.GreaterThanOrEqual(100)
```

#### Frequency

How often the trigger evaluates the query. See [Frequency Configuration](#frequency-configuration).

```go
Frequency: trigger.Minutes(5)   // Every 5 minutes
Frequency: trigger.Seconds(30)  // Every 30 seconds
```

#### Recipients

List of notification targets. See [Recipient Types](#recipient-types).

```go
Recipients: []trigger.Recipient{
    trigger.SlackChannel("#alerts"),
    trigger.PagerDutyService("api-team"),
    trigger.EmailAddress("team@example.com"),
}
```

#### Disabled

Whether the trigger is active. Use `true` to temporarily disable without deleting.

```go
Disabled: false  // Active
Disabled: true   // Disabled
```

---

## Direct Query Reference

Triggers reference existing `query.Query` variables directly, ensuring consistency and reusability.

### Pattern

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

### Multiple Triggers, One Query

```go
var LatencyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
    },
}

var HighLatencyWarning = trigger.Trigger{
    Name:      "High Latency Warning",
    Dataset:   "production",
    Query:     LatencyQuery,
    Threshold: trigger.GreaterThan(500),
    Frequency: trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#warnings"),
    },
}

var HighLatencyCritical = trigger.Trigger{
    Name:      "High Latency Critical",
    Dataset:   "production",
    Query:     LatencyQuery,  // Same query, different threshold
    Threshold: trigger.GreaterThan(1000),
    Frequency: trigger.Minutes(1),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),
    },
}
```

### Inline Query Definition

For single-use triggers, you can define the query inline:

```go
var OneTimeAlert = trigger.Trigger{
    Name:    "Memory Spike",
    Dataset: "production",
    Query: query.Query{
        Dataset:   "production",
        TimeRange: query.Minutes(15),
        Calculations: []query.Calculation{
            query.Max("memory_mb"),
        },
    },
    Threshold: trigger.GreaterThan(8000),
    Frequency: trigger.Minutes(5),
}
```

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

### GreaterThan

Fires when the metric exceeds the threshold value.

```go
Threshold: trigger.GreaterThan(500)    // > 500
Threshold: trigger.GreaterThan(0.05)   // > 0.05 (5%)
Threshold: trigger.GreaterThan(1000.5) // > 1000.5
```

**Use cases:**
- Latency exceeds target
- Error rate too high
- Memory usage too high
- Request count spike

### GreaterThanOrEqual

Fires when the metric is at least the threshold value.

```go
Threshold: trigger.GreaterThanOrEqual(100)  // >= 100
Threshold: trigger.GreaterThanOrEqual(0.95) // >= 0.95 (95%)
```

**Use cases:**
- Capacity at or above limit
- SLA threshold met or exceeded

### LessThan

Fires when the metric falls below the threshold value.

```go
Threshold: trigger.LessThan(0.99)   // < 0.99 (99%)
Threshold: trigger.LessThan(10)     // < 10
```

**Use cases:**
- Success rate too low
- Throughput below minimum
- Available resources running out

### LessThanOrEqual

Fires when the metric is at most the threshold value.

```go
Threshold: trigger.LessThanOrEqual(5)    // <= 5
Threshold: trigger.LessThanOrEqual(0.90) // <= 0.90 (90%)
```

**Use cases:**
- Capacity at or below minimum
- Performance degraded

### Threshold Examples

```go
// Latency alert
var HighLatency = trigger.Trigger{
    Name:      "High P99 Latency",
    Dataset:   "production",
    Query:     LatencyQuery,
    Threshold: trigger.GreaterThan(500),  // > 500ms
    Frequency: trigger.Minutes(5),
}

// Error rate alert
var HighErrorRate = trigger.Trigger{
    Name:      "High Error Rate",
    Dataset:   "production",
    Query:     ErrorRateQuery,
    Threshold: trigger.GreaterThan(0.05),  // > 5%
    Frequency: trigger.Minutes(5),
}

// Low success rate alert
var LowSuccessRate = trigger.Trigger{
    Name:      "Low Success Rate",
    Dataset:   "production",
    Query:     SuccessRateQuery,
    Threshold: trigger.LessThan(0.99),  // < 99%
    Frequency: trigger.Minutes(5),
}

// Capacity warning
var LowCapacity = trigger.Trigger{
    Name:      "Low Available Connections",
    Dataset:   "production",
    Query:     ConnectionPoolQuery,
    Threshold: trigger.LessThanOrEqual(10),  // <= 10
    Frequency: trigger.Minutes(1),
}
```

---

## Frequency Configuration

Frequency determines how often the trigger evaluates its query.

### Frequency Functions

| Function | Description | Converts to |
|----------|-------------|-------------|
| `Minutes(m)` | Evaluate every m minutes | `m * 60` seconds |
| `Seconds(s)` | Evaluate every s seconds | `s` seconds |

### Minutes

Evaluates the trigger every N minutes.

```go
Frequency: trigger.Minutes(1)   // Every 1 minute
Frequency: trigger.Minutes(5)   // Every 5 minutes
Frequency: trigger.Minutes(15)  // Every 15 minutes
Frequency: trigger.Minutes(30)  // Every 30 minutes
```

### Seconds

Evaluates the trigger every N seconds.

```go
Frequency: trigger.Seconds(30)  // Every 30 seconds
Frequency: trigger.Seconds(60)  // Every 60 seconds (1 minute)
Frequency: trigger.Seconds(300) // Every 300 seconds (5 minutes)
```

### Choosing Frequency

Consider these factors when setting frequency:

1. **Alert fatigue** - More frequent = more potential noise
2. **Response time** - Less frequent = slower detection
3. **Query cost** - More frequent = higher API usage
4. **Data granularity** - Match your data collection rate

**Recommendations:**

| Alert Type | Frequency | Rationale |
|------------|-----------|-----------|
| Critical production issues | 1-2 minutes | Fast detection required |
| High-priority warnings | 5 minutes | Balance speed and noise |
| General monitoring | 10-15 minutes | Sufficient for most use cases |
| Low-priority checks | 30+ minutes | Reduce alert fatigue |

### Frequency Examples

```go
// Critical alert - check frequently
var DatabaseDown = trigger.Trigger{
    Name:      "Database Connection Failed",
    Dataset:   "production",
    Query:     DBHealthQuery,
    Threshold: trigger.LessThan(1),
    Frequency: trigger.Minutes(1),  // Check every minute
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),
    },
}

// Standard monitoring - balanced check
var HighLatency = trigger.Trigger{
    Name:      "High P99 Latency",
    Dataset:   "production",
    Query:     LatencyQuery,
    Threshold: trigger.GreaterThan(500),
    Frequency: trigger.Minutes(5),  // Check every 5 minutes
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
    },
}

// Low-priority check - infrequent
var DiskUsageTrend = trigger.Trigger{
    Name:      "Disk Usage Increasing",
    Dataset:   "production",
    Query:     DiskUsageQuery,
    Threshold: trigger.GreaterThan(80),
    Frequency: trigger.Minutes(30),  // Check every 30 minutes
    Recipients: []trigger.Recipient{
        trigger.EmailAddress("ops@example.com"),
    },
}

// Fast response - seconds granularity
var HighTraffic = trigger.Trigger{
    Name:      "Traffic Spike",
    Dataset:   "production",
    Query:     TrafficQuery,
    Threshold: trigger.GreaterThan(10000),
    Frequency: trigger.Seconds(30),  // Check every 30 seconds
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#incidents"),
    },
}
```

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

### Slack

Send notifications to a Slack channel via Honeycomb's Slack integration.

```go
trigger.SlackChannel("#alerts")
trigger.SlackChannel("#incidents")
trigger.SlackChannel("#team-backend")
```

**Setup required:**
1. Configure Honeycomb Slack integration
2. Authorize access to channels
3. Use channel name with `#` prefix

### PagerDuty

Send notifications to PagerDuty for incident management.

```go
trigger.PagerDutyService("api-team")
trigger.PagerDutyService("PXXXXXX")  // PagerDuty service ID
trigger.PagerDutyService("platform-oncall")
```

**Setup required:**
1. Configure Honeycomb PagerDuty integration
2. Create or identify service ID in PagerDuty
3. Use service ID or integration key

### Email

Send notifications to email addresses.

```go
trigger.EmailAddress("team@example.com")
trigger.EmailAddress("oncall@company.io")
trigger.EmailAddress("alerts@internal.dev")
```

**Setup required:**
1. Verify email addresses in Honeycomb
2. Configure email notification settings

### Webhook

Send HTTP POST notifications to custom endpoints.

```go
trigger.WebhookURL("https://example.com/webhook")
trigger.WebhookURL("https://internal.company.io/alerts")
trigger.WebhookURL("https://hooks.slack.com/services/T00/B00/XXX")
```

**Setup required:**
1. Configure webhook endpoint to receive POST requests
2. Handle Honeycomb webhook payload format
3. Use HTTPS URLs

### Multiple Recipients

Triggers can notify multiple recipients simultaneously:

```go
var CriticalAlert = trigger.Trigger{
    Name:      "Critical System Failure",
    Dataset:   "production",
    Query:     SystemHealthQuery,
    Threshold: trigger.LessThan(1),
    Frequency: trigger.Minutes(1),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),           // Page on-call
        trigger.SlackChannel("#incidents"),            // Notify team
        trigger.EmailAddress("leadership@example.com"), // Escalate
        trigger.WebhookURL("https://status.example.com/api/incident"), // Update status page
    },
}
```

### Recipient Examples by Severity

```go
// Low severity - informational
var InfoAlert = trigger.Trigger{
    Name: "Daily Summary",
    Recipients: []trigger.Recipient{
        trigger.EmailAddress("team@example.com"),
    },
}

// Medium severity - team awareness
var WarningAlert = trigger.Trigger{
    Name: "Latency Warning",
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
    },
}

// High severity - immediate action
var CriticalAlert = trigger.Trigger{
    Name: "Service Down",
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),
        trigger.SlackChannel("#incidents"),
    },
}

// Custom integration
var CustomAlert = trigger.Trigger{
    Name: "Custom Event",
    Recipients: []trigger.Recipient{
        trigger.WebhookURL("https://internal.company.io/alerts"),
    },
}
```

---

## Complete Examples

### Example 1: High Latency Alert

Monitor P99 latency and alert when it exceeds 500ms.

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

// LatencyQuery monitors P99 latency across all endpoints
var LatencyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint"},
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
    },
    Filters: []query.Filter{
        query.Exists("duration_ms"),
    },
}

// HighLatencyAlert fires when P99 exceeds 500ms
var HighLatencyAlert = trigger.Trigger{
    Name:        "High P99 Latency",
    Description: "Alert when P99 latency exceeds 500ms for any endpoint",
    Dataset:     "production",
    Query:       LatencyQuery,
    Threshold:   trigger.GreaterThan(500),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
        trigger.PagerDutyService("api-team"),
    },
    Disabled: false,
}
```

### Example 2: Error Rate Spike

Monitor error rate and alert when it exceeds 5%.

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

// ErrorRateQuery calculates error rate by service
var ErrorRateQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Minutes(15),
    Breakdowns: []string{"service"},
    Calculations: []query.Calculation{
        query.Count(),
        query.CountDistinct("error"),
    },
    Filters: []query.Filter{
        query.NotEquals("http.status_code", "200"),
    },
}

// ErrorRateAlert fires when error rate exceeds 5%
var ErrorRateAlert = trigger.Trigger{
    Name:        "High Error Rate",
    Description: "Alert when error rate exceeds 5% across any service",
    Dataset:     "production",
    Query:       ErrorRateQuery,
    Threshold:   trigger.GreaterThan(0.05),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#incidents"),
        trigger.PagerDutyService("backend-oncall"),
        trigger.EmailAddress("backend-team@example.com"),
    },
    Disabled: false,
}
```

### Example 3: Database Connection Pool

Monitor database connection pool and alert when capacity is low.

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

// DBConnectionPoolQuery monitors available database connections
var DBConnectionPoolQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Minutes(5),
    Calculations: []query.Calculation{
        query.Min("db.pool.available"),
    },
    Filters: []query.Filter{
        query.Equals("db.name", "postgres"),
    },
}

// LowConnectionPoolAlert fires when available connections drop below 10
var LowConnectionPoolAlert = trigger.Trigger{
    Name:        "Database Connection Pool Low",
    Description: "Alert when available database connections fall below 10",
    Dataset:     "production",
    Query:       DBConnectionPoolQuery,
    Threshold:   trigger.LessThanOrEqual(10),
    Frequency:   trigger.Minutes(2),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("database-team"),
        trigger.SlackChannel("#database-alerts"),
    },
    Disabled: false,
}
```

### Example 4: Multi-Tier Alerting

Create warning and critical alerts for the same metric.

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

// MemoryUsageQuery monitors maximum memory usage
var MemoryUsageQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Minutes(10),
    Breakdowns: []string{"service"},
    Calculations: []query.Calculation{
        query.Max("memory_mb"),
    },
}

// MemoryWarning fires at 75% capacity (6GB of 8GB)
var MemoryWarning = trigger.Trigger{
    Name:        "Memory Usage Warning",
    Description: "Warning when memory usage exceeds 75%",
    Dataset:     "production",
    Query:       MemoryUsageQuery,
    Threshold:   trigger.GreaterThan(6000),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#warnings"),
    },
    Disabled: false,
}

// MemoryCritical fires at 90% capacity (7.2GB of 8GB)
var MemoryCritical = trigger.Trigger{
    Name:        "Memory Usage Critical",
    Description: "Critical alert when memory usage exceeds 90%",
    Dataset:     "production",
    Query:       MemoryUsageQuery,
    Threshold:   trigger.GreaterThan(7200),
    Frequency:   trigger.Minutes(2),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("platform-oncall"),
        trigger.SlackChannel("#incidents"),
    },
    Disabled: false,
}
```

### Example 5: Disabled Trigger

Create a disabled trigger for future use or testing.

```go
package alerts

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/trigger"
)

// ExperimentalFeatureQuery monitors usage of experimental feature
var ExperimentalFeatureQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("feature", "experimental_mode"),
    },
}

// ExperimentalFeatureAlert is disabled until feature is stable
var ExperimentalFeatureAlert = trigger.Trigger{
    Name:        "Experimental Feature Usage",
    Description: "Alert when experimental feature is used (disabled until stable)",
    Dataset:     "production",
    Query:       ExperimentalFeatureQuery,
    Threshold:   trigger.GreaterThan(100),
    Frequency:   trigger.Minutes(15),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#product-team"),
    },
    Disabled: true,  // Disabled until feature is stable
}
```

---

## JSON Output Format

wetwire-honeycomb serializes triggers to Honeycomb Trigger JSON format.

### Basic Trigger JSON

```go
var HighLatencyAlert = trigger.Trigger{
    Name:        "High P99 Latency",
    Description: "Alert when P99 exceeds 500ms",
    Dataset:     "production",
    Query:       LatencyQuery,
    Threshold:   trigger.GreaterThan(500),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
    },
    Disabled: false,
}
```

Generates:

```json
{
  "name": "High P99 Latency",
  "description": "Alert when P99 exceeds 500ms",
  "dataset": "production",
  "query": {
    "time_range": 3600,
    "calculations": [
      {
        "op": "P99",
        "column": "duration_ms"
      }
    ]
  },
  "threshold": {
    "op": ">",
    "value": 500
  },
  "frequency": 300,
  "recipients": [
    {
      "type": "slack",
      "target": "#alerts"
    }
  ],
  "disabled": false
}
```

### JSON Field Mapping

| Go Field | JSON Field | Type | Description |
|----------|-----------|------|-------------|
| `Name` | `name` | string | Trigger name |
| `Description` | `description` | string | Trigger description |
| `Dataset` | `dataset` | string | Dataset name |
| `Query` | `query` | object | Query specification (see Query JSON) |
| `Threshold.Op` | `threshold.op` | string | Operator (`>`, `>=`, `<`, `<=`) |
| `Threshold.Value` | `threshold.value` | number | Threshold value |
| `Frequency.Seconds` | `frequency` | number | Evaluation interval in seconds |
| `Recipients` | `recipients` | array | Notification targets |
| `Recipients[].Type` | `recipients[].type` | string | Recipient type |
| `Recipients[].Target` | `recipients[].target` | string | Target identifier |
| `Disabled` | `disabled` | boolean | Whether trigger is disabled |

### Complete JSON Example

```go
var ErrorRateAlert = trigger.Trigger{
    Name:        "High Error Rate",
    Description: "Alert when error rate exceeds 5%",
    Dataset:     "production",
    Query: query.Query{
        Dataset:   "production",
        TimeRange: query.Minutes(15),
        Breakdowns: []string{"service"},
        Calculations: []query.Calculation{
            query.Count(),
            query.CountDistinct("error"),
        },
        Filters: []query.Filter{
            query.NotEquals("http.status_code", "200"),
        },
    },
    Threshold: trigger.GreaterThan(0.05),
    Frequency: trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#incidents"),
        trigger.PagerDutyService("backend-oncall"),
        trigger.EmailAddress("backend-team@example.com"),
    },
    Disabled: false,
}
```

Generates:

```json
{
  "name": "High Error Rate",
  "description": "Alert when error rate exceeds 5%",
  "dataset": "production",
  "query": {
    "time_range": 900,
    "breakdowns": ["service"],
    "calculations": [
      {
        "op": "COUNT"
      },
      {
        "op": "COUNT_DISTINCT",
        "column": "error"
      }
    ],
    "filters": [
      {
        "column": "http.status_code",
        "op": "!=",
        "value": "200"
      }
    ]
  },
  "threshold": {
    "op": ">",
    "value": 0.05
  },
  "frequency": 300,
  "recipients": [
    {
      "type": "slack",
      "target": "#incidents"
    },
    {
      "type": "pagerduty",
      "target": "backend-oncall"
    },
    {
      "type": "email",
      "target": "backend-team@example.com"
    }
  ],
  "disabled": false
}
```

### Minimal JSON Example

```go
var MinimalAlert = trigger.Trigger{
    Name:    "Minimal Alert",
    Dataset: "production",
}
```

Generates:

```json
{
  "name": "Minimal Alert",
  "dataset": "production",
  "disabled": false
}
```

### JSON Omission Rules

Fields are omitted from JSON when:

- `Description` is empty string
- `Dataset` is empty string (use trigger-level or query-level)
- `Query` has no calculations and default time range
- `Threshold.Op` is empty
- `Frequency.Seconds` is 0
- `Recipients` is empty array

---

## CLI Commands

### Build Triggers

Generate JSON for all triggers:

```bash
wetwire-honeycomb build ./alerts -o triggers.json
```

Build specific triggers:

```bash
wetwire-honeycomb build ./alerts --resource trigger -o triggers.json
```

### List Triggers

List all discovered triggers:

```bash
wetwire-honeycomb list ./alerts --resource trigger
```

Output:

```
TRIGGER                PACKAGE  DATASET      THRESHOLD  FREQUENCY
HighLatencyAlert       alerts   production   >500       5m
ErrorRateAlert         alerts   production   >0.05      5m
LowConnectionPoolAlert alerts   production   <=10       2m
```

### Lint Triggers

Check triggers for issues:

```bash
wetwire-honeycomb lint ./alerts --resource trigger
```

See [LINT_RULES.md](LINT_RULES.md) for trigger-specific lint rules.

---

## Best Practices

### 1. Naming Conventions

Use descriptive names that indicate severity and metric:

```go
// GOOD: Clear severity and metric
var HighLatencyCritical = trigger.Trigger{
    Name: "Critical: P99 Latency Exceeds 1s",
}

var MemoryWarning = trigger.Trigger{
    Name: "Warning: Memory Usage Above 75%",
}

// BAD: Unclear severity
var Alert1 = trigger.Trigger{
    Name: "Alert",
}
```

### 2. Multi-Tier Alerting

Create multiple triggers for the same metric with different thresholds:

```go
// Warning at 75%
var MemoryWarning = trigger.Trigger{
    Threshold: trigger.GreaterThan(6000),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#warnings"),
    },
}

// Critical at 90%
var MemoryCritical = trigger.Trigger{
    Threshold: trigger.GreaterThan(7200),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("oncall"),
    },
}
```

### 3. Appropriate Frequency

Balance detection speed with alert fatigue:

```go
// Critical issues - check frequently
Frequency: trigger.Minutes(1)

// Standard monitoring - balanced
Frequency: trigger.Minutes(5)

// Low-priority checks - less frequent
Frequency: trigger.Minutes(30)
```

### 4. Query Reuse

Share queries between multiple triggers:

```go
var LatencyQuery = query.Query{...}

var WarningAlert = trigger.Trigger{
    Query:     LatencyQuery,
    Threshold: trigger.GreaterThan(500),
}

var CriticalAlert = trigger.Trigger{
    Query:     LatencyQuery,
    Threshold: trigger.GreaterThan(1000),
}
```

### 5. Meaningful Descriptions

Add context to help on-call engineers:

```go
Description: "Alert when P99 latency exceeds 500ms. Check database query performance and connection pool. Runbook: https://wiki.company.com/runbooks/latency"
```

### 6. Test Disabled First

Create new triggers in disabled state:

```go
var NewAlert = trigger.Trigger{
    Name:     "New Experimental Alert",
    Disabled: true,  // Enable after validation
}
```

### 7. Organize by Domain

Structure trigger files by domain or team:

```
alerts/
├── api/
│   ├── latency.go      # API latency alerts
│   └── errors.go       # API error alerts
├── database/
│   ├── connections.go  # DB connection alerts
│   └── queries.go      # DB query alerts
└── infrastructure/
    ├── memory.go       # Memory alerts
    └── cpu.go          # CPU alerts
```

### 8. Document Thresholds

Comment why specific thresholds were chosen:

```go
// Threshold based on SLA: 99% of requests under 500ms
Threshold: trigger.GreaterThan(500)

// Based on max pool size of 100 connections
Threshold: trigger.LessThanOrEqual(10)
```

---

## Common Patterns

### Pattern 1: SLA Monitoring

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
        trigger.EmailAddress("engineering@example.com"),
    },
}
```

### Pattern 2: Capacity Planning

```go
var CapacityWarning = trigger.Trigger{
    Name:        "Approaching Capacity Limit",
    Description: "Alert when resource usage exceeds 80% of capacity",
    Dataset:     "infrastructure",
    Query:       ResourceQuery,
    Threshold:   trigger.GreaterThan(0.80),
    Frequency:   trigger.Minutes(15),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#capacity"),
    },
}
```

### Pattern 3: Health Check

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

### Pattern 4: Anomaly Detection

```go
var TrafficSpike = trigger.Trigger{
    Name:        "Unusual Traffic Pattern",
    Description: "Alert when request rate is 3x normal baseline",
    Dataset:     "production",
    Query:       TrafficQuery,
    Threshold:   trigger.GreaterThan(30000),  // 3x baseline of 10k
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#traffic-alerts"),
    },
}
```

### Pattern 5: Security Monitoring

```go
var SuspiciousActivity = trigger.Trigger{
    Name:        "High Failed Login Attempts",
    Description: "Alert when failed login attempts exceed threshold",
    Dataset:     "security",
    Query:       FailedLoginsQuery,
    Threshold:   trigger.GreaterThan(100),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.PagerDutyService("security-team"),
        trigger.SlackChannel("#security-alerts"),
        trigger.EmailAddress("security@example.com"),
    },
}
```

---

## See Also

- [CLI Reference](CLI.md) - Complete command documentation
- [Lint Rules](LINT_RULES.md) - Trigger validation rules
- [Query Documentation](QUERY.md) - Query specification
- [Quick Start](QUICK_START.md) - Getting started guide
