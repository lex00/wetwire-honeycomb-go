// Package slos demonstrates Honeycomb SLO declarations.
//
// These SLOs define service level objectives with burn alert configurations
// for proactive reliability monitoring.
package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// APIAvailability defines a 99.9% availability SLO for the API service.
// Tracks successful requests (status < 500) over a 30-day rolling window.
var APIAvailability = slo.SLO{
	Name:        "API Availability",
	Description: "Track API availability - successful requests vs total requests",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.LT("http.status_code", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "production",
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
			Name:      "API Fast Burn - Critical",
			AlertType: slo.BudgetRate,
			Threshold: 10.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "platform-oncall"},
				{Type: "slack", Target: "#alerts-critical"},
			},
		},
		{
			Name:      "API Slow Burn - Warning",
			AlertType: slo.BudgetRate,
			Threshold: 5.0,
			Window:    slo.TimePeriod{Hours: 24},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-platform"},
			},
		},
	},
}

// CheckoutAvailability defines a 99.95% availability SLO for checkout flow.
// Higher target for business-critical checkout operations.
var CheckoutAvailability = slo.SLO{
	Name:        "Checkout Availability",
	Description: "Track checkout flow reliability - critical business path",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/checkout"),
				query.LT("http.status_code", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/checkout"),
			},
		},
	},
	Target:     slo.Percentage(99.95),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		slo.FastBurn(15.0),
		slo.SlowBurn(7.5),
	},
}

// AuthAvailability defines a 99.99% availability SLO for authentication.
// Highest availability target for authentication service.
var AuthAvailability = slo.SLO{
	Name:        "Authentication Availability",
	Description: "Track auth service availability - foundational service",
	Dataset:     "auth-service",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "auth-service",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Equals("result", "success"),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "auth-service",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(99.99),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "Auth Fast Burn - Page Immediately",
			AlertType: slo.BudgetRate,
			Threshold: 20.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "auth-oncall"},
				{Type: "slack", Target: "#alerts-auth-critical"},
			},
		},
		{
			Name:      "Auth Medium Burn - Investigate",
			AlertType: slo.BudgetRate,
			Threshold: 10.0,
			Window:    slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-auth"},
			},
		},
		{
			Name:      "Auth Slow Burn - Review",
			AlertType: slo.BudgetRate,
			Threshold: 5.0,
			Window:    slo.TimePeriod{Hours: 24},
			Recipients: []slo.Recipient{
				{Type: "email", Target: "auth-team@example.com"},
			},
		},
	},
}
