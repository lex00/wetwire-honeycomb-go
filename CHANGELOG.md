# Changelog

All notable changes to wetwire-honeycomb-go will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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
