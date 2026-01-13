# Changelog

All notable changes to wetwire-honeycomb-go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
