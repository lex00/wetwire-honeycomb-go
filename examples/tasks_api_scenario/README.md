# Task API Scenario

Honeycomb observability scenario for a Task API service.

## Running the Scenario

Scenarios use the **Claude CLI** (not the Anthropic API directly). The scenario runner invokes `claude` which handles its own authentication.

```bash
# From wetwire-core-go directory:
go run ./cmd/run_scenario /path/to/tasks_api_scenario [persona] [flags]

# Run single persona (default: intermediate)
go run ./cmd/run_scenario ./examples/tasks_api_scenario beginner --verbose

# Run all personas
go run ./cmd/run_scenario ./examples/tasks_api_scenario --all --verbose
```

### Personas

| Persona | Style | Timeout |
|---------|-------|---------|
| beginner | Conversational, minimal technical terms | 2 min |
| intermediate | Structured, uses domain terminology | 2 min |
| expert | Terse, assumes full domain knowledge | 2 min |

## Expected Output

| Type | Count | Resources |
|------|-------|-----------|
| Queries | 4 | RequestLatency, ErrorRate, SlowRequests, RequestThroughput |
| SLOs | 2 | Availability (99.9%), Latency (95% < 500ms) |
| Triggers | 2 | HighErrorRate, HighLatency |
| Boards | 1 | TasksAPIDashboard |

## File Structure

```
tasks_api_scenario/
├── scenario.yaml        # Scenario config (model, timeout, validation)
├── system_prompt.md     # Domain knowledge for Claude
├── prompt.md            # Default prompt
├── prompts/             # Persona-specific prompts
│   ├── beginner.md
│   ├── intermediate.md
│   └── expert.md
├── expected/            # Gold standard implementation
│   ├── queries/
│   ├── slos/
│   ├── triggers/
│   └── boards/
└── results/             # Generated output (gitignored)
```

## Validating Output

```bash
# List discovered resources
wetwire-honeycomb list ./examples/tasks_api_scenario/expected/...

# Lint resources
wetwire-honeycomb lint ./examples/tasks_api_scenario/expected/...

# Build to JSON
wetwire-honeycomb build ./examples/tasks_api_scenario/expected/...
```
