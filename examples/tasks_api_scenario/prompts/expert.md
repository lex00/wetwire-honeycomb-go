Dataset: `tasks-api`. Generate these JSON files:

**query-latency.json:**
- P99/P95/P50 on duration_ms, breakdown http.route+http.method

**query-errors.json:**
- COUNT where status>=400, breakdown route+status_code

**slo-availability.json:**
- 99.9%, 30d, SLI status<500

**slo-latency.json:**
- 95%, 7d, SLI duration<500ms

**trigger-error-rate.json:**
- Error rate >1%, 2m frequency, slack #alerts

**trigger-latency.json:**
- P99 >1000ms, 2m frequency, slack #alerts

**board-dashboard.json:**
- Tasks API Dashboard with latency and error panels
