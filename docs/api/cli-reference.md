---
layout: default
title: CLI Reference
---

# CLI Reference

This guide covers all command-line options and usage patterns for Copilot Agent Chain.

## Server Commands

### Start Server

```bash
./copilot-agent-chain
```

Starts the MCP server with default configuration.

**Output**:
```
2025-12-07T10:30:45.123Z	info	config	Loaded configuration	{"repo_root": ".", "log_level": "info"}
2025-12-07T10:30:45.234Z	info	discovery	Agent discovery started	{"path": "./.github/agents"}
2025-12-07T10:30:45.345Z	info	discovery	Discovered agents	{"count": 4, "agents": [...]}
2025-12-07T10:30:45.456Z	info	server	MCP server started	{"transport": "stdio"}
```

### Server with Custom Repository

```bash
export REPO_ROOT=/path/to/repo
./copilot-agent-chain
```

Scans for agents in `/path/to/repo/.github/agents/`.

### Server with Debug Logging

```bash
export LOG_LEVEL=debug
./copilot-agent-chain
```

Enables debug-level logging for troubleshooting.

### Server with Custom Timeout

```bash
export AGENT_TIMEOUT=60s
./copilot-agent-chain
```

Sets individual agent execution timeout to 60 seconds.

## MCP Tool Invocation

Once the server is running, use Copilot CLI to invoke tools:

### run_with_orchestrator

**Usage**:
```bash
copilot --agent=orchestrator --prompt "<your-prompt>"
```

**Options**:
- `--prompt` (required) — The request to orchestrate

**Examples**:
```bash
# Simple code review
copilot --agent=orchestrator --prompt "Review main.go"

# Complex request
copilot --agent=orchestrator --prompt "Review the authentication module for security, ensure good test coverage, and write documentation"

# Architecture discussion
copilot --agent=orchestrator --prompt "Should we add caching? What are the implications?"
```

### run_agent

**Usage**:
```bash
copilot --agent=<agent-name> --prompt "<your-prompt>"
```

**Examples**:
```bash
# Use specific agent
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go"

copilot --agent=architecture-advisor --prompt "Design a caching layer"

copilot --agent=test-generator --prompt "Generate tests for the evaluator"

copilot --agent=documentation-writer --prompt "Improve README documentation"
```

### list_agents

**Usage**:
```bash
copilot --agent=orchestrator --prompt "List all available agents"
```

**Output**:
```
Available agents:
1. code-reviewer - Specialized Go code reviewer
   Keywords: code-review, go, quality, testing, correctness, performance

2. architecture-advisor - System architecture specialist
   Keywords: architecture, design, pattern, structure, scale, integration

3. test-generator - Testing specialist
   Keywords: test, testing, coverage, unit-test, integration-test, mock, benchmark

4. documentation-writer - Documentation specialist
   Keywords: documentation, readme, guide, comment, explain, write, document
```

## Prompt Best Practices

### Make Prompts Specific

```bash
# Good
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go for correctness and performance bottlenecks"

# Poor
copilot --agent=code-reviewer --prompt "Review the code"
```

### Provide Context

```bash
# Good
copilot --agent=test-generator --prompt "Generate comprehensive unit tests for the prompt evaluator, focusing on vagueness detection heuristics and edge cases"

# Poor
copilot --agent=test-generator --prompt "Write tests"
```

### Ask Clarifying Questions

```bash
# Good
copilot --agent=architecture-advisor --prompt "Should we implement result caching? What are the architectural trade-offs, performance implications, and recommended approach?"

# Poor
copilot --agent=architecture-advisor --prompt "How do we cache results?"
```

### Use Multi-Agent for Complex Tasks

```bash
# Good - leverage orchestrator for multi-perspective review
copilot --agent=orchestrator --prompt "Comprehensively review the agent selection algorithm including code quality, architectural implications, test coverage, and documentation needs"

# Less ideal - would need to run agents separately
copilot --agent=code-reviewer --prompt "Review the agent selection algorithm"
```

## Environment Variables

### REPO_ROOT

**Default**: Current directory (`.`)

**Usage**:
```bash
export REPO_ROOT=/home/user/my-project
./copilot-agent-chain
```

**Purpose**: Location of repository containing `.github/agents/` directory.

### LOG_LEVEL

**Default**: `info`

**Valid Values**: `debug`, `info`, `warn`, `error`

