package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// RequestLatency tracks P99, P95, and P50 latency broken down by route and method
var RequestLatency = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.route", "http.method"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
	},
}

// ErrorRate counts errors (status >= 400) broken down by route and status code
var ErrorRate = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.route", "http.status_code"},
	Filters: []query.Filter{
		query.GTE("http.status_code", 400),
	},
	Calculations: []query.Calculation{
		query.Count(),
	},
}

// SlowRequests tracks requests that take longer than 500ms
var SlowRequests = query.Query{
	Dataset:   "tasks-api",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.route", "http.method"},
	Filters: []query.Filter{
		query.GT("duration_ms", 500),
	},
	Calculations: []query.Calculation{
		query.Count(),
		query.Avg("duration_ms"),
		query.Max("duration_ms"),
	},
}

// RequestThroughput measures request volume with 5-minute granularity
var RequestThroughput = query.Query{
	Dataset:     "tasks-api",
	TimeRange:   query.Hours(2),
	Granularity: 300, // 5 minutes in seconds
	Breakdowns:  []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
}
