package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/examples/ecommerce_scenario/expected/queries"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// HighLatencyAlert triggers when P99 latency exceeds 3 seconds for checkout flow services.
// This alert helps detect performance degradation in critical checkout paths.
//
// References:
//   - Query: CheckoutFlowLatency (queries/latency.go)
//   - Threshold: P99 duration > 3000ms
//   - Evaluation: Every 2 minutes
var HighLatencyAlert = trigger.Trigger{
	Name:        "Checkout High Latency",
	Description: "P99 latency exceeds 3 seconds in checkout flow services",
	Dataset:     "otel-demo",
	Query:       queries.CheckoutFlowLatency,
	Threshold:   trigger.GreaterThan(3000),
	Frequency:   trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#checkout-alerts"),
		trigger.EmailAddress("checkout-team@example.com"),
	},
	Disabled: false,
}
