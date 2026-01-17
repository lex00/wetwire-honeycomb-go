// Package boards provides example Honeycomb board declarations for a Task API service.
package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/examples/tasks_api_scenario/expected/queries"
)

// TasksAPIDashboard provides monitoring for the Task API service.
var TasksAPIDashboard = board.Board{
	Name:        "Task API Dashboard",
	Description: "Monitoring dashboard for the Task API service",
	Panels: []board.Panel{
		board.QueryPanel(
			queries.RequestLatency,
			board.WithTitle("Request Latency (P50/P95/P99)"),
			board.WithPosition(0, 0, 6, 4),
		),
		board.QueryPanel(
			queries.ErrorRate,
			board.WithTitle("Error Rate by Endpoint"),
			board.WithPosition(6, 0, 6, 4),
		),
		board.QueryPanel(
			queries.RequestThroughput,
			board.WithTitle("Request Throughput"),
			board.WithPosition(0, 4, 6, 4),
		),
		board.QueryPanel(
			queries.SlowRequests,
			board.WithTitle("Slow Requests (>500ms)"),
			board.WithPosition(6, 4, 6, 4),
		),
	},
	Tags: []board.Tag{
		{Key: "service", Value: "tasks-api"},
	},
}
