// Package sli demonstrates Service Level Indicator queries.
//
// These queries help track SLI/SLO metrics for reliability monitoring.
package sli

import "github.com/lex00/wetwire-honeycomb-go/query"

// AvailabilitySLI measures successful requests vs total requests.
// Core availability metric for SLO tracking.
var AvailabilitySLI = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(24),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.LT("http.status_code", 500),
	},
}

// LatencySLI measures requests meeting latency threshold.
// Track what percentage of requests complete within SLO.
var LatencySLI = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(24),
	Calculations: []query.Calculation{
		query.Count(),
		query.P99("duration_ms"),
	},
	Filters: []query.Filter{
		query.LT("duration_ms", 500),
	},
}

// ErrorBudgetByService tracks error budget consumption per service.
// Helps prioritize reliability improvements.
var ErrorBudgetByService = query.Query{
	Dataset:   "production",
	TimeRange: query.Days(7),
	Breakdowns: []string{"service.name"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 500),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// ThroughputSLI measures request throughput over time.
// Track capacity-related SLIs.
var ThroughputSLI = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
}

// ApdexScore approximates user satisfaction via latency thresholds.
// Satisfied: <200ms, Tolerating: 200-800ms, Frustrated: >800ms
var ApdexScore = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.LT("duration_ms", 200),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 25,
}

// CriticalEndpointHealth tracks SLIs for critical business endpoints.
// Focus monitoring on revenue-impacting paths.
var CriticalEndpointHealth = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
		query.P99("duration_ms"),
		query.P50("duration_ms"),
	},
	Filters: []query.Filter{
		query.In("http.route", []interface{}{
			"/api/v1/checkout",
			"/api/v1/payment",
			"/api/v1/auth/login",
		}),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 10,
}
