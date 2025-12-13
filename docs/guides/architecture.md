---
layout: default
title: Architecture Guide
nav_order: 2
parent: Guides
---

# Architecture Guide

This document describes the system architecture, design decisions, and component interactions in CoPilot OS.

## System Overview

CoPilot OS is a Model Context Protocol (MCP) server that orchestrates agent chains. It acts as a bridge between external callers (Copilot CLI, VS Code, etc.) and agents defined in a repository's `.github/agents/` directory.

```
External Caller (Copilot CLI)
        │
        │ (MCP Protocol)
        ▼
┌─────────────────────────────┐
│   MCP Server (hypermcp)     │
│  ┌───────────────────────┐  │
│  │   Orchestrator        │  │
│  │  (Go Component)       │  │
│  └───────────────────────┘  │
└─────────────────────────────┘
        │
        ├─► Agent Discovery (scan .github/agents/)
        │
        ├─► Prompt Evaluator (heuristics)
        │
        ├─► Agent Selection (keyword matching)
        │
        └─► CLI Invoker (spawn copilot subprocesses)
```

## Core Components

### 1. MCP Server (cmd/server/main.go)

**Responsibility**: Bootstrap and run the MCP server using hypermcp framework.

**Key Decisions**:
- Uses hypermcp for transport abstraction and standard infrastructure
- Stdio transport for simplicity (can evolve to HTTP if needed)
- Structured logging with zap for observability
- Configuration via environment variables for flexibility

**Exposes Tools**:
- `run_with_orchestrator` — Primary entry point
- `list_agents` — Discover available agents
- `evaluate_prompt` — Test prompt evaluation
- `run_agent` — Run specific agent (bypass orchestrator)

### 2. Orchestrator (internal/orchestrator/orchestrator.go)

**Responsibility**: Core logic for intelligent agent chaining.

The orchestrator executes these steps:

```
Input Prompt
    │
    ▼
1. Evaluate Clarity
    │ ├─► Detect vagueness
    │ ├─► Check for missing context
    │ └─► Auto-refine if needed
    │
    ▼
2. Select Agents
    │ ├─► Extract keywords from refined prompt
    │ ├─► Match against agent capabilities
    │ └─► Build chain (or use explicit chain)
    │
    ▼
3. Execute Chain
    │ ├─► Initialize context
    │ ├─► For each agent:
    │ │   ├─► Invoke via Copilot CLI
    │ │   ├─► Parse output
    │ │   ├─► Accumulate to context
    │ │   └─► Adapt chain if needed
    │ │
    │ └─► Synthesize final output
    │
    ▼
Final Result (JSON)
```

**Key Design Decisions**:
- **Single responsibility**: Orchestrator focuses on chaining, not agent definitions
- **Go component (not MCP agent)**: Ensures reliability and control; orchestration logic is in code, not in an LLM
- **Heuristic-based evaluation**: Rule-based prompt evaluation for predictability and performance
- **Keyword matching**: Simple, fast agent selection; evolves to semantic matching if needed
- **JSON context flow**: Structured, parseable context; easy to serialize/deserialize
- **Sequential execution**: Agents run in sequence (can parallelize in future if needed)

### 3. Prompt Evaluator (internal/prompt/evaluator.go)

**Responsibility**: Assess prompt clarity and auto-refine if needed.

**Heuristics**:
- **Vagueness detection**: Check for non-specific terms (module, thing, stuff, etc.)
- **Context completeness**: Detect missing specifics (no mention of scope, concerns, constraints)
- **Actionability**: Ensure prompt describes a task, not just a concept
- **Refinement suggestions**: Auto-improve prompt or request clarification

**Example**:
```
Input: "Review the auth module"
Issues:
  - "module" is too broad (which file? which aspect?)
  - Missing context (security? performance? design?)
Refined: "Review the authentication module for security vulnerabilities and performance bottlenecks"
Confidence: 0.85 (high enough to proceed)
```

**Key Decision**: Heuristic-based (rule-based) rather than LLM-based ensures:
- Fast evaluation (no API calls)
- Predictable behavior
- Easy to debug and tune

### 4. Agent Discovery (internal/agents/discovery.go)

