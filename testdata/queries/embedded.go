package queries

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// CustomQuery embeds a Query
type CustomQuery struct {
	query.Query
	Tags []string
}

// DatabaseQueries is a composite type with embedded queries
var DatabaseQueries = CustomQuery{
	Query: query.Query{
		Dataset:   "database",
		TimeRange: query.Hours(1),
		Breakdowns: []string{"query_type"},
		Calculations: []query.Calculation{
			query.Avg("execution_time"),
			query.Count(),
		},
	},
	Tags: []string{"database", "performance"},
}

// FunctionScopedQuery demonstrates a function-scoped query definition
func GetLatencyQuery() query.Query {
	return query.Query{
		Dataset:   "api",
		TimeRange: query.Hours(6),
		Breakdowns: []string{"region"},
		Calculations: []query.Calculation{
			query.P99("latency_ms"),
		},
	}
}
