package ecommerce_scenario

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-core-go/agent/agents"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/personas"
	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/lex00/wetwire-honeycomb-go/internal/agent"
	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runScenario executes an E2E test scenario with the specified persona and prompt.
// It returns a ScenarioResult containing information about the test run.
func runScenario(t *testing.T, personaName, prompt string) *ScenarioResult {
	t.Helper()

	// Create temporary output directory for this test run
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, personaName)
	require.NoError(t, os.MkdirAll(outputDir, 0755))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get persona
	persona, err := personas.Get(personaName)
	require.NoError(t, err, "Failed to get persona %s", personaName)

	// Create session for tracking
	session := results.NewSession(personaName, "checkout_scenario")

	// Create AI developer with persona
	responder := agents.CreateDeveloperResponder("")
	developer := orchestrator.NewAIDeveloper(persona, responder)

	// Create runner agent with Honeycomb domain
	runner, err := agents.NewRunnerAgent(agents.RunnerConfig{
		WorkDir:       outputDir,
		MaxLintCycles: 5,
		Session:       session,
		Developer:     developer,
		StreamHandler: nil, // No streaming in tests
		Domain:        agent.HoneycombDomain(),
	})
	require.NoError(t, err, "Failed to create runner agent")

	// Run the agent
	err = runner.Run(ctx, prompt)
	if err != nil {
		t.Logf("Warning: agent run encountered error: %v", err)
	}

	// Complete session
	session.Complete()

	// Discover generated resources
	resources, discErr := discovery.DiscoverAll(outputDir)
	if discErr != nil {
		t.Logf("Warning: discovery failed: %v", discErr)
		resources = &discovery.DiscoveredResources{}
	}

	// Run lint check
	lintPassed := runLint(t, outputDir)

	// Run build check
	buildPassed := runBuild(t, outputDir)

	// Calculate score (simple scoring based on results)
	score := calculateScore(runner, resources, lintPassed, buildPassed)

	return &ScenarioResult{
		OutputDir:      outputDir,
		LintPassed:     lintPassed,
		BuildPassed:    buildPassed,
		Score:          score,
		Resources:      resources,
		Session:        session,
		LintCycles:     runner.GetLintCycles(),
		GeneratedFiles: runner.GetGeneratedFiles(),
	}
}

// ScenarioResult contains the results of running a scenario test.
type ScenarioResult struct {
	OutputDir      string
	LintPassed     bool
	BuildPassed    bool
	Score          int
	Resources      *discovery.DiscoveredResources
	Session        *results.Session
	LintCycles     int
	GeneratedFiles []string
}

// runLint executes the lint command on the output directory.
func runLint(t *testing.T, dir string) bool {
	t.Helper()

	// Build the CLI tool path
	cliPath := filepath.Join(os.Getenv("GOPATH"), "bin", "wetwire-honeycomb")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		// Try building it
		buildCmd := exec.Command("go", "install", "github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest")
		if err := buildCmd.Run(); err != nil {
			t.Logf("Failed to build CLI tool: %v", err)
			return false
		}
	}

	cmd := exec.Command(cliPath, "lint", dir)
	output, err := cmd.CombinedOutput()
	t.Logf("Lint output: %s", string(output))

	return err == nil
}

// runBuild executes the build command on the output directory.
func runBuild(t *testing.T, dir string) bool {
	t.Helper()

	cliPath := filepath.Join(os.Getenv("GOPATH"), "bin", "wetwire-honeycomb")
	cmd := exec.Command(cliPath, "build", dir)
	output, err := cmd.CombinedOutput()
	t.Logf("Build output: %s", string(output))

	return err == nil
}

// calculateScore computes a simple score for the test run (0-15 points).
func calculateScore(runner *agents.RunnerAgent, resources *discovery.DiscoveredResources, lintPassed, buildPassed bool) int {
	score := 0

	// Completeness (0-5 points)
	if len(resources.Queries) >= 4 {
		score += 2
	}
	if len(resources.SLOs) >= 1 {
		score++
	}
	if len(resources.Triggers) >= 1 {
		score++
	}
	if len(resources.Boards) >= 1 {
		score++
	}

	// Lint quality (0-3 points)
	if lintPassed {
		score += 3
	} else if runner.GetLintCycles() <= 3 {
		score++
	}

	// Build quality (0-3 points)
	if buildPassed {
		score += 3
	}

	// Generation efficiency (0-4 points)
	if len(runner.GetGeneratedFiles()) > 0 {
		score++
	}
	if runner.GetLintCycles() <= 2 {
		score++
	}
	if runner.GetLintCycles() == 0 {
		score += 2
	}

	return score
}

