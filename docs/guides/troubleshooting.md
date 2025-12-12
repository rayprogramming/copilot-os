---
layout: default
title: Troubleshooting Guide
---

# Troubleshooting Guide

Common issues and solutions for Copilot Agent Chain.

## Server Issues

### Server Won't Start

**Symptom**: Command fails with error message

**Check 1: REPO_ROOT Directory**

```bash
# Error: "Cannot find .github/agents directory"

# Solution: Verify directory structure
ls -la $REPO_ROOT/.github/agents/

# Or explicitly set REPO_ROOT
export REPO_ROOT=$(pwd)
ls -la .github/agents/

# Verify correct path
pwd
```

**Check 2: Go Installation**

```bash
# Verify Go is installed
go version

# Should output: go version goX.XX.X ...
```

**Check 3: Dependencies**

```bash
# Verify dependencies are available
go mod download
go mod tidy
go mod verify
```

**Solution**:
```bash
# Clean rebuild
rm copilot-agent-chain
go clean -cache
go build -o copilot-agent-chain ./cmd/server

# Check for build errors
./copilot-agent-chain --help
```

### High Memory Usage

**Symptom**: Server uses excessive memory

**Causes**:
- Cache size too large
- Agent outputs are large
- Memory leak in code

**Solutions**:

```bash
# Reduce cache size
export CACHE_SIZE=100
export CACHE_TTL=15m

# Or disable caching
export CACHE_SIZE=0

# Monitor memory
top -p $(pgrep -f copilot-agent-chain)
```

### Server Crashes

**Symptom**: Server exits unexpectedly

**Troubleshoot**:

```bash
# Run with debug logging
export LOG_LEVEL=debug
./copilot-agent-chain

# Watch for error messages
# Look for "panic" or "fatal" in output

# Check system resources
free -h
df -h
```

**Common Causes**:
- Out of memory → Increase system memory or reduce cache
- Disk full → Clean up disk space
- File descriptor limits → `ulimit -n`

## Copilot CLI Issues

### Copilot CLI Not Found

**Symptom**: "copilot: command not found"

**Check Installation**:

```bash
# Verify installation
which copilot
copilot --version

# Output should be: copilot X.X.X
```

**Install Copilot CLI**:

```bash
# macOS
brew install gh-copilot

# Or download from GitHub
# https://github.com/github/gh-cli/releases

# Linux (example)
sudo apt-get install gh-copilot
```

**Add to PATH**:

```bash
# If installed but not in PATH
export PATH=$PATH:/usr/local/bin

# Or add to shell profile (~/.bashrc, ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.zshrc
source ~/.zshrc
```

### Authentication Errors

**Symptom**: "Not authenticated" or "Invalid credentials"

**Check Authentication Status**:

```bash
copilot auth status

# Should show: You are authenticated
```

**Re-Authenticate**:

```bash
# Logout first
copilot auth logout

# Login again
copilot auth login

# Verify authentication
copilot auth status
```

**Verify GitHub Access**:

```bash
# Check if you can access GitHub
gh auth status

# Should show: authenticated
```

### CLI Timeout

**Symptom**: Copilot CLI command times out

**Check Configuration**:

```bash
# Increase server-side timeout
export AGENT_TIMEOUT=120s

# Restart server
./copilot-agent-chain
```

**Try Simpler Request**:

```bash
# Try with shorter prompt
copilot --agent=code-reviewer --prompt "Review foo.go"

# If this works, longer prompts might need more time
```

## Agent Issues

### Agent Not Discovered

**Symptom**: Agent not found in agent list

**Check Directory**:

```bash
# Verify agent files exist
ls -la $REPO_ROOT/.github/agents/

# Should show files like: code-reviewer.md, architecture-advisor.md, etc.
```

**Check File Format**:

```bash
# Agent files must be .md (Markdown)
# With YAML frontmatter at the top

# Verify YAML frontmatter
head -5 $REPO_ROOT/.github/agents/code-reviewer.md

# Should show:
# ---
# name: code-reviewer
# description: ...
# keywords: [...]
# ---
```

**Debug Discovery**:

```bash
# Run with debug logging
export LOG_LEVEL=debug
./copilot-agent-chain

# Look for "Discovered agents" message
# Should list all found agents
```

**Solution**:
```bash
# Check REPO_ROOT is correct
echo $REPO_ROOT

# Create .github/agents directory if needed
mkdir -p $REPO_ROOT/.github/agents

# Add agent files
# (see Agent Guide for format)
```

### Agent Timeouts

**Symptom**: Agent execution exceeds timeout

**Check Timeout Setting**:

```bash
echo $AGENT_TIMEOUT

# Default is 30s - may be too short for complex analysis
```

**Increase Timeout**:

```bash
# Try 60 seconds
export AGENT_TIMEOUT=60s

# Restart server
./copilot-agent-chain

# Try again with same prompt
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go"
```

**Simplify Request**:

```bash
# More specific = faster
copilot --agent=code-reviewer --prompt "Review the executeChain function in orchestrator.go"

# Less specific = slower
copilot --agent=code-reviewer --prompt "Review all orchestrator code"
```

### Agent Returns No Output

**Symptom**: Agent runs but produces empty result

**Verify Agent is Responding**:

```bash
# Test with simple prompt
copilot --agent=code-reviewer --prompt "What is Go?"

# If no output, issue is with Copilot CLI or authentication
# If output appears, issue might be with specific file/prompt
```

