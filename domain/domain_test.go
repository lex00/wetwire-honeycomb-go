package domain

import (
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestHoneycombDomainImplementsInterface(t *testing.T) {
	// Compile-time check that HoneycombDomain implements Domain
	var _ coredomain.Domain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainImplementsListerDomain(t *testing.T) {
	// Compile-time check that HoneycombDomain implements ListerDomain
	var _ coredomain.ListerDomain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainImplementsGrapherDomain(t *testing.T) {
	// Compile-time check that HoneycombDomain implements GrapherDomain
	var _ coredomain.GrapherDomain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainName(t *testing.T) {
	d := &HoneycombDomain{}
	if d.Name() != "honeycomb" {
		t.Errorf("expected name 'honeycomb', got %q", d.Name())
	}
}

func TestHoneycombDomainVersion(t *testing.T) {
	d := &HoneycombDomain{}
	v := d.Version()
	if v == "" {
		t.Error("version should not be empty")
	}
}

func TestHoneycombDomainBuilder(t *testing.T) {
	d := &HoneycombDomain{}
	b := d.Builder()
	if b == nil {
		t.Error("builder should not be nil")
	}
}

func TestHoneycombDomainLinter(t *testing.T) {
	d := &HoneycombDomain{}
	l := d.Linter()
	if l == nil {
		t.Error("linter should not be nil")
	}
}

func TestHoneycombDomainInitializer(t *testing.T) {
	d := &HoneycombDomain{}
	i := d.Initializer()
	if i == nil {
		t.Error("initializer should not be nil")
	}
}

func TestHoneycombDomainValidator(t *testing.T) {
	d := &HoneycombDomain{}
	v := d.Validator()
	if v == nil {
		t.Error("validator should not be nil")
	}
}

func TestHoneycombDomainLister(t *testing.T) {
	d := &HoneycombDomain{}
	l := d.Lister()
	if l == nil {
		t.Error("lister should not be nil")
	}
}

func TestHoneycombDomainGrapher(t *testing.T) {
	d := &HoneycombDomain{}
	g := d.Grapher()
	if g == nil {
		t.Error("grapher should not be nil")
	}
}

func TestCreateRootCommand(t *testing.T) {
	cmd := CreateRootCommand(&HoneycombDomain{})
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}
	if cmd.Use != "wetwire-honeycomb" {
		t.Errorf("expected Use 'wetwire-honeycomb', got %q", cmd.Use)
	}
}
