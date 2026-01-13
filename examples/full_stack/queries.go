// Package full_stack demonstrates the complete Query→SLO→Trigger→Board chain.
//
// This example shows how to build a comprehensive monitoring setup where:
// - Queries define the base metrics
// - SLOs reference queries to track service level objectives
// - Triggers reference queries to alert on anomalies
// - Boards reference queries and SLOs to visualize the complete picture
package full_stack

import "github.com/lex00/wetwire-honeycomb-go/query"

// SuccessfulRequests counts all successful HTTP requests (status < 500).
// Used as the "good events" query in the APIAvailability SLO.
var SuccessfulRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.LT("http.status_code", 500),
		query.Exists("http.status_code"),
	},
}

// AllRequests counts all HTTP requests regardless of status.
// Used as the "total events" query in the APIAvailability SLO.
var AllRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.status_code"),
	},
}

// SlowRequests identifies requests exceeding latency thresholds.
// Used in the HighLatencyAlert trigger and PerformanceBoard.
var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "service.name"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.Avg("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 1000),
		query.Exists("http.route"),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 50,
}

// ErrorRate calculates the percentage of failed requests.
// Used in the ErrorRateAlert trigger and PerformanceBoard.
var ErrorRate = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "service.name"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 500),
		query.Exists("http.status_code"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 25,
}

// LatencyP99 tracks P99 latency across all endpoints.
// Used as the basis for the LatencySLO.
var LatencyP99 = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.route"),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 50,
}

// RequestThroughput tracks overall request volume.
// Used in the PerformanceBoard to monitor traffic patterns.
var RequestThroughput = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(24),
	Calculations: []query.Calculation{
		query.Rate("http.status_code"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.status_code"),
	},
	Granularity: 300, // 5-minute buckets
}