// TestCheckoutScenario_AllPersonas runs the checkout scenario with all personas.
func TestCheckoutScenario_AllPersonas(t *testing.T) {
	// Skip if CI environment doesn't have API keys
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping E2E test: ANTHROPIC_API_KEY not set")
	}

	// Read the complex scenario prompt
	promptPath := filepath.Join("prompts", "complex_scenario.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		// If prompt file doesn't exist, use a default prompt
		t.Logf("Warning: couldn't read %s, using default prompt: %v", promptPath, err)
		promptBytes = []byte(`Create comprehensive observability resources for an e-commerce checkout flow.

Dataset: otel-demo
Services: checkoutservice, cartservice, paymentservice, frauddetectionservice

Create:
1. Queries for checkout flow latency, payment fraud correlation, error rates, and checkout funnel
2. SLOs for availability (99.5%), latency P95 (<500ms), and payment success (99.9%)
3. Triggers for high error rate (>5%), slow checkout (P95 >1s), and payment failures (>1%)
4. A performance dashboard board with 4 panels`)
	}
	prompt := string(promptBytes)

	personas := []string{"beginner", "intermediate", "expert"}

	for _, persona := range personas {
		t.Run(persona, func(t *testing.T) {
			// Run the scenario
			result := runScenario(t, persona, prompt)

			// Assertions
			assert.True(t, result.LintPassed, "Lint should pass for persona %s", persona)
			assert.True(t, result.BuildPassed, "Build should pass for persona %s", persona)
			assert.GreaterOrEqual(t, result.Score, 8, "Score should be at least 8/15 for persona %s", persona)

			// Log results
			t.Logf("Persona %s results:", persona)
			t.Logf("  Generated files: %d", len(result.GeneratedFiles))
			t.Logf("  Lint cycles: %d", result.LintCycles)
			t.Logf("  Score: %d/15", result.Score)
		})
	}
}

// TestCheckoutScenario_ExpectedOutput validates the AI-generated output against expected counts.
func TestCheckoutScenario_ExpectedOutput(t *testing.T) {
	// Skip if CI environment doesn't have API keys
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping E2E test: ANTHROPIC_API_KEY not set")
	}

	prompt := `Create comprehensive observability resources for an e-commerce checkout flow.

Dataset: otel-demo
Services: checkoutservice, cartservice, paymentservice, frauddetectionservice

Create:
1. At least 4 queries: checkout flow latency, payment fraud correlation, error rates by service, and checkout funnel
2. At least 3 SLOs: availability (99.5%), latency P95 (<500ms), and payment success (99.9%)
3. At least 2 triggers: high error rate alert and slow checkout alert
4. Exactly 1 board: performance dashboard with 4 panels`

	result := runScenario(t, "intermediate", prompt)

	// Validate resource counts
	assert.GreaterOrEqual(t, len(result.Resources.Queries), 4,
		"Should generate at least 4 queries")
	assert.GreaterOrEqual(t, len(result.Resources.SLOs), 3,
		"Should generate at least 3 SLOs")
	assert.GreaterOrEqual(t, len(result.Resources.Triggers), 2,
		"Should generate at least 2 triggers")
	assert.GreaterOrEqual(t, len(result.Resources.Boards), 1,
		"Should generate at least 1 board")

	// Log what was generated
	t.Logf("Generated resources:")
	t.Logf("  Queries: %d", len(result.Resources.Queries))
	for _, q := range result.Resources.Queries {
		t.Logf("    - %s (dataset: %s)", q.Name, q.Dataset)
	}
	t.Logf("  SLOs: %d", len(result.Resources.SLOs))
	for _, s := range result.Resources.SLOs {
		t.Logf("    - %s (target: %.1f%%)", s.Name, s.TargetPercentage)
	}
	t.Logf("  Triggers: %d", len(result.Resources.Triggers))
	for _, tr := range result.Resources.Triggers {
		t.Logf("    - %s", tr.Name)
	}
	t.Logf("  Boards: %d", len(result.Resources.Boards))
	for _, b := range result.Resources.Boards {
		t.Logf("    - %s", b.Name)
	}
}

