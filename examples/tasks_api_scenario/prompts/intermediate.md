# Task API Observability

Dataset: `tasks-api` (CRUD API for /tasks endpoints)

## Create these files:

**expected/queries/queries.go:**
- RequestLatency: P99/P95/P50 by route/method
- ErrorRate: COUNT where status>=400, by route/status_code
- SlowRequests: COUNT/AVG/MAX where duration>500ms
- RequestThroughput: COUNT with 5m granularity

**expected/slos/slos.go:**
- Availability: 99.9% (status<500), 30-day window
- Latency: 95% under 500ms, 7-day window

**expected/triggers/triggers.go:**
- HighErrorRate: >1% errors, 2m frequency
- HighLatency: P99>1000ms, 2m frequency

**expected/boards/dashboard.go:**
- TasksAPIDashboard: 2x2 grid with all 4 queries

Use wetwire-honeycomb-go typed functions.
