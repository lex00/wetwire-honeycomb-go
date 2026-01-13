package serialize

import (
	"bytes"
	"encoding/json"

	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

// triggerJSON is the internal representation for Trigger JSON serialization.
type triggerJSON struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Dataset     string              `json:"dataset,omitempty"`
	Query       *queryJSON          `json:"query,omitempty"`
	Threshold   *thresholdJSON      `json:"threshold,omitempty"`
	Frequency   int                 `json:"frequency,omitempty"`
	Recipients  []triggerRecipientJSON `json:"recipients,omitempty"`
	Disabled    bool                `json:"disabled"`
}

type thresholdJSON struct {
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

type triggerRecipientJSON struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

// TriggerToJSON serializes a Trigger to Honeycomb Trigger JSON format.
func TriggerToJSON(t trigger.Trigger) ([]byte, error) {
	jt := toTriggerJSON(t)
	return json.Marshal(jt)
}

// TriggerToJSONPretty serializes a Trigger to indented JSON format.
func TriggerToJSONPretty(t trigger.Trigger) ([]byte, error) {
	jt := toTriggerJSON(t)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jt); err != nil {
		return nil, err
	}
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

func toTriggerJSON(t trigger.Trigger) triggerJSON {
	jt := triggerJSON{
		Name:        t.Name,
		Description: t.Description,
		Dataset:     t.Dataset,
		Disabled:    t.Disabled,
	}

	// Convert query
	if t.Query.Dataset != "" || len(t.Query.Calculations) > 0 || t.Query.TimeRange.TimeRange > 0 {
		q := toQueryJSON(t.Query)
		jt.Query = &q
	}

	// Convert threshold
	if t.Threshold.Op != "" {
		jt.Threshold = &thresholdJSON{
			Op:    string(t.Threshold.Op),
			Value: t.Threshold.Value,
		}
	}

	// Convert frequency
	if t.Frequency.Seconds > 0 {
		jt.Frequency = t.Frequency.Seconds
	}

	// Convert recipients
	if len(t.Recipients) > 0 {
		jt.Recipients = make([]triggerRecipientJSON, len(t.Recipients))
		for i, r := range t.Recipients {
			jt.Recipients[i] = triggerRecipientJSON{
				Type:   string(r.Type),
				Target: r.Target,
			}
		}
	}

	return jt
}
