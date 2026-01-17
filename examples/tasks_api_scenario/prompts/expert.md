Dataset: `tasks-api`. Create these files:

**expected/queries/queries.go:**
- RequestLatency: P99/P95/P50(duration_ms), breakdown http.route+http.method
- ErrorRate: COUNT where status>=400, breakdown route+status_code
- SlowRequests: COUNT/AVG/MAX where duration>500ms
- RequestThroughput: COUNT, 5m granularity

**expected/slos/slos.go:**
- Availability: 99.9%, 30d, SLI status<500, burn alerts (1h/2x, 6h/5x)
- Latency: 95%, 7d, SLI duration<500ms, burn alert (1h/10x)

**expected/triggers/triggers.go:**
- HighErrorRate: >1%, 2m frequency
- HighLatency: P99>1000ms, 2m frequency

**expected/boards/dashboard.go:**
- TasksAPIDashboard: 2x2 grid with all queries
