# E2E Demo: Ecommerce Checkout Scenario

This end-to-end demonstration showcases the complete workflow of wetwire-honeycomb-go, from natural language prompts to fully synthesized Honeycomb observability resources. It demonstrates the core value proposition: **prompts → Go code → JSON**.

## Overview

The ecommerce scenario simulates a realistic checkout flow monitoring setup across multiple microservices. This demo shows how wetwire-honeycomb can:

- Transform natural language requirements into type-safe Go code
- Generate queries, SLOs, triggers, and boards for a complex service architecture
- Validate generated code against best practices
- Produce production-ready Honeycomb Query JSON

The scenario focuses on an OpenTelemetry demo application with four interconnected services: `checkoutservice`, `cartservice`, `paymentservice`, and `frauddetectionservice`.

## Quick Start

### Run with intermediate persona

Test the scenario with the intermediate user persona, which generates code suitable for users with moderate Honeycomb experience:

```bash
# From the repository root
wetwire-honeycomb test --persona intermediate \
  "$(cat examples/ecommerce_scenario/prompts/complex_scenario.txt)"
```

### Run all personas

Compare outputs across beginner, intermediate, and expert personas to see how the tool adapts to different user experience levels:

```bash
wetwire-honeycomb test --all-personas \
  "$(cat examples/ecommerce_scenario/prompts/complex_scenario.txt)"
```

### View results

After running tests, results are stored in the `results/` directory organized by persona:

```bash
ls -la examples/ecommerce_scenario/results/
```

## Scenario Description

### Checkout Flow Architecture

The scenario models a typical ecommerce checkout flow with four critical services:

1. **checkoutservice** - Orchestrates the checkout process
2. **cartservice** - Manages shopping cart state
3. **paymentservice** - Processes payment transactions
4. **frauddetectionservice** - Validates transactions for fraud

### Generated Resources

Based on the scenario prompt, wetwire-honeycomb generates:

#### Queries
- **CheckoutFlowLatency** - P50/P95/P99 latency across all checkout services
- **PaymentFraudCorrelation** - Correlation between payment attempts and fraud detection results
- **ErrorRateByService** - Error rates broken down by individual service
- **CheckoutFunnel** - Conversion tracking through the checkout stages

#### SLOs (Service Level Objectives)
- **Availability SLO** - 99.9% uptime target for checkout flow
- **Latency SLO** - 95% of requests complete within 2000ms
- **Payment Success SLO** - 99.95% successful payment processing

#### Triggers
- Latency degradation alerts
- Error rate threshold alerts
- Fraud detection anomaly alerts

#### Boards
- Checkout flow performance dashboard
- Service health overview
- Payment and fraud correlation views

### Dataset Configuration

All queries target the `otel-demo` dataset, configured in `config/services.yaml`:

```yaml
checkout_flow:
  dataset: otel-demo
  services:
    - checkoutservice
    - cartservice
    - paymentservice
    - frauddetectionservice
  slo_targets:
    availability: 99.9%
    latency_p95: 2000ms
    payment_success: 99.95%
```

## File Structure

```
ecommerce_scenario/
├── README.md              # This file
├── config/                # Service and dataset configuration
│   ├── services.yaml      # Service definitions and SLO targets
│   ├── config_test.go     # Validation tests for configuration
│   ├── go.mod             # Go module for config tests
│   └── go.sum             # Go module checksums
├── prompts/               # Input prompts for generation
│   └── complex_scenario.txt  # Natural language scenario description
├── expected/              # Reference implementations
│   └── queries/           # Expected query implementations
│       └── queries_test.go   # Validation tests for expected queries
└── results/               # Generated outputs (created during test runs)
    ├── beginner/          # Beginner persona outputs
    ├── intermediate/      # Intermediate persona outputs
    └── expert/            # Expert persona outputs
```

### Directory Purposes

- **config/** - Defines the service architecture and SLO targets that inform code generation
- **prompts/** - Contains natural language descriptions of monitoring requirements
- **expected/** - Reference implementations that demonstrate correct patterns
- **results/** - Generated code and JSON outputs from test runs (gitignored)

## Comparing Results

### Scoring Metrics

The `wetwire-honeycomb test` command evaluates generated code across multiple dimensions:

1. **Correctness** - Does the code compile and follow Go syntax rules?
2. **Completeness** - Are all required resources (queries, SLOs, triggers, boards) generated?
3. **Best Practices** - Does the code use typed helpers instead of raw maps?
4. **Lint Compliance** - Does the code pass all WHC lint rules?
5. **Resource Accuracy** - Do the resources match the scenario requirements?

### Persona Differences

Each persona demonstrates different levels of sophistication:

#### Beginner Persona
- Generates simpler queries with basic calculations
- Uses straightforward filter conditions
- Includes extensive inline comments and documentation
- May use more verbose patterns for clarity

#### Intermediate Persona
- Balances simplicity with power
- Uses typed calculation helpers (e.g., `query.P99()`)
- Implements multi-dimensional breakdowns
- Includes appropriate filters and ordering

#### Expert Persona
- Generates optimized, concise code
- Leverages advanced query patterns
- Implements complex correlations and derivations
- Uses sophisticated filtering and aggregations

### Interpreting Scores

Higher scores indicate better alignment with:
- Honeycomb best practices
- Type-safe Go patterns
- wetwire-honeycomb idioms
- Production-ready code quality

Compare scores across personas to understand:
- How well each persona handles complex scenarios
- Which persona best matches your team's experience level
- Where generated code may need manual refinement

## Next Steps

After reviewing the demo:

1. **Customize the scenario** - Modify `config/services.yaml` to match your architecture
2. **Write your own prompts** - Create prompts in `prompts/` for your use cases
3. **Run tests** - Use `wetwire-honeycomb test` to generate code from your prompts
4. **Validate results** - Check generated code against expected patterns
5. **Refine prompts** - Iterate on prompts to improve generation quality
6. **Build JSON** - Use `wetwire-honeycomb build` to synthesize Honeycomb Query JSON
7. **Deploy** - Apply the generated JSON to your Honeycomb environment

## Running Tests

### Validate configuration

```bash
cd examples/ecommerce_scenario/config
go test -v
```

### Validate expected queries

```bash
cd examples/ecommerce_scenario/expected/queries
go test -v
```

### Run full E2E test

```bash
# From repository root
wetwire-honeycomb test --persona intermediate \
  "$(cat examples/ecommerce_scenario/prompts/complex_scenario.txt)"
```

## See Also

- [Full Stack Example](../full_stack/README.md) - Complete Query→SLO→Trigger→Board chain
- [CLI Reference](../../docs/CLI.md) - Full wetwire-honeycomb command documentation
- [Kiro Integration](../../docs/HONEYCOMB-KIRO-CLI.md) - AI-assisted query design
- [Quick Start](../../docs/QUICK_START.md) - Getting started guide
- [Lint Rules](../../docs/LINT_RULES.md) - All WHC validation rules
