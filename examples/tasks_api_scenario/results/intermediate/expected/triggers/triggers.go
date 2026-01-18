package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/examples/tasks_api_scenario/results/intermediate/expected/queries"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// HighErrorRate alerts when error count exceeds threshold
var HighErrorRate = trigger.Trigger{
	Name:        "High Error Rate",
	Description: "Alert when error rate is elevated",
	Query:       queries.ErrorRate,
	Frequency:   trigger.Minutes(2), // Check every 2 minutes
	Threshold:   trigger.GreaterThan(10),
	Recipients: []trigger.Recipient{
		{Type: "slack", Target: "#alerts"},
	},
}

// HighLatency alerts when P99 latency exceeds 1 second
var HighLatency = trigger.Trigger{
	Name:        "High Latency",
	Description: "Alert when P99 latency exceeds 1 second",
	Query:       queries.RequestLatency,
	Frequency:   trigger.Minutes(2), // Check every 2 minutes
	Threshold:   trigger.GreaterThan(1000), // 1 second in milliseconds
	Recipients: []trigger.Recipient{
		{Type: "slack", Target: "#alerts"},
	},
}
