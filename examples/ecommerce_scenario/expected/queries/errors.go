package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// ErrorRateByService tracks error counts and rates across all services.
// Use this to identify which services are experiencing the highest error rates.
var ErrorRateByService = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service.name", "http.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
		query.CountDistinct("trace.trace_id"),
		query.Avg("duration_ms"),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 50,
}

// CheckoutFunnel tracks request counts through each stage of the checkout flow.
// Use this to identify drop-off points and conversion rates in the checkout process.
var CheckoutFunnel = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.route", "service.name"},
	Calculations: []query.Calculation{
		query.Count(),
		query.CountDistinct("user.id"),
		query.CountDistinct("session.id"),
		query.P50("duration_ms"),
	},
	Filters: []query.Filter{
		query.In("http.route", []any{
			"/cart/view",
			"/cart/add",
			"/checkout/start",
			"/checkout/address",
			"/checkout/payment",
			"/checkout/complete",
		}),
		query.LT("http.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 100,
}
