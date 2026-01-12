// Package traces demonstrates distributed tracing queries.
//
// These queries help analyze trace data and service dependencies.
package traces

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowTraces finds the slowest distributed traces.
// Use root span duration for end-to-end latency.
var SlowTraces = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"name", "service.name"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("trace.parent_id", nil), // Root spans only
		query.GT("duration_ms", 1000),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 25,
}

// ServiceDependencies shows how services call each other.
// Useful for understanding service topology.
var ServiceDependencies = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"service.name", "peer.service"},
	Calculations: []query.Calculation{
		query.Count(),
		query.P99("duration_ms"),
	},
	Filters: []query.Filter{
		query.Exists("peer.service"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 50,
}

// SpansByService shows span distribution across services.
// Identifies services generating the most spans.
var SpansByService = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"service.name"},
	Calculations: []query.Calculation{
		query.Count(),
		query.Avg("duration_ms"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// TraceErrors finds traces with errors.
// Group by trace to see error propagation patterns.
var TraceErrors = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"trace.trace_id", "service.name"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("status_code", "ERROR"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 50,
}

// SpanDuration analyzes span duration by operation name.
// Helps identify slow operations within traces.
var SpanDuration = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"name"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P50("duration_ms"),
		query.Count(),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 30,
}
