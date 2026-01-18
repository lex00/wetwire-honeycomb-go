package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/examples/tasks_api_scenario/results/intermediate/expected/queries"
)

// TasksAPIDashboard provides a comprehensive view of Task API performance
var TasksAPIDashboard = board.Board{
	Name:        "Tasks API Dashboard",
	Description: "Overview of Tasks API performance metrics",
	Panels: []board.Panel{
		board.QueryPanel(
			queries.RequestLatency,
			board.WithTitle("Request Latency"),
			board.WithPosition(0, 0, 6, 4),
		),
		board.QueryPanel(
			queries.ErrorRate,
			board.WithTitle("Error Rate"),
			board.WithPosition(6, 0, 6, 4),
		),
		board.QueryPanel(
			queries.SlowRequests,
			board.WithTitle("Slow Requests"),
			board.WithPosition(0, 4, 6, 4),
		),
		board.QueryPanel(
			queries.RequestThroughput,
			board.WithTitle("Request Throughput"),
			board.WithPosition(6, 4, 6, 4),
		),
	},
}
