// Package boards demonstrates comprehensive dashboard composition for the ecommerce checkout flow.
//
// This package shows how to create unified monitoring dashboards that combine
// multiple queries and visualizations to provide end-to-end observability.
package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// CheckoutFlowLatency represents the latency query for checkout services.
// This should be defined in the queries package, but is included here as a reference.
var CheckoutFlowLatency = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service.name", "http.route"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.In("service.name", []any{
			"checkoutservice", "cartservice",
			"paymentservice", "frauddetectionservice",
		}),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 20,
}

// ErrorRateByService represents error counts grouped by service.
// This should be defined in the queries package, but is included here as a reference.
var ErrorRateByService = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service.name"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 500),
		query.In("service.name", []any{
			"checkoutservice", "cartservice",
			"paymentservice", "frauddetectionservice",
		}),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 10,
}

// RequestThroughput represents request volume over time.
var RequestThroughput = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(2),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.In("service.name", []any{
			"checkoutservice", "cartservice",
			"paymentservice", "frauddetectionservice",
		}),
	},
	Granularity: 300, // 5-minute buckets
}

// SLOBurnRate represents the burn rate status for checkout SLO.
var SLOBurnRate = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service.name"},
	Calculations: []query.Calculation{
		query.Count(),
		query.CountDistinct("trace.trace_id"),
	},
	Filters: []query.Filter{
		query.In("service.name", []any{
			"checkoutservice", "cartservice",
			"paymentservice", "frauddetectionservice",
		}),
		query.GTE("http.status_code", 500),
	},
	Orders: []query.Order{
		{Op: "COUNT", Order: "descending"},
	},
	Limit: 10,
}

// CheckoutDashboard provides end-to-end monitoring for the checkout flow.
//
// This board demonstrates how to compose multiple queries into a unified dashboard
// that provides comprehensive observability across the entire checkout service stack.
//
// Layout:
//   - Row 1: Dashboard overview (text panel)
//   - Row 2: Latency metrics (line chart) and Error rate (bar chart)
//   - Row 3: Request throughput (time series) and SLO burn rate (gauge)
//
// The board showcases:
//   - Query reuse: Multiple panels reference shared query definitions
//   - Type safety: All panels use typed builder functions
//   - Direct references: Queries are Go variables, not string IDs
//   - Preset filters: Board-level filters apply to all panels
//
// References:
//   - CheckoutFlowLatency query (latency.go)
//   - ErrorRateByService query (errors.go)
//   - RequestThroughput query (inline)
//   - SLOBurnRate query (inline)
var CheckoutDashboard = board.Board{
	Name:        "Checkout Flow Observability",
	Description: "End-to-end monitoring for checkout services including latency, errors, throughput, and SLO compliance",
	Panels: []board.Panel{
		// Row 1: Overview text
		board.TextPanel(
			"# Checkout Flow Monitoring\n\nThis dashboard provides comprehensive observability for the checkout service stack:\n\n- **Latency by Service**: Track P99/P95/P50 response times across checkout, cart, payment, and fraud detection services\n- **Error Rate by Service**: Monitor 5xx error counts to identify service reliability issues\n- **Request Throughput**: Observe traffic patterns and volume trends\n- **SLO Burn Rate**: Track error budget consumption to prevent SLO violations\n\nUse this board to:\n- Identify performance bottlenecks in the checkout flow\n- Detect service degradation early\n- Correlate errors across service boundaries\n- Monitor SLO compliance in real-time",
			board.WithTitle("Dashboard Overview"),
			board.WithPosition(0, 0, 12, 3),
		),

		// Row 2: Latency and errors
		board.QueryPanel(
			CheckoutFlowLatency,
			board.WithTitle("Latency by Service (P99/P95/P50)"),
			board.WithPosition(0, 3, 6, 4),
		),
		board.QueryPanel(
			ErrorRateByService,
			board.WithTitle("Error Rate by Service (5xx)"),
			board.WithPosition(6, 3, 6, 4),
		),

		// Row 3: Throughput and SLO burn
		board.QueryPanel(
			RequestThroughput,
			board.WithTitle("Request Throughput (5min buckets)"),
			board.WithPosition(0, 7, 6, 4),
		),
		board.QueryPanel(
			SLOBurnRate,
			board.WithTitle("SLO Burn Rate Status"),
			board.WithPosition(6, 7, 6, 4),
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
		{Key: "team", Value: "ecommerce"},
		{Key: "service", Value: "checkout"},
		{Key: "environment", Value: "production"},
		{Key: "demo", Value: "true"},
	},
}
