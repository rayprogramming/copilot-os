---
layout: default
title: Configuration Reference
---

# Configuration Reference

Complete reference for configuring Copilot Agent Chain.

## Environment Variables

All configuration is done via environment variables. There is no configuration file.

### REPO_ROOT

**Description**: Path to the repository containing `.github/agents/` directory.

**Type**: String (file path)

**Default**: Current working directory (`.`)

**Example**:
```bash
export REPO_ROOT=/home/user/my-project
```

**Notes**:
- Can be absolute or relative path
- Must contain `.github/agents/` directory
- Agents are scanned on server startup
- Use `$(pwd)` to reference current directory

### LOG_LEVEL

**Description**: Logging verbosity level.

**Type**: String

**Default**: `info`

**Valid Values**:
- `debug` — Very verbose, includes all internal operations
- `info` — Standard logging, shows startup and major events
- `warn` — Only warnings and errors
- `error` — Only errors

**Example**:
```bash
export LOG_LEVEL=debug
```

**When to Use Each Level**:
- `debug` — Troubleshooting, understanding flow
- `info` — Production, normal operation
- `warn` — Focus on issues
- `error` — Only when something breaks

### AGENT_TIMEOUT

**Description**: Maximum time to wait for individual agent execution.

**Type**: Duration (Go duration format)

**Default**: `30s`

**Valid Formats**:
- Seconds: `30s`, `45s`, `60s`
- Minutes: `1m`, `2m30s`, `5m`
- Mixed: `1m30s`, `2m15s`

**Example**:
```bash
export AGENT_TIMEOUT=60s
```

**When to Adjust**:
- Increase if agents regularly timeout
- Decrease to fail faster on hung agents
- Typical range: 20s-120s

**Note**: Each agent execution has this timeout independently.

### MCP_TRANSPORT

**Description**: Protocol for MCP (Model Context Protocol) communication.

**Type**: String

**Default**: `stdio`

**Valid Values**:
- `stdio` — Standard input/output (currently supported)
- `http` — HTTP transport (planned for future)

**Example**:
```bash
export MCP_TRANSPORT=stdio
```

**Notes**:
- `stdio` is best for Copilot CLI integration
- HTTP transport will enable remote agent orchestration
- No configuration needed for `stdio` transport

### CACHE_SIZE

**Description**: Maximum number of agent result entries to cache.

**Type**: Integer

**Default**: `1000`

**Example**:
```bash
export CACHE_SIZE=5000
```

**When to Adjust**:
- Higher values for frequent prompt repetition
- Lower values for memory-constrained environments
- Set to `0` to disable caching

**Notes**:
- Uses LRU (Least Recently Used) eviction
- Memory per entry varies by output size

### CACHE_TTL

**Description**: Time-to-live for cached agent results.

**Type**: Duration (Go duration format)

**Default**: `1h`

**Valid Formats**:
- Minutes: `15m`, `30m`, `45m`
- Hours: `1h`, `2h`, `24h`
- Mixed: `1h30m`, `2h15m`

**Example**:
```bash
export CACHE_TTL=2h
```

**When to Adjust**:
- Shorter TTL if agents' outputs change frequently
- Longer TTL if results are stable
- Set to `0` to disable expiry (keep until evicted by size)

## Complete Configuration Example

### Development Environment

```bash
export REPO_ROOT=$(pwd)
export LOG_LEVEL=debug
export AGENT_TIMEOUT=60s
export MCP_TRANSPORT=stdio
export CACHE_SIZE=100
export CACHE_TTL=30m

./copilot-agent-chain
```

### Production Environment

```bash
export REPO_ROOT=/var/lib/copilot-agent-chain/repo
export LOG_LEVEL=info
export AGENT_TIMEOUT=45s
export MCP_TRANSPORT=stdio
export CACHE_SIZE=5000
export CACHE_TTL=2h

./copilot-agent-chain
```

### High-Performance Environment

```bash
export REPO_ROOT=/opt/repo
export LOG_LEVEL=warn
export AGENT_TIMEOUT=30s
export MCP_TRANSPORT=stdio
export CACHE_SIZE=10000
export CACHE_TTL=4h

./copilot-agent-chain
```

## Agent Configuration

Individual agents are configured via `.github/agents/` markdown files, not via environment variables.

