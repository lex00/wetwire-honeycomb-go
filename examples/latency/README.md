# Latency Analysis Examples

This package demonstrates common latency analysis queries for Honeycomb.

## Queries

### SlowEndpoints

Identifies the slowest endpoints by P99 latency, grouped by route and service.

**Use case:** Find which endpoints need performance optimization.

```bash
wetwire-honeycomb build ./examples/latency
```

### LatencyDistribution

Shows the full latency distribution using a heatmap visualization.

**Use case:** Understand latency patterns beyond simple percentiles.

### LatencyByRegion

Compares latency metrics across different cloud regions.

**Use case:** Identify geographic performance issues or misconfigured deployments.

### SlowDatabaseQueries

Finds slow database operations by query statement.

**Use case:** Identify database bottlenecks and optimize queries.

## Output

Run `wetwire-honeycomb build` to generate Honeycomb Query JSON:

```json
{
  "time_range": 7200,
  "breakdowns": ["http.route", "service.name"],
  "calculations": [
    {"op": "P99", "column": "duration_ms"},
    {"op": "P95", "column": "duration_ms"},
    {"op": "P50", "column": "duration_ms"},
    {"op": "COUNT"}
  ],
  "filters": [
    {"column": "http.route", "op": "exists"}
  ],
  "orders": [
    {"op": "P99", "column": "duration_ms", "order": "descending"}
  ],
  "limit": 50
}
```
