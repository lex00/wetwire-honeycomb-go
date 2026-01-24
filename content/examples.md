---
title: "Examples"
---

This document provides an index of example queries and usage patterns for wetwire-honeycomb-go.

## Example Packages

| Package | Description | Queries |
|---------|-------------|---------|
| [latency](https://github.com/lex00/wetwire-honeycomb-go/tree/main/examples/latency/) | Latency analysis patterns | 4 |
| [errors](https://github.com/lex00/wetwire-honeycomb-go/tree/main/examples/errors/) | Error tracking and analysis | 5 |
| [traffic](https://github.com/lex00/wetwire-honeycomb-go/tree/main/examples/traffic/) | Traffic and capacity analysis | 5 |
| [sli](https://github.com/lex00/wetwire-honeycomb-go/tree/main/examples/sli/) | SLI/SLO metric tracking | 6 |
| [traces](https://github.com/lex00/wetwire-honeycomb-go/tree/main/examples/traces/) | Distributed trace analysis | 5 |

## Running Examples

Build all examples to generate Honeycomb Query JSON:

```bash
# Build all examples
wetwire-honeycomb build ./examples/...

# Build a specific package
wetwire-honeycomb build ./examples/latency

# Lint examples for issues
wetwire-honeycomb lint ./examples/...

# List all queries
wetwire-honeycomb list ./examples/...
```

## Quick Reference

### Latency Queries

| Query | Purpose |
|-------|---------|
| SlowEndpoints | Find slowest endpoints by P99 |
| LatencyDistribution | Heatmap of latency distribution |
| LatencyByRegion | Compare latency across regions |
| SlowDatabaseQueries | Identify slow DB operations |

### Error Queries

| Query | Purpose |
|-------|---------|
| ErrorsByService | Error counts per service |
| ErrorRate | Error counts per endpoint |
| ErrorsByType | Group by exception type |
| RecentErrors | Most recent errors |
| HTTPStatusCodes | Status code distribution |

### Traffic Queries

| Query | Purpose |
|-------|---------|
| RequestsByEndpoint | Request volume per endpoint |
| TrafficByService | Load distribution across services |
| UniqueUsers | Count distinct users |
| TrafficByMethod | GET/POST/PUT distribution |
| ThroughputByRegion | Request rate by region |

### SLI Queries

| Query | Purpose |
|-------|---------|
| AvailabilitySLI | Successful request count |
| LatencySLI | Requests within latency threshold |
| ErrorBudgetByService | Error budget consumption |
| ThroughputSLI | Request throughput |
| ApdexScore | User satisfaction approximation |
| CriticalEndpointHealth | Monitor key endpoints |

### Trace Queries

| Query | Purpose |
|-------|---------|
| SlowTraces | Slowest end-to-end traces |
| ServiceDependencies | Service call patterns |
| SpansByService | Span volume per service |
| TraceErrors | Traces containing errors |
| SpanDuration | Operation duration analysis |

## Common Patterns

### Percentile Comparison

```go
Calculations: []query.Calculation{
    query.P99("duration_ms"),
    query.P95("duration_ms"),
    query.P50("duration_ms"),
},
```

### Filtering to Recent Data

```go
TimeRange: query.Minutes(30),  // Last 30 minutes
TimeRange: query.Hours(1),     // Last hour
TimeRange: query.Days(7),      // Last 7 days (max)
```

### Grouping by Service

```go
Breakdowns: []string{"service.name", "http.route"},
```

### Counting Distinct Values

```go
Calculations: []query.Calculation{
    query.CountDistinct("user.id"),
},
```

### Filtering by Status Code

```go
Filters: []query.Filter{
    query.GTE("http.status_code", 500),  // Server errors
    query.GTE("http.status_code", 400),  // All errors
    query.LT("http.status_code", 400),   // Successes only
},
```

## Next Steps

1. Copy an example as a starting point
2. Customize for your dataset and fields
3. Run `wetwire-honeycomb lint` to check for issues
4. Run `wetwire-honeycomb build` to generate JSON
5. Use the JSON with Honeycomb's Query API
