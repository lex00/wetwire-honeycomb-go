# E2E Demo: Presenter Script

This is a timed presenter script for demonstrating wetwire-honeycomb's AI-assisted observability-as-code workflow.

**Total Duration:** 8-10 minutes

---

## Setup (Before Demo)

### Prerequisites

- Go 1.23+ installed
- wetwire-honeycomb CLI installed (`go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest`)
- Kiro CLI installed and authenticated (see [docs/HONEYCOMB-KIRO-CLI.md](../../docs/HONEYCOMB-KIRO-CLI.md))
- Terminal with split-pane view (recommended)

### Terminal Setup

Open 2-3 terminal panes:

1. **Main pane:** For running commands
2. **Watch pane:** For showing generated files (optional)
3. **Output pane:** For displaying results (optional)

### What to Have Ready

```bash
# Navigate to demo directory
cd examples/ecommerce_scenario

# Clear any previous output (if doing multiple demos)
rm -rf generated/ output/

# Have the prompt ready to cat
ls prompts/complex_scenario.txt

# Test that wetwire-honeycomb is installed
wetwire-honeycomb --version
```

---

## Demo Script (8-10 minutes)

### Opening (1 minute)

**SAY:**
> "Let's talk about observability-as-code. Most teams use Honeycomb's UI or Terraform to define queries, SLOs, and dashboards. But there's a problem."

**SHOW:** Open browser to Honeycomb MCP documentation (optional) or just state:

**SAY:**
> "Honeycomb's MCP server is read-only - AI agents can query data but can't create new queries or SLOs. And Terraform? It's powerful but not LLM-friendly. The syntax is complex and error-prone."

**TRANSITION:**
> "wetwire-honeycomb solves this. It's observability-as-code designed for agentic workflows - type-safe Go structs that compile to Honeycomb JSON."

---

### Show the Prompt (2 minutes)

**SAY:**
> "Here's our scenario: We're monitoring an e-commerce checkout flow. Payment processing, fraud detection, cart service - the whole critical path."

**RUN:**
```bash
cat prompts/complex_scenario.txt
```

**SAY (while scrolling through the prompt):**
> "This is natural language - the kind of requirements you'd get from a product manager or SRE team. We want:
> - 4 queries to track latency, errors, fraud correlation, and funnel analysis
> - 3 SLOs for availability, latency, and payment success
> - 2 triggers for high latency and error spikes
> - 1 dashboard to visualize everything"

**HIGHLIGHT (scroll to specific sections):**
> "Notice the specifics: P99/P95/P50 calculations, specific services, time ranges. This is the kind of domain knowledge that needs to translate into code."

---

### Run the Demo (3 minutes)

**SAY:**
> "Now watch this. We're going to use Kiro - Anthropic's agentic CLI - with wetwire-honeycomb's MCP server. The AI will generate type-safe Go code from this natural language prompt."

**RUN:**
```bash
wetwire-honeycomb test --provider kiro --persona intermediate "$(cat prompts/complex_scenario.txt)"
```

**SAY (as it runs):**
> "What's happening here:
> 1. The AI reads the prompt through Kiro's agent framework
> 2. It generates Go code using wetwire-honeycomb's type-safe builders
> 3. It runs the linter to catch issues
> 4. It fixes any problems and re-validates
> 5. It compiles the Go code to Honeycomb Query JSON"

**WATCH FOR:**
- Initial code generation messages
- Lint/fix cycles (if any)
- Success indicators
- Final score output

**SAY (when complete):**
> "There we go. The AI scored 85-100% based on how well it matched our requirements. Let's see what it created."

---

### Examine Output (2-3 minutes)

**RUN:**
```bash
# Show generated Go files
ls -la generated/

# Show one of the query files
cat generated/queries.go
```

**SAY (while showing queries.go):**
> "Look at this. Type-safe Go structs - not YAML, not JSON, not HCL. This is code that:
> - Compiles at build time
> - Gets caught by the linter before it ever reaches production
> - Is impossible to typo-squat because Go won't compile 'service.nam' instead of 'service.name'
> - Has inline documentation and examples"

