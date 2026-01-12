# Distributed Tracing Examples

This package demonstrates trace analysis queries for Honeycomb.

## Queries

### SlowTraces

Finds the slowest root spans (end-to-end traces).

**Use case:** Identify the slowest user-facing requests.

### ServiceDependencies

Maps service-to-service call patterns.

**Use case:** Understand service topology and dependencies.

### SpansByService

Shows span volume per service.

**Use case:** Identify which services generate the most trace data.

### TraceErrors

Finds traces containing errors.

**Use case:** Debug error propagation through the system.

### SpanDuration

Analyzes operation duration by span name.

**Use case:** Find the slowest operations within traces.

## Trace Analysis Patterns

### Finding Slow Traces

```go
// Start with SlowTraces to find candidates
// Then use trace.trace_id to drill down
```

### Service Dependency Mapping

```go
// ServiceDependencies query shows caller -> callee relationships
// Visualize as a service graph
```

### Error Investigation

```go
// 1. Use TraceErrors to find error traces
// 2. Filter by trace.trace_id for full trace view
// 3. Check SpanDuration for performance context
```

## Output

```json
{
  "time_range": 3600,
  "breakdowns": ["name", "service.name"],
  "calculations": [
    {"op": "P99", "column": "duration_ms"},
    {"op": "COUNT"}
  ],
  "filters": [
    {"column": "trace.parent_id", "op": "=", "value": null},
    {"column": "duration_ms", "op": ">", "value": 1000}
  ],
  "orders": [{"op": "P99", "column": "duration_ms", "order": "descending"}],
  "limit": 25
}
```

## OpenTelemetry Field Reference

Common fields used in these queries:
- `trace.trace_id` - Unique trace identifier
- `trace.parent_id` - Parent span (null for root)
- `name` - Span operation name
- `service.name` - Service that created the span
- `peer.service` - Called service (for client spans)
- `status_code` - Span status (OK, ERROR)
- `duration_ms` - Span duration in milliseconds
