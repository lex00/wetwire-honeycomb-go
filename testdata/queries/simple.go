package queries

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// SlowRequests finds requests that take too long
var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"endpoint", "service"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 500),
	},
}

// ErrorRate calculates error percentage by service
var ErrorRate = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(24),
	Breakdowns: []string{"service"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("status", "error"),
	},
}
