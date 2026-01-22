<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

This example demonstrates the complete Query→SLO→Trigger→Board chain in wetwire-honeycomb-go. It shows how to build a comprehensive monitoring setup where queries serve as the foundation for SLOs, triggers, and dashboards.

## Architecture

```
Queries (queries.go)
    ↓
    ├─→ SLOs (slos.go)
    │     └─→ Burn Alerts
    │
    ├─→ Triggers (triggers.go)
    │     └─→ Recipients (Slack, PagerDuty, Email)
    │
    └─→ Boards (boards.go)
          └─→ Panels (Query Panels, Text Panels, SLO Panels)
```

## Files

### queries.go

Defines the base metrics that power the entire monitoring stack:

- **SuccessfulRequests**: Counts successful HTTP requests (status < 500)
- **AllRequests**: Counts all HTTP requests regardless of status
- **SlowRequests**: Identifies requests exceeding 1000ms latency
- **ErrorRate**: Tracks 5xx error rates by endpoint
- **LatencyP99**: P99 latency across all endpoints
- **RequestThroughput**: Overall request volume and rate

### slos.go

Defines Service Level Objectives that reference queries:

- **APIAvailability**: 99.9% availability SLO over 30 days
  - Uses `SuccessfulRequests` as good events
  - Uses `AllRequests` as total events
  - Includes fast burn (1h) and slow burn (24h) alerts

- **LatencySLO**: 95% of requests under 1 second over 7 days
  - Defines inline queries for fast vs. all requests
  - Includes fast burn alert for latency degradation

### triggers.go

Defines alerts that reference queries:

- **HighLatencyAlert**: Triggers when P99 latency exceeds 2000ms
  - References `SlowRequests` query
  - Evaluates every 5 minutes
  - Notifies #performance and performance-team@example.com

- **ErrorRateAlert**: Triggers when error count exceeds 50 per minute
  - References `ErrorRate` query
  - Evaluates every 2 minutes
  - Notifies #oncall and PagerDuty

- **LowTrafficAlert**: Triggers when request rate drops below 100 req/sec
  - References `RequestThroughput` query
  - Evaluates every 10 minutes
  - Notifies #infrastructure

- **CriticalEndpointLatency**: Triggers for critical endpoints exceeding 500ms
  - Uses inline query with specific endpoint filters
  - Evaluates every 3 minutes
  - Notifies #oncall and PagerDuty

### boards.go

Defines dashboards that visualize queries and SLOs:

- **PerformanceBoard**: Comprehensive API performance dashboard
  - Text panel with overview
  - Query panels for error rate, slow requests, throughput, and latency
  - SLO status panel
  - Preset filters for production environment
  - Tags for team, environment, and service

- **IncidentResponseBoard**: Focused incident investigation dashboard
  - Critical metrics at the top (errors, latency, traffic)
  - Detailed latency breakdown
  - Investigation checklist
  - Optimized layout for rapid triage

## Usage

### List all resources

```bash
wetwire-honeycomb list examples/full_stack
```

### Build all resource types

```bash
# All resources
wetwire-honeycomb build examples/full_stack

# Specific resource type
wetwire-honeycomb build examples/full_stack --type query
wetwire-honeycomb build examples/full_stack --type slo
wetwire-honeycomb build examples/full_stack --type trigger
wetwire-honeycomb build examples/full_stack --type board
```

### Lint for issues

```bash
wetwire-honeycomb lint examples/full_stack
```

### View generated JSON

```bash
# Pretty-printed JSON
wetwire-honeycomb build examples/full_stack --type query --stdout --format pretty

# Compact JSON
wetwire-honeycomb build examples/full_stack --type query --stdout
```

## Key Concepts Demonstrated

### 1. Query Reuse

Queries defined in `queries.go` are referenced by:
- SLOs (e.g., `SuccessfulRequests` → `APIAvailability.SLI.GoodEvents`)
- Triggers (e.g., `SlowRequests` → `HighLatencyAlert.Query`)
- Boards (e.g., `ErrorRate` → `PerformanceBoard` panel)

### 2. Type Safety

All resources use typed constructors and helper functions:
- `query.P99()`, `query.Count()` for calculations
- `query.GT()`, `query.Exists()` for filters
- `trigger.GreaterThan()`, `trigger.Minutes()` for trigger config
- `slo.Percentage()`, `slo.Days()` for SLO config
- `board.QueryPanel()`, `board.TextPanel()` for board panels

### 3. Clear Dependencies

Comments document which resources reference which queries:
```go
// HighLatencyAlert triggers when P99 latency exceeds 2 seconds.
// References:
//   - Query: SlowRequests (queries.go)
var HighLatencyAlert = ...
```

### 4. Inline vs. Referenced Queries

- **Referenced**: Use top-level query variables (e.g., `SlowRequests`)
- **Inline**: Define query directly in SLO/trigger (useful for one-off queries)

Both approaches work and can be mixed as needed.

## Resource Counts

- **10 Queries**: Base metrics and inline queries
- **2 SLOs**: Availability and latency objectives
- **4 Triggers**: Alerts for various failure modes
- **2 Boards**: Performance monitoring and incident response

## Next Steps

1. **Customize datasets**: Change `"production"` to match your Honeycomb dataset names
2. **Adjust thresholds**: Modify latency/error thresholds based on your SLAs
3. **Update recipients**: Change Slack channels and PagerDuty services to your team's
4. **Add more queries**: Create additional queries for your specific use cases
5. **Build JSON**: Use `wetwire-honeycomb build` to generate Honeycomb-compatible JSON
6. **Deploy**: Use the generated JSON with Honeycomb API or Terraform

## Common Patterns

### Pattern 1: SLO with Burn Alerts

```go
var MySLO = slo.SLO{
    SLI: slo.SLI{
        GoodEvents:  MySuccessQuery,
        TotalEvents: MyTotalQuery,
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        {
            AlertType:  slo.BudgetRate,
            Threshold:  2.0,
            Window:     slo.TimePeriod{Hours: 1},
            Recipients: []slo.Recipient{...},
        },
    },
}
```

### Pattern 2: Trigger with Multiple Recipients

```go
var MyTrigger = trigger.Trigger{
    Query:     MyQuery,
    Threshold: trigger.GreaterThan(100),
    Frequency: trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#team"),
        trigger.PagerDutyService("service-id"),
        trigger.EmailAddress("team@example.com"),
    },
}
```

### Pattern 3: Board with Multiple Panel Types

```go
var MyBoard = board.Board{
    Panels: []board.Panel{
        board.TextPanel("# Overview", ...),
        board.QueryPanel(MyQuery, ...),
        board.SLOPanelByID("slo-id", ...),
    },
    PresetFilters: []board.Filter{...},
    Tags:          []board.Tag{...},
}
```
