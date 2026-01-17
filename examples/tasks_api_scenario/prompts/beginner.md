I need monitoring for my Task API. Dataset: `tasks-api`

Please create these files:

**expected/queries/queries.go:**
- RequestLatency - track response times
- ErrorRate - count errors
- SlowRequests - find slow requests
- RequestThroughput - track volume

**expected/slos/slos.go:**
- Availability - 99.9% uptime target
- Latency - 95% under 500ms

**expected/triggers/triggers.go:**
- HighErrorRate - alert on errors
- HighLatency - alert on slow responses

**expected/boards/dashboard.go:**
- TasksAPIDashboard - overview dashboard
