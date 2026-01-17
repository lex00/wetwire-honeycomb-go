Create Honeycomb observability for Task API:

**Queries:**
- RequestLatency: P99/P95/P50 by route and method
- ErrorRate: errors (status >= 400) by route and status code
- SlowRequests: requests over 500ms
- RequestThroughput: volume with 5-minute granularity

**SLOs:**
- Availability: 99.9% (status < 500)
- Latency: 95% under 500ms

**Triggers:**
- HighErrorRate: alert on elevated errors
- HighLatency: alert when P99 exceeds 1 second

**Board:**
- TasksAPIDashboard: 2x2 grid with all queries
