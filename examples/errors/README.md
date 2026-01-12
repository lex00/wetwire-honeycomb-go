# Error Tracking Examples

This package demonstrates error tracking and analysis queries for Honeycomb.

## Queries

### ErrorsByService

Counts errors grouped by service name. Quickly identifies which services have issues.

**Use case:** During incidents, identify the source of errors.

### ErrorRate

Calculates error counts by endpoint to identify problematic routes.

**Use case:** Track which API endpoints have the highest error rates.

### ErrorsByType

Groups errors by exception type to identify patterns.

**Use case:** Prioritize which error types to fix first.

### RecentErrors

Shows the most recent errors with full context.

**Use case:** Debug recent issues by seeing error patterns.

### HTTPStatusCodes

Distribution of all HTTP status codes.

**Use case:** Quick health check of API responses.

## Output

```json
{
  "time_range": 3600,
  "breakdowns": ["service.name"],
  "calculations": [{"op": "COUNT"}],
  "filters": [
    {"column": "http.status_code", "op": ">=", "value": 500}
  ],
  "orders": [{"op": "COUNT", "order": "descending"}],
  "limit": 20
}
```

## Dashboard Integration

These queries work well together in a dashboard:

1. **Overview panel:** HTTPStatusCodes for quick health
2. **Service breakdown:** ErrorsByService for drill-down
3. **Error details:** ErrorsByType for root cause analysis
