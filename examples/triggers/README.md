# Trigger Alert Examples

This package demonstrates various trigger patterns for alerting in Honeycomb.

## Triggers

### HighLatencyTrigger

Alerts when P99 latency exceeds 1000ms over a 5 minute window.

**Use case:** Monitor API endpoint performance and notify on-call team via PagerDuty when latency degrades.

**Recipients:** PagerDuty service, Slack channel

```bash
wetwire-honeycomb build ./examples/triggers
```

### ErrorRateTrigger

Alerts when error rate exceeds 5% over a 10 minute window.

**Use case:** Monitor application health and catch error spikes early.

**Recipients:** Slack channel, email

### SlowDatabaseTrigger

Alerts when P95 database query duration exceeds 500ms.

**Use case:** Monitor database performance and identify slow queries before they impact users.

**Recipients:** Slack channel, PagerDuty service

### LowTrafficTrigger

Alerts when request rate drops below 100 requests per minute.

**Use case:** Detect potential outages or traffic routing issues.

**Recipients:** PagerDuty service, Slack channel

### HighMemoryUsageTrigger

Alerts when average memory usage exceeds 85% over 10 minutes.

**Use case:** Monitor resource usage and prevent out-of-memory errors.

**Recipients:** Slack channel, email

### AuthenticationFailuresTrigger

Alerts when authentication failures exceed 50 per minute.

**Use case:** Monitor for suspicious login activity and potential security incidents.

**Recipients:** PagerDuty service, Slack channel, webhook

### ApiQuotaExceededTrigger

Alerts when API requests with 429 status exceed threshold.

**Use case:** Monitor API usage and identify clients hitting rate limits.

**Recipients:** Slack channel, email

### DeploymentErrorSpikeTrigger

Alerts when error count spikes in first 15 minutes after deployment.

**Use case:** Monitor deployment health during rollout window for immediate rollback decisions.

**Recipients:** PagerDuty service, Slack channel

## Output

Run `wetwire-honeycomb build` to generate Honeycomb Trigger JSON:

```json
{
  "queries": {
    "HighLatencyTrigger": {
      "time_range": 300,
      "calculations": [
        {"op": "P99", "column": "duration_ms"}
      ],
      "filters": [
        {"column": "http.route", "op": "exists"}
      ]
    }
  },
  "triggers": {
    "HighLatencyTrigger": {
      "name": "High Latency Alert",
      "description": "Fires when P99 latency exceeds 1000ms over a 5 minute window",
      "dataset": "production",
      "frequency": 120,
      "disabled": false
    }
  }
}
```

## Recipient Types

The examples demonstrate all four recipient types:

- **Slack:** `trigger.SlackChannel("#channel-name")`
- **PagerDuty:** `trigger.PagerDutyService("service-id")`
- **Email:** `trigger.EmailAddress("team@company.com")`
- **Webhook:** `trigger.WebhookURL("https://webhook.company.com/path")`

## Threshold Operators

Triggers support four comparison operators:

- **GreaterThan:** `trigger.GreaterThan(value)`
- **GreaterThanOrEqual:** `trigger.GreaterThanOrEqual(value)`
- **LessThan:** `trigger.LessThan(value)`
- **LessThanOrEqual:** `trigger.LessThanOrEqual(value)`

## Frequency

Evaluation frequency can be specified in minutes or seconds:

- **Minutes:** `trigger.Minutes(5)` - evaluates every 5 minutes
- **Seconds:** `trigger.Seconds(30)` - evaluates every 30 seconds
