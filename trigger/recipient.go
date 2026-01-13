package trigger

// Recipient represents a notification target for triggers.
type Recipient struct {
	// Type is the recipient type (slack, pagerduty, email, webhook)
	Type RecipientType

	// Target is the destination (channel, service ID, email address, URL)
	Target string
}

// RecipientType represents the type of notification recipient.
type RecipientType string

const (
	// Slack sends notifications to a Slack channel
	Slack RecipientType = "slack"

	// PagerDuty sends notifications to a PagerDuty service
	PagerDuty RecipientType = "pagerduty"

	// Email sends notifications to an email address
	Email RecipientType = "email"

	// Webhook sends notifications to a webhook URL
	Webhook RecipientType = "webhook"
)

// SlackChannel creates a Recipient for a Slack channel.
func SlackChannel(channel string) Recipient {
	return Recipient{Type: Slack, Target: channel}
}

// PagerDutyService creates a Recipient for a PagerDuty service.
func PagerDutyService(serviceID string) Recipient {
	return Recipient{Type: PagerDuty, Target: serviceID}
}

// EmailAddress creates a Recipient for an email address.
func EmailAddress(email string) Recipient {
	return Recipient{Type: Email, Target: email}
}

// WebhookURL creates a Recipient for a webhook URL.
func WebhookURL(url string) Recipient {
	return Recipient{Type: Webhook, Target: url}
}
