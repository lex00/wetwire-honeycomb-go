// Package slos demonstrates E2E SLO declarations for the ecommerce scenario.
package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// PaymentSuccess defines a 99.95% success rate SLO for payment processing.
// Tracks successful payment transactions - critical for revenue and customer trust.
var PaymentSuccess = slo.SLO{
	Name:        "Payment Processing Success Rate",
	Description: "Track payment success rate - 99.95% of payment attempts must succeed",
	Dataset:     "ecommerce-production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "ecommerce-production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/api/payment"),
				query.Equals("payment.status", "success"),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "ecommerce-production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Contains("http.route", "/api/payment"),
			},
		},
	},
	Target:     slo.Percentage(99.95),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:       "Fast Burn - Payment Success Critical",
			AlertType:  slo.BudgetRate,
			Threshold:  15.0,
			Window:     slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "payment-oncall"},
				{Type: "slack", Target: "#alerts-payment-critical"},
			},
		},
		{
			Name:       "Slow Burn - Payment Success Warning",
			AlertType:  slo.BudgetRate,
			Threshold:  7.5,
			Window:     slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-payment"},
			},
		},
	},
}
