// Package queries provides example Honeycomb query declarations for a Task API service.
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// RequestLatency tracks P50/P95/P99 latency across all endpoints.
// Use this to identify performance bottlenecks and latency distribution.
var RequestLatency = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "http.method"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
		query.Count(),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 100,
}

// ErrorRate tracks error counts by status code and endpoint.
// Use this to identify which endpoints are experiencing the most errors.
var ErrorRate = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "http.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 50,
}

// SlowRequests finds requests that exceed the 500ms latency threshold.
// Use this to investigate specific slow requests and their characteristics.
var SlowRequests = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "http.method", "http.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
		query.Avg("duration_ms"),
		query.Max("duration_ms"),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 500),
	},
	Orders: []query.Order{
		{Op: "MAX", Column: "duration_ms", Order: "descending"},
	},
	Limit: 50,
}

// RequestThroughput tracks request volume over time.
// Use this to understand traffic patterns and capacity needs.
var RequestThroughput = query.Query{
	Dataset:     "tasks-api",
	TimeRange:   query.Hours(2),
	Breakdowns:  []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Granularity: 300, // 5-minute buckets
}
