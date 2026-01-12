# Kiro CLI Integration

Use Kiro CLI with wetwire-honeycomb for AI-assisted query design and observability configuration.

## Prerequisites

- Go 1.23+ installed
- Kiro CLI installed ([installation guide](https://kiro.dev/docs/cli/installation/))
- AWS Builder ID or GitHub/Google account (for Kiro authentication)

---

## Step 1: Install wetwire-honeycomb

### Option A: Using Go (recommended)

```bash
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
```

### Option B: Pre-built binaries

Download from [GitHub Releases](https://github.com/lex00/wetwire-honeycomb-go/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/lex00/wetwire-honeycomb-go/releases/latest/download/wetwire-honeycomb-darwin-arm64
chmod +x wetwire-honeycomb-darwin-arm64
sudo mv wetwire-honeycomb-darwin-arm64 /usr/local/bin/wetwire-honeycomb

# macOS (Intel)
curl -LO https://github.com/lex00/wetwire-honeycomb-go/releases/latest/download/wetwire-honeycomb-darwin-amd64
chmod +x wetwire-honeycomb-darwin-amd64
sudo mv wetwire-honeycomb-darwin-amd64 /usr/local/bin/wetwire-honeycomb

# Linux (x86-64)
curl -LO https://github.com/lex00/wetwire-honeycomb-go/releases/latest/download/wetwire-honeycomb-linux-amd64
chmod +x wetwire-honeycomb-linux-amd64
sudo mv wetwire-honeycomb-linux-amd64 /usr/local/bin/wetwire-honeycomb
```

### Verify installation

```bash
wetwire-honeycomb --version
```

---

## Step 2: Install Kiro CLI

```bash
# Install Kiro CLI
curl -fsSL https://cli.kiro.dev/install | bash

# Verify installation
kiro-cli --version

# Sign in (opens browser)
kiro-cli login
```

---

## Step 3: Configure Kiro for wetwire-honeycomb

Run the design command with `--provider kiro` to auto-configure:

```bash
# Create a project directory
mkdir my-queries && cd my-queries

# Initialize Go module
go mod init my-queries

# Run design with Kiro provider (auto-installs configs on first run)
wetwire-honeycomb design --provider kiro "Create a query to find slow requests"
```

This automatically installs:

| File | Purpose |
|------|---------|
| `~/.kiro/agents/wetwire-honeycomb-runner.json` | Kiro agent configuration |
| `.kiro/mcp.json` | Project MCP server configuration |

### Manual configuration (optional)

The MCP server is provided as a subcommand `wetwire-honeycomb mcp`. If you prefer to configure manually:

**~/.kiro/agents/wetwire-honeycomb-runner.json:**
```json
{
  "name": "wetwire-honeycomb-runner",
  "description": "Query generator using wetwire-honeycomb",
  "prompt": "You are an observability assistant...",
  "model": "claude-sonnet-4",
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-honeycomb",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  },
  "tools": ["*"]
}
```

**.kiro/mcp.json:**
```json
{
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-honeycomb",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  }
}
```

> **Note:** The `cwd` field ensures MCP tools resolve paths correctly in your project directory. When using `wetwire-honeycomb design --provider kiro`, this is configured automatically.

---

## Step 4: Run Kiro with wetwire design

### Using the wetwire-honeycomb CLI

```bash
# Start Kiro design session
wetwire-honeycomb design --provider kiro "Create queries for P99 latency by endpoint"
```

This launches Kiro CLI with the wetwire-honeycomb-runner agent and your prompt.

### Using Kiro CLI directly

```bash
# Start chat with wetwire-honeycomb-runner agent
kiro-cli chat --agent wetwire-honeycomb-runner

# Or with an initial prompt
kiro-cli chat --agent wetwire-honeycomb-runner "Create a query to track error rates"
```

---

## Available MCP Tools

The wetwire-honeycomb MCP server exposes five tools to Kiro:

| Tool | Description | Example |
|------|-------------|---------|
| `wetwire_init` | Initialize a new project | `wetwire_init(path="./myqueries")` |
| `wetwire_lint` | Lint code for issues | `wetwire_lint(path="./queries/...")` |
| `wetwire_build` | Generate Query JSON | `wetwire_build(path="./queries/...", format="json")` |
| `wetwire_list` | List discovered queries | `wetwire_list(path="./queries/...")` |
| `wetwire_graph` | Generate dependency graph | `wetwire_graph(path="./queries/...", format="dot")` |

---

## Example Session

```
$ wetwire-honeycomb design --provider kiro "Create a query to find slow API requests with error tracking"

Installed Kiro agent config: ~/.kiro/agents/wetwire-honeycomb-runner.json
Installed project MCP config: .kiro/mcp.json
Starting Kiro CLI design session...

> I'll help you create a query to track slow API requests with error information.

Let me initialize the project and create the query code.

[Calling wetwire_init...]
[Calling wetwire_lint...]
[Calling wetwire_build...]

I've created the following files:
- queries/api_performance.go

The query includes:
- P99 latency tracking
- Error rate calculations
- Breakdown by endpoint and status code
- 2-hour time window

Would you like me to add any additional metrics or filters?
```

---

## Workflow

The Kiro agent follows this workflow:

1. **Explore** - Understand your requirements
2. **Plan** - Design the query structure
3. **Implement** - Generate Go code using wetwire-honeycomb patterns
4. **Lint** - Run `wetwire_lint` to check for issues
5. **Build** - Run `wetwire_build` to generate Query JSON

---

## Using Generated Queries

After Kiro generates your query code:

```bash
# Build the Query JSON
wetwire-honeycomb build ./queries > queries.json

# Use with Honeycomb Query API
curl -X POST https://api.honeycomb.io/1/queries/my-dataset \
  -H "X-Honeycomb-Team: $HONEYCOMB_API_KEY" \
  -H "Content-Type: application/json" \
  -d @queries.json

# Or create derived columns programmatically
# (user's responsibility to integrate with Honeycomb API)
```

---

## Troubleshooting

### MCP server not found

```
Mcp error: -32002: No such file or directory
```

**Solution:** Ensure `wetwire-honeycomb` is in your PATH:

```bash
which wetwire-honeycomb

# If not found, add to PATH or reinstall
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
```

### Kiro CLI not found

```
kiro-cli not found in PATH
```

**Solution:** Install Kiro CLI:

```bash
curl -fsSL https://cli.kiro.dev/install | bash
```

### Authentication issues

```
Error: Not authenticated
```

**Solution:** Sign in to Kiro:

```bash
kiro-cli login
```

---

## Known Limitations

### Automated Testing

When using `wetwire-honeycomb test --provider kiro`, tests run in non-interactive mode (`--no-interactive`). This means:

- The agent runs autonomously without waiting for user input
- Persona simulation is limited - all personas behave similarly
- The agent won't ask clarifying questions

For true persona simulation with multi-turn conversations, use the Anthropic provider:

```bash
wetwire-honeycomb test --provider anthropic --persona expert "Create a query for SLI tracking"
```

### Interactive Design Mode

Interactive design mode (`wetwire-honeycomb design --provider kiro`) works fully as expected:

- Real-time conversation with the agent
- Agent can ask clarifying questions
- Lint loop executes as specified in the agent prompt

---

## See Also

- [CLI Reference](CLI.md) - Full wetwire-honeycomb CLI documentation
- [Quick Start](QUICK_START.md) - Getting started with wetwire-honeycomb
- [Kiro CLI Installation](https://kiro.dev/docs/cli/installation/) - Official installation guide
- [Kiro CLI Docs](https://kiro.dev/docs/cli/) - Official Kiro documentation
