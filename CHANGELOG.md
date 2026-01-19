# Changelog

All notable changes to wetwire-honeycomb-go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **LintOpts.Fix and LintOpts.Disable support** (#117)
  - `opts.Fix` support in Linter.Lint() (auto-fix not yet implemented, returns message)
  - `opts.Disable` support to skip specified rule IDs (e.g., `["WHC001", "WHC002"]`)
  - `LintAllWithConfig()` function for configurable linting with disabled rules
  - `LintBoardsWithRules()`, `LintSLOsWithRules()`, `LintTriggersWithRules()` helper functions

### Changed
- **Renamed `internal/discovery` to `internal/discover`** for consistent naming (#112)
- **Lint Severity type migration** to wetwire-core-go/lint (#111)
  - Upgraded wetwire-core-go to v1.16.0 for shared Severity type
  - Added Severity type alias and constants (SeverityError, SeverityWarning, SeverityInfo)
  - Replaced string severity values with typed constants across all lint rules
  - Updated LintResult.Severity and LintConfig.SeverityOverrides to use Severity type
  - Updated FilterBySeverity to accept Severity parameter instead of string
- **MCP server migration** to domain.BuildMCPServer() (#101)
  - Upgraded wetwire-core-go to v1.13.0 for auto-generated MCP server support
  - Replaced 833-line manual MCP implementation with minimal 40-line version
  - Removed manual mcpRegister* functions in favor of domain.BuildMCPServer()
  - Removed `import` command - HoneycombDomain doesn't implement ImporterDomain
  - MCP server now auto-generates all standard tools (init, build, lint, list, graph)

### Added
- **CLI enhancements** for full resource support (#54)
  - `list` command now shows all resources (queries, boards, SLOs, triggers)
  - `lint` command now validates all resource types with WHC001-056 rules
  - `build` command now outputs all resources with grouped JSON structure
  - New `--type` flag to filter resources by type
  - `LintBoards()`, `LintSLOs()`, `LintTriggers()`, `LintAll()` functions in lint package
- **Lint rules WHC030-056** for boards, SLOs, and triggers (#53)
  - Board rules: WHC030 (no panels), WHC034 (exceeds 20 panels)
  - SLO rules: WHC040 (missing name), WHC044 (target out of range), WHC047 (no burn alerts)
  - Trigger rules: WHC050 (missing name), WHC053 (no recipients), WHC054 (frequency under 1min), WHC056 (disabled)
  - New types: `BoardRule`, `SLORule`, `TriggerRule` with Check functions
- **Serialization extensions** for boards, SLOs, and triggers (#52)
  - `BoardToJSON()` / `BoardToJSONPretty()` for board serialization
  - `SLOToJSON()` / `SLOToJSONPretty()` for SLO serialization
  - `TriggerToJSON()` / `TriggerToJSONPretty()` for trigger serialization
  - Export `board.PanelConfig` for serialization access
- **Discovery extensions** for boards, SLOs, and triggers (#51)
  - `DiscoverBoards()` finds board.Board declarations with panel count, query/SLO refs
  - `DiscoverSLOs()` finds slo.SLO declarations with target, time period, burn alerts
  - `DiscoverTriggers()` finds trigger.Trigger declarations with threshold, frequency, recipients
  - `DiscoverAll()` discovers all resource types in one call
  - `DiscoveredResources.TotalCount()` helper for total resource count
- **Trigger package** (`trigger/`) for type-safe trigger declarations (#50)
  - `Trigger` struct with Name, Description, Dataset, Query, Threshold, Frequency, Recipients, Disabled
  - `Threshold` with comparison operators (`GT`, `GTE`, `LT`, `LTE`)
  - `Frequency` with `Minutes()` and `Seconds()` builders
  - `Recipient` with `SlackChannel()`, `PagerDutyService()`, `EmailAddress()`, `WebhookURL()` builders
- **SLO package** (`slo/`) for type-safe SLO declarations (#49)
  - `SLO` struct with Name, Description, Dataset, SLI, Target, TimePeriod, BurnAlerts
  - `SLI` struct with direct query references (GoodEvents, TotalEvents)
  - `Target` and `TimePeriod` types with `Percentage()` and `Days()` builders
  - `BurnAlert` with `FastBurn()` and `SlowBurn()` helpers
  - `Recipient` type for notification targets
- **Board package** (`board/`) for type-safe board declarations (#48)
  - `Board` struct with Name, Description, Panels, PresetFilters, Tags
  - `Panel` interface with Query, SLO, and Text panel implementations
  - `PanelOption` pattern for `WithTitle` and `WithPosition`
- Query type definitions with full Honeycomb Query API support
- AST-based query discovery for Go source files
- JSON serialization with proper Honeycomb Query format
- Query builder with topological ordering
- Comprehensive lint rules (WHC001-WHC010)
- CLI commands: build, lint, import, validate, list, init, graph
- Round-trip testing infrastructure
- Test coverage: Discovery 96.7%, Lint 100%

### Documentation
- CLAUDE.md with AI assistant guidelines
- CLI.md with complete command reference
- FAQ.md with common questions
- LINT_RULES.md with all rule documentation


## [1.6.0] - 2026-01-19

### Changed

- **Claude CLI as default provider for design command**
  - No API key required - uses existing Claude authentication
  - Falls back to Anthropic API if Claude CLI not installed
  - Updated wetwire-core-go to v1.17.1

## [0.1.0] - 2026-01-12

### Added
- Initial scaffold with core infrastructure
- Query package with types for calculations, filters, breakdowns
- Discovery package for AST-based query detection
- Serialize package for JSON output
- Builder package for query construction
- Lint package with 10 validation rules
- Basic CLI structure

[Unreleased]: https://github.com/lex00/wetwire-honeycomb-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/lex00/wetwire-honeycomb-go/releases/tag/v0.1.0
