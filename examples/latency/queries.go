// Package latency demonstrates latency analysis queries.
//
// These queries help identify slow endpoints and track performance over time.
package latency

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowEndpoints finds the slowest endpoints by P99 latency.
// Use this to identify which endpoints need optimization.
var SlowEndpoints = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.route", "service.name"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
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

// LatencyDistribution shows the full latency distribution for a specific endpoint.
// Use this to understand latency patterns beyond percentiles.
var LatencyDistribution = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Heatmap("duration_ms"),
	},
	Filters: []query.Filter{
		query.Equals("http.route", "/api/v1/users"),
	},
}

// LatencyByRegion compares latency across different regions.
// Useful for identifying geographic performance issues.
var LatencyByRegion = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"cloud.region"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Avg("duration_ms"),
		query.Count(),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 20,
}

// SlowDatabaseQueries identifies slow database operations.
// Helps find database bottlenecks.
var SlowDatabaseQueries = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"db.statement", "db.name"},
	Calculations: []query.Calculation{
		query.P99("db.duration_ms"),
		query.Avg("db.duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("db.statement"),
		query.GT("db.duration_ms", 100),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "db.duration_ms", Order: "descending"},
	},
	Limit: 25,
}