**Usage**:
```bash
export LOG_LEVEL=debug
./copilot-agent-chain
```

**Purpose**: Control logging verbosity for troubleshooting.

### AGENT_TIMEOUT

**Default**: `30s`

**Format**: Duration (e.g., `30s`, `1m`, `90s`)

**Usage**:
```bash
export AGENT_TIMEOUT=60s
./copilot-agent-chain
```

**Purpose**: Maximum time to wait for agent execution.

### MCP_TRANSPORT

**Default**: `stdio`

**Valid Values**: `stdio`, `http` (pending)

**Usage**:
```bash
export MCP_TRANSPORT=stdio
./copilot-agent-chain
```

**Purpose**: Transport protocol for MCP communication.

## Common Workflows

### Code Review Workflow

1. Start server:
   ```bash
   export REPO_ROOT=$(pwd)
   ./copilot-agent-chain
   ```

2. Run code review in another terminal:
   ```bash
   copilot --agent=code-reviewer --prompt "Review internal/agents/discovery.go for correctness and efficiency"
   ```

3. Get architecture feedback:
   ```bash
   copilot --agent=architecture-advisor --prompt "Is the current agent discovery design scalable?"
   ```

4. Generate tests:
   ```bash
   copilot --agent=test-generator --prompt "Generate comprehensive tests for the discovery process"
   ```

### Feature Development Workflow

1. Get architecture review:
   ```bash
   copilot --agent=architecture-advisor --prompt "Design a new feature to support parallel agent execution"
   ```

2. Implement feature
3. Review implementation:
   ```bash
   copilot --agent=code-reviewer --prompt "Review my implementation in internal/orchestrator/chain.go"
   ```

4. Generate tests:
   ```bash
   copilot --agent=test-generator --prompt "Generate tests for the parallel execution implementation"
   ```

5. Document feature:
   ```bash
   copilot --agent=documentation-writer --prompt "Update documentation for the new parallel execution feature"
   ```

### Complete Module Review Workflow

```bash
# Run comprehensive review
copilot --agent=orchestrator --prompt "Comprehensively review the agent discovery module including code quality, architecture, test coverage, and documentation. Suggest improvements in all areas."
```

This single command chains all agents for a complete review.

## Troubleshooting

### Server Won't Start

**Problem**: "Cannot find .github/agents directory"

**Solution**:
```bash
# Verify REPO_ROOT is set correctly
echo $REPO_ROOT

# Check directory exists
ls -la $REPO_ROOT/.github/agents/

# Or explicitly set it
export REPO_ROOT=$(pwd)
./copilot-agent-chain
```

### Copilot CLI Not Found

**Problem**: "copilot: command not found"

**Solution**:
```bash
# Verify installation
which copilot
copilot --version

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin
```

### Agent Timeout

**Problem**: "Agent execution timeout"

**Solution**:
```bash
# Increase timeout
export AGENT_TIMEOUT=120s
./copilot-agent-chain

# Or try again with smaller prompt
copilot --agent=code-reviewer --prompt "Review foo.go"
```

### No Output from Agent

**Problem**: Agent runs but produces no output

**Solution**:
```bash
# Check with verbose logging
export LOG_LEVEL=debug
./copilot-agent-chain

# Try with a simpler prompt
copilot --agent=code-reviewer --prompt "What is this file about?"
```

### Authentication Issues

**Problem**: "Not authenticated" or "Invalid credentials"

**Solution**:
```bash
# Re-authenticate
copilot auth logout
copilot auth login
copilot auth status

# Restart server
export REPO_ROOT=$(pwd)
./copilot-agent-chain
```

## Advanced Usage

### Chaining Custom Agents

If you've created custom agents in `.github/agents/`, use them the same way:

```bash
copilot --agent=your-custom-agent --prompt "Your request here"
```

### Debugging Agent Selection

View which agents the orchestrator would select:

```bash
copilot --agent=orchestrator --prompt "Debug: What agents would be selected for: 'Review authentication module for security'?"
```

### Performance Profiling

Run with debug logging to see execution times:

```bash
export LOG_LEVEL=debug
./copilot-agent-chain &

# Make requests
copilot --agent=orchestrator --prompt "Review internal/prompt/evaluator.go"

# Check logs for timing information
```

---

**Previous:** [API Documentation](index.md) | **Next:** [Configuration Reference](configuration.md)
