# CopilotOS

**Copilot Orchestration System** - A Go-based MCP (Model Context Protocol) server that orchestrates intelligent agent chains using GitHub Copilot CLI. The server automatically evaluates user prompts, selects optimal agents from your codebase's `.github/agents/` directory, chains them together with context flow, and returns synthesized results.

**📖 Full Documentation**: See the `docs/` folder for comprehensive guides

## Features

- **Intelligent Orchestration**: Built-in Go orchestrator evaluates prompts, auto-refines unclear requests, and intelligently chains agents
- **Agent Discovery**: Automatically loads agents from your repository's `.github/agents/` directory
- **Smart Agent Selection**: Keyword-based capability matching to select optimal agents for each task
- **Context Flow**: Accumulates results as JSON, passes rich context between agents
- **MCP Integration**: Full Model Context Protocol support via hypermcp framework
- **Development Agents**: Includes Code Reviewer, Architecture Advisor, Test Generator, and Documentation Writer agents for repository development
- **Zero Configuration**: Works with any repository that has `.github/agents/` agent definitions

## Quick Start

For detailed instructions, see the **[Getting Started Guide](docs/guides/getting-started.md)**.

### Prerequisites

- Go 1.21 or later
- GitHub Copilot CLI installed and authenticated
- A repository with agent definitions in `.github/agents/`

### Installation

```bash
git clone https://github.com/rayprogramming/copilot-os.git
cd copilot-os
go build -o copilot-os ./cmd/server
```

### Running the Server

```bash
# Run from your repository root (agents discovered from ./.github/agents/)
./copilot-os

# Or point to a specific repository
export REPO_ROOT=/path/to/your/repo
./copilot-os
```

The server will start listening on stdio and automatically discover agents from the repository's `.github/agents/` directory.

### Using from Copilot CLI

```bash
# Run with orchestrator (default, auto-chains agents)
copilot --agent=orchestrator --prompt "Review the authentication module for security issues"

# Run a specific agent
copilot --agent=code-reviewer --prompt "Check the CLI invocation logic"

# List available agents
copilot --agent=orchestrator --prompt "List all available agents"
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│          Copilot CLI / External Caller              │
└─────────────────┬───────────────────────────────────┘
                  │ (MCP Protocol)
┌─────────────────▼───────────────────────────────────┐
│         MCP Server (hypermcp)                       │
│  ┌──────────────────────────────────────────────┐  │
│  │  Orchestrator (Go Component)                 │  │
│  │  ┌────────────────────────────────────────┐  │  │
│  │  │ 1. Prompt Evaluator (heuristics)       │  │  │
│  │  │    - Detect clarity issues             │  │  │
│  │  │    - Auto-refine if needed             │  │  │
│  │  └────────────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────────────┐  │  │
│  │  │ 2. Agent Discovery & Selection         │  │  │
│  │  │    - Scan .github/agents/*.md          │  │  │
│  │  │    - Keyword matching for selection    │  │  │
│  │  └────────────────────────────────────────┘  │  │
│  │  ┌────────────────────────────────────────┐  │  │
│  │  │ 3. Agent Chain Executor                │  │  │
│  │  │    - Invoke agents via Copilot CLI     │  │  │
│  │  │    - Accumulate context as JSON        │  │  │
│  │  │    - Adapt chain based on results      │  │  │
│  │  └────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                  │
      ┌───────────┴──────────────┬──────────────┐
      │                          │              │
   [Agent 1]              [Agent 2]        [Agent N]
 (Code Reviewer)      (Architecture)     (Custom Agent)
```

## Project Structure

```
copilot-os/
├── cmd/
│   └── server/
│       └── main.go              # Server entry point
├── internal/
│   ├── orchestrator/            # Orchestration logic
│   │   ├── orchestrator.go      # Main orchestrator component
│   │   └── chain.go             # Agent chain execution
│   ├── agents/                  # Agent management
│   │   ├── discovery.go         # Agent discovery & registry
│   │   └── types.go             # Agent metadata types
│   ├── prompt/                  # Prompt evaluation
│   │   └── evaluator.go         # Clarity & auto-refinement
│   ├── cli/                     # Copilot CLI invocation
│   │   └── invoker.go           # Subprocess management
│   └── config/
│       └── config.go            # Configuration handling
├── .github/
│   └── agents/                  # Development agents
│       ├── code-reviewer.md
│       ├── architecture-advisor.md
│       ├── test-generator.md
│       └── documentation-writer.md
├── tests/                       # Unit & integration tests
├── examples/                    # Usage examples
├── DEVELOPMENT.md               # Local setup & contributing
├── ARCHITECTURE.md              # System design details
├── AGENTS.md                    # Agent documentation
└── go.mod                       # Go module definition
```

## Documentation

