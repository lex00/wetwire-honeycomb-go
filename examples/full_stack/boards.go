package full_stack

import "github.com/lex00/wetwire-honeycomb-go/board"

// PerformanceBoard provides a comprehensive view of API performance metrics.
// Combines query panels and SLO panels to visualize the complete monitoring picture.
//
// Layout:
//   - Row 1: Overview text and key metrics
//   - Row 2: Error rate and slow requests
//   - Row 3: Request throughput and latency distribution
//   - Row 4: SLO tracking panels
//
// References:
//   - ErrorRate query (queries.go)
//   - SlowRequests query (queries.go)
//   - RequestThroughput query (queries.go)
//   - LatencyP99 query (queries.go)
//   - APIAvailability SLO (slos.go)
//   - LatencySLO (slos.go)
var PerformanceBoard = board.Board{
	Name:        "API Performance Dashboard",
	Description: "Comprehensive monitoring dashboard tracking API performance, availability, and SLO compliance",
	Panels: []board.Panel{
		// Row 1: Overview
		board.TextPanel(
			"# API Performance Overview\n\nThis dashboard tracks key performance indicators for the production API:\n- **Availability**: Success rate and error trends\n- **Latency**: P99 response times across endpoints\n- **Throughput**: Request volume and rate\n- **SLOs**: Compliance with service level objectives",
			board.WithTitle("Dashboard Overview"),
			board.WithPosition(0, 0, 12, 2),
		),

		// Row 2: Error tracking
		board.QueryPanel(
			ErrorRate,
			board.WithTitle("Error Rate by Endpoint"),
			board.WithPosition(0, 2, 6, 4),
		),
		board.QueryPanel(
			SlowRequests,
			board.WithTitle("Slow Requests (>1s)"),
			board.WithPosition(6, 2, 6, 4),
		),

		// Row 3: Traffic and latency
		board.QueryPanel(
			RequestThroughput,
			board.WithTitle("Request Throughput (5min buckets)"),
			board.WithPosition(0, 6, 6, 4),
		),
		board.QueryPanel(
			LatencyP99,
			board.WithTitle("P99 Latency by Endpoint"),
			board.WithPosition(6, 6, 6, 4),
		),

		// Row 4: SLO panels
		// Note: SLO panels typically use SLOPanelByID with actual SLO IDs from Honeycomb
		// For this example, we demonstrate the panel structure
		board.TextPanel(
			"## SLO Tracking\n\n**API Availability SLO**: 99.9% success rate over 30 days\n\n**Latency SLO**: 95% of requests under 1 second over 7 days\n\nRefer to slos.go for full SLO definitions.",
			board.WithTitle("SLO Status"),
			board.WithPosition(0, 10, 12, 3),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "service.name",
			Operation: "=",
			Value:     "api-service",
		},
	},
	Tags: []board.Tag{
		{Key: "team", Value: "platform"},
		{Key: "environment", Value: "production"},
		{Key: "service", Value: "api"},
	},
}

// IncidentResponseBoard provides focused views for incident investigation.
// Optimized layout for troubleshooting during active incidents.
//
// References:
//   - ErrorRate query (queries.go)
//   - SlowRequests query (queries.go)
//   - LatencyP99 query (queries.go)
var IncidentResponseBoard = board.Board{
	Name:        "Incident Response Dashboard",
	Description: "Focused dashboard for rapid incident investigation and triage",
	Panels: []board.Panel{
		// Critical metrics at the top
		board.QueryPanel(
			ErrorRate,
			board.WithTitle("Current Error Rate"),
			board.WithPosition(0, 0, 4, 4),
		),
		board.QueryPanel(
			SlowRequests,
			board.WithTitle("Slow Endpoints"),
			board.WithPosition(4, 0, 4, 4),
		),
		board.QueryPanel(
			RequestThroughput,
			board.WithTitle("Traffic Volume"),
			board.WithPosition(8, 0, 4, 4),
		),

		// Detailed breakdown
		board.QueryPanel(
			LatencyP99,
			board.WithTitle("Latency by Route"),
			board.WithPosition(0, 4, 12, 4),
		),

		// Investigation notes
		board.TextPanel(
			"## Incident Investigation Checklist\n\n1. Check error rate trends\n2. Identify affected endpoints\n3. Review latency distribution\n4. Compare traffic patterns\n5. Check recent deployments\n6. Review error logs in context",
			board.WithTitle("Investigation Guide"),
			board.WithPosition(0, 8, 12, 3),
		),
	},
	PresetFilters: []board.Filter{
		{
			Column:    "service.name",
			Operation: "=",
			Value:     "api-service",
		},
	},
	Tags: []board.Tag{
		{Key: "team", Value: "oncall"},
		{Key: "purpose", Value: "incident-response"},
	},
}