### Agent File Format

Each agent is a markdown file with YAML frontmatter:

```yaml
---
name: agent-name
description: What this agent does
keywords: [keyword1, keyword2, keyword3]
---

## Agent Instructions

Detailed instructions for the agent go here.
```

### Agent Properties

#### name
- **Type**: String
- **Required**: Yes
- **Constraint**: Alphanumeric + hyphens, lowercase
- **Example**: `code-reviewer`

#### description
- **Type**: String
- **Required**: Yes
- **Purpose**: Display in agent listings
- **Example**: "Specialized Go code reviewer"

#### keywords
- **Type**: Array of strings
- **Required**: Yes
- **Purpose**: Used by orchestrator for agent selection
- **Example**: `[code-review, go, quality, testing]`

## Runtime Behavior

### On Startup
1. Load configuration from environment variables
2. Verify REPO_ROOT exists
3. Scan `.github/agents/` for agent definitions
4. Load and validate agent metadata
5. Initialize logger with LOG_LEVEL
6. Start MCP server on stdio transport

### During Execution
- Each agent invocation respects AGENT_TIMEOUT
- Results are cached with CACHE_TTL expiry
- Logs are written at LOG_LEVEL granularity

## Validation

The server validates configuration on startup:

| Property | Validation |
|----------|-----------|
| REPO_ROOT | Must exist and be readable |
| LOG_LEVEL | Must be valid level |
| AGENT_TIMEOUT | Must be positive duration |
| CACHE_SIZE | Must be non-negative integer |
| CACHE_TTL | Must be non-negative duration |

Invalid configuration causes startup failure with helpful error message.

## Performance Tuning

### For Fast Response Times
```bash
export AGENT_TIMEOUT=20s
export LOG_LEVEL=warn
export CACHE_SIZE=10000
export CACHE_TTL=4h
```

### For Detailed Debugging
```bash
export LOG_LEVEL=debug
export AGENT_TIMEOUT=120s
export CACHE_SIZE=100
```

### For Memory Efficiency
```bash
export CACHE_SIZE=100
export CACHE_TTL=15m
export LOG_LEVEL=warn
```

## Security Considerations

### REPO_ROOT Security
- Only set REPO_ROOT to trusted repositories
- The server will scan `.github/agents/` for malicious instructions
- Agent files should be reviewed before execution

### Prompt Sanitization
- User prompts are passed to Copilot CLI
- The server does not execute arbitrary code from prompts
- However, use caution with untrusted user input

### Authentication
- Copilot CLI credentials are used from user's configuration
- The server does not store or transmit credentials
- Authentication is handled entirely by Copilot CLI

## Docker Configuration

If running in Docker, mount the repository and set environment variables:

```bash
docker run -it \
  -v /path/to/repo:/repo \
  -e REPO_ROOT=/repo \
  -e LOG_LEVEL=info \
  -e AGENT_TIMEOUT=60s \
  copilot-agent-chain:latest
```

## Kubernetes Configuration

ConfigMap for Kubernetes deployment:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: copilot-agent-chain-config
data:
  REPO_ROOT: /repo
  LOG_LEVEL: info
  AGENT_TIMEOUT: "60s"
  MCP_TRANSPORT: "stdio"
  CACHE_SIZE: "5000"
  CACHE_TTL: "2h"
```

## Common Configuration Mistakes

### Mistake 1: REPO_ROOT Not Set
```bash
# Wrong - uses current working directory
./copilot-agent-chain

# Right - explicitly set it
export REPO_ROOT=/path/to/repo
./copilot-agent-chain
```

### Mistake 2: Invalid Duration Format
```bash
# Wrong
export AGENT_TIMEOUT=30  # No unit!

# Right
export AGENT_TIMEOUT=30s
```

### Mistake 3: Timeout Too Short
```bash
# May cause timeouts
export AGENT_TIMEOUT=5s

# More reasonable
export AGENT_TIMEOUT=30s
```

### Mistake 4: Excessive Logging
```bash
# May slow down execution
export LOG_LEVEL=debug
export CACHE_SIZE=100

# Better for production
export LOG_LEVEL=info
export CACHE_SIZE=5000
```

---

**Previous:** [CLI Reference](cli-reference.md) | **Next:** [Examples](../examples/index.md)
