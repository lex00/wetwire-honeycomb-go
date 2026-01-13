package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// AuthServiceLatency tracks latency metrics for the auth service.
var AuthServiceLatency = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "auth-service"),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 15,
}

// AuthServiceErrors tracks error rates for the auth service.
var AuthServiceErrors = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "auth-service"),
		query.GTE("http.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 10,
}

// AuthServiceThroughput tracks request volume for the auth service.
var AuthServiceThroughput = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "auth-service"),
	},
	Granularity: 300, // 5-minute buckets
}

// AuthServiceLatencyDistribution shows the latency heatmap for the auth service.
var AuthServiceLatencyDistribution = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Calculations: []query.Calculation{
		query.Heatmap("duration_ms"),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "auth-service"),
	},
}

// AuthServiceExceptions tracks recent exceptions in the auth service.
var AuthServiceExceptions = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"exception.type", "http.route"},
	Calculations: []query.Calculation{
		query.Count(),
		query.CountDistinct("trace.trace_id"),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "auth-service"),
		query.Exists("exception.type"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// AuthServiceHealthBoard provides a comprehensive health dashboard for the auth service.
// This board demonstrates how to compose multiple queries into a service health dashboard.
var AuthServiceHealthBoard = board.Board{
	Name:        "Auth Service Health Dashboard",
	Description: "Comprehensive health monitoring for auth-service",
	Panels: []board.Panel{
		board.TextPanel(
			"# Auth Service Health\n\nMonitoring latency, errors, throughput, and resource utilization.",
			board.WithTitle("Service Overview"),
			board.WithPosition(0, 0, 12, 2),
		),
		board.QueryPanel(
			AuthServiceLatency,
			board.WithTitle("Endpoint Latency Breakdown"),
			board.WithPosition(0, 2, 6, 4),
		),
		board.QueryPanel(
			AuthServiceErrors,
			board.WithTitle("Error Rates (4xx/5xx)"),
			board.WithPosition(6, 2, 6, 4),
		),
		board.QueryPanel(
			AuthServiceThroughput,
			board.WithTitle("Request Throughput Over Time"),
			board.WithPosition(0, 6, 6, 3),
		),
		board.QueryPanel(
			AuthServiceLatencyDistribution,
			board.WithTitle("Latency Distribution Heatmap"),
			board.WithPosition(6, 6, 6, 3),
		),
		board.QueryPanel(
			AuthServiceExceptions,
			board.WithTitle("Recent Exceptions"),
			board.WithPosition(0, 9, 12, 4),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "service.name",
			Operation: "=",
			Value:     "auth-service",
		},
	},
	Tags: []board.Tag{
		{Key: "service", Value: "auth-service"},
		{Key: "environment", Value: "production"},
	},
}

// APIGatewayTopEndpoints tracks top endpoints by request volume.
var APIGatewayTopEndpoints = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route", "http.method"},
	Calculations: []query.Calculation{
		query.Count(),
		query.P95("duration_ms"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 30,
}

// APIGatewayRateLimited tracks rate-limited requests.
var APIGatewayRateLimited = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("http.status_code", 429),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// APIGatewayAuthFailures tracks authentication failures.
var APIGatewayAuthFailures = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"http.route"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("http.status_code", 401),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 15,
}

// APIGatewayLatencyPercentiles tracks gateway latency percentiles over time.
var APIGatewayLatencyPercentiles = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(4),
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
	},
	Granularity: 300, // 5-minute buckets
}

// APIGatewayHealthBoard provides monitoring for API gateway health.
// Focuses on endpoint-level metrics, rate limiting, and authentication patterns.
var APIGatewayHealthBoard = board.Board{
	Name:        "API Gateway Health",
	Description: "API Gateway monitoring including endpoints, rate limits, and auth",
	Panels: []board.Panel{
		board.TextPanel(
			"# API Gateway Monitoring\n\nTrack endpoint performance, rate limiting, authentication, and quota usage.",
			board.WithTitle("Gateway Overview"),
			board.WithPosition(0, 0, 12, 2),
		),
		board.QueryPanel(
			APIGatewayTopEndpoints,
			board.WithTitle("Top Endpoints by Request Volume"),
			board.WithPosition(0, 2, 6, 4),
		),
		board.QueryPanel(
			APIGatewayRateLimited,
			board.WithTitle("Rate Limited Requests (429)"),
			board.WithPosition(6, 2, 6, 4),
		),
		board.QueryPanel(
			APIGatewayAuthFailures,
			board.WithTitle("Authentication Failures (401)"),
			board.WithPosition(0, 6, 6, 3),
		),
		board.QueryPanel(
			APIGatewayLatencyPercentiles,
			board.WithTitle("Gateway Latency Percentiles Over Time"),
			board.WithPosition(6, 6, 6, 3),
		),
	},
	Tags: []board.Tag{
		{Key: "component", Value: "api-gateway"},
		{Key: "team", Value: "platform"},
	},
}

// PaymentServiceDownstreamCalls tracks downstream service dependencies for payment service.
var PaymentServiceDownstreamCalls = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"downstream.service", "downstream.endpoint"},
	Calculations: []query.Calculation{
		query.Count(),
		query.P95("downstream.duration_ms"),
		query.Avg("downstream.duration_ms"),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "payment-service"),
		query.Exists("downstream.service"),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 25,
}

// PaymentServiceDownstreamErrors tracks downstream service errors for payment service.
var PaymentServiceDownstreamErrors = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"downstream.service", "downstream.status_code"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "payment-service"),
		query.GTE("downstream.status_code", 400),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 20,
}

// PaymentServiceDownstreamLatency tracks downstream service latency for payment service.
var PaymentServiceDownstreamLatency = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"downstream.service"},
	Calculations: []query.Calculation{
		query.P99("downstream.duration_ms"),
	},
	Filters: []query.Filter{
		query.Equals("service.name", "payment-service"),
		query.Exists("downstream.service"),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "downstream.duration_ms", Order: "descending"},
	},
	Limit: 20,
	Granularity: 300, // 5-minute buckets
}

// PaymentServiceDependenciesBoard monitors payment service dependencies.
// Tracks service dependencies, downstream latency, and inter-service errors.
var PaymentServiceDependenciesBoard = board.Board{
	Name:        "Payment Service Microservice Dependencies",
	Description: "Track payment-service dependencies and inter-service communication",
	Panels: []board.Panel{
		board.QueryPanel(
			PaymentServiceDownstreamCalls,
			board.WithTitle("Downstream Service Dependencies"),
			board.WithPosition(0, 0, 6, 4),
		),
		board.QueryPanel(
			PaymentServiceDownstreamErrors,
			board.WithTitle("Downstream Service Errors"),
			board.WithPosition(6, 0, 6, 4),
		),
		board.QueryPanel(
			PaymentServiceDownstreamLatency,
			board.WithTitle("Downstream Service Latency (P99)"),
			board.WithPosition(0, 4, 12, 3),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "service.name",
			Operation: "=",
			Value:     "payment-service",
		},
	},
	Tags: []board.Tag{
		{Key: "service", Value: "payment-service"},
		{Key: "focus", Value: "dependencies"},
	},
}