**Responsibility**: Scan repository for agent definitions and build registry.

**Process**:
```
1. Scan $REPO_ROOT/.github/agents/ for *.md files
2. For each file:
   ├─► Parse YAML frontmatter:
   │   ├─► name
   │   ├─► description
   │   ├─► keywords
   │   └─► (other metadata)
   └─► Validate and add to registry
3. Return registry for lookups
```

**Agent Metadata Format** (YAML frontmatter in `.md`):
```yaml
---
name: code-reviewer
description: Specialized Go code reviewer
keywords: [code-review, go, quality, testing, correctness, performance]
---
(markdown content)
```

**Key Decision**: Load agents once at startup (or periodically), cache registry for fast lookups. This ensures consistent behavior during a request.

### 5. Agent Selection (internal/agents/selection.go)

**Responsibility**: Match user prompt to available agents via keyword matching.

**Algorithm**:
```
1. Extract keywords from refined prompt (nouns, verbs, domain terms)
2. For each available agent:
   ├─► Calculate match score:
   │   ├─► Exact keyword matches (weight: 2.0)
   │   ├─► Partial matches (weight: 1.0)
   │   └─► Synonym matches (weight: 0.5)
   │
   └─► Accumulate scores
3. Rank agents by score
4. Select top N agents (default: 2-3) to form chain
```

**Example**:
```
Prompt: "Review authentication module for security and add tests"
Keywords: [authentication, security, testing, review, add]

Agent Scores:
  - code-reviewer: 3.0 (matches: review, security)
  - test-generator: 2.5 (matches: testing, add)
  - architecture-advisor: 1.0 (matches: authentication)

Selected Chain: [code-reviewer, test-generator]
```

**Key Decision**: Keyword-based (not semantic) ensures:
- Fast selection
- Transparent and debuggable
- Handles domain-specific keywords well

### 6. CLI Invoker (internal/cli/invoker.go)

**Responsibility**: Spawn and manage Copilot CLI subprocesses.

**Process**:
```
1. Build command: copilot --agent=<name> --prompt="<prompt>"
2. Set up I/O (stdin, stdout, stderr)
3. Execute subprocess with timeout
4. Parse JSON output
5. Handle errors and retries
6. Return structured result
```

**Error Handling**:
- Timeout: Retry once, then fail gracefully
- CLI not found: Return helpful error
- Agent failure: Capture stderr, include in context
- Connection errors: Retry with backoff

**Key Decision**: Subprocess invocation (not embedding Copilot) allows:
- Using user's authenticated Copilot CLI
- Isolation from server crashes
- Easy to test with mocks

### 7. Context Accumulation

**Responsibility**: Thread context through agent chain.

**Structure**:
```json
{
  "originalPrompt": "user's initial request",
  "refinedPrompt": "orchestrator's refined version",
  "evaluationFeedback": {
    "isClear": true,
    "changes": ["clarified scope"]
  },
  "agentResults": [
    {
      "agent": "code-reviewer",
      "output": { ... agent response ... },
      "executedAt": "2024-01-15T10:30:00Z"
    },
    {
      "agent": "test-generator",
      "output": { ... },
      "context": {
        "previousAgent": "code-reviewer",
        "focusAreas": ["identified_issues"]
      }
    }
  ],
  "finalOutput": "synthesized result from all agents"
}
```

**Key Decision**: JSON accumulation ensures:
- Structured, parseable context
- Easy to serialize/transmit
- Can be analyzed or stored
- Transparent to observers

## Design Trade-offs

### 1. Go Orchestrator vs. Agent-based Orchestrator

**Choice**: Go component (not an LLM agent)

**Rationale**:
- **Reliability**: Control flow is in code, not in prompt/LLM
- **Performance**: No extra LLM call for routing
- **Debuggability**: Logic is explicit and testable
- **Cost**: No additional API calls

**Trade-off**: Less flexible if needs become complex (but evolves easily)

### 2. Heuristic Evaluation vs. LLM-based Evaluation

**Choice**: Heuristics (rule-based)

**Rationale**:
- **Speed**: No API call
- **Cost**: Free
- **Predictability**: Consistent behavior
- **Debuggability**: Rules are explicit