**HIGHLIGHT specific code:**
```go
var CheckoutFlowLatency = query.Query{
    Dataset:   "otel-demo",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.P99("duration_ms"),  // <- Type-safe builder
        query.P95("duration_ms"),
    },
    Filters: []query.Filter{
        query.In("service.name", []any{  // <- Type-safe filter
            "checkoutservice",
            "paymentservice",
        }),
    },
}
```

**SAY:**
> "Notice `query.P99()` and `query.In()` - these are typed functions, not strings. The LLM can't hallucinate a calculation type that doesn't exist."

**RUN:**
```bash
# Show compiled JSON output
cat output/queries.json | head -40
```

**SAY:**
> "And here's the Honeycomb Query JSON that gets sent to the API. Ready to deploy."

---

### Competitive Positioning (2 minutes)

**SAY:**
> "Why does this matter? Let's compare approaches."

| Approach | AI-Friendly? | Type-Safe? | Compile-Time Validation? |
|----------|--------------|------------|-------------------------|
| **Honeycomb MCP** | Read-only | N/A | N/A |
| **Terraform** | ⚠️ Complex | ❌ Strings | ❌ Runtime only |
| **wetwire-honeycomb** | ✅ Yes | ✅ Yes | ✅ Yes |

**SAY:**
> "Honeycomb MCP is great for querying but can't create new configs. Terraform works but it's too complex for LLMs to reliably generate - you end up with syntax errors, invalid references, and runtime failures."

