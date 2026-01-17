// Package triggers provides example Honeycomb trigger declarations for a Task API service.
package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/examples/tasks_api_scenario/expected/queries"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// HighErrorRate triggers when error rate exceeds 1%.
// Evaluates every 2 minutes to provide rapid detection of service failures.
var HighErrorRate = trigger.Trigger{
	Name:        "High Error Rate",
	Description: "Error rate exceeds 1% threshold",
	Dataset:     "tasks-api",
	Query:       queries.ErrorRate,
	Threshold:   trigger.GreaterThan(1.0),
	Frequency:   trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#alerts"),
	},
	Disabled: false,
}

// HighLatency triggers when P99 latency exceeds 1000ms.
// Helps detect performance degradation before it impacts users.
var HighLatency = trigger.Trigger{
	Name:        "High Latency",
	Description: "P99 latency exceeds 1000ms threshold",
	Dataset:     "tasks-api",
	Query:       queries.RequestLatency,
	Threshold:   trigger.GreaterThan(1000),
	Frequency:   trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#alerts"),
	},
	Disabled: false,
}