**Trade-off**: May miss nuanced clarity issues (addressed via logging and continuous improvement)

### 3. Sequential Execution vs. Parallel

**Choice**: Sequential (for now)

**Rationale**:
- **Simplicity**: Easier to implement and reason about
- **Context passing**: Later agents can use earlier results
- **Resource efficiency**: Don't spawn all agents at once
- **Determinism**: Easier to debug

**Evolution**: Can parallelize agents that don't depend on each other

### 4. Keyword Matching vs. Semantic Matching

**Choice**: Keyword matching

**Rationale**:
- **Speed**: No embeddings or LLM call
- **Transparency**: Easy to debug why agent was selected
- **Cost**: Free
- **Robustness**: Works well for domain-specific agents

**Evolution**: Can add semantic matching using embeddings if needed

### 5. Subprocess Invocation vs. Library Integration

**Choice**: Subprocess (spawn Copilot CLI)

**Rationale**:
- **User authentication**: Uses their configured credentials
- **Isolation**: Server doesn't crash if CLI has issues
- **Independence**: No tight coupling to Copilot internals
- **Testing**: Easy to mock

**Trade-off**: Subprocess overhead (acceptable for agent workflows)

## Extension Points

### 1. Add New Agent Selection Strategy

Current: Keyword matching

Future options:
- Semantic matching (embeddings)
- User-provided explicit chains
- LLM-based routing agent
- Rule-based (if X then use agent Y)

**Implementation**: Add `Selector` interface, multiple implementations.

### 2. Support Parallel Agent Execution

Current: Sequential

Future:
- Spawn independent agents in parallel
- Wait for all to complete
- Merge results in final synthesis

**Implementation**: Use goroutines with channels for coordination.

### 3. Add Agent Chaining Strategies

Current: Simple sequence

Future:
- Branching (if condition, take different path)
- Loops (iterate on results)
- Conditional routing (based on intermediate results)

**Implementation**: Extend `Chain` type with conditional logic.

### 4. Integrate Result Caching

Current: No caching

Future:
- Cache agent outputs by (agent, prompt hash)
- Reuse across multiple runs
- Invalidation strategies

**Implementation**: Use hypermcp's built-in cache layer.

## Testing Strategy

### Unit Tests
- Prompt evaluator rules
- Agent discovery parsing
- Agent selection scoring
- Context accumulation

### Integration Tests
- End-to-end orchestration flow (with mocked CLI)
- Agent discovery with real `.md` files
- Full request/response cycle

### Manual Tests
- Against real Copilot CLI
- With various prompt types
- With different agent combinations

## Performance Considerations

### Agent Discovery
- Load once at startup → O(1) lookups
- Cache agent registry in memory
- Lazy load agent content if needed

### Prompt Evaluation
- Heuristics are O(n) in prompt length (acceptable)
- Pre-compile regex patterns if needed

### Agent Selection
- Keyword matching is O(agents × keywords) (acceptable)
- Cache selection results if same prompt repeated

### CLI Invocation
- Subprocess overhead ~100ms (acceptable for workflow)
- Consider connection pooling in future

## Security Considerations

### Input Validation
- Validate agent names (alphanumeric + hyphens)
- Sanitize prompts before passing to CLI
- Limit prompt size to prevent DoS

### Subprocess Execution
- Don't pass unsanitized user input directly to shell
- Use exec (not shell invocation) to prevent injection
- Set resource limits on subprocess

### Repository Access
- Scan only `.github/agents/` directory (not arbitrary dirs)
- Validate agent file format before parsing
- Don't execute arbitrary code from agent files

## Future Enhancements

### Short-term
- Add more comprehensive tests
- Improve prompt evaluator heuristics
- Add result caching layer
- Support agent parameters

### Medium-term
- Semantic agent selection
- Parallel agent execution
- Branching and conditional logic
- Agent dependency management

### Long-term
- Distributed agent execution
- Agent composition and reuse
- Performance optimization and scaling
- Multi-language agent support

---

**Previous:** [Development Guide](development.md) | **Next:** [Agent Guide](agents.md)