**SAY:**
> "wetwire-honeycomb is designed for agentic workflows:
> - **Flat struct declarations** - LLMs understand `var X = Struct{}` patterns
> - **Type-safe builders** - Prevent hallucinations (can't invent `query.P999()`)
> - **Compile-time validation** - Catch errors before deployment
> - **Domain-specific lints** - Encode observability expertise (WHC001-WHC005 rules)"

---

### Closing (1 minute)

**SAY:**
> "So to recap: We took a natural language prompt, used an AI agent with wetwire-honeycomb, and generated production-ready observability code in under 2 minutes. The code is type-safe, linted, and ready to deploy to Honeycomb."

**TRANSITION to Q&A:**
> "Questions?"

---

## Talking Points

### Key Differentiators

1. **"Flat struct declarations are LLM-friendly"**
   - Simple `var X = Struct{Field: Value}` patterns
   - No nested definitions or complex inheritance
   - Top-level declarations are easy to discover and modify

2. **"Type-safe builders prevent hallucinations"**
   - `query.P99()` exists or code won't compile
   - Can't typo `query.P999()` or `query.Percentile99()`
   - Auto-complete in IDEs shows available options

3. **"Compile-time validation catches errors before runtime"**
   - Invalid field references = compilation error
   - Type mismatches caught immediately
   - No "deploy and pray" workflows

4. **"Domain-specific lint rules encode expertise"**
   - WHC001: Use typed calculations, not strings
   - WHC002: Use direct filter functions, not raw maps
   - WHC003: Validate dataset references
   - WHC004: Check time range validity
   - WHC005: Ensure calculation names are unique

### When to Use This Demo

- **AI/ML conferences** - Showcase agentic workflows
- **Observability talks** - Infrastructure-as-code for monitoring
- **SRE/DevOps teams** - Reduce toil in query management
- **Sales/partnerships** - Honeycomb + AI integration story

### Elevator Pitch (30 seconds)

> "wetwire-honeycomb is observability-as-code for AI agents. It generates Honeycomb queries, SLOs, and dashboards from natural language using type-safe Go structs. Unlike Terraform or raw JSON, it's designed for LLMs - with typed builders that prevent hallucinations and compile-time validation that catches errors before deployment."

---

## Q&A Prep

### Common Questions

**Q: "Why Go? Why not Python or TypeScript?"**

A: Go compiles to a single binary, has excellent type safety, and is widely used in infrastructure tooling. The simplicity of Go structs makes it easy for LLMs to generate correct code. Plus, Go's tooling (go build, go test, go vet) provides robust validation.

**Q: "Does this replace Terraform?"**

A: No, it's complementary. Use Terraform for infrastructure (compute, networking, databases). Use wetwire-honeycomb for observability configs (queries, SLOs, triggers). You can even call wetwire-honeycomb from Terraform using local-exec provisioners.

**Q: "What if I don't use Kiro? Can I use other AI providers?"**

A: Yes! The `--provider` flag supports:
- `kiro` (Anthropic's CLI)
- `anthropic` (direct API)
- `openai` (OpenAI API)
- `bedrock` (AWS Bedrock)

You can also generate code manually and use wetwire-honeycomb as a pure CLI tool without AI.

**Q: "How does this handle secrets? API keys?"**

A: wetwire-honeycomb generates code and JSON - it doesn't execute queries or store credentials. You manage Honeycomb API keys separately (environment variables, secret managers, etc.) when you actually deploy the JSON to Honeycomb.

**Q: "Can I import existing Honeycomb queries into wetwire-honeycomb?"**

A: Yes! Use `wetwire-honeycomb import` to convert Honeycomb Query JSON into Go code:

```bash
wetwire-honeycomb import --input existing_query.json --output queries.go
```

**Q: "What about SLOs, Triggers, and Boards?"**

A: Full support! This demo shows queries, but wetwire-honeycomb also supports:
- **SLOs** (`slo.SLO` type)
- **Triggers** (`trigger.Trigger` type)
- **Boards** (`board.Board` type with panels)

See `examples/full_stack/` for complete examples.

**Q: "How do you handle breaking changes in Honeycomb's API?"**

A: wetwire-honeycomb tracks Honeycomb's Query Specification. When Honeycomb adds new features:
1. We add typed builders (e.g., `query.P999()` for new percentiles)
2. Update linter rules if needed
3. Publish a new version

Existing code continues to work because Go is compiled.

**Q: "What's the test command doing under the hood?"**

A: The `test` command:
1. Spins up the wetwire-honeycomb MCP server
2. Connects the AI provider (Kiro/Anthropic/etc.) to the MCP server
3. Sends the prompt to the AI
4. The AI uses MCP tools to generate code, lint, and validate
5. Scores the output against expected results (if `expected/` exists)

It's essentially an end-to-end integration test for AI-generated code.

### Edge Cases to Be Aware Of

1. **Large prompts**: Very complex scenarios (20+ queries) may hit token limits. Break them into multiple prompts.

2. **Dataset validation**: wetwire-honeycomb doesn't validate that datasets exist in your Honeycomb account - it just generates JSON. You'll get errors when you try to use the query if the dataset doesn't exist.

3. **Field name typos**: The linter can't validate that `duration_ms` exists in your dataset's schema. It only validates Go syntax and structure.

4. **Rate limiting**: When using AI providers, be aware of API rate limits (especially for OpenAI/Anthropic direct APIs).

5. **Non-deterministic output**: LLMs may generate slightly different code each time. That's expected. As long as it passes linting and compiles, it's valid.

---

## Troubleshooting

### Demo fails with "kiro-cli not found"

**Fix:**
```bash
curl -fsSL https://cli.kiro.dev/install | bash
kiro-cli login
```

### AI generates invalid Go code

**Fix:**
- Check the prompt is clear and specific
- Try a different persona (`--persona expert` vs `--persona intermediate`)
- Run `wetwire-honeycomb lint` manually to see specific errors

### Compilation errors

**Fix:**
```bash
cd generated/
go mod init ecommerce_scenario
go mod tidy
go build
```

### MCP server not starting

**Fix:**
- Check that `wetwire-honeycomb mcp` runs without errors
- Verify Kiro agent config at `~/.kiro/agents/wetwire-honeycomb-runner.json`
- Check logs in `~/.kiro/logs/`

---

## Additional Resources

- [Full Documentation](../../README.md)
- [Kiro CLI Integration Guide](../../docs/HONEYCOMB-KIRO-CLI.md)
- [Lint Rules Reference](../../docs/LINT_RULES.md)
- [Example Queries](../../examples/)
- [FAQ](../../docs/FAQ.md)

---

**Pro Tips for Presenters:**

1. **Practice the timing** - Run through the script 2-3 times to hit the 8-10 minute mark
2. **Have a backup** - Pre-generate code in case of network issues
3. **Customize the prompt** - Adapt the scenario to your audience (e.g., financial services, healthcare)
4. **Show the linter** - Run `wetwire-honeycomb lint` manually to demo the validation
5. **Zoom in** - Terminal font should be 16-18pt for visibility
6. **Record it** - This makes a great async demo video

---

Last updated: 2026-01-14
