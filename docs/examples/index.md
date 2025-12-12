---
layout: default
title: Examples & Scenarios
---

# Examples & Scenarios

This section provides real-world usage examples and integration patterns for Copilot Agent Chain.

## Quick Examples

### Example 1: Simple Code Review

Review a specific Go file for code quality and correctness:

```bash
copilot --agent=code-reviewer --prompt "Review cmd/server/main.go for error handling and resource cleanup"
```

**Expected Output**:
- Code quality feedback
- Error handling analysis
- Resource cleanup verification
- Specific improvement suggestions

### Example 2: Architecture Discussion

Get architectural perspective on a design decision:

```bash
copilot --agent=architecture-advisor --prompt "Should we add result caching to improve performance? What are the architectural implications and trade-offs?"
```

**Expected Output**:
- Architectural analysis
- Trade-off discussion
- Design pattern recommendations
- Integration approach

### Example 3: Generate Tests

Create comprehensive tests for a function:

```bash
copilot --agent=test-generator --prompt "Generate comprehensive unit tests for the PromptEvaluator.Evaluate function including edge cases and error scenarios"
```

**Expected Output**:
- Table-driven test code
- Tests for normal cases
- Edge case coverage
- Error path testing

### Example 4: Improve Documentation

Get help writing better documentation:

```bash
copilot --agent=documentation-writer --prompt "Write a comprehensive guide explaining how the orchestrator works, including the agent selection algorithm and context flow"
```

**Expected Output**:
- Well-structured documentation
- Clear explanations
- Helpful examples
- Diagrams or ASCII art

## Complex Workflows

### Workflow 1: Complete Code Review Process

Review code from multiple perspectives in one command:

```bash
copilot --agent=orchestrator --prompt "Comprehensively review internal/agents/discovery.go including: code quality and correctness, architectural implications, test coverage with generated tests, and documentation improvements needed"
```

This single command:
1. ✅ Code Reviewer — Analyzes code quality
2. ✅ Architecture Advisor — Evaluates design
3. ✅ Test Generator — Identifies test gaps
4. ✅ Documentation Writer — Suggests docs

### Workflow 2: Feature Design & Implementation Review

Design, implement, and review a new feature:

1. **Design Phase**:
```bash
copilot --agent=architecture-advisor --prompt "Design a caching layer for agent results. What should we cache, how long should we keep results, and how should it integrate with the current architecture?"
```

2. **Implementation Phase**:
```bash
# Write implementation code
# Then review it...
```

3. **Review Phase**:
```bash
copilot --agent=orchestrator --prompt "Comprehensively review the new caching implementation in internal/cache/ including correctness, performance, test coverage, and documentation"
```

### Workflow 3: Performance Optimization

Find and fix performance issues:

1. **Identify Issues**:
```bash
copilot --agent=code-reviewer --prompt "Identify performance bottlenecks in internal/agents/selection.go. What could be optimized?"
```

2. **Design Solution**:
```bash
copilot --agent=architecture-advisor --prompt "Design a more efficient agent selection algorithm. What data structures and algorithms should we use?"
```

3. **Implement & Review**:
```bash
copilot --agent=orchestrator --prompt "Review the optimized agent selection implementation for correctness, performance improvement, test coverage, and documentation"
```

## Integration Patterns

### Pattern 1: Pre-Commit Code Review

Review changes before committing:

```bash
#!/bin/bash
# pre-commit-hook.sh

# Start server if not running
if ! pgrep -f "copilot-agent-chain" > /dev/null; then
    export REPO_ROOT=$(pwd)
    ./copilot-agent-chain &
    sleep 2
fi

# Review staged files
STAGED_FILES=$(git diff --cached --name-only | grep '\.go$')

for file in $STAGED_FILES; do
    echo "Reviewing $file..."
    copilot --agent=code-reviewer --prompt "Review the staged changes in $file"
done
```

### Pattern 2: PR Review Automation

Automatically review pull requests:

```bash
#!/bin/bash
# pr-review.sh

PR_URL=$1
CHANGED_FILES=$(gh pr diff $PR_URL --name-only)

echo "=== Architecture Review ===" 
copilot --agent=architecture-advisor --prompt "Review the PR for architectural impacts: $PR_URL"

echo "=== Code Quality Review ==="
for file in $CHANGED_FILES; do
    copilot --agent=code-reviewer --prompt "Review changes in $file from PR $PR_URL"
done

echo "=== Test Coverage Review ==="
copilot --agent=test-generator --prompt "Review test coverage for PR $PR_URL and identify missing tests"

echo "=== Documentation Review ==="
copilot --agent=documentation-writer --prompt "Review documentation updates in PR $PR_URL"
```

### Pattern 3: Documentation Generation Workflow

Generate comprehensive project documentation:

```bash
#!/bin/bash
# generate-docs.sh

echo "Generating documentation..."

echo "=== README ==="
copilot --agent=documentation-writer --prompt "Write a comprehensive README.md for this project including overview, quick start, features, architecture, and contributing guidelines"

echo "=== Development Guide ==="
copilot --agent=documentation-writer --prompt "Write a DEVELOPMENT.md guide covering local setup, building, testing, and contributing"

echo "=== Architecture Guide ==="
copilot --agent=documentation-writer --prompt "Write an ARCHITECTURE.md explaining system design, components, and design decisions"

echo "=== API Documentation ==="
copilot --agent=documentation-writer --prompt "Create comprehensive API documentation for all MCP tools and interfaces"
```

## Real-World Use Cases

### Use Case 1: Code Review as a Service

Use agents to automatically review code in a CI/CD pipeline:

