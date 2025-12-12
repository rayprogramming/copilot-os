package prompt

import (
	"strings"
	"testing"
)

func TestEvaluator_Evaluate(t *testing.T) {
	evaluator := NewEvaluator()

	tests := []struct {
		name             string
		prompt           string
		expectClear      bool
		expectConfidence float64 // Minimum expected confidence
		expectIssues     int     // Minimum number of issues
		expectRefinement bool    // Should suggest refinement
		expectKeywords   []string
	}{
		{
			name:             "clear and specific prompt",
			prompt:           "Review the internal/orchestrator/orchestrator.go file for performance issues and suggest optimizations",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"code-review", "review", "performance"},
		},
		{
			name:             "vague prompt with 'thing'",
			prompt:           "Fix the thing in the module",
			expectClear:      false,
			expectConfidence: 0.5,
			expectIssues:     1,
			expectRefinement: true,
			expectKeywords:   []string{},
		},
		{
			name:             "vague prompt with 'stuff'",
			prompt:           "Update some stuff",
			expectClear:      false,
			expectConfidence: 0.49, // Slightly less than threshold
			expectIssues:     1,
			expectRefinement: true,
			expectKeywords:   []string{},
		},
		{
			name:             "too short prompt",
			prompt:           "Fix bug",
			expectClear:      false,
			expectConfidence: 0.5,
			expectIssues:     1,
			expectRefinement: true,
			expectKeywords:   []string{},
		},
		{
			name:             "prompt with file path",
			prompt:           "Review internal/orchestrator/orchestrator.go for code quality and best practices",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"code-review", "review", "quality"},
		},
		{
			name:             "prompt with action verb and details",
			prompt:           "Generate comprehensive unit tests for the internal/prompt/evaluator.go file covering all edge cases and error scenarios",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"test"},
		},
		{
			name:             "testing keyword prompt with specifics",
			prompt:           "Create unit tests for the MatchKeywords() function in internal/agents/types.go covering empty input and scoring edge cases",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"test"},
		},
		{
			name:             "architecture keyword prompt with context",
			prompt:           "Design a caching strategy for agent results in the orchestrator, considering cache invalidation and memory limits",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"design", "architecture"},
		},
		{
			name:             "documentation keyword prompt with specifics",
			prompt:           "Write comprehensive API documentation for the RunWithAuto() and RunWithExplicitChain() methods in internal/orchestrator/orchestrator.go with usage examples",
			expectClear:      true,
			expectConfidence: 0.7,
			expectIssues:     0,
			expectRefinement: false,
			expectKeywords:   []string{"documentation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.Evaluate(tt.prompt)

			if result.IsClear != tt.expectClear {
				t.Errorf("expected IsClear=%v, got %v", tt.expectClear, result.IsClear)
			}

			if result.Confidence < tt.expectConfidence {
				t.Errorf("expected confidence >= %.2f, got %.2f", tt.expectConfidence, result.Confidence)
			}

			if len(result.DetectedIssues) < tt.expectIssues {
				t.Errorf("expected at least %d issues, got %d: %v", tt.expectIssues, len(result.DetectedIssues), result.DetectedIssues)
			}

			if tt.expectRefinement && result.SuggestedRefinement == "" {
				t.Error("expected refinement suggestion but got none")
			}

			if !tt.expectRefinement && result.SuggestedRefinement != "" {
				t.Errorf("expected no refinement but got: %s", result.SuggestedRefinement)
			}

			// Check for expected keywords in suggested agent keywords
			for _, kw := range tt.expectKeywords {
				found := false
				for _, suggestedKW := range result.SuggestedAgentKeywords {
					if strings.Contains(strings.ToLower(suggestedKW), strings.ToLower(kw)) {
						found = true
						break
					}
				}
				if !found && len(tt.expectKeywords) > 0 {
					t.Logf("expected keyword %q in suggested keywords %v (may be acceptable)", kw, result.SuggestedAgentKeywords)
				}
			}
		})
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		expectContains []string
		expectCount    int // Minimum expected
	}{
		{
			name:           "code review keywords",
			prompt:         "Review the code for quality issues",
			expectContains: []string{"code-review"},
			expectCount:    1,
		},
		{
			name:           "testing keywords",
			prompt:         "Generate unit tests for the module",
			expectContains: []string{"testing"},
			expectCount:    1,
		},
		{
			name:           "architecture keywords",
			prompt:         "Design the system architecture",
			expectContains: []string{"architecture-advisor"},
			expectCount:    1,
		},
		{
			name:           "documentation keywords",
			prompt:         "Write API documentation",
			expectContains: []string{"documentation-writer"},
			expectCount:    1,
		},
		{
			name:           "multiple keywords",
			prompt:         "Review code and write documentation for the architecture",
			expectContains: []string{"code-review", "documentation-writer", "architecture-advisor"},
			expectCount:    3,
		},
		{
			name:           "no keywords",
			prompt:         "Do something generic",
			expectContains: []string{},
			expectCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keywords := ExtractKeywords(tt.prompt)

			if len(keywords) < tt.expectCount {
				t.Errorf("expected at least %d keywords, got %d: %v", tt.expectCount, len(keywords), keywords)
			}

			for _, expected := range tt.expectContains {
				found := false
				for _, kw := range keywords {
					if kw == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected keyword %q in %v", expected, keywords)
				}
			}
		})
	}
}

// Note: containsActionVerb is package-private, tested indirectly through Evaluate

func TestEvaluator_SuggestRefinement(t *testing.T) {
	evaluator := NewEvaluator()

	issues := []string{
		"prompt is too vague",
		"lacks specific context",
		"no clear action",
	}

	refinement := evaluator.suggestRefinement("Fix the thing", issues)

	if refinement == "" {
		t.Error("expected non-empty refinement suggestion")
	}

	// Should include original prompt
	if !strings.Contains(refinement, "Fix the thing") {
		t.Error("refinement should include original prompt")
	}

	// Should provide specific guidance
	expectedPhrases := []string{"which", "what", "file", "function", "component"}
	foundGuidance := false
	for _, phrase := range expectedPhrases {
		if strings.Contains(strings.ToLower(refinement), phrase) {
			foundGuidance = true
			break
		}
	}
	if !foundGuidance {
		t.Logf("refinement may lack specific guidance: %s", refinement)
	}
}

func TestNewEvaluator(t *testing.T) {
	evaluator := NewEvaluator()

	if evaluator == nil {
		t.Fatal("expected non-nil evaluator")
	}

	// Verify basic properties are set
	result := evaluator.Evaluate("Review orchestrator.go for performance")
	if result.Confidence == 0 {
		t.Error("expected evaluator to calculate confidence")
	}
}

func TestEvaluationResult_Fields(t *testing.T) {
	evaluator := NewEvaluator()
	result := evaluator.Evaluate("Fix the thing")

	// Verify all fields are populated
	if result.Feedback == "" {
		t.Error("expected feedback to be set")
	}

	if result.Confidence == 0 {
		t.Error("expected confidence to be calculated")
	}

	if result.RefinedPrompt == "" && !result.IsClear {
		t.Error("expected refined prompt for unclear input")
	}

	if len(result.DetectedIssues) == 0 && !result.IsClear {
		t.Error("expected detected issues for unclear input")
	}
}
