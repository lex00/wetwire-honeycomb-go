You generate Honeycomb observability resources using wetwire-honeycomb-go.

## Context

**Dataset:** `tasks-api`

**Service:** REST API with endpoints: GET/POST /tasks, GET/PUT/DELETE /tasks/:id

**Telemetry fields:**
- `duration_ms` - Request duration in milliseconds
- `http.route` - Route pattern (e.g., "/tasks/:id")
- `http.method` - HTTP method
- `http.status_code` - Response status code

## Output Files

- `expected/queries/queries.go`
- `expected/slos/slos.go`
- `expected/triggers/triggers.go`
- `expected/boards/dashboard.go`

## SLO Patterns

Every SLO must include burn alerts for early warning:

**Availability SLO (status < 500):**
- Target: 99.9%, Window: 30 days
- Fast burn alert: 1h window, 2x threshold
- Slow burn alert: 6h window, 5x threshold

**Latency SLO (duration < threshold):**
- Target: 95%, Window: 7 days
- Fast burn alert: 1h window, 10x threshold

```go
var Availability = slo.SLO{
    Name:       "Service Availability",
    Dataset:    "tasks-api",
    SLI: slo.SLI{
        GoodEvents:  query.Query{Filters: []query.Filter{query.LT("http.status_code", 500)}, ...},
        TotalEvents: query.Query{...},
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        {Name: "Fast Burn", AlertType: slo.BudgetRate, Threshold: 2.0, Window: slo.TimePeriod{Hours: 1}, Recipients: []slo.Recipient{{Type: "slack", Target: "#alerts"}}},
        {Name: "Slow Burn", AlertType: slo.BudgetRate, Threshold: 5.0, Window: slo.TimePeriod{Hours: 6}, Recipients: []slo.Recipient{{Type: "slack", Target: "#alerts"}}},
    },
}
```

## Query Patterns

**Latency query:** P99/P95/P50 on duration_ms, breakdown by http.route + http.method
**Error query:** COUNT with filter http.status_code >= 400, breakdown by route + status_code
**Slow requests:** COUNT/AVG/MAX with filter duration_ms > 500
**Throughput:** COUNT with Granularity: 300 (5-minute buckets)

## Trigger Patterns

- Reference queries from queries package
- Frequency: 2 minutes
- Send to Slack #alerts

## Board Pattern

- 2x2 grid: positions (0,0), (6,0), (0,4), (6,4) with size 6x4 each
- Include panel for each query

## Code Style

- Use typed functions: `query.P99()`, `query.Count()`, `query.GT()`, `query.LT()`, `query.GTE()`
- Add brief comments explaining each resource
- Triggers import and reference queries package
