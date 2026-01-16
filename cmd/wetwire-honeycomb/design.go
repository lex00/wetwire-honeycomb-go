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
	"github.com/lex00/wetwire-core-go/agent/results"
	coremcp "github.com/lex00/wetwire-core-go/mcp"
	anthropicprovider "github.com/lex00/wetwire-core-go/providers/anthropic"
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
// It creates a unified agent with MCP tools for code generation.
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
	humanDeveloper := &humanDeveloperAdapter{
		reader: reader,
	}

	// Create stream handler if streaming enabled
	var streamHandler agents.StreamHandler
	if stream {
		streamHandler = func(text string) {
			fmt.Print(text)
		}
	}

	// Create MCP server with Honeycomb tools
	mcpServer := coremcp.NewServer(coremcp.Config{
		Name:    "wetwire-honeycomb-design",
		Version: "1.0.0",
	})

	// Register standard wetwire tools for design mode
	mcpRegisterDesignTools(mcpServer, outputDir)

	// Create Anthropic provider
	anthropicProvider, err := anthropicprovider.New(anthropicprovider.Config{})
	if err != nil {
		return fmt.Errorf("creating provider: %w", err)
	}

	// Create unified agent with MCP server
	designAgent, err := agents.NewAgent(agents.AgentConfig{
		Provider:      anthropicProvider,
		MCPServer:     agents.NewMCPServerAdapter(mcpServer),
		Session:       session,
		Developer:     humanDeveloper,
		StreamHandler: streamHandler,
		SystemPrompt:  agent.HoneycombSystemPrompt(),
	})
	if err != nil {
		return fmt.Errorf("creating agent: %w", err)
	}

	fmt.Println("Starting AI-assisted design session...")
	fmt.Println("The AI will ask questions and generate Honeycomb query code.")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Run the agent
	if err := designAgent.Run(ctx, prompt); err != nil {
		return fmt.Errorf("design session failed: %w", err)
	}

	fmt.Println("\n--- Session Complete ---")

	return nil
}

// humanDeveloperAdapter adapts orchestrator.HumanDeveloper to agents.Developer interface.
type humanDeveloperAdapter struct {
	reader *bufio.Reader
}

// Respond implements agents.Developer interface.
func (h *humanDeveloperAdapter) Respond(ctx context.Context, message string) (string, error) {
	fmt.Printf("\n%s\n> ", message)
	answer, err := h.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

// mcpRegisterDesignTools registers tools needed for design mode.
// These are the same tools as MCP server plus file operations.
func mcpRegisterDesignTools(server *coremcp.Server, workDir string) {
	// Register standard Honeycomb tools
	mcpRegisterStandardTools(server)

	// Add file write/read tools for design mode
	server.RegisterToolWithSchema("wetwire_write", "Write content to a file", func(ctx context.Context, args map[string]any) (string, error) {
		return coremcp.DefaultFileWriteHandler(ctx, args)
	}, coremcp.WriteSchema)

	server.RegisterToolWithSchema("wetwire_read", "Read content from a file", func(ctx context.Context, args map[string]any) (string, error) {
		return coremcp.DefaultFileReadHandler(ctx, args)
	}, coremcp.ReadSchema)
}

// runDesignKiro starts a Kiro CLI chat session for interactive design.
func runDesignKiro(prompt string) error {
	fmt.Println("Starting Kiro CLI chat session...")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Launch interactive kiro session
	return kiro.LaunchChat("wetwire-honeycomb-runner", prompt)
}