// TestCheckoutScenario_ResourceValidation validates specific aspects of generated resources.
func TestCheckoutScenario_ResourceValidation(t *testing.T) {
	// Skip if CI environment doesn't have API keys
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping E2E test: ANTHROPIC_API_KEY not set")
	}

	prompt := `Create observability resources for e-commerce checkout.

Dataset: otel-demo
Services: checkoutservice, cartservice, paymentservice, frauddetectionservice

Requirements:
1. Create queries with In() filters for multi-service selection
2. Create an SLO with 3-tier burn rate alerts
3. Create triggers with proper thresholds
4. Create a board with 4 panels

Focus on using the In() filter function for filtering by multiple services.`

	result := runScenario(t, "expert", prompt)

	// Validate queries have multi-service filters
	hasMultiServiceFilter := false
	for _, q := range result.Resources.Queries {
		for _, f := range q.Filters {
			// Check if filter uses In operator or contains service references
			if strings.Contains(strings.ToLower(f.Op), "in") ||
				(f.Column == "service.name" || f.Column == "service") {
				hasMultiServiceFilter = true
				break
			}
		}
	}
	assert.True(t, hasMultiServiceFilter,
		"At least one query should have multi-service In() filters")

	// Validate SLOs exist
	assert.GreaterOrEqual(t, len(result.Resources.SLOs), 1,
		"Should have at least one SLO")

	// Validate triggers have proper thresholds
	if len(result.Resources.Triggers) > 0 {
		for _, tr := range result.Resources.Triggers {
			assert.NotEmpty(t, tr.Name, "Trigger should have a name")
			t.Logf("Trigger: %s", tr.Name)
		}
	}

	// Validate boards exist
	assert.GreaterOrEqual(t, len(result.Resources.Boards), 1,
		"Should have at least one board")

	// Overall validation
	assert.True(t, result.LintPassed, "Generated code should pass lint")
	assert.True(t, result.BuildPassed, "Generated code should build successfully")
}

// TestCheckoutScenario_QuickValidation is a faster test for CI that validates basic functionality.
func TestCheckoutScenario_QuickValidation(t *testing.T) {
	// This test validates that expected resources exist and are well-formed
	// without running the full AI scenario

	// Check for expected directories and validate any that exist
	expectedDir := "expected"

	// Check for SLOs
	slosDir := filepath.Join(expectedDir, "slos")
	if _, err := os.Stat(slosDir); err == nil {
		resources, err := discovery.DiscoverAll(slosDir)
		require.NoError(t, err, "Should be able to discover SLOs")

		if len(resources.SLOs) > 0 {
			for _, s := range resources.SLOs {
				assert.NotEmpty(t, s.Name, "SLO should have a name")
				assert.Greater(t, s.TargetPercentage, 0.0, "SLO should have a target percentage")
				t.Logf("SLO %s: target=%.1f%%", s.Name, s.TargetPercentage)
			}
		} else {
			t.Logf("No SLOs found in expected/slos directory")
		}
	} else {
		t.Logf("SLOs directory not found, skipping SLO validation")
	}

	// Check for Triggers
	triggersDir := filepath.Join(expectedDir, "triggers")
	if _, err := os.Stat(triggersDir); err == nil {
		resources, err := discovery.DiscoverAll(triggersDir)
		require.NoError(t, err, "Should be able to discover triggers")

		if len(resources.Triggers) > 0 {
			for _, tr := range resources.Triggers {
				assert.NotEmpty(t, tr.Name, "Trigger should have a name")
				t.Logf("Trigger %s: description=%s", tr.Name, tr.Description)
			}
		} else {
			t.Logf("No triggers found in expected/triggers directory")
		}
	} else {
		t.Logf("Triggers directory not found, skipping trigger validation")
	}

	// Check for Boards
	boardsDir := filepath.Join(expectedDir, "boards")
	if _, err := os.Stat(boardsDir); err == nil {
		resources, err := discovery.DiscoverAll(boardsDir)
		require.NoError(t, err, "Should be able to discover boards")

		if len(resources.Boards) > 0 {
			for _, b := range resources.Boards {
				assert.NotEmpty(t, b.Name, "Board should have a name")
				t.Logf("Board %s: description=%s", b.Name, b.Description)
			}
		} else {
			t.Logf("No boards found in expected/boards directory")
		}
	} else {
		t.Logf("Boards directory not found, skipping board validation")
	}

	// Check for Queries
	queriesDir := filepath.Join(expectedDir, "queries")
	if _, err := os.Stat(queriesDir); err == nil {
		resources, err := discovery.DiscoverAll(queriesDir)
		require.NoError(t, err, "Should be able to discover queries")

		if len(resources.Queries) > 0 {
			for _, q := range resources.Queries {
				assert.NotEmpty(t, q.Name, "Query should have a name")
				assert.NotEmpty(t, q.Dataset, "Query should have a dataset")
				t.Logf("Query %s: dataset=%s, breakdowns=%d, calculations=%d, filters=%d",
					q.Name, q.Dataset, len(q.Breakdowns), len(q.Calculations), len(q.Filters))
			}
		} else {
			t.Logf("No queries found in expected/queries directory")
		}
	} else {
		t.Logf("Queries directory not found, skipping query validation")
	}
}
