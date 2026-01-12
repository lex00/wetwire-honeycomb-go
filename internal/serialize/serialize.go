// Package serialize provides JSON serialization for Honeycomb queries.
package serialize

import (
	"bytes"
	"encoding/json"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

// queryJSON is the internal representation for JSON serialization.
// It matches the Honeycomb Query API format with snake_case keys.
type queryJSON struct {
	TimeRange         int               `json:"time_range,omitempty"`
	StartTime         int               `json:"start_time,omitempty"`
	EndTime           int               `json:"end_time,omitempty"`
	Breakdowns        []string          `json:"breakdowns,omitempty"`
	Calculations      []calculationJSON `json:"calculations,omitempty"`
	Filters           []filterJSON      `json:"filters,omitempty"`
	FilterCombination string            `json:"filter_combination,omitempty"`
	Orders            []orderJSON       `json:"orders,omitempty"`
	Limit             int               `json:"limit,omitempty"`
	Granularity       int               `json:"granularity,omitempty"`
}

type calculationJSON struct {
	Op     string `json:"op"`
	Column string `json:"column,omitempty"`
}

type filterJSON struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  any    `json:"value,omitempty"`
}

type orderJSON struct {
	Column string `json:"column,omitempty"`
	Op     string `json:"op,omitempty"`
	Order  string `json:"order"`
}

// ToJSON serializes a Query to Honeycomb Query JSON format.
func ToJSON(q query.Query) ([]byte, error) {
	jq := toQueryJSON(q)
	return json.Marshal(jq)
}

// ToJSONPretty serializes a Query to indented JSON format.
func ToJSONPretty(q query.Query) ([]byte, error) {
	jq := toQueryJSON(q)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jq); err != nil {
		return nil, err
	}
	// Remove trailing newline from Encode
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

func toQueryJSON(q query.Query) queryJSON {
	jq := queryJSON{
		TimeRange:         q.TimeRange.TimeRange,
		StartTime:         q.TimeRange.StartTime,
		EndTime:           q.TimeRange.EndTime,
		Breakdowns:        q.Breakdowns,
		FilterCombination: q.FilterCombination,
		Limit:             q.Limit,
		Granularity:       q.Granularity,
	}

	// Convert calculations
	if len(q.Calculations) > 0 {
		jq.Calculations = make([]calculationJSON, len(q.Calculations))
		for i, c := range q.Calculations {
			jq.Calculations[i] = calculationJSON{
				Op:     c.Op,
				Column: c.Column,
			}
		}
	}

	// Convert filters
	if len(q.Filters) > 0 {
		jq.Filters = make([]filterJSON, len(q.Filters))
		for i, f := range q.Filters {
			jq.Filters[i] = filterJSON{
				Column: f.Column,
				Op:     f.Op,
				Value:  f.Value,
			}
		}
	}

	// Convert orders
	if len(q.Orders) > 0 {
		jq.Orders = make([]orderJSON, len(q.Orders))
		for i, o := range q.Orders {
			jq.Orders[i] = orderJSON{
				Column: o.Column,
				Op:     o.Op,
				Order:  o.Order,
			}
		}
	}

	return jq
}
