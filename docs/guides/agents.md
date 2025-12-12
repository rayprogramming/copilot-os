---
layout: default
title: Agent Guide
---

# Development Agents Guide

This document describes the development agents included in this repository and how to use them effectively.

## Overview

The repository includes four specialized development agents that help with different aspects of the codebase:

1. **Code Reviewer** — Code quality and correctness
2. **Architecture Advisor** — System design and patterns
3. **Test Generator** — Test creation and coverage
4. **Documentation Writer** — Documentation and comments

All agents are defined in `.github/agents/` and are automatically discovered by the server.

## Code Reviewer

**File**: `.github/agents/code-reviewer.md`

**Expertise**:
- Go idioms and best practices
- Code correctness and error handling
- Performance and resource management
- Design patterns and anti-patterns
- Testing strategies

**When to Use**:
- Reviewing new code for quality
- Identifying performance bottlenecks
- Catching edge cases or error handling gaps
- Ensuring Go conventions are followed
- Validating refactoring changes

**Example Usage**:

```bash
# Review a specific file
copilot --agent=code-reviewer --prompt "Review internal/orchestrator/orchestrator.go for correctness and performance"

# Review a specific function
copilot --agent=code-reviewer --prompt "Check the agent selection algorithm in internal/agents/selection.go for efficiency"

# Review code changes
copilot --agent=code-reviewer --prompt "Review the new context accumulation logic for memory leaks and efficiency"

# Quality assessment
copilot --agent=code-reviewer --prompt "Assess the error handling in the CLI invoker for robustness"
```

**What to Expect**:
- Detailed feedback on code quality
- Specific issues with line references (when available)
- Performance concerns and optimization suggestions
- Go idiom violations and corrections
- Edge cases or error scenarios missed

**Tips**:
- Be specific about the concern (performance, correctness, style)
- Provide context about the code's purpose
- Ask about tradeoffs for design choices
- Request examples when suggesting improvements

## Architecture Advisor

**File**: `.github/agents/architecture-advisor.md`

**Expertise**:
- System architecture and design
- Module organization and boundaries
- MCP server patterns
- Scalability and reliability
- Design patterns and tradeoffs

**When to Use**:
- Discussing system design decisions
- Planning new features or refactors
- Evaluating architectural tradeoffs
- Organizing code into modules
- Ensuring system coherence

**Example Usage**:

```bash
# Design discussion
copilot --agent=architecture-advisor --prompt "Should the orchestrator be a Go component or an MCP agent? What are the tradeoffs?"

# Feature planning
copilot --agent=architecture-advisor --prompt "Design a caching layer for agent results. How should it integrate with the current architecture?"

# Module organization
copilot --agent=architecture-advisor --prompt "How should we organize the agent discovery, selection, and execution into packages?"

# Scalability
copilot --agent=architecture-advisor --prompt "The server currently executes agents sequentially. How could we support parallel execution while maintaining context flow?"

# Integration design
copilot --agent=architecture-advisor --prompt "Design how the server should integrate with Copilot CLI. What are the key integration points?"
```

**What to Expect**:
- Architectural trade-offs explained
- Design patterns and principles applied
- Scalability and reliability considerations
- Coupling and cohesion analysis
- Specific suggestions for module organization
- Diagrams or detailed descriptions of components

**Tips**:
- Present the problem before asking for a solution
- Discuss constraints (performance, complexity, time)
- Ask about tradeoffs explicitly
- Request diagrams for complex designs
- Reference existing code/patterns for context

## Test Generator

**File**: `.github/agents/test-generator.md`

**Expertise**:
- Table-driven testing in Go
- Unit and integration testing
- Edge cases and error scenarios
- Mocking and test fixtures
- Test coverage and benchmarks

**When to Use**:
- Generating tests for new code
- Improving test coverage
- Creating integration tests
- Writing benchmark tests
- Testing error paths and edge cases

**Example Usage**:

