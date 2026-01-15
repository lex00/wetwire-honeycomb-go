// Package slos demonstrates E2E SLO declarations for the ecommerce scenario.
package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// CheckoutLatency defines a 95% latency SLO for checkout operations.
// Ensures that 95% of checkout requests complete within 2 seconds to
// maintain excellent user experience during the purchase flow.
var CheckoutLatency = slo.SLO{
	Name:        "Checkout Service Latency",
	Description: "Track checkout latency - 95% of requests must complete under 2000ms",
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
				query.LT("duration_ms", 2000),
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
	Target:     slo.Percentage(95.0),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn - Checkout Latency",
			AlertType:  slo.BudgetRate,
			Threshold:  10.0,
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-checkout"},
			},
		},
		{
			Name:       "Slow Burn - Checkout Latency",
			AlertType:  slo.BudgetRate,
			Threshold:  5.0,
			Window:     slo.TimePeriod{Hours: 24},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-checkout"},
			},
		},
	},
}
