# Task API Scenario

Honeycomb observability scenario for a Task API service.

## What This Demonstrates

**wetwire** is an AI-native infrastructure architecture. Instead of bolting AI onto existing tools, the entire system is designed to make AI agents reliable:

- **Typed abstraction layer** — Agents work with Go structs, not raw JSON. Types encode valid structures. Agents can't hallucinate field names.
- **Deterministic synthesis** — Same declarations produce the same output. Every time.
- **Semantic-level authoring** — Agents operate on *what you want*, not formatting details.

This scenario tests whether different AI personas (beginner, intermediate, expert) can reliably generate the same correct observability resources. The typed API constrains output to valid structures, producing consistent results regardless of how the request is phrased.

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
├── scenario.yaml        # Config: model, timeout, validation rules
├── system_prompt.md     # Claude's context (API patterns, SLO guidelines)
├── prompts/             # User prompts (one per persona)
│   ├── beginner.md      # Conversational style
│   ├── intermediate.md  # Structured, default
│   └── expert.md        # Terse, technical
├── expected/            # Reference implementation (gold standard)
│   ├── queries/
│   ├── slos/
│   ├── triggers/
│   └── boards/
└── results/             # Scenario output (gitignored)
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
