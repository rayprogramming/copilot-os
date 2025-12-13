---
layout: home
title: Copilot Agent Chain
nav_order: 1
---

# Copilot Agent Chain Documentation

Welcome to the **Copilot Agent Chain** documentation! This site provides comprehensive guidance for using, developing, and extending the intelligent agent orchestration system.

## What is Copilot Agent Chain?

**Copilot Agent Chain** is a Go-based MCP (Model Context Protocol) server that orchestrates intelligent agent chains using GitHub Copilot CLI. The server automatically evaluates user prompts, selects optimal agents from your codebase's `.github/agents/` directory, chains them together with context flow, and returns synthesized results.

### Key Features

- **🤖 Intelligent Orchestration** — Built-in Go orchestrator evaluates prompts, auto-refines unclear requests, and intelligently chains agents
- **🔍 Agent Discovery** — Automatically loads agents from your repository's `.github/agents/` directory
- **🎯 Smart Agent Selection** — Keyword-based capability matching to select optimal agents for each task
- **🔗 Context Flow** — Accumulates results as JSON, passes rich context between agents
- **📦 MCP Integration** — Full Model Context Protocol support via hypermcp framework
- **🛠️ Development Agents** — Includes Code Reviewer, Architecture Advisor, Test Generator, and Documentation Writer agents
- **⚙️ Zero Configuration** — Works with any repository that has `.github/agents/` agent definitions

## Getting Started

Start with the **[Getting Started Guide](guides/getting-started.md)** to:
- Install Copilot Agent Chain
- Configure your repository
- Run your first agent chain

## Documentation Structure

### 📖 Guides

- **[Getting Started](guides/getting-started.md)** — Installation, setup, and quick start
- **[Development Guide](guides/development.md)** — Local development, building, testing, and debugging
- **[Architecture Guide](guides/architecture.md)** — System design, component interactions, and extending the system
- **[Agent Guide](guides/agents.md)** — Understanding and using the built-in development agents

### 📚 API Reference

- **[API Documentation](api/index.md)** — Complete reference for MCP tools and interfaces
- **[CLI Reference](api/cli-reference.md)** — Command-line usage and options
- **[Configuration Reference](api/configuration.md)** — Environment variables and settings

### 💡 Examples

- **[Examples & Scenarios](examples/index.md)** — Real-world usage examples and integration patterns

### 🔧 Additional Resources

- **[Troubleshooting](guides/troubleshooting.md)** — Common issues and solutions
- **[Contributing](guides/contributing.md)** — How to contribute to Copilot Agent Chain
- **[Project Status](status/implementation.md)** — Project completion status and test coverage

## System Architecture

```
┌─────────────────────────────────────┐
│    Copilot CLI / External Caller    │
└─────────────────┬───────────────────┘
                  │ (MCP Protocol)
┌─────────────────▼───────────────────┐
│      MCP Server (hypermcp)          │
│  ┌─────────────────────────────┐   │
│  │   Orchestrator (Go)         │   │
│  │ • Prompt Evaluator          │   │
│  │ • Agent Selection           │   │
│  │ • Chain Executor            │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
         ▼      ▼      ▼      ▼
      Agent   Agent   Agent   Custom
      (Code)  (Arch)  (Test)  (Agents)
```

## Quick Example

```bash
# Point to your repository
export REPO_ROOT=/path/to/your/repo

# Build the server
go build -o copilot-agent-chain ./cmd/server

# Run the server
./copilot-agent-chain

# Use it with Copilot CLI
copilot --agent=orchestrator --prompt "Review the authentication module for security"
```

## Key Concepts

### Agents
Specialized tools defined in `.github/agents/` that perform specific tasks (code review, architecture advice, test generation, documentation).

### Orchestrator
The intelligent core that evaluates prompts, selects appropriate agents, chains them together, and synthesizes results.

### Context Flow
Results are accumulated as JSON and passed between agents, enabling sophisticated multi-step workflows.

### Smart Selection
Agents are matched to tasks using keyword-based capability matching with intelligent scoring.

## Requirements

- Go 1.21 or later
- GitHub Copilot CLI installed and authenticated
- A repository with agent definitions in `.github/agents/`

## Contributing

We welcome contributions! See the **[Contributing Guide](guides/contributing.md)** for details on how to:
- Report issues
- Propose features
- Submit pull requests
- Set up your development environment

## License

Copilot Agent Chain is licensed under the MIT License. See the LICENSE file in the repository for details.

## Support

- 📖 Check the [Troubleshooting Guide](guides/troubleshooting.md)
- 🐛 Report issues on [GitHub](https://github.com/rayprogramming/copilot-agent-chain/issues)
- 💬 Start a discussion on [GitHub Discussions](https://github.com/rayprogramming/copilot-agent-chain/discussions)

---

**Last Updated:** December 7, 2025  
**Version:** 1.0.0  
**Status:** ✅ Production Ready