Complete documentation is available in the `docs/` folder and at **[https://rayprogramming.github.io/copilot-os](https://rayprogramming.github.io/copilot-os)**.

### Quick Links

- **[Getting Started](docs/guides/getting-started.md)** — Installation, setup, and first use
- **[Development Guide](docs/guides/development.md)** — Local setup, building, testing, and contributing
- **[Architecture Guide](docs/guides/architecture.md)** — System design and design decisions  
- **[Agent Guide](docs/guides/agents.md)** — Understanding and using development agents
- **[API Reference](docs/api/index.md)** — Complete MCP tools and interface documentation
- **[Examples](docs/examples/index.md)** — Real-world usage scenarios
- **[Troubleshooting](docs/guides/troubleshooting.md)** — Common issues and solutions
- **[Contributing](docs/guides/contributing.md)** — How to contribute to the project

## How It Works

### 1. Prompt Evaluation

When a prompt is received, the orchestrator evaluates clarity:

```
Input: "Review the auth module"
Evaluation:
  - Detected vagueness: "module" is non-specific
  - Missing context: No mention of security concerns or performance
  - Suggested refinement: "Review the authentication module for security vulnerabilities and performance bottlenecks"
```

### 2. Agent Selection

The orchestrator matches refined prompt against available agents:

```
Refined Prompt Keywords: [authentication, security, performance]
Available Agents:
  - code-reviewer (keywords: code-review, go, quality, testing, correctness, performance)
  - test-generator (keywords: testing, test-generation, unit-tests, coverage)
  - architecture-advisor (keywords: architecture, design, patterns, modularity)

Selected Chain: [code-reviewer, test-generator]
Rationale: code-reviewer matches on security/performance, test-generator for validation
```

### 3. Agent Chain Execution

Agents execute sequentially with context flow:

```
Agent 1 (code-reviewer):
  Input: {
    originalPrompt: "Review the auth module",
    refinedPrompt: "Review authentication module for security and performance",
    currentContext: {}
  }
  Output: {
    issues: [...],
    recommendations: [...],
    riskLevel: "medium"
  }

Agent 2 (test-generator):
  Input: {
    originalPrompt: "Review the auth module",
    refinedPrompt: "...",
    previousAgentOutput: { issues, recommendations, riskLevel },
    currentContext: { focusAreas: ["identified_issues"] }
  }
  Output: {
    testCases: [...],
    coverage: "85%",
    gaps: [...]
  }

Final Synthesis:
  {
    originalPrompt: "Review the auth module",
    refinedPrompt: "...",
    agentResults: [ code-reviewer output, test-generator output ],
    finalOutput: "Comprehensive summary integrating all findings"
  }
```

## Usage Examples

### Example 1: Full Review with Auto-Chaining

```bash
copilot --agent=orchestrator --prompt "Check the prompt evaluator for quality"
```

The orchestrator will:
1. Evaluate the prompt (likely clear enough)
2. Auto-select: `[code-reviewer, test-generator]`
3. Chain them with context
4. Return comprehensive feedback

### Example 2: Specific Agent

```bash
copilot --agent=code-reviewer --prompt "Review internal/cli/invoker.go for subprocess safety"
```

### Example 3: Complex Task

```bash
copilot --agent=orchestrator --prompt "Design a resilient retry mechanism for Copilot CLI calls that fails gracefully when agents are unavailable"
```

The orchestrator will:
1. Detect need for architecture focus
2. Auto-select: `[architecture-advisor, test-generator, documentation-writer]`
3. Chain them to produce a complete design with examples and docs

## Configuration

Configuration via environment variables:

- `REPO_ROOT` — Path to repository with agents (default: current directory). When using VS Code MCP integration, use `${workspaceFolder}` to automatically reference the current workspace.
- `LOG_LEVEL` — Logging level: debug, info, warn, error (default: info)
- `CACHE_ENABLED` — Enable result caching (default: true)
- `COPILOT_CLI_TIMEOUT` — Timeout for Copilot CLI calls in seconds (default: 300)

## Dependencies

- `github.com/rayprogramming/hypermcp` — MCP server framework
- `go.uber.org/zap` — Structured logging
- Standard Go library for CLI invocation and file I/O

## API Reference

### MCP Tools

#### `run_with_orchestrator`

Runs the intelligent orchestrator on a prompt.

**Parameters:**
- `prompt` (string): User request/task
- `agentChain` ([]string, optional): Explicit agent chain. If omitted, orchestrator auto-selects.

**Response:**
- JSON object with: `originalPrompt`, `refinedPrompt`, `agentResults[]`, `finalOutput`

#### `list_agents`

Lists all available agents discovered from repository.

**Parameters:** None

**Response:**
- Array of agent objects with: `name`, `description`, `keywords`

#### `evaluate_prompt`

Evaluates prompt clarity without executing agents.

**Parameters:**
- `prompt` (string): Prompt to evaluate

**Response:**
- Object with: `isClear` (bool), `feedback` (string), `suggestedRefinement` (string)

#### `run_agent`

Runs a specific agent (bypass orchestrator).

**Parameters:**
- `agent` (string): Agent name
- `prompt` (string): Agent prompt
- `context` (object, optional): Previous context

**Response:**
- Agent output as returned by Copilot CLI

## Testing

The project includes comprehensive unit and integration tests with ~90% coverage:

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/agents -v
go test ./internal/prompt -v
go test ./tests -v
```

**Test Coverage**: ~90% across core packages (36 tests, all passing)

For detailed testing information, see the [Testing Guide](docs/guides/development.md#testing).

## Development Agents

This repository includes four specialized development agents for code review, architecture guidance, testing, and documentation:

- **Code Reviewer** — Go code quality and correctness
- **Architecture Advisor** — System design and patterns
- **Test Generator** — Test creation and coverage
- **Documentation Writer** — Documentation and comments

See the [Agent Guide](docs/guides/agents.md) for usage examples.

## Contributing

This repository uses its own development agents for code review, architecture guidance, testing, and documentation. 

**To contribute:**
1. See the [Contributing Guide](docs/guides/contributing.md) for development workflow
2. See the [Development Guide](docs/guides/development.md) for setup and building
3. Use the development agents to review your changes before submitting a PR

## License

MIT

## Support

For issues, questions, or contributions, please open an issue on GitHub.
