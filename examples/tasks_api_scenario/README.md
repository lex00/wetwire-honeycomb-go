# Task API Scenario

Example Honeycomb observability resources for a hypothetical Task API service.

## Dataset

`tasks-api`

## Resources

| Type | Count | Description |
|------|-------|-------------|
| Queries | 4 | RequestLatency, ErrorRate, SlowRequests, RequestThroughput |
| SLOs | 2 | Availability (99.9%), Latency (95% < 500ms) |
| Triggers | 2 | HighErrorRate, HighLatency |
| Boards | 1 | TasksAPIDashboard |

## Usage

```bash
# List all resources
wetwire-honeycomb list ./examples/tasks_api_scenario/expected/...

# Lint resources
wetwire-honeycomb lint ./examples/tasks_api_scenario/expected/...

# Build to JSON
wetwire-honeycomb build ./examples/tasks_api_scenario/expected/...
```
