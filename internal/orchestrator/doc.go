// Package orchestrator provides intelligent agent chain orchestration for the
// CopilotOS server.
//
// This package handles:
//   - Automatic Orchestration: Evaluate prompts, select agents, execute chains
//   - Explicit Chaining: Execute specific agents in a defined order
//   - Context Management: Pass accumulated context between agents
//   - Result Synthesis: Combine outputs from multiple agents
//   - Error Recovery: Handle agent failures gracefully
//
// # Orchestration Modes
//
// The orchestrator supports two execution modes:
//
// 1. Automatic Mode (RunWithAuto):
//   - Evaluates prompt clarity
//   - Refines unclear prompts
//   - Extracts keywords
//   - Selects optimal agents
//   - Executes agent chain
//   - Synthesizes results
//
// 2. Explicit Mode (RunWithExplicitChain):
//   - Uses specified agent names
//   - Executes in given order
//   - No automatic selection
//   - Full control over chain
//
// # Automatic Orchestration Workflow
//
// The automatic orchestration follows this workflow:
//
//  1. Prompt Evaluation:
//     - Analyze prompt clarity
//     - Calculate confidence score
//     - Identify issues
//
//  2. Prompt Refinement (if needed):
//     - Refine unclear prompts
//     - Expand vague terms
//     - Add specificity
//
//  3. Keyword Extraction:
//     - Extract significant terms
//     - Remove stop words
//     - Normalize keywords
//
//  4. Agent Selection:
//     - Match keywords to agent capabilities
//     - Rank agents by relevance score
//     - Select top N agents (default: 2)
//
//  5. Chain Execution:
//     - Invoke first agent with refined prompt
//     - Capture result and add to context
//     - Pass context to next agent
//     - Continue until all agents complete
//
//  6. Result Synthesis:
//     - Combine all agent outputs
//     - Structure as JSON context state
//     - Include metadata (timing, rationale)
//
// # Context Flow
//
// Context accumulates as the chain executes:
//
// Initial Context:
//
//	{
//	  "original_prompt": "Review authentication code",
//	  "refined_prompt": "Review authentication code in auth.go for security issues"
//	}
//
// After Agent 1 (Code Reviewer):
//
//	{
//	  "original_prompt": "...",
//	  "refined_prompt": "...",
//	  "agent_results": [
//	    {
//	      "agent": "code-reviewer",
//	      "output": "Found 3 security issues...",
//	      "success": true
//	    }
//	  ]
//	}
//
// After Agent 2 (Test Generator):
//
//	{
//	  "original_prompt": "...",
//	  "refined_prompt": "...",
//	  "agent_results": [
//	    { "agent": "code-reviewer", ... },
//	    {
//	      "agent": "test-generator",
//	      "output": "Generated 5 test cases...",
//	      "success": true
//	    }
//	  ],
//	  "final_output": "Combined synthesis..."
//	}
//
// # Agent Selection Algorithm
//
// The agent selection algorithm uses keyword matching:
//
//  1. Extract keywords from refined prompt
//  2. For each agent:
//     - Calculate match score with agent keywords
//     - Score based on direct matches and partial matches
//  3. Rank agents by score (descending)
//  4. Select top N agents
//  5. If no matches, fall back to top general-purpose agents
//
// Scoring Formula:
//   - Direct keyword match: +1.0 point
//   - Partial keyword match (substring): +0.5 point
//   - Normalize by agent keyword count
//   - Clamp to [0.0, 1.0]
//
// # Error Handling
//
// The orchestrator handles errors at multiple levels:
//
// 1. Prompt Evaluation Errors:
//   - Log but continue (use original prompt)
//
// 2. Agent Selection Errors:
//   - Fall back to default agents
//   - Log selection rationale
//
// 3. Agent Invocation Errors:
//   - Capture error in result
//   - Continue to next agent
//   - Include error in final context
//
// 4. Critical Errors:
//   - Return immediately
//   - Provide partial context
//
// Usage Example (Automatic Mode)
//
//	// Create orchestrator
//	orch := orchestrator.NewOrchestrator(registry, invoker, logger)
//
//	// Run automatic orchestration
//	state, err := orch.RunWithAuto(ctx, "Review authentication code")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access results
//	fmt.Printf("Selected agents: %v\n", state.SelectedAgents)
//	fmt.Printf("Final output: %s\n", state.FinalOutput)
//
// Usage Example (Explicit Mode)
//
//	// Specify exact agent chain
//	agents := []string{"code-reviewer", "test-generator"}
//
//	// Run explicit chain
//	state, err := orch.RunWithExplicitChain(ctx, "Review and test auth.go", agents)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Context State Structure
//
// The ContextState captures the entire execution:
//   - OriginalPrompt: User's original prompt
//   - RefinedPrompt: Automatically refined version
//   - EvaluationFeedback: Prompt evaluation details
//   - AgentResults: Results from each agent
//   - FinalOutput: Synthesized final result
//   - SelectedAgents: Names of agents executed
//   - SelectionRationale: Why these agents were chosen
//   - TotalDuration: Total execution time in milliseconds
//
// # Performance Considerations
//
// Orchestration performance depends on:
//   - Number of agents in chain (more agents = longer time)
//   - Agent execution time (varies by complexity)
//   - CLI overhead (process spawning, parsing)
//   - Network latency (if using remote agents)
//
// Typical timings:
//   - Prompt evaluation: <10ms
//   - Agent selection: <50ms
//   - Agent invocation: 2-30 seconds each
//   - Result synthesis: <100ms
//
// For better performance:
//   - Limit agent chain length (2-3 agents)
//   - Use explicit mode when agents are known
//   - Cache agent results at a higher level
//   - Run independent agents in parallel (future enhancement)
//
// # Thread Safety
//
// The Orchestrator is safe for concurrent use. Multiple goroutines can
// run orchestrations simultaneously without additional synchronization.
package orchestrator
