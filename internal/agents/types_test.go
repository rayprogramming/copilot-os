package agents

import (
	"testing"
)

func TestRegistry_Add(t *testing.T) {
	registry := NewRegistry()

	agent := &Agent{
		Name:        "test-agent",
		Description: "Test agent",
		Keywords:    []string{"test", "example"},
	}

	registry.Add(agent)

	retrieved := registry.Get("test-agent")
	if retrieved == nil {
		t.Fatal("expected agent to be retrievable")
	}

	if retrieved.Name != "test-agent" {
		t.Errorf("expected name 'test-agent', got %q", retrieved.Name)
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	registry := NewRegistry()

	retrieved := registry.Get("nonexistent")
	if retrieved != nil {
		t.Error("expected nil for nonexistent agent")
	}
}

func TestRegistry_All(t *testing.T) {
	registry := NewRegistry()

	agents := []*Agent{
		{Name: "agent1", Description: "First", Keywords: []string{"one"}},
		{Name: "agent2", Description: "Second", Keywords: []string{"two"}},
		{Name: "agent3", Description: "Third", Keywords: []string{"three"}},
	}

	for _, agent := range agents {
		registry.Add(agent)
	}

	all := registry.All()
	if len(all) != 3 {
		t.Errorf("expected 3 agents, got %d", len(all))
	}

	// Check insertion order is preserved
	if all[0].Name != "agent1" {
		t.Errorf("expected first agent 'agent1', got %q", all[0].Name)
	}
}

func TestRegistry_MatchKeywords(t *testing.T) {
	registry := NewRegistry()

	agents := []*Agent{
		{
			Name:        "code-reviewer",
			Description: "Reviews code",
			Keywords:    []string{"code-review", "go", "quality"},
		},
		{
			Name:        "test-generator",
			Description: "Generates tests",
			Keywords:    []string{"testing", "unit-tests", "go"},
		},
		{
			Name:        "documentation-writer",
			Description: "Writes docs",
			Keywords:    []string{"documentation", "readme"},
		},
	}

	for _, agent := range agents {
		registry.Add(agent)
	}

	tests := []struct {
		name             string
		keywords         []string
		expectedCount    int
		expectedFirst    string
		expectedContains []string
	}{
		{
			name:             "single keyword match",
			keywords:         []string{"testing"},
			expectedCount:    1,
			expectedFirst:    "test-generator",
			expectedContains: []string{"test-generator"},
		},
		{
			name:             "multiple keyword matches",
			keywords:         []string{"go"},
			expectedCount:    2,
			expectedFirst:    "code-reviewer", // First added
			expectedContains: []string{"code-reviewer", "test-generator"},
		},
		{
			name:             "exact match prioritized",
			keywords:         []string{"code-review"},
			expectedCount:    1,
			expectedFirst:    "code-reviewer",
			expectedContains: []string{"code-reviewer"},
		},
		{
			name:             "no matches",
			keywords:         []string{"nonexistent"},
			expectedCount:    0,
			expectedFirst:    "",
			expectedContains: []string{},
		},
		{
			name:             "multiple input keywords",
			keywords:         []string{"testing", "documentation"},
			expectedCount:    2,
			expectedFirst:    "test-generator", // Higher score
			expectedContains: []string{"test-generator", "documentation-writer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := registry.MatchKeywords(tt.keywords)

			if len(matches) != tt.expectedCount {
				t.Errorf("expected %d matches, got %d", tt.expectedCount, len(matches))
			}

			if tt.expectedCount > 0 && matches[0].Name != tt.expectedFirst {
				t.Errorf("expected first match %q, got %q", tt.expectedFirst, matches[0].Name)
			}

			// Verify all expected agents are present
			matchNames := make(map[string]bool)
			for _, match := range matches {
				matchNames[match.Name] = true
			}

			for _, expected := range tt.expectedContains {
				if !matchNames[expected] {
					t.Errorf("expected to find agent %q in matches", expected)
				}
			}
		})
	}
}

func TestRegistry_MatchKeywords_Scoring(t *testing.T) {
	registry := NewRegistry()

	agents := []*Agent{
		{
			Name:     "exact-match",
			Keywords: []string{"testing"},
		},
		{
			Name:     "partial-match",
			Keywords: []string{"test"},
		},
	}

	for _, agent := range agents {
		registry.Add(agent)
	}

	// Exact match should score higher
	matches := registry.MatchKeywords([]string{"testing"})
	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}

	if matches[0].Name != "exact-match" {
		t.Errorf("expected exact match to rank first, got %q", matches[0].Name)
	}
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	if registry == nil {
		t.Fatal("expected non-nil registry")
	}

	if len(registry.All()) != 0 {
		t.Error("expected new registry to be empty")
	}
}

func TestRegistry_DuplicateAdd(t *testing.T) {
	registry := NewRegistry()

	agent1 := &Agent{Name: "test", Description: "First"}
	agent2 := &Agent{Name: "test", Description: "Second"}

	err1 := registry.Add(agent1)
	if err1 != nil {
		t.Fatalf("first add failed: %v", err1)
	}

	err2 := registry.Add(agent2)
	if err2 == nil {
		t.Error("expected error on duplicate add, got nil")
	}

	// Should keep first agent
	retrieved := registry.Get("test")
	if retrieved.Description != "First" {
		t.Errorf("expected description 'First', got %q", retrieved.Description)
	}

	// Should only have one agent
	all := registry.All()
	if len(all) != 1 {
		t.Errorf("expected 1 agent after duplicate add attempt, got %d", len(all))
	}
}
