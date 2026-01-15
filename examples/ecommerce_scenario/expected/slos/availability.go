// Package slos demonstrates E2E SLO declarations for the ecommerce scenario.
//
// These SLOs define service level objectives with multi-tier burn alert
// configurations for comprehensive reliability monitoring.
package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// CheckoutAvailability defines a 99.9% availability SLO for checkout service.
// Tracks successful checkout requests over a 30-day rolling window with
// 3-tier burn alerts for proactive incident detection.
var CheckoutAvailability = slo.SLO{
	Name:        "Checkout Service Availability",
	Description: "Track checkout service reliability - 99.9% of requests must succeed (status < 500)",
	Dataset:     "ecommerce-production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "ecommerce-production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/api/checkout"),
				query.LT("http.status_code", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "ecommerce-production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/api/checkout"),
			},
		},
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn - Checkout Availability Critical",
			AlertType:  slo.BudgetRate,
			Threshold:  2.0, // Alert if burning budget at 2x normal rate
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "checkout-oncall"},
				{Type: "slack", Target: "#alerts-checkout-critical"},
			},
		},
		{
			Name:       "Slow Burn - Checkout Availability Warning",
			AlertType:  slo.BudgetRate,
			Threshold:  5.0, // Alert if burning budget at 5x normal rate over longer window
			Window:     slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-checkout"},
			},
		},
		{
			Name:       "Exhaustion Forecast - Checkout Availability",
			AlertType:  slo.ExhaustionTime,
			Threshold:  24, // Alert if budget will be exhausted within 24 hours
			Window:     slo.Days(1),
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-checkout"},
				{Type: "email", Target: "checkout-team@example.com"},
			},
		},
	},
}
