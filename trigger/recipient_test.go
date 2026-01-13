package trigger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipient_BasicFields(t *testing.T) {
	r := Recipient{
		Type:   Slack,
		Target: "#alerts",
	}

	assert.Equal(t, Slack, r.Type)
	assert.Equal(t, "#alerts", r.Target)
}

func TestRecipientType_Constants(t *testing.T) {
	assert.Equal(t, RecipientType("slack"), Slack)
	assert.Equal(t, RecipientType("pagerduty"), PagerDuty)
	assert.Equal(t, RecipientType("email"), Email)
	assert.Equal(t, RecipientType("webhook"), Webhook)
}

func TestSlackChannel_Builder(t *testing.T) {
	tests := []struct {
		name    string
		channel string
	}{
		{"alerts", "#alerts"},
		{"oncall", "#oncall"},
		{"engineering", "#engineering"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := SlackChannel(tt.channel)
			assert.Equal(t, Slack, r.Type)
			assert.Equal(t, tt.channel, r.Target)
		})
	}
}

func TestPagerDutyService_Builder(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
	}{
		{"api-team", "api-team"},
		{"platform", "platform"},
		{"service-123", "service-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := PagerDutyService(tt.serviceID)
			assert.Equal(t, PagerDuty, r.Type)
			assert.Equal(t, tt.serviceID, r.Target)
		})
	}
}

func TestEmailAddress_Builder(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"team", "team@example.com"},
		{"alerts", "alerts@company.io"},
		{"oncall", "oncall@internal.dev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := EmailAddress(tt.email)
			assert.Equal(t, Email, r.Type)
			assert.Equal(t, tt.email, r.Target)
		})
	}
}

func TestWebhookURL_Builder(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"example", "https://example.com/webhook"},
		{"internal", "https://internal.company.io/alerts"},
		{"custom", "https://hooks.custom.dev/abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := WebhookURL(tt.url)
			assert.Equal(t, Webhook, r.Type)
			assert.Equal(t, tt.url, r.Target)
		})
	}
}

func TestMultipleRecipients(t *testing.T) {
	recipients := []Recipient{
		SlackChannel("#alerts"),
		PagerDutyService("api-team"),
		EmailAddress("team@example.com"),
		WebhookURL("https://example.com/webhook"),
	}

	assert.Len(t, recipients, 4)

	expectedTypes := []RecipientType{Slack, PagerDuty, Email, Webhook}
	for i, r := range recipients {
		assert.Equal(t, expectedTypes[i], r.Type)
	}
}
