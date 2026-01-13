package full_stack

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// HighLatencyAlert triggers when P99 latency exceeds 2 seconds.
// Evaluates every 5 minutes and sends notifications to the performance team.
//
// References:
//   - Query: SlowRequests (queries.go)
var HighLatencyAlert = trigger.Trigger{
	Name:        "High Latency Alert",
	Description: "Alerts when P99 latency exceeds 2000ms, indicating performance degradation",
	Dataset:     "production",
	Query:       SlowRequests,
	Threshold:   trigger.GreaterThan(2000),
	Frequency:   trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#performance"),
		trigger.EmailAddress("performance-team@example.com"),
	},
	Disabled: false,
}

// ErrorRateAlert triggers when error rate exceeds 5% of total traffic.
// Evaluates every 2 minutes for rapid incident detection.
//
// References:
//   - Query: ErrorRate (queries.go)
var ErrorRateAlert = trigger.Trigger{
	Name:        "High Error Rate Alert",
	Description: "Alerts when 5xx error count exceeds 50 requests per minute",
	Dataset:     "production",
	Query:       ErrorRate,
	Threshold:   trigger.GreaterThan(50),
	Frequency:   trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#oncall"),
		trigger.PagerDutyService("api-oncall-service"),
	},
	Disabled: false,
}

// LowTrafficAlert triggers when request volume drops unexpectedly.
// Helps detect upstream issues or traffic routing problems.
//
// References:
//   - Query: RequestThroughput (queries.go)
var LowTrafficAlert = trigger.Trigger{
	Name:        "Low Traffic Alert",
	Description: "Alerts when request rate drops below 100 req/sec, indicating potential traffic routing issues",
	Dataset:     "production",
	Query:       RequestThroughput,
	Threshold:   trigger.LessThan(100),
	Frequency:   trigger.Minutes(10),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#infrastructure"),
	},
	Disabled: false,
}

// CriticalEndpointLatency triggers on latency issues for critical endpoints.
// More aggressive threshold (500ms) for high-priority API routes.
var CriticalEndpointLatency = trigger.Trigger{
	Name:        "Critical Endpoint Latency",
	Description: "Alerts when critical endpoints exceed 500ms P99 latency",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(15),
		Breakdowns: []string{"http.route"},
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
			query.Count(),
		},
		Filters: []query.Filter{
			query.In("http.route", []any{
				"/api/v1/checkout",
				"/api/v1/payment",
				"/api/v1/auth/login",
			}),
			query.GT("duration_ms", 500),
		},
		Orders: []query.Order{
			{Op: "P99", Column: "duration_ms", Order: "descending"},
		},
		Limit: 10,
	},
	Threshold: trigger.GreaterThan(500),
	Frequency: trigger.Minutes(3),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#oncall"),
		trigger.PagerDutyService("critical-api-service"),
	},
	Disabled: false,
}
