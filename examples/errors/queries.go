// Package errors demonstrates error tracking queries.
//
// These queries help identify and analyze errors in your system.
package errors

import "github.com/lex00/wetwire-honeycomb-go/query"

// ErrorsByService shows error counts grouped by service.
// Use this to identify which services have the most errors.
var ErrorsByService = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
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

// ErrorRate calculates error rate as percentage of all requests.
// Group by endpoint to see which routes are most problematic.
var ErrorRate = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 25,
}

// ErrorsByType groups errors by exception type or error message.
// Helps identify the most common error patterns.
var ErrorsByType = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"exception.type", "service.name"},
	Calculations: []query.Calculation{
		query.Count(),
		query.CountDistinct("trace.trace_id"),
	},
	Filters: []query.Filter{
		query.Exists("exception.type"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 30,
}

// RecentErrors shows the most recent errors with context.
// Useful for debugging recent issues.
var RecentErrors = query.Query{
	Dataset:   "production",
	TimeRange: query.Minutes(30),
	Breakdowns: []string{"exception.message", "service.name", "http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("exception.message"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 50,
}

// HTTPStatusCodes shows distribution of HTTP status codes.
// Quick overview of response status patterns.
var HTTPStatusCodes = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.status_code"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}
