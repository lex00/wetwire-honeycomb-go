// Command design provides AI-assisted query design for Honeycomb.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lex00/wetwire-core-go/agent/agents"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/lex00/wetwire-honeycomb-go/internal/agent"
	"github.com/lex00/wetwire-honeycomb-go/internal/kiro"
	"github.com/spf13/cobra"
)

// newDesignCmd creates the "design" subcommand for AI-assisted query design.
// It uses the Anthropic API for interactive code generation.
func newDesignCmd() *cobra.Command {
	var outputDir string
	var maxLintCycles int
	var stream bool
	var provider string

	cmd := &cobra.Command{
		Use:   "design [prompt]",
		Short: "AI-assisted query design",
		Long: `Start an interactive AI-assisted session to design and generate Honeycomb queries.

The AI agent will:
1. Ask clarifying questions about your requirements
2. Generate Go code using wetwire-honeycomb patterns
3. Run the linter and fix any issues
4. Build the Query JSON

Example:
    wetwire-honeycomb design "Show me P99 latency by endpoint for the last 2 hours"
    wetwire-honeycomb design "Create an SLO dashboard for API latency"
    wetwire-honeycomb design "Find slow database queries with error tracking"`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := strings.Join(args, " ")
			if prompt == "" {
				return fmt.Errorf("prompt is required")
			}
			return runDesign(prompt, outputDir, maxLintCycles, stream, provider)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated files")
	cmd.Flags().IntVarP(&maxLintCycles, "max-lint-cycles", "l", 5, "Maximum lint/fix cycles")
	cmd.Flags().BoolVarP(&stream, "stream", "s", true, "Stream AI responses")
	cmd.Flags().StringVar(&provider, "provider", "anthropic", "AI provider: 'anthropic' or 'kiro'")

	return cmd
}

// runDesign starts an AI-assisted design session using the specified provider.
// It creates a runner agent that generates code, runs the linter, and fixes issues.
func runDesign(prompt, outputDir string, maxLintCycles int, stream bool, provider string) error {
	// Handle kiro provider
	if provider == "kiro" {
		return runDesignKiro(prompt)
	}

	// Default to anthropic provider
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nInterrupted, cleaning up...")
		cancel()
	}()

	// Create session for tracking
	session := results.NewSession("human", "design")

	// Create human developer (reads from stdin)
	reader := bufio.NewReader(os.Stdin)
	developer := orchestrator.NewHumanDeveloper(func() (string, error) {
		return reader.ReadString('\n')
	})

	// Create stream handler if streaming enabled
	var streamHandler agents.StreamHandler
	if stream {
		streamHandler = func(text string) {
			fmt.Print(text)
		}
	}

	// Create runner agent with Honeycomb domain
	runner, err := agents.NewRunnerAgent(agents.RunnerConfig{
		WorkDir:       outputDir,
		MaxLintCycles: maxLintCycles,
		Session:       session,
		Developer:     developer,
		StreamHandler: streamHandler,
		Domain:        agent.HoneycombDomain(),
	})
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	fmt.Println("Starting AI-assisted design session...")
	fmt.Println("The AI will ask questions and generate Honeycomb query code.")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Run the agent
	if err := runner.Run(ctx, prompt); err != nil {
		return fmt.Errorf("design session failed: %w", err)
	}

	// Print summary
	fmt.Println("\n--- Session Summary ---")
	fmt.Printf("Generated files: %d\n", len(runner.GetGeneratedFiles()))
	for _, f := range runner.GetGeneratedFiles() {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Printf("Lint cycles: %d\n", runner.GetLintCycles())
	fmt.Printf("Lint passed: %v\n", runner.LintPassed())

	return nil
}

// runDesignKiro starts a Kiro CLI chat session for interactive design.
func runDesignKiro(prompt string) error {
	fmt.Println("Starting Kiro CLI chat session...")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Launch interactive kiro session
	return kiro.LaunchChat("wetwire-honeycomb-runner", prompt)
}
