// Package triggers demonstrates various trigger patterns for alerting.
//
// These triggers monitor key metrics and send notifications when thresholds are exceeded.
package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// HighLatencyTrigger alerts when P99 latency exceeds 1000ms.
// Monitors API endpoint performance and notifies on-call via PagerDuty.
var HighLatencyTrigger = trigger.Trigger{
	Name:        "High Latency Alert",
	Description: "Fires when P99 latency exceeds 1000ms over a 5 minute window",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(5),
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
		Filters: []query.Filter{
			query.Exists("http.route"),
		},
	},
	Threshold: trigger.GreaterThan(1000),
	Frequency: trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.PagerDutyService("P123ABC"),
		trigger.SlackChannel("#alerts-latency"),
	},
	Disabled: false,
}

// ErrorRateTrigger alerts when error rate exceeds 5% of total requests.
// Monitors application health and notifies engineering team via Slack.
var ErrorRateTrigger = trigger.Trigger{
	Name:        "High Error Rate",
	Description: "Fires when error rate exceeds 5% over a 10 minute window",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(10),
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.GreaterThanOrEqual("http.status_code", 500),
		},
	},
	Threshold: trigger.GreaterThan(5),
	Frequency: trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#engineering"),
		trigger.EmailAddress("oncall@company.com"),
	},
	Disabled: false,
}

// SlowDatabaseTrigger alerts when database query latency is too high.
// Monitors database performance and notifies database team.
var SlowDatabaseTrigger = trigger.Trigger{
	Name:        "Slow Database Queries",
	Description: "Fires when P95 database query duration exceeds 500ms",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(15),
		Breakdowns: []string{"db.name"},
		Calculations: []query.Calculation{
			query.P95("db.duration_ms"),
			query.Count(),
		},
		Filters: []query.Filter{
			query.Exists("db.statement"),
			query.GT("db.duration_ms", 100),
		},
		Orders: []query.Order{
			{Op: "P95", Column: "db.duration_ms", Order: "descending"},
		},
		Limit: 10,
	},
	Threshold: trigger.GreaterThan(500),
	Frequency: trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#database-team"),
		trigger.PagerDutyService("DB456XYZ"),
	},
	Disabled: false,
}

// LowTrafficTrigger alerts when request volume drops significantly.
// Monitors for potential outages or traffic routing issues.
var LowTrafficTrigger = trigger.Trigger{
	Name:        "Low Traffic Alert",
	Description: "Fires when request rate drops below 100 requests per minute",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(5),
		Calculations: []query.Calculation{
			query.Rate("request.id"),
		},
	},
	Threshold: trigger.LessThan(100),
	Frequency: trigger.Minutes(1),
	Recipients: []trigger.Recipient{
		trigger.PagerDutyService("P789DEF"),
		trigger.SlackChannel("#ops-alerts"),
	},
	Disabled: false,
}

// HighMemoryUsageTrigger alerts when memory consumption is too high.
// Monitors resource usage and helps prevent out-of-memory issues.
var HighMemoryUsageTrigger = trigger.Trigger{
	Name:        "High Memory Usage",
	Description: "Fires when average memory usage exceeds 85% over 10 minutes",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(10),
		Breakdowns: []string{"service.name", "host.name"},
		Calculations: []query.Calculation{
			query.Avg("memory.usage_percent"),
		},
		Filters: []query.Filter{
			query.Exists("memory.usage_percent"),
		},
		Orders: []query.Order{
			{Op: "AVG", Column: "memory.usage_percent", Order: "descending"},
		},
		Limit: 20,
	},
	Threshold: trigger.GreaterThanOrEqual(85),
	Frequency: trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#infrastructure"),
		trigger.EmailAddress("infra-team@company.com"),
	},
	Disabled: false,
}

// AuthenticationFailuresTrigger alerts on suspicious login activity.
// Monitors security events and notifies security team.
var AuthenticationFailuresTrigger = trigger.Trigger{
	Name:        "High Authentication Failures",
	Description: "Fires when authentication failures exceed 50 per minute",
	Dataset:     "security",
	Query: query.Query{
		Dataset:   "security",
		TimeRange: query.Minutes(5),
		Breakdowns: []string{"user.ip"},
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.Equals("auth.status", "failed"),
			query.Contains("http.route", "/auth"),
		},
		Orders: []query.Order{
			{Op: "COUNT", Order: "descending"},
		},
		Limit: 50,
	},
	Threshold: trigger.GreaterThan(50),
	Frequency: trigger.Minutes(2),
	Recipients: []trigger.Recipient{
		trigger.PagerDutyService("SEC123"),
		trigger.SlackChannel("#security-alerts"),
		trigger.WebhookURL("https://security-system.company.com/webhook"),
	},
	Disabled: false,
}

// ApiQuotaExceededTrigger alerts when API quota is being exceeded.
// Monitors API usage and helps prevent rate limiting issues.
var ApiQuotaExceededTrigger = trigger.Trigger{
	Name:        "API Quota Exceeded",
	Description: "Fires when API requests with 429 status exceed threshold",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(10),
		Breakdowns: []string{"api.client_id"},
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.Equals("http.status_code", 429),
			query.Exists("api.client_id"),
		},
		Orders: []query.Order{
			{Op: "COUNT", Order: "descending"},
		},
		Limit: 25,
	},
	Threshold: trigger.GreaterThan(100),
	Frequency: trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#api-team"),
		trigger.EmailAddress("api-team@company.com"),
	},
	Disabled: false,
}

// DeploymentErrorSpikeTrigger alerts on errors immediately after deployment.
// Monitors deployment health during rollout window.
var DeploymentErrorSpikeTrigger = trigger.Trigger{
	Name:        "Post-Deployment Error Spike",
	Description: "Fires when error count spikes in first 15 minutes after deployment",
	Dataset:     "production",
	Query: query.Query{
		Dataset:   "production",
		TimeRange: query.Minutes(15),
		Breakdowns: []string{"service.version", "service.name"},
		Calculations: []query.Calculation{
			query.Count(),
			query.P99("duration_ms"),
		},
		Filters: []query.Filter{
			query.GTE("http.status_code", 500),
			query.Exists("deployment.id"),
		},
		Orders: []query.Order{
			{Op: "COUNT", Order: "descending"},
		},
		Limit: 15,
	},
	Threshold: trigger.GreaterThan(25),
	Frequency: trigger.Minutes(3),
	Recipients: []trigger.Recipient{
		trigger.PagerDutyService("DEPLOY456"),
		trigger.SlackChannel("#deployments"),
	},
	Disabled: false,
}
