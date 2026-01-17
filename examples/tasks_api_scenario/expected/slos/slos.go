// Package slos provides example Honeycomb SLO declarations for a Task API service.
package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// Availability defines a 99.9% availability SLO for the Task API.
// Tracks successful requests (status < 500) over a 30-day rolling window.
var Availability = slo.SLO{
	Name:        "Task API Availability",
	Description: "99.9% of requests must succeed (status < 500)",
	Dataset:     "tasks-api",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.LT("http.status_code", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "Fast Burn - Availability",
			AlertType: slo.BudgetRate,
			Threshold: 2.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts"},
			},
		},
		{
			Name:      "Slow Burn - Availability",
			AlertType: slo.BudgetRate,
			Threshold: 5.0,
			Window:    slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts"},
			},
		},
	},
}

// Latency defines a 95% latency SLO for the Task API.
// 95% of requests must complete within 500ms over a 7-day rolling window.
var Latency = slo.SLO{
	Name:        "Task API Latency",
	Description: "95% of requests must complete under 500ms",
	Dataset:     "tasks-api",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.LT("duration_ms", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(95.0),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "Fast Burn - Latency",
			AlertType: slo.BudgetRate,
			Threshold: 10.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts"},
			},
		},
	},
}
