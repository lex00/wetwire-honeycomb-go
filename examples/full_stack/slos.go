package full_stack

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// APIAvailability defines a 99.9% availability SLO over a 30-day window.
// Tracks the ratio of successful requests (status < 500) to all requests.
//
// References:
//   - GoodEvents: SuccessfulRequests query
//   - TotalEvents: AllRequests query
var APIAvailability = slo.SLO{
	Name:        "API Availability",
	Description: "Tracks the percentage of successful HTTP requests (status < 500) over a 30-day rolling window. Target: 99.9%",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents:  SuccessfulRequests,
		TotalEvents: AllRequests,
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn - API Availability",
			AlertType:  slo.BudgetRate,
			Threshold:  2.0, // Alert if burning budget at 2x normal rate
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#oncall"},
				{Type: "pagerduty", Target: "api-oncall-service"},
			},
		},
		{
			Name:       "Slow Burn - API Availability",
			AlertType:  slo.BudgetRate,
			Threshold:  1.5, // Alert if burning budget at 1.5x normal rate
			Window:     slo.TimePeriod{Hours: 24},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#api-team"},
			},
		},
	},
}

// LatencySLO defines a latency SLO ensuring 95% of requests complete under 1000ms.
// Uses a 7-day rolling window for more responsive latency tracking.
//
// References:
//   - GoodEvents: Fast requests (< 1000ms)
//   - TotalEvents: All requests
var LatencySLO = slo.SLO{
	Name:        "API Latency P95 < 1s",
	Description: "Ensures 95% of API requests complete within 1 second over a 7-day rolling window.",
	Dataset:     "production",
	SLI: slo.SLI{
		// Good events: requests completing under 1000ms
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.LT("duration_ms", 1000),
				query.Exists("http.route"),
			},
		},
		// Total events: all requests
		TotalEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Exists("http.route"),
			},
		},
	},
	Target:     slo.Percentage(95.0),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn - Latency SLO",
			AlertType:  slo.BudgetRate,
			Threshold:  3.0, // More aggressive threshold for latency issues
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#performance"},
			},
		},
	},
}
