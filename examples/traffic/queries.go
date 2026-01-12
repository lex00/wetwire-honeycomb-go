// Package traffic demonstrates traffic analysis queries.
//
// These queries help understand request patterns and traffic distribution.
package traffic

import "github.com/lex00/wetwire-honeycomb-go/query"

// RequestsByEndpoint shows request volume per endpoint.
// Identifies the most frequently called endpoints.
var RequestsByEndpoint = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.route"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 25,
}

// TrafficByService shows request distribution across services.
// Useful for capacity planning.
var TrafficByService = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(6),
	Breakdowns: []string{"service.name"},
	Calculations: []query.Calculation{
		query.Count(),
		query.P99("duration_ms"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// UniqueUsers counts distinct users over time.
// Tracks active user count.
var UniqueUsers = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(24),
	Calculations: []query.Calculation{
		query.CountDistinct("user.id"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("user.id"),
	},
}

// TrafficByMethod shows request counts grouped by HTTP method.
// Identifies the mix of read vs write operations.
var TrafficByMethod = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"http.method"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Exists("http.method"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 10,
}

// ThroughputByRegion shows request rate across regions.
// Helps with geographic load balancing decisions.
var ThroughputByRegion = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"cloud.region"},
	Calculations: []query.Calculation{
		query.Count(),
		query.Sum("http.response_content_length"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 15,
}
