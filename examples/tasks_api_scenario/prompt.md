# Task API Observability

Create Honeycomb observability resources for a Task API service.

## Requirements

### Queries (4)

1. **Request Latency** - Track P50/P95/P99 latency by endpoint
2. **Error Rate** - Count errors by status code and endpoint
3. **Slow Requests** - Find requests exceeding 500ms
4. **Request Throughput** - Track request volume over time

### SLOs (2)

1. **Availability** - 99.9% of requests must succeed (status < 500)
2. **Latency** - 95% of requests must complete under 500ms

### Triggers (2)

1. **High Error Rate** - Alert when error rate exceeds 1%
2. **High Latency** - Alert when P99 exceeds 1000ms

### Board (1)

1. **Dashboard** - Overview with panels for latency, errors, throughput, and slow requests

## Expected Outputs

- `expected/queries/queries.go`
- `expected/slos/slos.go`
- `expected/triggers/triggers.go`
- `expected/boards/dashboard.go`
