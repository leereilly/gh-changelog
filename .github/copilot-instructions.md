# Project Overview

gh-changelog is a GitHub CLI extension that allows users to view the latest GitHub Changelog entries directly in their terminal. The tool fetches and parses the RSS feed from GitHub's changelog blog and presents it in a clean, terminal-friendly format.

# Technology Stack

- **Language**: Go 1.21+
- **Platform**: GitHub CLI extension
- **Key Libraries**: Standard library (encoding/xml, net/http, flag)
- **Testing**: Go's built-in testing framework
- **Distribution**: GitHub CLI extension system

# Coding Standards

- Follow standard Go conventions and idioms
- Use gofmt for code formatting
- Prefer explicit error handling over panic
- Use descriptive variable and function names
- Keep functions focused and single-purpose
- Use standard library packages when possible
- Minimize external dependencies

## Go-Specific Guidelines

- Use named return values sparingly, only when they improve clarity
- Prefer composition over inheritance
- Use interfaces for abstractions
- Handle errors explicitly - do not ignore them
- Use defer for cleanup operations
- Follow Go's error wrapping conventions with fmt.Errorf

## Code Organization

- Main business logic in main.go
- Tests in main_test.go
- Keep the codebase simple and maintainable
- Avoid over-engineering for this CLI tool

# Testing Strategy

- Write unit tests for all parsing and formatting functions
- Use table-driven tests where appropriate
- Test edge cases and error conditions
- Mock external HTTP calls in tests
- Run tests with `go test -v`
- Ensure tests are deterministic and isolated

## Test Coverage

- Focus on critical functionality: feed parsing, formatting, and display logic
- Test error handling paths
- Validate RSS XML parsing correctness
- Test relative time formatting (e.g., "2 days ago")

# Development Workflow

## Building

```bash
go build -o gh-changelog
```

## Testing

```bash
go test -v
```

## Running Locally

```bash
./gh-changelog
./gh-changelog --pretty
./gh-changelog view <id>
./gh-changelog open <id>
```

# CLI Interface Guidelines

- Use the flag package for command-line argument parsing
- Provide clear error messages to stderr
- Use consistent formatting for output
- Support both list and detail views
- Implement --pretty flag for enhanced formatting
- Exit with appropriate exit codes (0 for success, 1 for errors)

# External Dependencies

- Minimize external dependencies
- Prefer standard library solutions
- Only add dependencies for significant functionality gaps
- Document rationale for any new dependencies added

# RSS Feed Processing

- Parse RSS from https://github.blog/changelog/feed/
- Handle XML parsing errors gracefully
- Display items in reverse chronological order (newest first)
- Support both brief and detailed views
- Convert HTML content to readable plain text for terminal display

# Error Handling

- Always check and handle errors from HTTP requests
- Provide meaningful error messages to users
- Use fmt.Fprintf(os.Stderr, ...) for error output
- Exit with non-zero status codes on failures
- Handle malformed RSS feeds gracefully

# Documentation

- Keep README.md up to date with usage examples
- Document all public functions with comments
- Include examples in documentation
- Keep installation and usage instructions clear and concise
