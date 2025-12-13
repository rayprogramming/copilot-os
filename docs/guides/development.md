---
layout: default
title: Development Guide
nav_order: 3
parent: Guides
---

# Development Guide

This guide covers local setup, building, testing, and contributing to the CoPilot OS repository.

## Prerequisites

- Go 1.21 or later
- GitHub Copilot CLI (latest version)
- Git
- Make (optional, for using Makefile)

## Local Setup

### 1. Clone the Repository

```bash
git clone https://github.com/rayprogramming/copilot-agent-chain.git
cd copilot-agent-chain
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Authenticate with Copilot CLI

```bash
copilot auth login
copilot auth status
```

### 4. Verify Setup

```bash
go run ./cmd/server --help
```

## Building

### Development Build

```bash
go build -o copilot-agent-chain ./cmd/server
```

### Optimized Build

```bash
go build -ldflags="-s -w" -o copilot-agent-chain ./cmd/server
```

## Running

### Run Local Server

```bash
# Without pointing to specific repo (uses current directory)
./copilot-agent-chain

# Point to specific repository
REPO_ROOT=/path/to/repo ./copilot-agent-chain

# With debug logging
LOG_LEVEL=debug ./copilot-agent-chain
```

### Test Against This Repository

```bash
# Use this repo's agents
export REPO_ROOT=$(pwd)
./copilot-agent-chain
```

In another terminal:

```bash
copilot --agent=code-reviewer --prompt "Review the main.go file"
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Generate Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Test

```bash
go test -run TestPromptEvaluator ./internal/prompt/...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Benchmarks

```bash
go test -bench=. -benchmem ./...
```

## Code Quality

### Linting

```bash
# Using golangci-lint (recommended)
golangci-lint run

# Or individual tools
go vet ./...
go fmt ./...
```

### Format Code

```bash
go fmt ./...
gofmt -s -w .
```

### Vet Code

```bash
go vet ./...
```

## Using Development Agents

This repository includes development agents in `.github/agents/` to help with code review, architecture guidance, testing, and documentation.

### Code Review

```bash
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go for correctness and performance"
```

### Architecture Advice

```bash
copilot --agent=architecture-advisor --prompt "Should we refactor the agent discovery to support remote agent repositories?"
```

### Generate Tests

```bash
copilot --agent=test-generator --prompt "Generate comprehensive tests for the prompt evaluator"
```

### Improve Documentation

```bash
copilot --agent=documentation-writer --prompt "Write a detailed guide on how the agent chain execution works"
```

### Full Review (Auto-Chaining)

```bash
copilot --agent=orchestrator --prompt "Perform a comprehensive review of the new agent selection algorithm including architecture, testing, and documentation"
```

## Contributing

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

Follow Go best practices:
- Write clear, self-documenting code
- Add appropriate comments for non-obvious logic
- Keep functions small and focused
- Handle errors explicitly

### 3. Test Your Changes

```bash
# Run tests
go test ./...

# Run linters
golangci-lint run

# Format code
go fmt ./...
```

### 4. Use Development Agents for Review

Before submitting a PR, use the agents to review your changes:

```bash
# Code quality review
copilot --agent=code-reviewer --prompt "Review my changes in internal/orchestrator/ for correctness and performance"

# Architecture review
copilot --agent=architecture-advisor --prompt "Does the new agent selection strategy align with the overall architecture?"

# Test coverage
copilot --agent=test-generator --prompt "Identify untested code paths in my implementation"

# Documentation
copilot --agent=documentation-writer --prompt "Improve comments and documentation for my changes"
```

### 5. Commit Changes

```bash
git add .
git commit -m "feat: description of changes"
```

Follow conventional commit format: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `perf:`, etc.

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Create a pull request on GitHub with a clear description of changes.

## Project Structure Quick Reference

```
cmd/server/                 # Server entry point
internal/
  ├── orchestrator/         # Orchestration logic
  ├── agents/               # Agent management
  ├── prompt/               # Prompt evaluation
  ├── cli/                  # Copilot CLI integration
  └── config/               # Configuration
tests/                      # Unit & integration tests
examples/                   # Usage examples
.github/agents/             # Development agent definitions
```

## Common Tasks

### Add a New Package

```bash
# Create package directory
mkdir -p internal/newpackage

# Add types.go for type definitions
# Add implementations.go for logic
# Add implementations_test.go for tests
```

### Add a New Tool to MCP Server

See `internal/orchestrator/orchestrator.go` for examples of registering MCP tools.

### Debug Agent Discovery

```bash
LOG_LEVEL=debug ./copilot-agent-chain
```

Check logs for agent discovery details.

### Test Agent Chaining

```bash
# Create a test prompt that should trigger multiple agents
copilot --agent=orchestrator --prompt "Design and implement tests for a new caching layer that improves performance"
```

## Troubleshooting

### Copilot CLI Not Found

```bash
# Verify installation
which copilot
copilot --version

# Add to PATH if needed
export PATH=$PATH:/path/to/copilot
```

### Server Fails to Start

```bash
# Check logs
LOG_LEVEL=debug ./copilot-agent-chain

# Verify go.mod and dependencies
go mod tidy
go mod verify
```

### Agent Not Discovered

```bash
# Check REPO_ROOT is set correctly
echo $REPO_ROOT

# Verify agent files exist
ls -la $REPO_ROOT/.github/agents/

# Check file format (should be .md with YAML frontmatter)
head -20 $REPO_ROOT/.github/agents/code-reviewer.md
```

### Tests Failing

```bash
# Run with verbose output
go test -v ./...

# Check for race conditions
go test -race ./...

# Look for specific test failure
go test -run TestName -v ./path/to/package
```

## Performance Optimization

### Profile the Server

```bash
# CPU profiling
go run -cpuprofile=cpu.prof ./cmd/server &
# Make some requests
# Kill the server
go tool pprof cpu.prof

# Memory profiling
go run -memprofile=mem.prof ./cmd/server &
# Make some requests
# Kill the server
go tool pprof mem.prof
```

### Benchmark Agent Discovery

```bash
go test -bench=BenchmarkAgentDiscovery -benchmem ./internal/agents/...
```

## Release Process

### Create a Release

```bash
# Tag a version
git tag -a v1.2.3 -m "Release version 1.2.3"
git push origin v1.2.3

# Build release binaries
go build -ldflags="-s -w" -o copilot-agent-chain ./cmd/server
```

## Documentation

### Update Documentation When:
- Changing public APIs
- Adding new features
- Fixing bugs that required investigation
- Changing architecture or design decisions

Use the `documentation-writer` agent:

```bash
copilot --agent=documentation-writer --prompt "Update the documentation with information about the new feature"
```

## Questions or Issues?

Open an issue on GitHub with:
1. Description of the problem
2. Steps to reproduce
3. Expected vs actual behavior
4. Relevant logs or error messages
5. Your environment (Go version, OS, etc.)

---

**Previous:** [Getting Started](getting-started.md) | **Next:** [Architecture Guide](architecture.md)
