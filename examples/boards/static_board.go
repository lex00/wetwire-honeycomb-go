// Package boards demonstrates board declarations using the board package.
//
// These examples show how to compose queries into comprehensive dashboards.
package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// ProductionOverviewBoard provides a comprehensive view of production system health.
// Combines latency, error rates, and throughput metrics on a single dashboard.
var ProductionOverviewBoard = board.Board{
	Name:        "Production Overview",
	Description: "High-level health metrics for production environment",
	Panels: []board.Panel{
		board.TextPanel(
			"# Production System Health\n\nThis board provides a comprehensive overview of production system health, including latency, errors, and throughput.",
			board.WithTitle("Overview"),
			board.WithPosition(0, 0, 12, 2),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(4),
				Breakdowns: []string{"http.route"},
				Calculations: []query.Calculation{
					query.P99("duration_ms"),
					query.P95("duration_ms"),
					query.P50("duration_ms"),
					query.Count(),
				},
				Orders: []query.Order{
					{Op: "P99", Column: "duration_ms", Order: "descending"},
				},
				Limit: 20,
			},
			board.WithTitle("Endpoint Latency (P50/P95/P99)"),
			board.WithPosition(0, 2, 6, 4),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(4),
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
				Limit: 15,
			},
			board.WithTitle("Errors by Service (5xx)"),
			board.WithPosition(6, 2, 6, 4),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(4),
				Calculations: []query.Calculation{
					query.Count(),
				},
				Granularity: 300, // 5-minute buckets
			},
			board.WithTitle("Request Throughput"),
			board.WithPosition(0, 6, 6, 3),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(4),
				Breakdowns: []string{"http.status_code"},
				Calculations: []query.Calculation{
					query.Count(),
				},
				Orders: []query.Order{
					{Op: "COUNT", Order: "descending"},
				},
				Limit: 10,
			},
			board.WithTitle("HTTP Status Code Distribution"),
			board.WithPosition(6, 6, 6, 3),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "service.name",
			Operation: "exists",
			Value:     nil,
		},
	},
	Tags: []board.Tag{
		{Key: "team", Value: "platform"},
		{Key: "environment", Value: "production"},
	},
}

// DatabasePerformanceBoard tracks database query performance and health.
// Focus on slow queries, connection pool metrics, and database errors.
var DatabasePerformanceBoard = board.Board{
	Name:        "Database Performance",
	Description: "Database query performance and health monitoring",
	Panels: []board.Panel{
		board.TextPanel(
			"# Database Performance Monitoring\n\nTrack slow queries, connection pool health, and database errors.",
			board.WithTitle("Database Overview"),
			board.WithPosition(0, 0, 12, 2),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(2),
				Breakdowns: []string{"db.statement", "db.name"},
				Calculations: []query.Calculation{
					query.P99("db.duration_ms"),
					query.P95("db.duration_ms"),
					query.Avg("db.duration_ms"),
					query.Count(),
				},
				Filters: []query.Filter{
					query.Exists("db.statement"),
					query.GT("db.duration_ms", 50),
				},
				Orders: []query.Order{
					{Op: "P99", Column: "db.duration_ms", Order: "descending"},
				},
				Limit: 25,
			},
			board.WithTitle("Slow Database Queries"),
			board.WithPosition(0, 2, 8, 4),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(2),
				Breakdowns: []string{"db.name"},
				Calculations: []query.Calculation{
					query.Count(),
				},
				Filters: []query.Filter{
					query.Exists("db.name"),
				},
				Orders: []query.Order{
					{Op: "COUNT", Order: "descending"},
				},
				Limit: 10,
			},
			board.WithTitle("Queries by Database"),
			board.WithPosition(8, 2, 4, 4),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(2),
				Calculations: []query.Calculation{
					query.Avg("db.duration_ms"),
					query.P99("db.duration_ms"),
				},
				Filters: []query.Filter{
					query.Exists("db.duration_ms"),
				},
				Granularity: 300, // 5-minute buckets
			},
			board.WithTitle("Database Query Latency Over Time"),
			board.WithPosition(0, 6, 12, 3),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "db.name",
			Operation: "exists",
			Value:     nil,
		},
	},
	Tags: []board.Tag{
		{Key: "team", Value: "platform"},
		{Key: "component", Value: "database"},
	},
}

// ErrorAnalysisBoard provides detailed error tracking and analysis.
// Track error rates, error types, and affected endpoints.
var ErrorAnalysisBoard = board.Board{
	Name:        "Error Analysis",
	Description: "Comprehensive error tracking and analysis dashboard",
	Panels: []board.Panel{
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(6),
				Calculations: []query.Calculation{
					query.Count(),
				},
				Filters: []query.Filter{
					query.GTE("http.status_code", 500),
				},
				Granularity: 300, // 5-minute buckets
			},
			board.WithTitle("5xx Error Rate Over Time"),
			board.WithPosition(0, 0, 6, 3),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(6),
				Calculations: []query.Calculation{
					query.Count(),
				},
				Filters: []query.Filter{
					query.GTE("http.status_code", 400),
					query.LT("http.status_code", 500),
				},
				Granularity: 300, // 5-minute buckets
			},
			board.WithTitle("4xx Error Rate Over Time"),
			board.WithPosition(6, 0, 6, 3),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(6),
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
				Limit: 20,
			},
			board.WithTitle("Errors by Exception Type"),
			board.WithPosition(0, 3, 6, 4),
		),
		board.QueryPanel(
			query.Query{
				Dataset:   "production",
				TimeRange: query.Hours(6),
				Breakdowns: []string{"http.route", "http.status_code"},
				Calculations: []query.Calculation{
					query.Count(),
				},
				Filters: []query.Filter{
					query.GTE("http.status_code", 400),
				},
				Orders: []query.Order{
					{Op: "COUNT", Order: "descending"},
				},
				Limit: 30,
			},
			board.WithTitle("Errors by Endpoint"),
			board.WithPosition(6, 3, 6, 4),
		),
	},
	Tags: []board.Tag{
		{Key: "team", Value: "sre"},
		{Key: "focus", Value: "reliability"},
	},
}
