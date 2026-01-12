package queries

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// AdvancedQuery demonstrates a complex query with multiple features
var AdvancedQuery = query.Query{
	Dataset:           "logs",
	TimeRange:         query.Last24Hours(),
	Breakdowns:        []string{"service", "environment", "status_code"},
	FilterCombination: "AND",
	Calculations: []query.Calculation{
		query.Count(),
		query.P99("duration_ms"),
		query.Avg("response_size"),
		query.Max("memory_usage"),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 1000),
		query.Equals("environment", "production"),
		query.Contains("path", "/api/"),
	},
	Limit: 100,
}

// MetricsQuery demonstrates various calculation types
var MetricsQuery = query.Query{
	Dataset:   "metrics",
	TimeRange: query.Minutes(30),
	Breakdowns: []string{"host"},
	Calculations: []query.Calculation{
		query.Sum("bytes_sent"),
		query.Avg("cpu_usage"),
		query.P95("latency"),
		query.CountDistinct("user_id"),
	},
	Filters: []query.Filter{
		query.LT("cpu_usage", 80),
		query.Exists("user_id"),
	},
}

// ErrorTracking demonstrates filter combinations
var ErrorTracking = query.Query{
	Dataset:   "errors",
	TimeRange: query.Days(7),
	Breakdowns: []string{"error_type", "service"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.NotEquals("status_code", 200),
		query.DoesNotContain("path", "/health"),
	},
}
