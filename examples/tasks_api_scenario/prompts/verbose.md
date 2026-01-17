Create Honeycomb observability resources for a Task API service (CRUD operations on /tasks endpoints).

Dataset: `tasks-api`

## Create these files:

**expected/queries/queries.go:**
- RequestLatency: P99/P95/P50 percentiles on duration_ms, breakdown by http.route and http.method
- ErrorRate: COUNT of requests where http.status_code >= 400, breakdown by route and status_code
- SlowRequests: COUNT/AVG/MAX of duration_ms where duration > 500ms
- RequestThroughput: COUNT over time with 5-minute granularity

**expected/slos/slos.go:**
- Availability: 99.9% target over 30-day window, SLI is requests with status < 500
- Latency: 95% target over 7-day window, SLI is requests with duration < 500ms
- Include burn rate alerts for both SLOs

**expected/triggers/triggers.go:**
- HighErrorRate: Alert when error rate exceeds 1%, check every 2 minutes
- HighLatency: Alert when P99 latency exceeds 1000ms, check every 2 minutes

**expected/boards/dashboard.go:**
- TasksAPIDashboard: Dashboard with panels for all 4 queries in a 2x2 grid layout
