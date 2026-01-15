package queries

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestQueriesExist(t *testing.T) {
	tests := []struct {
		name  string
		query query.Query
	}{
		{"CheckoutFlowLatency", CheckoutFlowLatency},
		{"PaymentFraudCorrelation", PaymentFraudCorrelation},
		{"ErrorRateByService", ErrorRateByService},
		{"CheckoutFunnel", CheckoutFunnel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.query.Dataset == "" {
				t.Errorf("%s: Dataset is empty", tt.name)
			}
			if tt.query.TimeRange.TimeRange == 0 && tt.query.TimeRange.StartTime == 0 {
				t.Errorf("%s: TimeRange is not set", tt.name)
			}
			if len(tt.query.Calculations) == 0 {
				t.Errorf("%s: No calculations defined", tt.name)
			}
		})
	}
}

func TestCheckoutFlowLatency(t *testing.T) {
	q := CheckoutFlowLatency

	if q.Dataset != "otel-demo" {
		t.Errorf("Expected dataset 'otel-demo', got '%s'", q.Dataset)
	}

	if len(q.Breakdowns) != 2 {
		t.Errorf("Expected 2 breakdowns, got %d", len(q.Breakdowns))
	}

	if len(q.Calculations) != 4 {
		t.Errorf("Expected 4 calculations, got %d", len(q.Calculations))
	}

	if len(q.Filters) == 0 {
		t.Error("Expected at least one filter (In filter)")
	}

	if len(q.Orders) == 0 {
		t.Error("Expected at least one order clause")
	}
}

func TestPaymentFraudCorrelation(t *testing.T) {
	q := PaymentFraudCorrelation

	if q.Dataset != "otel-demo" {
		t.Errorf("Expected dataset 'otel-demo', got '%s'", q.Dataset)
	}

	if len(q.Calculations) == 0 {
		t.Error("Expected calculations for correlation query")
	}
}

func TestErrorRateByService(t *testing.T) {
	q := ErrorRateByService

	if q.Dataset != "otel-demo" {
		t.Errorf("Expected dataset 'otel-demo', got '%s'", q.Dataset)
	}

	if len(q.Breakdowns) == 0 {
		t.Error("Expected breakdown by service")
	}

	if len(q.Calculations) == 0 {
		t.Error("Expected calculations for error rate")
	}
}

func TestCheckoutFunnel(t *testing.T) {
	q := CheckoutFunnel

	if q.Dataset != "otel-demo" {
		t.Errorf("Expected dataset 'otel-demo', got '%s'", q.Dataset)
	}

	if len(q.Breakdowns) == 0 {
		t.Error("Expected breakdown by stage or route")
	}

	if len(q.Calculations) == 0 {
		t.Error("Expected calculations for funnel counts")
	}
}
