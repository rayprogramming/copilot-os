# Example Prompts for CoPilot OS

This directory contains example prompts and expected behaviors for the orchestrator.

## Example 1: Code Review Request

**Prompt:**
```
Review the orchestrator.go file for performance issues and potential memory leaks
```

**Expected Behavior:**
- ✅ Prompt is clear and specific
- ✅ Confidence: High (>0.7)
- ✅ Selected Agent: code-reviewer
- ✅ Keywords extracted: code-review, performance, review

**Orchestrator Flow:**
1. Evaluates prompt → Clear, no refinement needed
2. Extracts keywords → ["code-review", "performance"]
3. Matches agents → code-reviewer (high score)
4. Executes: `copilot --agent=code-reviewer --prompt="..."`
5. Returns code review with performance analysis

---

## Example 2: Test Generation Request

**Prompt:**
```
Generate comprehensive unit tests for the agent registry matching algorithm including edge cases
```

**Expected Behavior:**
- ✅ Prompt is clear and specific
- ✅ Confidence: High (>0.7)
- ✅ Selected Agent: test-generator
- ✅ Keywords extracted: testing, unit-test

**Orchestrator Flow:**
1. Evaluates prompt → Clear, actionable
2. Extracts keywords → ["testing", "unit-test"]
3. Matches agents → test-generator (exact match)
4. Executes: `copilot --agent=test-generator --prompt="..."`
5. Returns test code with edge cases

---

## Example 3: Architecture Discussion

**Prompt:**
```
Design a caching strategy for agent results considering cache invalidation and memory constraints
```

**Expected Behavior:**
- ✅ Prompt is clear and specific
- ✅ Confidence: High (>0.7)
- ✅ Selected Agent: architecture-advisor
- ✅ Keywords extracted: architecture, design

**Orchestrator Flow:**
1. Evaluates prompt → Clear with context
2. Extracts keywords → ["architecture", "design"]
3. Matches agents → architecture-advisor
4. Executes: `copilot --agent=architecture-advisor --prompt="..."`
5. Returns architectural recommendations

---

## Example 4: Documentation Request

**Prompt:**
```
Write API documentation for the orchestrator public methods including usage examples
```

**Expected Behavior:**
- ✅ Prompt is clear and specific
- ✅ Confidence: High (>0.7)
- ✅ Selected Agent: documentation-writer
- ✅ Keywords extracted: documentation, api-docs

**Orchestrator Flow:**
1. Evaluates prompt → Clear, specific scope
2. Extracts keywords → ["documentation", "api-docs"]
3. Matches agents → documentation-writer
4. Executes: `copilot --agent=documentation-writer --prompt="..."`
5. Returns formatted API documentation

---

## Example 5: Multi-Agent Chain (Automatic)

**Prompt:**
```
Perform a comprehensive review of the prompt evaluator including code quality, test coverage, and documentation completeness
```

**Expected Behavior:**
- ✅ Prompt is clear and comprehensive
- ✅ Confidence: High (>0.7)
- ✅ Selected Agents: code-reviewer, test-generator, documentation-writer (chained)
- ✅ Keywords extracted: code-review, testing, documentation

**Orchestrator Flow:**
1. Evaluates prompt → Clear, multi-faceted
2. Extracts keywords → ["code-review", "testing", "documentation"]
3. Matches agents → Multiple matches (top 2-3)
4. Executes chain:
   - code-reviewer → analyzes quality
   - test-generator → reviews test coverage (receives code review context)
   - documentation-writer → checks docs (receives both contexts)
5. Synthesizes final output combining all agent results

---

## Example 6: Vague Prompt (Auto-Refined)

**Prompt:**
```
Fix the thing in the module
```

**Expected Behavior:**
- ❌ Prompt is vague and unclear
- ⚠️  Confidence: Low (<0.7)
- ✅ Auto-refinement triggered
- ✅ Refined prompt: "Fix the thing in the module. Please specify: which module, which component, what issue, and what file?"

**Orchestrator Flow:**
1. Evaluates prompt → Unclear, vague language detected
2. Detects issues: ["contains vague term 'thing'", "lacks specific context"]
3. Generates refinement suggestion
4. Could either:
   - Return refinement to user for clarification, OR
   - Proceed with best-effort agent selection based on limited keywords

---

## Example 7: Explicit Agent Chain

**Prompt:**
```
Review and improve the CLI invoker
```

**Explicit Chain:**
```json
{
  "prompt": "Review and improve the CLI invoker",
  "agents": ["code-reviewer", "test-generator", "documentation-writer"]
}
```

**Expected Behavior:**
- ✅ Uses explicitly specified agent chain
- ✅ Evaluates prompt for context
- ✅ Executes agents in specified order with context flow

**Orchestrator Flow:**
1. Evaluates prompt → Provides context
2. Validates agents exist in registry
3. Executes chain in order:
   - code-reviewer → Reviews CLI code
   - test-generator → Generates tests (with review context)
   - documentation-writer → Documents improvements (with full context)
4. Returns synthesized output

---

## Usage in MCP Client

### Via run_with_orchestrator tool:
```json
{
  "tool": "run_with_orchestrator",
  "arguments": {
    "prompt": "Review the orchestrator for performance issues"
  }
}
```

### Via run_with_orchestrator with explicit agents:
```json
{
  "tool": "run_with_orchestrator",
  "arguments": {
    "prompt": "Comprehensive module review",
    "explicit_agents": ["code-reviewer", "test-generator"]
  }
}
```

### Via evaluate_prompt tool (no execution):
```json
{
  "tool": "evaluate_prompt",
  "arguments": {
    "prompt": "Fix the thing"
  }
}
```

Response:
```json
{
  "is_clear": false,
  "confidence": 0.3,
  "feedback": "Prompt is too vague",
  "suggested_refinement": "Fix the thing. Please specify: which component, what issue...",
  "detected_issues": ["contains vague term 'thing'", "lacks context"]
}
```

---

## Testing the Examples

You can test these examples locally by:

1. **Running the server:**
   ```bash
   ./server
   ```

2. **Using a test MCP client** (or Copilot CLI if integrated)

3. **Checking logs** to see evaluation, keyword extraction, and agent selection

4. **Verifying context flow** in multi-agent chains by examining the accumulated context in logs

---

## Expected Output Structure

All orchestrator executions return a `ContextState` object:

```json
{
  "original_prompt": "Review orchestrator.go",
  "refined_prompt": "Review orchestrator.go for performance issues and code quality",
  "evaluation_feedback": {
    "is_clear": true,
    "confidence": 0.85,
    "feedback": "Prompt is clear and actionable"
  },
  "selected_agents": ["code-reviewer"],
  "selection_rationale": "Selected 'code-reviewer' based on keywords: code-review, review",
  "agent_results": [
    {
      "agent": "code-reviewer",
      "success": true,
      "output": "...",
      "duration": "2.3s"
    }
  ],
  "final_output": "Code review completed. Found 3 potential issues...",
  "total_duration": "2.5s"
}
```
