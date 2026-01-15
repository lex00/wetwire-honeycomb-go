// Package queries demonstrates e-commerce observability queries.
//
// These queries analyze checkout flow performance, payment fraud correlation,
// and service health for an OpenTelemetry demo application.
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// CheckoutFlowLatency tracks P99/P95/P50 latency across checkout services.
// Use this to identify performance bottlenecks in the critical checkout path.
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
			"checkoutservice",
			"cartservice",
			"paymentservice",
			"frauddetectionservice",
		}),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 100,
}

// PaymentFraudCorrelation analyzes traces where fraud detection impacts payment latency.
// Helps understand the relationship between fraud checks and payment processing time.
var PaymentFraudCorrelation = query.Query{
	Dataset:   "otel-demo",
	TimeRange: query.Hours(4),
	Breakdowns: []string{"trace.parent_id", "service.name", "fraud.score_bucket"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Avg("duration_ms"),
		query.Count(),
		query.Sum("fraud.detection_time_ms"),
	},
	Filters: []query.Filter{
		query.In("service.name", []any{
			"paymentservice",
			"frauddetectionservice",
		}),
		query.Exists("fraud.score"),
		query.GT("duration_ms", 100),
	},
	Orders: []query.Order{
		{Op: "SUM", Column: "fraud.detection_time_ms", Order: "descending"},
	},
	Limit: 50,
}
