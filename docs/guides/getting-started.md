---
layout: default
title: Getting Started
nav_order: 1
parent: Guides
---

# Getting Started Guide

This guide walks you through installing CoPilot OS, setting up your repository, and running your first agent chain.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.21 or later** — Download from [golang.org](https://golang.org/dl/)
- **GitHub Copilot CLI** — Install with `brew install gh-copilot` (macOS) or from [GitHub CLI releases](https://github.com/github/gh-cli)
- **Git** — For version control
- **GitHub Copilot subscription** — Required for using Copilot CLI

### Verify Prerequisites

```bash
# Check Go version
go version

# Check Copilot CLI is installed
copilot --version

# Verify Copilot CLI authentication
copilot auth status
```

## Installation

### Step 1: Clone the Repository

```bash
git clone https://github.com/rayprogramming/copilot-agent-chain.git
cd copilot-agent-chain
```

### Step 2: Install Dependencies

```bash
go mod download
go mod verify
```

### Step 3: Authenticate with Copilot CLI

If you haven't already, authenticate with GitHub Copilot:

```bash
copilot auth login
```

Verify your authentication:

```bash
copilot auth status
```

### Step 4: Build the Server

```bash
go build -o copilot-agent-chain ./cmd/server
```

You should see no errors. The binary `copilot-agent-chain` is now ready to use.

## Configuration

The server is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REPO_ROOT` | Current directory | Path to the repository containing `.github/agents/` |
| `LOG_LEVEL` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `AGENT_TIMEOUT` | `30s` | Timeout for individual agent execution |
| `MCP_TRANSPORT` | `stdio` | MCP transport: `stdio` or `http` |

### Example Configuration

```bash
export REPO_ROOT=/path/to/my/repo
export LOG_LEVEL=debug
export AGENT_TIMEOUT=60s
```

## Running Your First Agent Chain

### Option 1: Use with GitHub Copilot in VS Code

The MCP server integrates directly with GitHub Copilot Chat in VS Code, allowing you to access the agent orchestration capabilities through the Copilot interface.

For detailed instructions on using MCP with GitHub Copilot, see the [official GitHub documentation](https://docs.github.com/en/copilot/how-tos/provide-context/use-mcp/extend-copilot-chat-with-mcp).

#### Prerequisites

- VS Code with GitHub Copilot extension installed
- GitHub Copilot subscription
- Built `copilot-agent-chain` binary

#### Configuration

Configure the MCP server for GitHub Copilot by creating or editing the MCP settings file:

**Location:**
- **macOS/Linux**: `~/.config/github-copilot/mcp_settings.json`
- **Windows**: `%APPDATA%\github-copilot\mcp_settings.json`

**Configuration:**

```json
{
  "mcpServers": {
    "copilot-agent-chain": {
      "command": "/absolute/path/to/copilot-agent-chain",
      "args": [],
      "env": {
        "REPO_ROOT": "${workspaceFolder}",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

**Important Configuration Notes:**
- Use **absolute path** for `command` — Replace `/absolute/path/to/copilot-agent-chain` with the full path to your built binary
- Use `${workspaceFolder}` for `REPO_ROOT` — VS Code will automatically resolve this to your current workspace
- If you need to use a different repository, you can specify an absolute path instead: `"REPO_ROOT": "/absolute/path/to/your/repo"`
- Adjust `LOG_LEVEL` as needed: `debug`, `info`, `warn`, or `error`
- Ensure the binary has executable permissions: `chmod +x /path/to/copilot-agent-chain`

#### Example Configuration

With the binary in `/usr/local/bin` and using the current workspace:

```json
{
  "mcpServers": {
    "copilot-agent-chain": {
      "command": "/usr/local/bin/copilot-agent-chain",
      "args": [],
      "env": {
        "REPO_ROOT": "${workspaceFolder}",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

**Note:** The `${workspaceFolder}` variable is automatically resolved by VS Code to the current workspace directory. This makes your configuration portable across different machines and users.

#### Restart and Verify

1. **Restart VS Code** to load the new configuration
2. **Open GitHub Copilot Chat** (View > Open View > GitHub Copilot Chat)
3. **Verify the connection**: The MCP tools should now be available in Copilot Chat
4. **Check available tools**: `run_with_orchestrator`, `list_agents`, `evaluate_prompt`, `run_agent`

#### Using the Agent Chain in Copilot Chat

Once configured, you can use the orchestrator tools directly in GitHub Copilot Chat:

**List available agents:**
```
@workspace Use the list_agents tool to see all available agents
```

**Run with orchestrator:**
```
@workspace Use run_with_orchestrator to review the authentication module for security issues
```

**Run a specific agent:**
```
@workspace Use run_agent with agent="code-reviewer" to check internal/orchestrator/orchestrator.go
```

The orchestrator will automatically evaluate your prompt, select appropriate agents, chain them together, and return synthesized results.

### Option 2: Run the Server Directly

```bash
# Set repository root
export REPO_ROOT=/path/to/your/repo

# Run the server
./copilot-agent-chain
```

You should see output like:

```
2025-12-07T10:30:45.123Z	info	config	Loaded configuration	{"repo_root": "/path/to/your/repo", "log_level": "info"}
2025-12-07T10:30:45.234Z	info	discovery	Agent discovery started	{"path": "/path/to/your/repo/.github/agents"}
2025-12-07T10:30:45.345Z	info	discovery	Discovered agents	{"count": 4, "agents": ["code-reviewer", "architecture-advisor", "test-generator", "documentation-writer"]}
2025-12-07T10:30:45.456Z	info	server	MCP server started	{"transport": "stdio"}
```

### Option 3: Use with Copilot CLI

In another terminal, run:

```bash
# Use the orchestrator (auto-selects and chains agents)
copilot --agent=orchestrator --prompt "Review the main.go file for security issues"

# Use a specific agent
copilot --agent=code-reviewer --prompt "Check internal/orchestrator/orchestrator.go for performance"

# List available agents
copilot --agent=orchestrator --prompt "List all available agents"
```

## Quick Examples

### Example 1: Code Review

```bash
copilot --agent=code-reviewer --prompt "Review cmd/server/main.go for error handling and resource cleanup"
```

Expected output: Detailed feedback on code quality, potential issues, and improvement suggestions.

### Example 2: Architecture Feedback

```bash
copilot --agent=architecture-advisor --prompt "Should we extract the context accumulation logic into a separate package?"
```

Expected output: Architectural analysis, design patterns, scalability considerations.

### Example 3: Generate Tests

```bash
copilot --agent=test-generator --prompt "Generate comprehensive tests for internal/agents/discovery.go"
```

Expected output: Test code that can be added to your test suite.

### Example 4: Document Code

```bash
copilot --agent=documentation-writer --prompt "Create a developer guide for the orchestrator package"
```

Expected output: Markdown documentation explaining the orchestrator design and usage.

## Understanding Agent Chains

When you use the **orchestrator** agent (the default), the system:

1. **Evaluates your prompt** — Checks clarity, specificity, and intent
2. **Selects agents** — Matches keywords to agent capabilities
3. **Chains agents** — Executes them in sequence, passing context
4. **Synthesizes results** — Combines outputs into a coherent response

For example, this prompt:

```bash
copilot --agent=orchestrator --prompt "Review the authentication module, ensure it follows best practices, and generate tests"
```

Might trigger:
1. **Code Reviewer** — Review authentication code
2. **Architecture Advisor** — Check design patterns
3. **Test Generator** — Create test cases

## Customizing Agents

Agents are defined in `.github/agents/` as markdown files with YAML frontmatter. Each file specifies:

- **name** — Agent identifier
- **description** — What the agent does
- **keywords** — Topics it handles
- **system_prompt** — Instructions for Copilot

See the [Agent Guide](agents.md) for details on creating custom agents.

## Troubleshooting

### Server fails to start

**Problem:** "Cannot find .github/agents directory"

**Solution:** Ensure `REPO_ROOT` points to a directory with a `.github/agents/` folder.

```bash
# Check directory structure
ls -la $REPO_ROOT/.github/agents/

# Or explicitly set REPO_ROOT
export REPO_ROOT=$(pwd)
```

### Copilot CLI not found

**Problem:** "copilot: command not found"

**Solution:** Install Copilot CLI:

```bash
# macOS
brew install gh-copilot

# Linux/Windows - Download from GitHub
# https://github.com/github/gh-cli/releases
```

### Authentication errors

**Problem:** "Not authenticated" or "Invalid credentials"

**Solution:** Re-authenticate with Copilot:

```bash
copilot auth logout
copilot auth login
copilot auth status
```

### Agent timeout

**Problem:** Agents take longer than the configured timeout

**Solution:** Increase the timeout:

```bash
export AGENT_TIMEOUT=120s
./copilot-agent-chain
```

### GitHub Copilot MCP server not connecting

**Problem:** "Server failed to start" or MCP tools not appearing in Copilot Chat

**Solution:** Check the following:

1. **Verify the configuration file location:**
   - macOS/Linux: `~/.config/github-copilot/mcp_settings.json`
   - Windows: `%APPDATA%\github-copilot\mcp_settings.json`

2. **Verify the binary path is absolute:**
   ```json
   {
     "mcpServers": {
       "copilot-agent-chain": {
         "command": "/absolute/path/to/copilot-agent-chain",
         ...
       }
     }
   }
   ```

3. **Ensure the binary is executable:**
   ```bash
   chmod +x /path/to/copilot-agent-chain
   ```

4. **Check your workspace has the required directory structure:**
   ```bash
   ls .github/agents/
   ```
   The workspace should contain a `.github/agents/` directory with agent definition files.

5. **Verify the JSON syntax is valid:** Use a JSON validator or check for missing commas, brackets, or quotes

6. **Restart VS Code completely** after making configuration changes

7. **Check GitHub Copilot's output logs** (View > Output or `Cmd/Ctrl + Shift + U`, then select "GitHub Copilot" from the dropdown) for detailed error messages

8. **Try running the binary manually** to verify it works:
   ```bash
   REPO_ROOT=/path/to/your/repo /path/to/copilot-agent-chain
   ```

### MCP tools not appearing in Copilot Chat

**Problem:** Server starts but tools aren't available in GitHub Copilot Chat

**Solution:**

1. **Verify GitHub Copilot extension is active** in VS Code
2. **Reload VS Code window** (Cmd/Ctrl + Shift + P → "Reload Window")
3. **Check agent discovery** by running the binary manually and verifying "Discovered agents" appears in the output
4. **Ensure `.github/agents/` directory exists** and contains valid agent definition files (`.md` files with YAML frontmatter)
5. **Check the configuration syntax** matches the example format exactly

## Next Steps

- **[Development Guide](development.md)** — Set up a local development environment
- **[Architecture Guide](architecture.md)** — Understand how the system works
- **[Agent Guide](agents.md)** — Learn about agents and create custom ones
- **[API Reference](../api/index.md)** — Complete API documentation
- **[Examples](../examples/index.md)** — Real-world usage scenarios

## Getting Help

- 📖 See [Troubleshooting Guide](troubleshooting.md) for common issues
- 🐛 Report bugs on [GitHub Issues](https://github.com/rayprogramming/copilot-agent-chain/issues)
- 💬 Ask questions on [GitHub Discussions](https://github.com/rayprogramming/copilot-agent-chain/discussions)
- 📚 Read full documentation in the `docs/` folder

---

**Next:** [Development Guide](development.md)
