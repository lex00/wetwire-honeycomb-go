package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// Availability SLO - 99.9% of requests must have status code < 500
var Availability = slo.SLO{
	Name:    "Service Availability",
	Dataset: "tasks-api",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Days(30),
			Filters: []query.Filter{
				query.LT("http.status_code", 500),
			},
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Days(30),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn",
			AlertType:  slo.BudgetRate,
			Threshold:  2.0,
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{{Type: "slack", Target: "#alerts"}},
		},
		{
			Name:       "Slow Burn",
			AlertType:  slo.BudgetRate,
			Threshold:  5.0,
			Window:     slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{{Type: "slack", Target: "#alerts"}},
		},
	},
}

// Latency SLO - 95% of requests must complete in under 500ms
var Latency = slo.SLO{
	Name:    "Request Latency",
	Dataset: "tasks-api",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Days(7),
			Filters: []query.Filter{
				query.LT("duration_ms", 500),
			},
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "tasks-api",
			TimeRange: query.Days(7),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(95),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn",
			AlertType:  slo.BudgetRate,
			Threshold:  10.0,
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{{Type: "slack", Target: "#alerts"}},
		},
	},
}
