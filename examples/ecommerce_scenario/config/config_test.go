package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	CheckoutFlow struct {
		Dataset    string   `yaml:"dataset"`
		Services   []string `yaml:"services"`
		SLOTargets struct {
			Availability   string `yaml:"availability"`
			LatencyP95     string `yaml:"latency_p95"`
			PaymentSuccess string `yaml:"payment_success"`
		} `yaml:"slo_targets"`
	} `yaml:"checkout_flow"`
}

func TestServicesYAMLValid(t *testing.T) {
	// Read the YAML file
	data, err := os.ReadFile("services.yaml")
	if err != nil {
		t.Fatalf("Failed to read services.yaml: %v", err)
	}

	// Parse the YAML
	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse services.yaml: %v", err)
	}

	// Validate structure
	if config.CheckoutFlow.Dataset == "" {
		t.Error("Dataset is empty")
	}

	if len(config.CheckoutFlow.Services) == 0 {
		t.Error("No services defined")
	}

	// Validate expected services
	expectedServices := []string{"checkoutservice", "cartservice", "paymentservice", "frauddetectionservice"}
	if len(config.CheckoutFlow.Services) != len(expectedServices) {
		t.Errorf("Expected %d services, got %d", len(expectedServices), len(config.CheckoutFlow.Services))
	}

	for i, service := range expectedServices {
		if i >= len(config.CheckoutFlow.Services) {
			t.Errorf("Missing service: %s", service)
			continue
		}
		if config.CheckoutFlow.Services[i] != service {
			t.Errorf("Expected service %s at index %d, got %s", service, i, config.CheckoutFlow.Services[i])
		}
	}

	// Validate SLO targets exist
	if config.CheckoutFlow.SLOTargets.Availability == "" {
		t.Error("Availability SLO target is empty")
	}
	if config.CheckoutFlow.SLOTargets.LatencyP95 == "" {
		t.Error("Latency P95 SLO target is empty")
	}
	if config.CheckoutFlow.SLOTargets.PaymentSuccess == "" {
		t.Error("Payment success SLO target is empty")
	}

	// Validate dataset
	if config.CheckoutFlow.Dataset != "otel-demo" {
		t.Errorf("Expected dataset 'otel-demo', got '%s'", config.CheckoutFlow.Dataset)
	}
}
