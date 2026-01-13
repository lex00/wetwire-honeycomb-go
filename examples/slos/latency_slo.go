package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// APILatency defines a 95% latency SLO for API endpoints.
// Tracks requests completing within 500ms threshold.
var APILatency = slo.SLO{
	Name:        "API Latency P95",
	Description: "Ensure 95% of API requests complete within 500ms",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.LT("duration_ms", 500),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
		},
	},
	Target:     slo.Percentage(95.0),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		slo.FastBurn(10.0),
		slo.SlowBurn(5.0),
	},
}

// DatabaseQueryLatency defines a 99% latency SLO for database operations.
// Ensures database queries complete within 100ms for optimal performance.
var DatabaseQueryLatency = slo.SLO{
	Name:        "Database Query Latency P99",
	Description: "Track database query performance - 99% under 100ms",
	Dataset:     "backend",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "backend",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Equals("span.kind", "database"),
				query.LT("duration_ms", 100),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "backend",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Equals("span.kind", "database"),
			},
		},
	},
	Target:     slo.Percentage(99.0),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "DB Latency Fast Burn",
			AlertType: slo.BudgetRate,
			Threshold: 15.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "backend-oncall"},
			},
		},
		{
			Name:      "DB Latency Slow Burn",
			AlertType: slo.BudgetRate,
			Threshold: 7.5,
			Window:    slo.TimePeriod{Hours: 24},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-backend"},
			},
		},
	},
}

// CheckoutLatency defines a 99.5% latency SLO for checkout operations.
// Critical user-facing flow requires aggressive latency targets.
var CheckoutLatency = slo.SLO{
	Name:        "Checkout Flow Latency",
	Description: "Track checkout latency - 99.5% under 1000ms",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.In("http.route", []any{
					"/api/v1/checkout",
					"/api/v1/payment/process",
					"/api/v1/order/confirm",
				}),
				query.LT("duration_ms", 1000),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.In("http.route", []any{
					"/api/v1/checkout",
					"/api/v1/payment/process",
					"/api/v1/order/confirm",
				}),
			},
		},
	},
	Target:     slo.Percentage(99.5),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "Checkout Latency Critical",
			AlertType: slo.BudgetRate,
			Threshold: 20.0,
			Window:    slo.TimePeriod{Hours: 1},
			Recipients: []slo.Recipient{
				{Type: "pagerduty", Target: "checkout-oncall"},
				{Type: "slack", Target: "#alerts-checkout"},
			},
		},
		{
			Name:      "Checkout Latency Warning",
			AlertType: slo.BudgetRate,
			Threshold: 10.0,
			Window:    slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-checkout"},
			},
		},
	},
}

// SearchLatency defines a 90% latency SLO for search operations.
// Lower target acceptable for non-critical search functionality.
var SearchLatency = slo.SLO{
	Name:        "Search Latency P90",
	Description: "Track search performance - 90% under 2000ms",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.StartsWith("http.route", "/api/v1/search"),
				query.LT("duration_ms", 2000),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "production",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.StartsWith("http.route", "/api/v1/search"),
			},
		},
	},
	Target:     slo.Percentage(90.0),
	TimePeriod: slo.Days(7),
	BurnAlerts: []slo.BurnAlert{
		{
			Name:      "Search Latency Degradation",
			AlertType: slo.BudgetRate,
			Threshold: 8.0,
			Window:    slo.TimePeriod{Hours: 6},
			Recipients: []slo.Recipient{
				{Type: "slack", Target: "#alerts-search"},
			},
		},
	},
}

// UploadLatency defines a 95% latency SLO for file upload operations.
// Tracks upload completion time for user experience monitoring.
var UploadLatency = slo.SLO{
	Name:        "File Upload Latency",
	Description: "Track upload performance - 95% under 5000ms",
	Dataset:     "storage-service",
	SLI: slo.SLI{
		GoodEvents: query.Query{
			Dataset:   "storage-service",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Equals("operation", "upload"),
				query.LT("duration_ms", 5000),
			},
		},
		TotalEvents: query.Query{
			Dataset:   "storage-service",
			TimeRange: query.Hours(1),
			Calculations: []query.Calculation{
				query.Count(),
			},
			Filters: []query.Filter{
				query.Equals("operation", "upload"),
			},
		},
	},
	Target:     slo.Percentage(95.0),
	TimePeriod: slo.Days(30),
	BurnAlerts: []slo.BurnAlert{
		slo.FastBurn(12.0),
		slo.SlowBurn(6.0),
	},
}
