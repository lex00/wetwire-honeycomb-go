package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/examples/ecommerce_scenario/expected/queries"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// ErrorRateAlert triggers when error rate exceeds 1% in checkout services.
// This alert provides rapid detection of service failures or degradation.
//
// References:
//   - Query: ErrorRateByService (queries/errors.go)
//   - Threshold: Error rate > 1%
//   - Evaluation: Every 2 minutes
var ErrorRateAlert = trigger.Trigger{
	Name:        "High Error Rate",
	Description: "Error rate exceeds 1% threshold in checkout services",
	Dataset:     "otel-demo",
	Query:       queries.ErrorRateByService,
	Threshold:   trigger.GreaterThan(1.0),
	Frequency:   trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#oncall"),
		trigger.PagerDutyService("checkout-oncall-service"),
	},
	Disabled: false,
}