**Check File Path**:

```bash
# Make sure file exists
ls -la internal/orchestrator/orchestrator.go

# Use absolute path if needed
copilot --agent=code-reviewer --prompt "Review /home/user/repo/internal/orchestrator/orchestrator.go"
```

**Debug with Logging**:

```bash
export LOG_LEVEL=debug
./copilot-agent-chain

# Look for agent invocation logs
# Check stderr for errors from agent
```

## Performance Issues

### Slow Agent Execution

**Symptom**: Agents take 30+ seconds to respond

**Check Network**:

```bash
# Verify internet connectivity
ping github.com

# Check GitHub status
# https://www.githubstatus.com
```

**Check Load**:

```bash
# Monitor system resources
top
free
```

**Optimize Configuration**:

```bash
# Reduce cache lookups
export CACHE_SIZE=100

# Lower log level
export LOG_LEVEL=warn

# Reduce timeout if acceptable
export AGENT_TIMEOUT=30s
```

### Inconsistent Performance

**Symptom**: Sometimes fast, sometimes slow

**Likely Causes**:
- Network latency
- Copilot API response time
- System load
- Cache hits/misses

**Optimize**:

```bash
# Enable caching
export CACHE_SIZE=5000
export CACHE_TTL=2h

# Run orchestrator during off-peak hours
# Monitor with debug logging
export LOG_LEVEL=debug
```

## Prompt Issues

### Prompt Gets Refined

**Symptom**: Orchestrator refines your vague prompt

**Example**:
```bash
# You input:
copilot --agent=orchestrator --prompt "Review the module"

# Orchestrator refines to:
# "Review the authentication module for security, correctness, and performance"
```

**Why**: Orchestrator detects vagueness and auto-refines

**Solution**: Be more specific in your prompt

```bash
# Instead of:
copilot --agent=orchestrator --prompt "Review the code"

# Use:
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go for correctness and performance"
```

### Prompt Gets Rejected

**Symptom**: Orchestrator says prompt is too unclear

**Example**:
```bash
# Very vague prompt
copilot --agent=orchestrator --prompt "Review stuff"

# Result: "Confidence too low, cannot proceed"
```

**Solution**: Provide context and specificity

```bash
# Provide more details
copilot --agent=orchestrator --prompt "Review the agent selection algorithm in internal/agents/selection.go for correctness, efficiency, and maintainability"
```

## Result Issues

### Result Is Truncated

**Symptom**: Output ends abruptly or is incomplete

**Causes**:
- Agent timeout
- Network interrupted
- Output too large

**Solutions**:

```bash
# Increase timeout
export AGENT_TIMEOUT=120s

# Use more specific prompt (shorter response)
copilot --agent=code-reviewer --prompt "Review only the calculateScore function in selection.go"

# Check for error messages in logs
export LOG_LEVEL=debug
```

### Result Is Irrelevant

**Symptom**: Agent output doesn't match your request

**Check Prompt Clarity**:

```bash
# Vague prompt = vague response
copilot --agent=orchestrator --prompt "What about testing?"

# Specific prompt = specific response
copilot --agent=test-generator --prompt "Generate comprehensive unit tests for the PromptEvaluator.Evaluate function"
```

**Check Agent Selection**:

```bash
# Verify correct agent was selected
copilot --agent=orchestrator --prompt "List all available agents"

# Use specific agent if orchestrator chose wrong one
copilot --agent=code-reviewer --prompt "Review internal/agents/selection.go"
```

## Getting Help

### Collect Debugging Information

Before reporting an issue, gather:

```bash
# System info
uname -a
go version
copilot --version

# Configuration
echo $REPO_ROOT
echo $LOG_LEVEL
echo $AGENT_TIMEOUT

# Logs
export LOG_LEVEL=debug
./copilot-agent-chain > server.log 2>&1 &
# Make your request...
tail -100 server.log
```

### Report an Issue

When reporting issues on GitHub, include:

1. **Error message** — Full error text
2. **Reproduction steps** — Exact commands to reproduce
3. **Environment** — OS, Go version, Copilot CLI version
4. **Configuration** — Environment variables set
5. **Logs** — Debug logs showing the issue

**Example Report**:
```
**Title**: Agent timeout when reviewing large files

**Description**: 
Server times out when reviewing files over 1000 lines

**Steps to Reproduce**:
1. Set AGENT_TIMEOUT=30s
2. Run: copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go"
3. Wait 30 seconds
4. Get timeout error

**Environment**:
- OS: Linux 5.10
- Go: 1.21.3
- Copilot CLI: 2.0.1

**Expected**: Agent completes review within 60 seconds
**Actual**: Times out at 30 seconds
```

## Quick Reference

| Issue | Quick Fix |
|-------|-----------|
| Server won't start | `export REPO_ROOT=$(pwd)` |
| Copilot CLI not found | `brew install gh-copilot` |
| Authentication error | `copilot auth login` |
| Agent timeout | `export AGENT_TIMEOUT=120s` |
| Agent not discovered | `ls $REPO_ROOT/.github/agents/` |
| Slow response | Reduce cache size or increase timeout |
| Vague output | Be more specific in prompt |
| No output | Check authentication and agent logs |

---

**Next:** [Contributing Guide](contributing.md)
