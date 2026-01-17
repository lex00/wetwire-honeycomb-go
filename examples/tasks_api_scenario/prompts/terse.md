Dataset: `tasks-api`

**expected/queries/queries.go:** RequestLatency (P99/P95/P50), ErrorRate, SlowRequests (>500ms), RequestThroughput

**expected/slos/slos.go:** Availability (99.9%), Latency (95%<500ms)

**expected/triggers/triggers.go:** HighErrorRate (>1%), HighLatency (P99>1s)

**expected/boards/dashboard.go:** TasksAPIDashboard
