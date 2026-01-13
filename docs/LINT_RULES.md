# Lint Rules

This document describes all lint rules for wetwire-honeycomb-go.

---

## Overview

wetwire-honeycomb linter checks query declarations for:
- Required fields (dataset, time range, calculations)
- Honeycomb API constraints and limits
- Best practices for query performance
- Common mistakes and anti-patterns

### Rule Naming

Rules follow the format `WHC<NNN>` where:
- `W` = wetwire
- `HC` = Honeycomb (domain code)
- `<NNN>` = Rule number (001, 002, etc.)

### Severity Levels

| Level | Description | Exit Code |
|-------|-------------|-----------|
| **error** | Must fix - violates Honeycomb constraints or missing required fields | 1 |
| **warning** | Should fix - potential performance or usability issue | 1 |

---

## Rule Index

| Rule | Description | Severity |
|------|-------------|----------|
| [WHC001](#whc001-missing-dataset) | Missing dataset | error |
| [WHC002](#whc002-missing-time-range) | Missing time range | error |
| [WHC003](#whc003-empty-calculations) | Empty calculations | error |
| [WHC004](#whc004-breakdown-without-order) | Breakdown without order | warning |
| [WHC005](#whc005-high-cardinality-breakdown) | High cardinality breakdown | warning |
| [WHC006](#whc006-invalid-calculation-for-column-type) | Invalid calculation for column type | error |
| [WHC007](#whc007-invalid-filter-operator) | Invalid filter operator | error |
| [WHC008](#whc008-missing-limit-with-breakdowns) | Missing limit with breakdowns | warning |
| [WHC009](#whc009-time-range-exceeds-7-days) | Time range exceeds 7 days | error |
| [WHC010](#whc010-excessive-filter-count) | Excessive filter count | warning |
| [WHC011](#whc011-circular-dependency) | Circular dependency | error |
| [WHC012](#whc012-secret-in-filter) | Secret in filter | error |
| [WHC013](#whc013-sensitive-column-exposure) | Sensitive column exposure | warning |
| [WHC014](#whc014-hardcoded-credentials) | Hardcoded credentials | error |
| [WHC020](#whc020-inline-calculation-definition) | Inline calculation definition | warning |
| [WHC021](#whc021-inline-filter-definition) | Inline filter definition | warning |
| [WHC022](#whc022-raw-map-literal) | Raw map literal | warning |
| [WHC023](#whc023-deeply-nested-configuration) | Deeply nested configuration | warning |
| **Board Rules** | | |
| [WHC030](#whc030-board-has-no-panels) | Board has no panels | error |
| [WHC034](#whc034-board-exceeds-panel-limit) | Board exceeds panel limit | warning |
| **SLO Rules** | | |
| [WHC040](#whc040-slo-missing-name) | SLO missing name | error |
| [WHC044](#whc044-target-out-of-range) | Target out of range | error |
| [WHC047](#whc047-slo-no-burn-alerts) | SLO no burn alerts | warning |
| **Trigger Rules** | | |
| [WHC050](#whc050-trigger-missing-name) | Trigger missing name | error |
| [WHC053](#whc053-trigger-no-recipients) | Trigger no recipients | warning |
| [WHC054](#whc054-trigger-frequency-under-1-minute) | Trigger frequency under 1 minute | warning |
| [WHC056](#whc056-trigger-is-disabled) | Trigger is disabled | info |

---

## Rules

### WHC001: Missing dataset

**Severity:** error

**Description:**

Every query must specify a dataset. This is a required field for the Honeycomb Query API.

**Why:**

- Honeycomb requires a target dataset for all queries
- Queries without a dataset will fail at the API level

**Bad:**

```go
var MyQuery = query.Query{
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    // Missing Dataset field
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC002: Missing time range

**Severity:** error

**Description:**

Every query must specify a time range, either as a relative duration or absolute start/end times.

**Why:**

- Honeycomb requires a time window for all queries
- Queries without a time range will fail at the API level

**Bad:**

```go
var MyQuery = query.Query{
    Dataset: "production",
    Calculations: []query.Calculation{
        query.Count(),
    },
    // Missing TimeRange field
}
```

**Good:**

```go
// Relative time range
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Calculations: []query.Calculation{
        query.Count(),
    },
}

// Absolute time range
var AnotherQuery = query.Query{
    Dataset: "production",
    TimeRange: query.Absolute(startTime, endTime),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC003: Empty calculations

**Severity:** error

**Description:**

Every query must have at least one calculation. A query without calculations would return no useful data.

**Why:**

- Calculations define what aggregations to compute
- Honeycomb requires at least one calculation per query

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint"},
    // Missing Calculations
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint"},
    Calculations: []query.Calculation{
        query.Count(),
        query.P99("duration_ms"),
    },
}
```

---

### WHC004: Breakdown without order

**Severity:** warning

**Description:**

Queries with breakdowns should specify an ordering to ensure consistent, predictable results.

**Why:**

- Without explicit ordering, result order may vary between query executions
- Makes it harder to compare results over time
- Dashboard displays may be inconsistent

**Triggers:**

This rule triggers when a query has one or more breakdowns specified.

**Example:**

```go
// This query will trigger WHC004
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint", "service"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    // No Orders field specified
}
```

**Fix:**

Add explicit ordering to your query:

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"endpoint", "service"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    Orders: []query.Order{
        {Op: "COUNT", Order: "descending"},
    },
}
```

---

### WHC005: High cardinality breakdown

**Severity:** warning

**Description:**

Queries with a limit greater than 100 may return high-cardinality results that are difficult to visualize and analyze.

**Why:**

- High cardinality results are harder to interpret
- May cause performance issues in Honeycomb UI
- Dashboard visualizations become cluttered

**Threshold:** Limit > 100

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user_id"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    Limit: 1000, // High cardinality
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user_id"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    Limit: 100, // Reasonable cardinality
}
```

---

### WHC006: Invalid calculation for column type

**Severity:** error

**Description:**

Numeric calculations (percentiles, SUM, AVG, etc.) should not be used on columns that appear to be string fields.

**Why:**

- Percentile and sum operations require numeric values
- Using numeric operations on string columns will cause query errors
- Early detection prevents runtime failures

**Detection:**

Uses heuristic pattern matching to detect likely string columns:
- Column names containing: `name`, `message`, `error`, `status`, `endpoint`, `path`, `url`, `type`, `service`, `env`, `environment`
- Excludes columns with numeric suffixes: `_ms`, `_bytes`, `_count`

**Invalid Calculations for String Columns:**

`P50`, `P75`, `P90`, `P95`, `P99`, `P999`, `SUM`, `AVG`, `MIN`, `MAX`, `HEATMAP`

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("endpoint"),      // Percentile on string field
        query.Sum("error_message"), // Sum on string field
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"),       // Percentile on numeric field
        query.CountDistinct("endpoint"), // Count distinct on string field
    },
}
```

---

### WHC007: Invalid filter operator

**Severity:** error

**Description:**

Filter operators must be valid Honeycomb filter operators.

**Why:**

- Invalid operators will cause API errors
- Typos in operators are common mistakes

**Valid Operators:**

| Operator | Description |
|----------|-------------|
| `=` | Equals |
| `!=` | Not equals |
| `>` | Greater than |
| `>=` | Greater than or equal |
| `<` | Less than |
| `<=` | Less than or equal |
| `contains` | String contains |
| `does-not-contain` | String does not contain |
| `exists` | Field exists |
| `does-not-exist` | Field does not exist |
| `starts-with` | String starts with |
| `in` | Value in list |
| `not-in` | Value not in list |

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        {Column: "status", Op: "==", Value: "error"},  // Invalid: use "="
        {Column: "path", Op: "like", Value: "/api/*"}, // Invalid: use "contains" or "starts-with"
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("status", "error"),
        query.StartsWith("path", "/api/"),
    },
}
```

---

### WHC008: Missing limit with breakdowns

**Severity:** warning

**Description:**

Queries with breakdowns should specify a limit to prevent returning too many results.

**Why:**

- Without a limit, queries may return thousands of groups
- Large result sets impact performance
- Honeycomb applies a default limit, which may not match expectations

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user_id", "endpoint"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    // No Limit specified
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user_id", "endpoint"},
    Calculations: []query.Calculation{
        query.Count(),
    },
    Limit: 100,
}
```

---

### WHC009: Time range exceeds 7 days

**Severity:** error

**Description:**

Honeycomb has a maximum time range of 7 days for queries. Queries exceeding this limit will fail.

**Why:**

- Honeycomb API enforces a 7-day maximum time range
- Queries exceeding this will return an error
- Large time ranges have significant performance implications

**Threshold:** 604,800 seconds (7 days)

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(30), // 30 days exceeds limit
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Days(7), // Maximum allowed
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC010: Excessive filter count

**Severity:** warning

**Description:**

Queries with more than 50 filters may have performance issues and are difficult to maintain.

**Why:**

- Many filters increase query complexity
- Performance degrades with filter count
- Queries become hard to understand and debug

**Threshold:** > 50 filters

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        // 51+ filters...
    },
}
```

**Good:**

Consider restructuring queries with many filters:
- Use `in` operator for multiple value matches
- Split into multiple focused queries
- Use derived columns in Honeycomb

---

### WHC011: Circular dependency

**Severity:** error

**Description:**

Detects self-referential patterns where a query name appears in its own filter or calculation columns.

**Why:**

- Self-referential queries indicate a logical error
- May cause unexpected behavior or confusing results
- Often indicates copy-paste errors from other queries

**Detection:**

Checks if the query variable name (case-insensitive) appears in:
- Filter column names
- Calculation column names

**Bad:**

```go
var LatencyMetrics = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("LatencyMetrics_value"), // Self-reference to query name
    },
}
```

**Good:**

```go
var LatencyMetrics = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"), // Actual column name
    },
}
```

---

### WHC012: Secret in filter

**Severity:** error

**Description:**

Detects potential secrets, tokens, or credentials in filter values that should not be hardcoded.

**Why:**

- Secrets in source code can be exposed in version control
- Hardcoded credentials are a security vulnerability
- API keys and tokens should be managed separately

**Detected Patterns:**

| Pattern | Example |
|---------|---------|
| OpenAI keys | `sk-...`, `sk_live_...`, `sk_test_...` |
| Generic tokens | Values containing `token`, `bearer`, `secret` |
| API keys | Values containing `api_key`, `apikey`, `api-key` |
| Passwords | Values containing `password`, `passwd` |
| Stripe keys | `sk_live_...`, `sk_test_...` |
| Access credentials | Values containing `access_key`, `auth_token`, `credential` |

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Equals("api.key", "sk-1234567890abcdef"), // Secret in filter!
    },
}
```

**Good:**

```go
// Use environment variables or configuration for secrets
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        query.Exists("api.key"), // Check existence without exposing value
    },
}
```

---

### WHC013: Sensitive column exposure

**Severity:** warning

**Description:**

Warns when breakdown columns might expose personally identifiable information (PII) or sensitive data.

**Why:**

- Breaking down by sensitive fields exposes individual values
- PII in query results may violate privacy regulations
- Sensitive data should be aggregated, not grouped

**Detected Column Patterns:**

| Category | Patterns |
|----------|----------|
| Authentication | `password`, `passwd`, `secret`, `auth_token`, `api_key`, `access_token`, `private_key` |
| Financial | `credit_card`, `card_number`, `cvv`, `pin` |
| Personal | `ssn`, `social_security` |

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user.password", "credit_card_number"}, // Sensitive!
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Breakdowns: []string{"user.id", "transaction_type"}, // Non-sensitive
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC014: Hardcoded credentials

**Severity:** error

**Description:**

Detects hardcoded credentials or authentication parameters in dataset names.

**Why:**

- Dataset names may appear in logs and error messages
- Credentials in dataset names are easily exposed
- Indicates misconfiguration of data sources

**Detected Patterns:**

Patterns in dataset names like:
- `password=...`
- `key=...`
- `token=...`
- `secret=...`
- `apikey=...`, `api_key=...`, `api-key=...`
- `access_key=...`, `access-key=...`
- `auth=...`
- `credential=...`

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production?password=secret123", // Credential in dataset!
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production", // Clean dataset name
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC020: Inline calculation definition

**Severity:** warning

**Description:**

Detects inline Calculation definitions that should be extracted to named variables for better readability and reusability.

**Why:**

- Named variables improve code readability
- Calculations can be reused across multiple queries
- Makes testing easier with standalone calculation definitions
- Follows Go best practices for configuration

**Threshold:** More than 3 inline calculation definitions

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        {Op: "COUNT"},
        {Op: "P50", Column: "duration_ms"},
        {Op: "P95", Column: "duration_ms"},
        {Op: "P99", Column: "duration_ms"},
        {Op: "AVG", Column: "response_size"},
    },
}
```

**Good:**

```go
// Named calculations for reusability
var (
    CountCalc       = query.Count()
    P50Latency      = query.P50("duration_ms")
    P95Latency      = query.P95("duration_ms")
    P99Latency      = query.P99("duration_ms")
    AvgResponseSize = query.Avg("response_size")
)

var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        CountCalc,
        P50Latency,
        P95Latency,
        P99Latency,
        AvgResponseSize,
    },
}
```

---

### WHC021: Inline filter definition

**Severity:** warning

**Description:**

Detects inline Filter definitions that should be extracted to named variables for better readability and reusability.

**Why:**

- Named filters can be reused across multiple queries
- Improves code organization and maintainability
- Makes complex filter logic easier to test
- Follows Go best practices for configuration

**Threshold:** More than 3 inline filter definitions

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        {Column: "status", Op: "=", Value: "error"},
        {Column: "env", Op: "=", Value: "prod"},
        {Column: "service", Op: "=", Value: "api"},
        {Column: "duration_ms", Op: ">", Value: 1000},
    },
}
```

**Good:**

```go
// Named filters for reusability
var (
    ErrorFilter      = query.Equals("status", "error")
    ProductionFilter = query.Equals("env", "prod")
    APIServiceFilter = query.Equals("service", "api")
    SlowRequestFilter = query.GT("duration_ms", 1000)
)

var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    Filters: []query.Filter{
        ErrorFilter,
        ProductionFilter,
        APIServiceFilter,
        SlowRequestFilter,
    },
}
```

---

### WHC022: Raw map literal

**Severity:** warning

**Description:**

Detects raw map literals used instead of typed query builders. Using typed builders provides better type safety and IDE support.

**Why:**

- Typed builders provide compile-time type checking
- IDE autocomplete works with typed structs
- Reduces runtime errors from typos in field names
- Follows Go best practices for type safety

**Bad:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    // Raw map literal - no type safety
    Metadata: map[string]interface{}{
        "owner": "team-a",
        "environment": "production",
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    // Use typed struct instead of map
    Metadata: query.Metadata{
        Owner:       "team-a",
        Environment: "production",
    },
}
```

---

### WHC023: Deeply nested configuration

**Severity:** warning

**Description:**

Detects deeply nested query configurations that should be flattened for better readability.

**Why:**

- Deep nesting makes code harder to read and maintain
- Increases cognitive load when reviewing code
- Makes debugging more difficult
- Often indicates code that could be refactored

**Threshold:** Nesting depth greater than 4 levels

**Bad:**

```go
var MyQuery = query.Query{
    Dataset: "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Calculation{
            Op: "HEATMAP",
            Column: "duration_ms",
            Config: query.CalculationConfig{
                Buckets: query.BucketConfig{
                    Min: 0,
                    Max: 10000,
                    Count: 50,
                    Distribution: query.Distribution{
                        Type: "logarithmic",
                        Base: 2,
                    },
                },
            },
        },
    },
}
```

**Good:**

```go
// Extract nested configurations to named variables
var logDistribution = query.Distribution{
    Type: "logarithmic",
    Base: 2,
}

var bucketConfig = query.BucketConfig{
    Min:          0,
    Max:          10000,
    Count:        50,
    Distribution: logDistribution,
}

var heatmapConfig = query.CalculationConfig{
    Buckets: bucketConfig,
}

var LatencyHeatmap = query.Calculation{
    Op:     "HEATMAP",
    Column: "duration_ms",
    Config: heatmapConfig,
}

var MyQuery = query.Query{
    Dataset:      "production",
    TimeRange:    query.Hours(1),
    Calculations: []query.Calculation{LatencyHeatmap},
}
```

---

## Board Rules

### WHC030: Board has no panels

**Severity:** error

**Description:**

Every board must have at least one panel. Empty boards serve no purpose.

**Why:**

- Honeycomb boards require panels to display content
- Empty boards provide no value to users

**Bad:**

```go
var EmptyBoard = board.Board{
    Name: "Empty Dashboard",
    // No Panels
}
```

**Good:**

```go
var ServiceBoard = board.Board{
    Name: "Service Dashboard",
    Panels: []board.Panel{
        board.QueryPanel(ServiceLatency),
        board.QueryPanel(ErrorRate),
    },
}
```

---

### WHC034: Board exceeds panel limit

**Severity:** warning

**Description:**

Boards with more than 20 panels may have performance issues and are harder to navigate.

**Why:**

- Large boards take longer to load
- Too many panels make it difficult to find relevant data
- Consider splitting into multiple focused boards

**Bad:**

```go
var MegaBoard = board.Board{
    Name: "Everything Dashboard",
    Panels: []board.Panel{
        // 25+ panels...
    },
}
```

**Good:**

```go
// Split into focused boards
var LatencyBoard = board.Board{
    Name: "Latency Overview",
    Panels: []board.Panel{/* 5-8 latency panels */},
}

var ErrorBoard = board.Board{
    Name: "Error Analysis",
    Panels: []board.Panel{/* 5-8 error panels */},
}
```

---

## SLO Rules

### WHC040: SLO missing name

**Severity:** error

**Description:**

Every SLO must have a name for identification in Honeycomb.

**Why:**

- SLOs without names are hard to identify
- Required for proper organization and alerting

**Bad:**

```go
var MySLO = slo.SLO{
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  SuccessfulRequests,
        TotalEvents: AllRequests,
    },
    Target: slo.Percentage(99.9),
    // Missing Name
}
```

**Good:**

```go
var APIAvailability = slo.SLO{
    Name:    "API Availability",
    Dataset: "production",
    SLI: slo.SLI{
        GoodEvents:  SuccessfulRequests,
        TotalEvents: AllRequests,
    },
    Target: slo.Percentage(99.9),
}
```

---

### WHC044: Target out of range

**Severity:** error

**Description:**

SLO target percentage must be between 0 and 100 (exclusive of 0, inclusive of 100).

**Why:**

- 0% target makes no sense (no reliability expected)
- Values over 100% are mathematically impossible

**Bad:**

```go
var BadSLO = slo.SLO{
    Name:   "Impossible SLO",
    Target: slo.Percentage(101), // Over 100%
}

var ZeroSLO = slo.SLO{
    Name:   "Zero Target",
    Target: slo.Percentage(0), // No reliability
}
```

**Good:**

```go
var RealisticSLO = slo.SLO{
    Name:   "API Availability",
    Target: slo.Percentage(99.9), // 99.9% availability
}
```

---

### WHC047: SLO no burn alerts

**Severity:** warning

**Description:**

SLOs without burn alerts won't notify you when the error budget is being consumed too quickly.

**Why:**

- Burn alerts provide early warning of SLO violations
- Without alerts, you may only discover issues after the SLO is breached

**Bad:**

```go
var MySLO = slo.SLO{
    Name:   "API Availability",
    Target: slo.Percentage(99.9),
    // No BurnAlerts configured
}
```

**Good:**

```go
var MySLO = slo.SLO{
    Name:   "API Availability",
    Target: slo.Percentage(99.9),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(),  // 14.4x burn rate, 1hr window
        slo.SlowBurn(),  // 1x burn rate, 24hr window
    },
}
```

---

## Trigger Rules

### WHC050: Trigger missing name

**Severity:** error

**Description:**

Every trigger must have a name for identification in alerts.

**Why:**

- Unnamed triggers are hard to identify in alert notifications
- Required for proper organization and debugging

**Bad:**

```go
var MyTrigger = trigger.Trigger{
    Dataset:   "production",
    Query:     SlowRequests,
    Threshold: trigger.GreaterThan(500),
    // Missing Name
}
```

**Good:**

```go
var HighLatencyAlert = trigger.Trigger{
    Name:      "High P99 Latency",
    Dataset:   "production",
    Query:     SlowRequests,
    Threshold: trigger.GreaterThan(500),
}
```

---

### WHC053: Trigger no recipients

**Severity:** warning

**Description:**

Triggers without recipients won't notify anyone when they fire.

**Why:**

- Alerts need to reach someone to be useful
- Consider adding Slack, PagerDuty, email, or webhook recipients

**Bad:**

```go
var SilentAlert = trigger.Trigger{
    Name:      "Silent Alert",
    Query:     ErrorRate,
    Threshold: trigger.GreaterThan(0.05),
    // No Recipients - alert fires but nobody knows
}
```

**Good:**

```go
var EffectiveAlert = trigger.Trigger{
    Name:      "Error Rate Alert",
    Query:     ErrorRate,
    Threshold: trigger.GreaterThan(0.05),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
        trigger.PagerDutyService("production-oncall"),
    },
}
```

---

### WHC054: Trigger frequency under 1 minute

**Severity:** warning

**Description:**

Trigger frequencies under 1 minute may cause excessive alerting and API load.

**Why:**

- Very frequent triggers can overwhelm alert channels
- May not provide actionable information
- Consider if such high frequency is really needed

**Bad:**

```go
var TooFrequent = trigger.Trigger{
    Name:      "Overeager Alert",
    Query:     SlowRequests,
    Frequency: trigger.Seconds(30), // Every 30 seconds
}
```

**Good:**

```go
var ReasonableAlert = trigger.Trigger{
    Name:      "Latency Alert",
    Query:     SlowRequests,
    Frequency: trigger.Minutes(5), // Every 5 minutes
}
```

---

### WHC056: Trigger is disabled

**Severity:** info

**Description:**

The trigger is explicitly disabled and will not fire alerts.

**Why:**

- This is informational, not necessarily a problem
- Disabled triggers may be intentional (maintenance, testing)
- Consider removing if no longer needed

**Example:**

```go
var DisabledTrigger = trigger.Trigger{
    Name:     "Temporarily Disabled",
    Query:    SlowRequests,
    Disabled: true, // Explicitly disabled
}
```

---

## Disabling Rules

### Command Line

Check specific rules only:

```bash
wetwire-honeycomb lint --rules WHC001,WHC002 ./queries/...
```

### Inline Comments

Rules cannot currently be disabled inline. This feature is planned for a future release.

---

## See Also

- [CLI Reference](CLI.md) - Complete command documentation
- [FAQ](FAQ.md) - Common questions
- [Honeycomb Query Best Practices](https://docs.honeycomb.io/working-with-your-data/queries/) - Official guide