```bash
#!/bin/bash
# ci-review.sh - runs in GitHub Actions

# Setup
export REPO_ROOT=$(pwd)
export LOG_LEVEL=info

# Start server
go build -o copilot-agent-chain ./cmd/server
./copilot-agent-chain &
SERVER_PID=$!
sleep 2

# Get changed files
CHANGED_FILES=$(git diff origin/main --name-only | grep '\.go$')

# Review each file
FAILED=0
for file in $CHANGED_FILES; do
    echo "Reviewing $file..."
    if ! copilot --agent=code-reviewer --prompt "Review $file and flag any critical issues"; then
        FAILED=1
    fi
done

# Cleanup
kill $SERVER_PID

exit $FAILED
```

### Use Case 2: Architecture Review Board

Get architectural perspective on proposed changes:

```bash
#!/bin/bash
# architecture-review.sh

echo "=== Proposed Change ==="
echo "$1"

echo ""
echo "=== Architecture Board Review ==="
copilot --agent=architecture-advisor --prompt "As an architecture board, review this proposed change: $1. Evaluate impacts, recommend alternatives, and identify risks."
```

### Use Case 3: Test Coverage Improvement

Identify and fix test coverage gaps:

```bash
#!/bin/bash
# improve-coverage.sh

echo "Analyzing test coverage..."

# Get coverage report
go test -coverprofile=coverage.out ./...
UNCOVERED_FILES=$(go tool cover -func=coverage.out | grep -E "^.*\s+\d+\.\d+%" | grep " [0-9][0-9]\.[0-9]" | cut -d: -f1 | sort | uniq)

echo "Files needing test coverage:"
for file in $UNCOVERED_FILES; do
    echo "  - $file"
    
    # Generate tests
    copilot --agent=test-generator --prompt "Analyze $file and generate tests to improve coverage from the current level to at least 80%"
done
```

## Agent Output Examples

### Code Reviewer Output Example

```
REVIEW: internal/orchestrator/orchestrator.go

Issues Found:
1. [Line 45] Error not properly checked: result := json.Unmarshal(...)
   - Should check and handle JSON unmarshaling errors
   - Risk: Silent failures in context accumulation

2. [Line 78] Missing timeout context on subprocess invocation
   - Current: exec.Command() without timeout
   - Recommended: Use context.WithTimeout()
   
3. [Line 156] Resource leak: file handle not closed
   - Suggestion: Use defer f.Close()

Performance:
✓ Good: Using sync.Once for singleton initialization
⚠ Consider: Caching agent registry to avoid rescanning

Recommendations:
1. Add comprehensive error handling
2. Implement context timeouts
3. Add resource cleanup
4. Consider result caching
```

### Architecture Advisor Output Example

```
ARCHITECTURE REVIEW: Agent Selection Algorithm

Current Design:
- Keyword matching with weighted scoring
- O(agents × keywords) complexity
- Simple, transparent, maintainable

Architectural Assessment:

Strengths:
✓ Low latency (no LLM calls)
✓ Transparent and debuggable
✓ Works well for domain-specific agents
✓ Easy to extend with new agents

Potential Improvements:
1. Semantic Matching (optional enhancement)
   - Use embeddings for better matching
   - Trade-off: Latency vs. accuracy
   
2. Machine Learning Routing (future)
   - Learn agent effectiveness over time
   - Adapt selection based on success rates

3. User-Provided Explicit Chains (planned)
   - Allow users to specify agent order
   - More control, less automation

Design Trade-offs:
- Chose keyword matching over semantic for speed
- Sequential execution vs. parallel (choose based on needs)
- Rule-based over LLM routing (for reliability)

Recommendation: Current design is sound. Consider semantic matching as enhancement.
```

## Performance Scenarios

### Scenario 1: Large File Review

```bash
# Review a large file (100+ lines)
copilot --agent=code-reviewer --prompt "Carefully review internal/orchestrator/orchestrator.go analyzing code structure, error handling, resource management, and performance"

# Expected execution time: 10-20 seconds
# Result: Comprehensive review with specific feedback
```

### Scenario 2: Multiple Agents Chain

```bash
# Full orchestrator with multiple agents
copilot --agent=orchestrator --prompt "Perform a comprehensive review of the new caching feature including code quality, architecture, tests, and documentation"

# Expected execution time: 30-60 seconds
# Result: Multi-perspective analysis from all agents
```

### Scenario 3: Large Codebase

```bash
# Review across multiple files
copilot --agent=code-reviewer --prompt "Review all files in internal/agents/ for consistency, error handling, and Go best practices"

# Expected execution time: 15-30 seconds
# Result: Analysis of consistency and patterns across module
```

## Troubleshooting Examples

### Issue: Agent Timeout

```bash
# Problem: Agent execution exceeds 30 second default timeout
# Solution: Increase timeout and use more specific prompt

export AGENT_TIMEOUT=120s
./copilot-agent-chain

# Use more specific prompt
copilot --agent=code-reviewer --prompt "Review only the Evaluate function in internal/prompt/evaluator.go"
```

### Issue: Vague Prompt

```bash
# Problem: Orchestrator refines vague prompt
# Solution: Provide more specific request

# Vague prompt - might get refined
copilot --agent=orchestrator --prompt "Review the authentication code"

# Specific prompt - no refinement needed
copilot --agent=code-reviewer --prompt "Review the authentication module in internal/auth/main.go for security vulnerabilities"
```

### Issue: Agent Not Found

```bash
# Problem: Requested agent doesn't exist
copilot --agent=security-reviewer --prompt "Review for security issues"

# Solution: Use available agents
copilot --agent=code-reviewer --prompt "Review for security vulnerabilities"

# Or list available agents
copilot --agent=orchestrator --prompt "List all available agents"
```

---

**Previous:** [Agent Guide](../guides/agents.md) | **Next:** [Troubleshooting](../guides/troubleshooting.md)
