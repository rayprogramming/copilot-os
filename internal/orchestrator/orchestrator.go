package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"github.com/rayprogramming/copilot-os/internal/agents"
	"github.com/rayprogramming/copilot-os/internal/cli"
	"github.com/rayprogramming/copilot-os/internal/prompt"
	"go.uber.org/zap"
)

// ContextState represents the accumulated state throughout agent execution.
type ContextState struct {
	OriginalPrompt     string                  `json:"original_prompt"`
	RefinedPrompt      string                  `json:"refined_prompt"`
	EvaluationFeedback prompt.EvaluationResult `json:"evaluation_feedback"`
	AgentResults       []cli.InvocationResult  `json:"agent_results"`
	FinalOutput        string                  `json:"final_output"`
	SelectedAgents     []string                `json:"selected_agents"`
	SelectionRationale string                  `json:"selection_rationale"`
	TotalDuration      int64                   `json:"total_duration_ms"`
}

// Orchestrator orchestrates agent chains intelligently.
type Orchestrator struct {
	registry  *agents.Registry
	invoker   *cli.Invoker
	evaluator *prompt.Evaluator
	logger    *zap.Logger
}

// NewOrchestrator creates a new orchestrator.
func NewOrchestrator(registry *agents.Registry, invoker *cli.Invoker, logger *zap.Logger) *Orchestrator {
	return &Orchestrator{
		registry:  registry,
		invoker:   invoker,
		evaluator: prompt.NewEvaluator(),
		logger:    logger,
	}
}

// RunWithAuto automatically evaluates the prompt, selects agents, and executes the chain.
func (o *Orchestrator) RunWithAuto(ctx context.Context, userPrompt string) (*ContextState, error) {
	state := &ContextState{
		OriginalPrompt: userPrompt,
		AgentResults:   []cli.InvocationResult{},
	}

	// Step 1: Evaluate prompt
	o.logger.Debug("evaluating prompt", zap.String("prompt", userPrompt))
	evaluation := o.evaluator.Evaluate(userPrompt)
	state.EvaluationFeedback = evaluation

	refinedPrompt := evaluation.RefinedPrompt
	if !evaluation.IsClear {
		o.logger.Info("prompt refined",
			zap.String("original", userPrompt),
			zap.String("refined", refinedPrompt),
		)
	}
	state.RefinedPrompt = refinedPrompt

	// Step 2: Extract keywords and select agents
	keywords := o.extractKeywords(refinedPrompt)
	selectedAgents := o.selectAgents(keywords, 2) // Select up to 2 agents by default

	if len(selectedAgents) == 0 {
		o.logger.Warn("no agents selected, trying broader search")
		// If no agents matched, select top agents
		selectedAgents = o.selectTopAgents(3)
	}

	state.SelectedAgents = o.agentNames(selectedAgents)
	state.SelectionRationale = o.buildRationale(keywords, selectedAgents)

	o.logger.Info("agents selected",
		zap.Strings("agents", state.SelectedAgents),
		zap.String("rationale", state.SelectionRationale),
	)

	// Step 3: Execute agent chain
	finalOutput, results, err := o.executeChain(ctx, refinedPrompt, selectedAgents, ContextState{})
	if err != nil {
		o.logger.Error("chain execution failed", zap.Error(err))
		return state, err
	}

	state.AgentResults = results
	state.FinalOutput = finalOutput

	return state, nil
}

// RunWithExplicitChain executes agents in a specific order.
func (o *Orchestrator) RunWithExplicitChain(ctx context.Context, userPrompt string, agentNames []string) (*ContextState, error) {
	state := &ContextState{
		OriginalPrompt: userPrompt,
		RefinedPrompt:  userPrompt,
		SelectedAgents: agentNames,
		AgentResults:   []cli.InvocationResult{},
	}

	// Get agent objects
	selectedAgents := make([]*agents.Agent, 0)
	for _, name := range agentNames {
		agent := o.registry.Get(name)
		if agent == nil {
			return state, fmt.Errorf("agent %q not found", name)
		}
		selectedAgents = append(selectedAgents, agent)
	}

	// Evaluate prompt (but don't change it)
	evaluation := o.evaluator.Evaluate(userPrompt)
	state.EvaluationFeedback = evaluation

	// Execute chain
	finalOutput, results, err := o.executeChain(ctx, userPrompt, selectedAgents, ContextState{})
	if err != nil {
		return state, err
	}

	state.AgentResults = results
	state.FinalOutput = finalOutput

	return state, nil
}

// executeChain executes a sequence of agents with context flow.
func (o *Orchestrator) executeChain(ctx context.Context, prompt string, agents []*agents.Agent, initialContext ContextState) (string, []cli.InvocationResult, error) {
	results := []cli.InvocationResult{}
	contextState := initialContext

	for _, agent := range agents {
		select {
		case <-ctx.Done():
			return "", results, ctx.Err()
		default:
		}

		// Build agent prompt with context
		agentPrompt := o.buildAgentPrompt(prompt, agent, contextState)

		o.logger.Debug("invoking agent",
			zap.String("agent", agent.Name),
			zap.String("prompt", agentPrompt),
		)

		// Invoke agent
		result, err := o.invoker.InvokeAgent(ctx, agent.Name, agentPrompt)
		if err != nil {
			o.logger.Error("agent invocation error",
				zap.String("agent", agent.Name),
				zap.Error(err),
			)
			// Continue to next agent instead of failing
			result = &cli.InvocationResult{
				Agent:    agent.Name,
				Success:  false,
				Error:    err.Error(),
				ExitCode: 1,
			}
		}

		results = append(results, *result)

		// If agent succeeded, include output in context
		if result.Success && result.Output != nil {
			contextState.AgentResults = append(contextState.AgentResults, *result)
		}
	}

	// Synthesize final output
	finalOutput := o.synthesizeOutput(contextState, results)

	return finalOutput, results, nil
}

