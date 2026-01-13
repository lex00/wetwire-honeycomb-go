package serialize

import (
	"bytes"
	"encoding/json"

	"github.com/lex00/wetwire-honeycomb-go/slo"
)

// sloJSON is the internal representation for SLO JSON serialization.
type sloJSON struct {
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	Dataset          string          `json:"dataset,omitempty"`
	SLI              *sliJSON        `json:"sli,omitempty"`
	TargetPerMillion int             `json:"target_per_million,omitempty"`
	TimePeriodDays   int             `json:"time_period_days,omitempty"`
	BurnAlerts       []burnAlertJSON `json:"burn_alerts,omitempty"`
}

type sliJSON struct {
	GoodEvents  *queryJSON `json:"good_events,omitempty"`
	TotalEvents *queryJSON `json:"total_events,omitempty"`
}

type burnAlertJSON struct {
	Name       string         `json:"name,omitempty"`
	AlertType  string         `json:"alert_type"`
	Threshold  float64        `json:"threshold"`
	WindowHours int           `json:"window_hours,omitempty"`
	Recipients []recipientJSON `json:"recipients,omitempty"`
}

type recipientJSON struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

// SLOToJSON serializes an SLO to Honeycomb SLO JSON format.
func SLOToJSON(s slo.SLO) ([]byte, error) {
	js := toSLOJSON(s)
	return json.Marshal(js)
}

// SLOToJSONPretty serializes an SLO to indented JSON format.
func SLOToJSONPretty(s slo.SLO) ([]byte, error) {
	js := toSLOJSON(s)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(js); err != nil {
		return nil, err
	}
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

func toSLOJSON(s slo.SLO) sloJSON {
	js := sloJSON{
		Name:        s.Name,
		Description: s.Description,
		Dataset:     s.Dataset,
	}

	// Convert target percentage to per-million
	if s.Target.Percentage > 0 {
		js.TargetPerMillion = int(s.Target.Percentage * 10000)
	}

	// Convert time period
	if s.TimePeriod.Days > 0 {
		js.TimePeriodDays = s.TimePeriod.Days
	}

	// Convert SLI
	if s.SLI.GoodEvents.Dataset != "" || s.SLI.TotalEvents.Dataset != "" ||
		len(s.SLI.GoodEvents.Calculations) > 0 || len(s.SLI.TotalEvents.Calculations) > 0 {
		goodEventsQuery := toQueryJSON(s.SLI.GoodEvents)
		totalEventsQuery := toQueryJSON(s.SLI.TotalEvents)
		js.SLI = &sliJSON{
			GoodEvents:  &goodEventsQuery,
			TotalEvents: &totalEventsQuery,
		}
	}

	// Convert burn alerts
	if len(s.BurnAlerts) > 0 {
		js.BurnAlerts = make([]burnAlertJSON, len(s.BurnAlerts))
		for i, ba := range s.BurnAlerts {
			js.BurnAlerts[i] = toBurnAlertJSON(ba)
		}
	}

	return js
}

func toBurnAlertJSON(ba slo.BurnAlert) burnAlertJSON {
	jba := burnAlertJSON{
		Name:      ba.Name,
		AlertType: string(ba.AlertType),
		Threshold: ba.Threshold,
	}

	// Convert window to hours
	if ba.Window.Hours > 0 {
		jba.WindowHours = ba.Window.Hours
	} else if ba.Window.Days > 0 {
		jba.WindowHours = ba.Window.Days * 24
	}

	// Convert recipients
	if len(ba.Recipients) > 0 {
		jba.Recipients = make([]recipientJSON, len(ba.Recipients))
		for i, r := range ba.Recipients {
			jba.Recipients[i] = recipientJSON{
				Type:   r.Type,
				Target: r.Target,
			}
		}
	}

	return jba
}
