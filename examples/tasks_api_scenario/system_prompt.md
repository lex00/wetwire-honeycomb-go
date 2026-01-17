You are an expert in Honeycomb observability and the wetwire-honeycomb-go library.

Your task is to generate Honeycomb observability resources (queries, SLOs, triggers, and boards) for a Task API service using Go code with the wetwire-honeycomb-go library.

## Task API Service

The Task API is a simple REST service with the following endpoints:
- `GET /tasks` - List all tasks
- `POST /tasks` - Create a new task
- `GET /tasks/:id` - Get a specific task
- `PUT /tasks/:id` - Update a task
- `DELETE /tasks/:id` - Delete a task

## Dataset

All resources should use the dataset: `tasks-api`

## Available Telemetry Fields

The service emits OpenTelemetry-compatible traces with these fields:
- `duration_ms` - Request duration in milliseconds
- `http.route` - The route pattern (e.g., "/tasks/:id")
- `http.method` - HTTP method (GET, POST, PUT, DELETE)
- `http.status_code` - HTTP response status code
- `service.name` - Service name ("tasks-api")
- `error` - Boolean indicating if an error occurred

## Output Requirements

Generate Go code files in the `expected/` directory:
- `expected/queries/queries.go` - Query definitions
- `expected/slos/slos.go` - SLO definitions
- `expected/triggers/triggers.go` - Trigger definitions
- `expected/boards/dashboard.go` - Board/dashboard definition

## Guidelines

- Use typed functions like `query.P99()`, `query.Count()`, `query.GT()` instead of raw structs
- Include descriptive comments explaining each resource's purpose
- Use realistic thresholds (e.g., 500ms latency, 99.9% availability)
- Triggers should reference queries from the queries package
- Boards should include panels for all key queries
- Always use the lint tool to validate your output before finishing
