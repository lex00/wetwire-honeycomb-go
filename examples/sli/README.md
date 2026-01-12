# SLI/SLO Examples

This package demonstrates Service Level Indicator queries for Honeycomb.

## Queries

### AvailabilitySLI

Counts successful requests (status < 500) for availability calculation.

**SLO Example:** 99.9% availability = max 43 minutes downtime/month

### LatencySLI

Counts requests meeting latency threshold (< 500ms).

**SLO Example:** 95% of requests complete in < 500ms

### ErrorBudgetByService

Tracks error counts per service over 7 days.

**Use case:** Identify which services are consuming error budget fastest.

### ThroughputSLI

Simple request count for capacity SLI tracking.

**Use case:** Alert when throughput drops below expected baseline.

### ApdexScore

Approximates user satisfaction via latency thresholds.

**Calculation:** (Satisfied + Tolerating/2) / Total

### CriticalEndpointHealth

Focused monitoring for revenue-critical endpoints.

**Use case:** Prioritized alerting for business-critical paths.

## SLO Dashboard Pattern

Combine these queries for a complete SLO dashboard:

```
+------------------+------------------+
|  Availability    |  Latency SLI     |
|  (AvailabilitySLI)  (LatencySLI)    |
+------------------+------------------+
|  Error Budget by Service           |
|  (ErrorBudgetByService)            |
+------------------------------------|
|  Critical Endpoint Health          |
|  (CriticalEndpointHealth)          |
+------------------------------------+
```

## Error Budget Calculations

Use these queries to calculate remaining error budget:

1. Query total requests with AvailabilitySLI
2. Query error requests (http.status_code >= 500)
3. Calculate: ErrorBudget = (1 - SLO) * TotalRequests - ActualErrors

## Output

```json
{
  "time_range": 86400,
  "calculations": [{"op": "COUNT"}],
  "filters": [
    {"column": "http.status_code", "op": "<", "value": 500}
  ]
}
```
