package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rayprogramming/copilot-os/internal/agents"
	"github.com/rayprogramming/copilot-os/internal/orchestrator"
	"github.com/rayprogramming/copilot-os/internal/prompt"
	"go.uber.org/zap"
)

// TestIntegration_PromptToAgentSelection tests the full flow from prompt evaluation to agent selection
func TestIntegration_PromptToAgentSelection(t *testing.T) {
	// Setup test agents
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, ".github", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test agent files
	testAgents := map[string]string{
		"code-reviewer.md": `---
name: code-reviewer
description: Reviews code quality
keywords: [code-review, quality, review, refactor]
---`,
		"test-generator.md": `---
name: test-generator
description: Generates tests
keywords: [testing, test, unit-test, coverage]
---`,
		"documentation-writer.md": `---
name: documentation-writer
description: Writes documentation
keywords: [documentation, docs, readme, api-docs]
---`,
	}

	for filename, content := range testAgents {
		path := filepath.Join(agentsDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Initialize components
	logger := zap.NewNop()
	discovery := agents.NewDiscovery(tmpDir, logger)
	if err := discovery.Discover(); err != nil {
		t.Fatalf("failed to discover agents: %v", err)
	}

	registry := discovery.Registry()
	evaluator := prompt.NewEvaluator()

	// Test cases
	tests := []struct {
		name                 string
		prompt               string
		expectClear          bool
		expectSelectedAgents int
		expectContainsAgent  string
	}{
		{
			name:                 "code review request",
			prompt:               "Review the orchestrator code for quality issues",
			expectClear:          true,
			expectSelectedAgents: 1,
			expectContainsAgent:  "code-reviewer",
		},
		{
			name:                 "testing request",
			prompt:               "Generate unit tests for the agent registry",
			expectClear:          true,
			expectSelectedAgents: 1,
			expectContainsAgent:  "test-generator",
		},
		{
			name:                 "documentation request",
			prompt:               "Write API documentation for the orchestrator",
			expectClear:          true,
			expectSelectedAgents: 1,
			expectContainsAgent:  "documentation-writer",
		},
		{
			name:                 "multiple agent match",
			prompt:               "Review code and generate tests for the module",
			expectClear:          true,
			expectSelectedAgents: 2,
			expectContainsAgent:  "code-reviewer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Evaluate prompt
			evalResult := evaluator.Evaluate(tt.prompt)

			if evalResult.IsClear != tt.expectClear {
				t.Errorf("expected IsClear=%v, got %v", tt.expectClear, evalResult.IsClear)
			}

			// Extract keywords
			keywords := prompt.ExtractKeywords(tt.prompt)
			if len(keywords) == 0 && tt.expectSelectedAgents > 0 {
				t.Logf("warning: no keywords extracted from prompt %q", tt.prompt)
			}

			// Match agents
			matches := registry.MatchKeywords(keywords)

			if len(matches) < tt.expectSelectedAgents {
				t.Errorf("expected at least %d agent matches, got %d", tt.expectSelectedAgents, len(matches))
			}

			// Verify expected agent is in matches
			if tt.expectContainsAgent != "" {
				found := false
				for _, match := range matches {
					if match.Name == tt.expectContainsAgent {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected agent %q in matches, got: %v", tt.expectContainsAgent, agentNames(matches))
				}
			}
		})
	}
}

// TestIntegration_OrchestratorWithoutCLI tests orchestrator logic without actual CLI invocation
func TestIntegration_OrchestratorWithoutCLI(t *testing.T) {
	// Setup test agents
	registry := agents.NewRegistry()
	registry.Add(&agents.Agent{
		Name:        "test-agent",
		Description: "Test agent",
		Keywords:    []string{"test"},
	})

	logger := zap.NewNop()

	// Note: We can't fully test orchestrator without CLI invoker
	// This test validates the orchestrator can be instantiated
	// and its registry is accessible

	orch := orchestrator.NewOrchestrator(registry, nil, logger)
	if orch == nil {
		t.Fatal("expected non-nil orchestrator")
	}

	// List agents through orchestrator
	availableAgents := orch.ListAgents()
	if len(availableAgents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(availableAgents))
	}

	if availableAgents[0].Name != "test-agent" {
		t.Errorf("expected agent name 'test-agent', got %q", availableAgents[0].Name)
	}
}

// TestIntegration_AgentDiscoveryToRegistry tests full agent discovery flow
func TestIntegration_AgentDiscoveryToRegistry(t *testing.T) {
	// Use actual project structure if running in repo
	if _, err := os.Stat("../.github/agents"); err == nil {
		logger := zap.NewNop()
		discovery := agents.NewDiscovery("..", logger)

		if err := discovery.Discover(); err != nil {
			t.Fatalf("failed to discover agents: %v", err)
		}

		registry := discovery.Registry()
		allAgents := registry.All()

		// Should discover the 4 development agents
		expectedAgents := []string{"code-reviewer", "architecture-advisor", "test-generator", "documentation-writer"}

		if len(allAgents) < len(expectedAgents) {
			t.Errorf("expected at least %d agents, got %d", len(expectedAgents), len(allAgents))
		}

		for _, expected := range expectedAgents {
			agent := registry.Get(expected)
			if agent == nil {
				t.Errorf("expected to find agent %q", expected)
			} else {
				t.Logf("found agent: %s - %s", agent.Name, agent.Description)
			}
		}

		// Test JSON export
		jsonData, err := discovery.ExportAgentsJSON()
		if err != nil {
			t.Errorf("failed to export agents as JSON: %v", err)
		}

		if len(jsonData) == 0 {
			t.Error("expected non-empty JSON export")
		}
	} else {
		t.Skip("skipping integration test - not in project directory")
	}
}

// Helper function
func agentNames(agents []*agents.Agent) []string {
	names := make([]string, len(agents))
	for i, a := range agents {
		names[i] = a.Name
	}
	return names
}
