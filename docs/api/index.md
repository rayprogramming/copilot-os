---
layout: default
title: API Reference
nav_order: 3
---

# API Reference

This section provides complete documentation for the CoPilot OS MCP server's tools, interfaces, and APIs.

## MCP Tools

The server exposes the following MCP tools for use via the Copilot CLI or other MCP clients:

### 1. run_with_orchestrator

**Purpose**: Run the intelligent orchestrator to evaluate a prompt, select agents, and chain them together.

**Parameters**:
- `prompt` (string, required) — The user's request to orchestrate

**Returns**:
```json
{
  "originalPrompt": "user's initial request",
  "refinedPrompt": "orchestrator's refined version (if refined)",
  "evaluationFeedback": {
    "isClear": true|false,
    "confidence": 0.0-1.0,
    "issues": ["issue1", "issue2"],
    "suggestions": "refinement suggestions if needed"
  },
  "selectedAgents": ["agent-name1", "agent-name2"],
  "agentResults": [
    {
      "agent": "agent-name",
      "output": {...},
      "executedAt": "2025-12-07T10:30:00Z"
    }
  ],
  "finalOutput": "synthesized result",
  "executionTime": "5.234s"
}
```

**Example**:
```bash
copilot --agent=orchestrator --prompt "Review the authentication module for security issues"
```

### 2. list_agents

**Purpose**: List all available agents and their capabilities.

**Parameters**: None

**Returns**:
```json
{
  "agents": [
    {
      "name": "code-reviewer",
      "description": "Specialized Go code reviewer",
      "keywords": ["code-review", "go", "quality", "testing", "correctness", "performance"]
    },
    ...
  ],
  "count": 4
}
```

**Example**:
```bash
copilot --agent=orchestrator --prompt "List all available agents"
```

### 3. run_agent

**Purpose**: Run a specific agent directly (bypass orchestrator).

**Parameters**:
- `agentName` (string, required) — Name of the agent to run
- `prompt` (string, required) — The prompt to send to the agent

**Returns**:
```json
{
  "agent": "agent-name",
  "prompt": "user prompt",
  "output": {...},
  "executedAt": "2025-12-07T10:30:00Z",
  "executionTime": "3.456s"
}
```

**Example**:
```bash
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go for performance"
```

### 4. evaluate_prompt

**Purpose**: Evaluate a prompt for clarity and get refinement suggestions (for debugging).

**Parameters**:
- `prompt` (string, required) — The prompt to evaluate

**Returns**:
```json
{
  "originalPrompt": "user prompt",
  "isClear": true|false,
  "confidence": 0.0-1.0,
  "issues": ["vagueness", "missing-context"],
  "refinedPrompt": "improved version if needed",
  "keywords": ["extracted", "keywords"],
  "selectedAgents": ["best-match-agents"]
}
```

**Example**:
```bash
copilot --agent=orchestrator --prompt "Debug this tool: evaluate_prompt('Review the module')"
```

## Go API

If you're extending the server, here are the key Go interfaces:

### Orchestrator Interface

```go
type Orchestrator interface {
    // RunWithAuto evaluates the prompt and automatically chains agents
    RunWithAuto(ctx context.Context, prompt string) (*OrchestrationResult, error)
    
    // RunWithExplicitChain runs a specific chain of agents
    RunWithExplicitChain(ctx context.Context, prompt string, agents []string) (*OrchestrationResult, error)
}
```

### Agent Registry Interface

```go
type Registry interface {
    // Add registers an agent
    Add(agent *Agent) error
    
    // Get retrieves an agent by name
    Get(name string) (*Agent, error)
    
    // All returns all agents
    All() []*Agent
    
    // MatchKeywords finds agents matching keywords
    MatchKeywords(keywords []string, maxResults int) []*Agent
}
```

### Prompt Evaluator Interface

```go
type Evaluator interface {
    // Evaluate analyzes prompt clarity
    Evaluate(prompt string) *EvaluationResult
    
    // ExtractKeywords extracts domain keywords
    ExtractKeywords(prompt string) []string
    
    // SuggestRefinement generates a refined version
    SuggestRefinement(prompt string) string
}
```

### CLI Invoker Interface

```go
type Invoker interface {
    // InvokeAgent runs an agent via Copilot CLI
    InvokeAgent(ctx context.Context, agentName, prompt string) (*AgentResult, error)
    
    // IsAvailable checks if Copilot CLI is available
    IsAvailable() bool
    
    // CheckAuth verifies Copilot CLI authentication
    CheckAuth() error
}
```

## Configuration

Configure the server via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REPO_ROOT` | Current directory | Path to repository with `.github/agents/` |
| `LOG_LEVEL` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `AGENT_TIMEOUT` | `30s` | Timeout for individual agent execution |
| `MCP_TRANSPORT` | `stdio` | MCP transport: `stdio` (HTTP support pending) |
| `CACHE_SIZE` | `1000` | Size of agent result cache (entries) |
| `CACHE_TTL` | `1h` | Time-to-live for cached results |

## Response Formats

### Successful Response

All tools return responses in this format:

```json
{
  "status": "success",
  "data": {... tool-specific data ...},
  "executionTime": "3.456s",
  "timestamp": "2025-12-07T10:30:00Z"
}
```

### Error Response

When an error occurs:

```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional context"
  },
  "timestamp": "2025-12-07T10:30:00Z"
}
```

**Common Error Codes**:
- `AGENT_NOT_FOUND` — Requested agent doesn't exist
- `CLI_NOT_AVAILABLE` — Copilot CLI not installed or authenticated
- `EXECUTION_TIMEOUT` — Agent execution exceeded timeout
- `INVALID_PROMPT` — Prompt validation failed
- `ORCHESTRATION_FAILED` — Orchestration process failed
- `UNKNOWN_ERROR` — Unexpected error

## Data Types

### Agent

```json
{
  "name": "code-reviewer",
  "description": "Specialized Go code reviewer",
  "keywords": ["code-review", "go", "quality"],
  "path": "/path/to/.github/agents/code-reviewer.md"
}
```

### EvaluationResult

```json
{
  "isClear": true,
  "confidence": 0.85,
  "issues": [],
  "suggestions": null,
  "keywords": ["review", "code", "quality"]
}
```

### AgentResult

```json
{
  "agent": "code-reviewer",
  "output": {...},
  "executedAt": "2025-12-07T10:30:00Z",
  "executionTime": "3.456s"
}
```

## Rate Limiting

Currently, there are no built-in rate limits. Future versions will include:
- Per-agent rate limits
- Global request rate limiting
- Priority queuing

## Authentication

Authentication is handled through the user's Copilot CLI configuration. Ensure:

1. Copilot CLI is installed
2. User is authenticated: `copilot auth login`
3. GitHub Copilot subscription is active

## Versioning

The API follows semantic versioning:
- **Major** — Breaking changes (new major version)
- **Minor** — New features (backwards compatible)
- **Patch** — Bug fixes

Current version: **1.0.0**

---

**Next:** [CLI Reference](cli-reference.md)