```bash
# Generate tests for a new function
copilot --agent=test-generator --prompt "Generate comprehensive tests for the PromptEvaluator.Evaluate function in internal/prompt/evaluator.go"

# Identify untested code
copilot --agent=test-generator --prompt "Review internal/orchestrator/orchestrator.go and identify untested code paths that need tests"

# Edge case tests
copilot --agent=test-generator --prompt "Generate tests for edge cases in the agent selection algorithm (empty agents, no matches, etc.)"

# Integration tests
copilot --agent=test-generator --prompt "Design integration tests for the full orchestration flow from prompt input to final output"

# Benchmark tests
copilot --agent=test-generator --prompt "Write benchmark tests for the agent discovery scanning performance"
```

**What to Expect**:
- Table-driven test code (Go idiom)
- Tests for normal cases, edge cases, and errors
- Mocking strategies for external dependencies
- Clear test names describing what's tested
- Commentary on test structure and strategy
- Coverage analysis and gaps

**Tips**:
- Provide the function signature and purpose
- Ask specifically about what scenarios to test
- Request tests for error cases you're concerned about
- Ask for help understanding how to test async code
- Request benchmark tests for performance-critical code

## Documentation Writer

**File**: `.github/agents/documentation-writer.md`

**Expertise**:
- Technical writing and clarity
- README and guide creation
- API documentation
- Architecture documentation
- Code comments and inline docs

**When to Use**:
- Writing or improving documentation
- Creating architecture guides
- Explaining complex systems
- Improving code comments
- Creating user guides and tutorials

**Example Usage**:

```bash
# Improve README
copilot --agent=documentation-writer --prompt "Improve the README.md with better examples and clearer architecture explanation"

# Write architecture docs
copilot --agent=documentation-writer --prompt "Write a detailed ARCHITECTURE.md explaining the orchestrator design and agent chain flow"

# Code comments
copilot --agent=documentation-writer --prompt "Improve comments in internal/orchestrator/orchestrator.go to explain non-obvious logic"

# API documentation
copilot --agent=documentation-writer --prompt "Create API documentation for the MCP tools (run_with_orchestrator, list_agents, etc.)"

# Development guide
copilot --agent=documentation-writer --prompt "Write a DEVELOPMENT.md guide for local setup, building, testing, and contributing"

# Inline comments
copilot --agent=documentation-writer --prompt "Add helpful comments to internal/prompt/evaluator.go explaining the heuristic rules"
```

**What to Expect**:
- Well-structured documentation
- Clear explanations of complex concepts
- Examples and use cases
- Helpful inline comments
- Diagrams or ASCII art where useful
- Step-by-step guides

**Tips**:
- Specify the audience level (beginners, experienced devs)
- Provide context about what needs to be documented
- Ask for specific sections or structures
- Request examples for clarity
- Ask for diagrams if describing complex flows

## Using Agents Together (Orchestrator)

The intelligent orchestrator automatically selects and chains agents based on your prompt. This is often the best way to get comprehensive results:

```bash
# Full code review with tests and docs
copilot --agent=orchestrator --prompt "Perform a comprehensive review of the orchestrator including code quality, architecture, tests, and documentation needs"

# Feature planning with all perspectives
copilot --agent=orchestrator --prompt "Design a new feature to support parallel agent execution with context synchronization"

# Complete module review
copilot --agent=orchestrator --prompt "Review and improve the agent selection module including code quality, architecture, test coverage, and documentation"

# Investigation and planning
copilot --agent=orchestrator --prompt "Investigate the current prompt evaluator heuristics and propose improvements including refactoring, tests, and updated documentation"
```

**What Happens**:
1. Orchestrator evaluates your prompt
2. Auto-refines if needed
3. Selects optimal agents (e.g., code-reviewer + test-generator)
4. Executes agents sequentially with context flow
5. Synthesizes final comprehensive result

## Agent Selection Reference

The orchestrator uses keyword matching to select agents. Here's what triggers each agent:

### Code Reviewer Selected By
- Keywords: review, quality, correctness, bug, issue, fix, check, error, performance, refactor
- Works well for: Code quality, bug investigation, performance analysis, best practices

### Architecture Advisor Selected By
- Keywords: architecture, design, pattern, structure, organize, scale, integration, module
- Works well for: System design, planning, module organization, scalability decisions

### Test Generator Selected By
- Keywords: test, testing, coverage, unit-test, integration-test, mock, benchmark, edge-case
- Works well for: Test creation, coverage improvement, test strategy