// buildAgentPrompt constructs the prompt for an agent, including context.
//
// Context Accumulation Strategy:
//
// This function implements a sequential context flow pattern where each agent
// in the chain receives:
//  1. The original/refined user prompt (base context)
//  2. Agent-specific context (who they are, what they do)
//  3. Results from all previous agents (accumulated context)
//
// Prompt Structure:
//
//	<base-prompt>
//
//	[Agent Context: You are the <agent-name>. <agent-description>]
//
//	[Previous Agent Results:]
//	- Agent 1 (<name>): <output>
//	- Agent 2 (<name>): <output>
//	...
//	[Consider these results in your response]
//
// Context Flow Example:
//
//	Agent 1 (Code Reviewer): Receives only base prompt
//	Agent 2 (Test Generator): Receives base prompt + Agent 1's findings
//	Agent 3 (Documentation Writer): Receives base prompt + Agent 1 & 2 results
//
// This approach enables:
//   - Sequential refinement: later agents build on earlier findings
//   - Collaborative analysis: agents can reference each other's work
//   - Comprehensive results: final output combines all perspectives
//
// Considerations:
//   - Context grows with each agent (may hit token limits)
//   - Agent order matters (earlier agents influence later ones)
//   - All agents see all previous results (no selective context)
//
// Future enhancements:
//   - Selective context: only pass relevant previous results
//   - Context summarization: compress older results
//   - Parallel execution: run independent agents concurrently
func (o *Orchestrator) buildAgentPrompt(basePrompt string, agent *agents.Agent, contextState ContextState) string {
	// Start with base prompt
	agentPrompt := basePrompt

	// Add agent-specific context to help the agent understand its role
	agentPrompt += fmt.Sprintf("\n\n[Agent Context: You are the %s. %s]", agent.Name, agent.Description)

	// Accumulate previous agent results to enable context flow
	if len(contextState.AgentResults) > 0 {
		agentPrompt += "\n\n[Previous Agent Results:]"
		for i, prevResult := range contextState.AgentResults {
			// Include each previous agent's output in order
			agentPrompt += fmt.Sprintf("\n- Agent %d (%s): %s", i+1, prevResult.Agent, string(prevResult.Output))
		}
		// Instruct the agent to consider previous results
		agentPrompt += "\n[Consider these results in your response]"
	}

	return agentPrompt
}

// synthesizeOutput combines all agent results into a final output.
func (o *Orchestrator) synthesizeOutput(contextState ContextState, results []cli.InvocationResult) string {
	var output strings.Builder

	output.WriteString("=== Agent Chain Results ===\n\n")

	for i, result := range results {
		output.WriteString(fmt.Sprintf("## Agent %d: %s\n", i+1, result.Agent))

		if result.Success {
			output.WriteString(fmt.Sprintf("Status: ✓ Success\n"))
			if result.Output != nil {
				output.WriteString(fmt.Sprintf("Output:\n%s\n\n", string(result.Output)))
			}
		} else {
			output.WriteString(fmt.Sprintf("Status: ✗ Failed\n"))
			if result.Error != "" {
				output.WriteString(fmt.Sprintf("Error: %s\n\n", result.Error))
			}
		}
	}

	output.WriteString("\n=== Summary ===\n")
	output.WriteString(fmt.Sprintf("Total Agents Executed: %d\n", len(results)))
	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}
	output.WriteString(fmt.Sprintf("Successful Executions: %d\n", successCount))

	return output.String()
}

// selectAgents selects agents based on keywords (keyword matching).
func (o *Orchestrator) selectAgents(keywords []string, maxCount int) []*agents.Agent {
	matched := o.registry.MatchKeywords(keywords)
	if len(matched) > maxCount {
		matched = matched[:maxCount]
	}
	return matched
}

// selectTopAgents selects the top N agents by default.
func (o *Orchestrator) selectTopAgents(count int) []*agents.Agent {
	agents := o.registry.All()
	if len(agents) > count {
		return agents[:count]
	}
	return agents
}

// extractKeywords extracts keywords from the prompt for agent selection.
func (o *Orchestrator) extractKeywords(p string) []string {
	return prompt.ExtractKeywords(p)
}

// buildRationale creates a human-readable rationale for agent selection.
func (o *Orchestrator) buildRationale(keywords []string, selectedAgents []*agents.Agent) string {
	if len(selectedAgents) == 0 {
		return "No agents matched the prompt keywords"
	}

	var rationale strings.Builder
	rationale.WriteString(fmt.Sprintf("Selected based on keywords: %s. ", strings.Join(keywords, ", ")))
	rationale.WriteString("Agents: ")

	for i, agent := range selectedAgents {
		if i > 0 {
			rationale.WriteString(", ")
		}
		rationale.WriteString(agent.Name)
	}

	return rationale.String()
}

// agentNames extracts names from agent objects.
func (o *Orchestrator) agentNames(agents []*agents.Agent) []string {
	names := make([]string, len(agents))
	for i, a := range agents {
		names[i] = a.Name
	}
	return names
}

// ListAgents returns all available agents.
func (o *Orchestrator) ListAgents() []*agents.Agent {
	return o.registry.All()
}