### Documentation Writer Selected By
- Keywords: documentation, readme, guide, comment, explain, write, document, api-doc
- Works well for: Documentation creation, improving clarity, writing guides

## Tips for Effective Agent Use

### 1. Be Specific
**Good**: "Review the agent selection algorithm in internal/agents/selection.go for correctness and efficiency"
**Poor**: "Review the code"

### 2. Provide Context
**Good**: "The prompt evaluator currently uses heuristics. Should we add LLM-based evaluation? What are the tradeoffs?"
**Poor**: "Improve the evaluator"

### 3. Ask for What You Need
**Good**: "Generate tests for the orchestrator focusing on error handling and agent chain failures"
**Poor**: "Write tests"

### 4. Use Orchestrator for Complex Tasks
**Good**: "Perform a comprehensive review of the new feature" (orchestrator chains multiple agents)
**Poor**: Running each agent separately

### 5. Reference Specific Files
**Good**: "Review internal/orchestrator/orchestrator.go focusing on the chain execution loop"
**Poor**: "Review the code"

## Customizing Agents

All agents are defined as Markdown files in `.github/agents/`. You can:

1. **Modify existing agents**: Edit the instructions in the `.md` file
2. **Add new agents**: Create a new `.md` file with YAML frontmatter and instructions
3. **Specialize agents**: Add domain-specific knowledge or constraints

**Agent File Format**:
```yaml
---
name: your-agent-name
description: Short description of the agent
keywords: [keyword1, keyword2, keyword3]
---

Your agent instructions go here. Explain:
- What the agent knows about
- What it's good at
- What it should focus on
- Any constraints or guidelines
```

The orchestrator will automatically discover and use new agents.

## Creating Custom Agents

To create a new agent, follow these steps:

### 1. Create the Agent File

Create a new markdown file in `.github/agents/` named `your-agent-name.md`:

```markdown
---
name: your-agent-name
description: A brief description of what this agent does
keywords: [keyword1, keyword2, keyword3]
---

## Overview

Describe the agent's purpose and expertise here.

## Guidelines

Provide specific instructions for how the agent should approach tasks.
```

### 2. Define Keywords

Choose keywords that represent what the agent specializes in. These are used for:
- Agent discovery and listing
- Automatic selection by the orchestrator
- Helping users find the right agent

### 3. Write Agent Instructions

Provide clear, detailed instructions that:
- Explain the agent's expertise
- Describe what it should focus on
- Specify any constraints or guidelines
- Include examples where helpful

### 4. Test the Agent

```bash
# Verify agent is discovered
copilot --agent=orchestrator --prompt "List all available agents"

# Use the agent
copilot --agent=your-agent-name --prompt "Your test prompt"
```

## Examples with Expected Results

### Example 1: Code Review

**Prompt**:
```
copilot --agent=code-reviewer --prompt "Review internal/cli/invoker.go for subprocess safety and error handling"
```

**Expected Result**:
- Analysis of subprocess invocation safety
- Error handling comprehensiveness
- Timeout handling
- Resource cleanup
- Specific improvement suggestions

### Example 2: Architecture Discussion

**Prompt**:
```
copilot --agent=architecture-advisor --prompt "Should we add result caching? What are the architectural implications?"
```

**Expected Result**:
- Trade-offs of caching vs. fresh results
- Cache invalidation strategies
- Integration approach
- Performance implications
- Recommended design

### Example 3: Test Coverage

**Prompt**:
```
copilot --agent=test-generator --prompt "Generate tests for the agent discovery process including scanning, parsing, and error cases"
```

**Expected Result**:
- Table-driven test code
- Tests for file discovery
- YAML parsing tests
- Error condition tests
- Validation tests

### Example 4: Complete Review

**Prompt**:
```
copilot --agent=orchestrator --prompt "Comprehensively review the context accumulation logic including quality, architecture, tests, and documentation"
```

**Expected Result**:
- Code review feedback
- Architectural assessment
- Test strategy and generated tests
- Documentation suggestions
- Integrated recommendations

---

**Previous:** [Architecture Guide](architecture.md) | **Next:** [Examples](../examples/index.md)
