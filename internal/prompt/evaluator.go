package prompt

import (
	"regexp"
	"strings"
)

// EvaluationResult holds the result of prompt evaluation.
type EvaluationResult struct {
	IsClear                bool     `json:"is_clear"`
	Feedback               string   `json:"feedback"`
	SuggestedRefinement    string   `json:"suggested_refinement,omitempty"`
	Confidence             float64  `json:"confidence"`
	DetectedIssues         []string `json:"detected_issues"`
	RefinedPrompt          string   `json:"refined_prompt"`
	SuggestedAgentKeywords []string `json:"suggested_agent_keywords,omitempty"`
}

// Evaluator evaluates prompt clarity and suggests refinements.
type Evaluator struct {
	vaguenessPatterns   []*regexp.Regexp
	minimumLength       int
	minimumWords        int
	confidenceThreshold float64
}

// NewEvaluator creates a new prompt evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		vaguenessPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\b(thing|stuff|something)\b`),
			regexp.MustCompile(`\b(check|look|review|fix)\s+(it|this|that)\b`),
		},
		minimumLength:       10,
		minimumWords:        2,
		confidenceThreshold: 0.7,
	}
}

// Evaluate evaluates the clarity and quality of a prompt.
//
// This function analyzes a prompt using multiple heuristics to determine
// whether it's clear enough for effective agent execution.
//
// Evaluation Criteria:
//
//  1. Length Checks:
//     - Minimum length: 10 characters
//     - Minimum words: 2 words
//     - Penalty: -0.2 confidence per violation
//
//  2. Vagueness Detection:
//     - Checks for: "thing", "stuff", "something"
//     - Checks for: "it", "this", "that" without context
//     - Penalty: -0.2 confidence
//
//  3. Specificity Bonuses:
//     - File references (.go, /path/to): +0.1 confidence
//     - Function references (func, ()): +0.1 confidence
//     - Action verbs (review, test, etc.): +0.1 confidence
//
//  4. Confidence Score:
//     - Base: 0.7
//     - Range: 0.0 to 1.0
//     - Threshold: 0.7 for "clear"
//
// Returns:
//
//	EvaluationResult with:
//	- IsClear: true if confidence ≥ threshold
//	- Feedback: Human-readable explanation
//	- SuggestedRefinement: Improved version (if not clear)
//	- Confidence: Calculated score (0.0-1.0)
//	- DetectedIssues: List of specific problems
//	- RefinedPrompt: Refined version to use
//	- SuggestedAgentKeywords: Keywords for agent selection
//
// Example:
//
//	Input: "check it"
//	Output: IsClear=false, Confidence=0.3, Issues=["too short", "vague terms"]
//
//	Input: "Review auth.go for security vulnerabilities"
//	Output: IsClear=true, Confidence=0.9, Issues=[]
func (e *Evaluator) Evaluate(prompt string) EvaluationResult {
	result := EvaluationResult{
		DetectedIssues: []string{},
		RefinedPrompt:  prompt,
	}

	prompt = strings.TrimSpace(prompt)

	// Check for empty prompt
	if prompt == "" {
		result.IsClear = false
		result.Feedback = "Prompt is empty. Please provide a task description."
		result.Confidence = 0.0
		return result
	}

	// Start with base confidence
	result.Confidence = 0.7

	// Check length
	if len(prompt) < e.minimumLength {
		result.DetectedIssues = append(result.DetectedIssues, "Prompt is too short and lacks context")
		result.Confidence -= 0.2
	}

	// Check word count
	words := strings.Fields(prompt)
	if len(words) < e.minimumWords {
		result.DetectedIssues = append(result.DetectedIssues, "Prompt lacks sufficient detail")
		result.Confidence -= 0.2
	}

	// Check for vagueness
	for _, pattern := range e.vaguenessPatterns {
		if pattern.MatchString(strings.ToLower(prompt)) {
			result.DetectedIssues = append(result.DetectedIssues, "Prompt uses vague terms (e.g., 'it', 'this', 'module')")
			result.Confidence -= 0.2
			break
		}
	}

	// Check for specificity indicators (bonuses)
	hasSpecificFile := strings.Contains(prompt, ".go") || strings.Contains(prompt, "/")
	hasFunction := strings.Contains(prompt, "func") || strings.Contains(prompt, "()")
	hasActionVerb := containsActionVerb(prompt)

	if hasSpecificFile {
		result.Confidence += 0.1
	}
	if hasFunction {
		result.Confidence += 0.1
	}
	if hasActionVerb {
		result.Confidence += 0.1
	}

	// Ensure confidence is between 0 and 1
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}

	// Determine clarity
	result.IsClear = result.Confidence >= e.confidenceThreshold

	// Generate feedback
	if result.IsClear {
		result.Feedback = "Prompt is clear and actionable."
		result.RefinedPrompt = prompt // No refinement needed
	} else {
		result.Feedback = "Prompt could be improved for clarity and specificity."
		result.SuggestedRefinement = e.suggestRefinement(prompt, result.DetectedIssues)
		result.RefinedPrompt = result.SuggestedRefinement
	}

	// Extract suggested keywords for agent selection
	result.SuggestedAgentKeywords = ExtractKeywords(prompt)

	return result
}

// suggestRefinement generates a suggested refinement of the prompt.
func (e *Evaluator) suggestRefinement(prompt string, issues []string) string {
	refined := prompt

	// Add specificity if vague
	if len(issues) > 0 && strings.Contains(issues[0], "vague") {
		// Add common specificity hints
		if !strings.Contains(refined, " for ") && !strings.Contains(refined, " including ") {
			refined += " for correctness and best practices"
		}
	}

	// Add context if too short
	if len(issues) > 0 && strings.Contains(issues[0], "too short") {
		if !strings.Contains(refined, " including ") {
			refined += " including error handling and edge cases"
		}
	}

	return refined
}

// containsActionVerb checks if the prompt contains common action verbs.
//
// Action verbs indicate that the prompt is task-oriented and specific.
// Prompts with action verbs tend to be clearer and more actionable.
//
// Recognized action verbs include:
//   - Analysis: review, analyze, check, validate, verify
//   - Creation: generate, design, create, implement
//   - Improvement: refactor, improve, optimize, fix, debug
//   - Documentation: explain, document, describe
//
// The function performs case-insensitive substring matching.
//
// Examples:
//
//	"Review the authentication code" → true (contains "review")
//	"Check for bugs in auth.go" → true (contains "check")
//	"The code is broken" → false (no action verb)
//
// This is used as a positive indicator in prompt confidence scoring.
// Prompts with action verbs typically receive a +0.1 confidence boost.
func containsActionVerb(prompt string) bool {
	// List of common action verbs in development contexts
	actions := []string{
		"review", "analyze", "check", "test", "generate", "design", "create",
		"refactor", "improve", "optimize", "fix", "debug", "explain",
		"implement", "architect", "validate", "verify",
	}
	promptLower := strings.ToLower(prompt)
	for _, action := range actions {
		if strings.Contains(promptLower, action) {
			return true
		}
	}
	return false
}

// ExtractKeywords extracts potential keywords from the prompt for agent selection.
//
// Keyword Extraction Heuristics:
//
// This function uses domain-specific pattern matching to identify relevant
// keywords for agent selection. The algorithm:
//
//  1. Define domain patterns (regex) → agent keyword mappings
//  2. Lowercase the prompt for case-insensitive matching
//  3. Test each pattern against the prompt
//  4. Collect matching agent keywords
//  5. Remove duplicates to avoid over-weighting
//
// Domain Mappings:
//   - Code Quality: "code", "review", "bug", "fix" → ["code-review", "quality"]
//   - Testing: "test", "coverage", "mock" → ["test-generator", "testing"]
//   - Architecture: "design", "pattern", "scale" → ["architecture-advisor", "design"]
//   - Documentation: "doc", "readme", "guide" → ["documentation-writer", "docs"]
//
// The patterns use regex alternation (|) to match any of the terms.
// This is a simple but effective heuristic that provides good agent selection
// for common development tasks.
//
// Limitations:
//   - No stemming or lemmatization ("reviewing" won't match "review")
//   - No synonym expansion ("inspect" won't match "review")
//   - Fixed patterns (not learned from data)
//
// Future improvements could use:
//   - Natural language processing (NLP) for better term extraction
//   - Machine learning models for keyword classification
//   - User feedback to refine patterns
//
// Example:
//
//	Input: "Review the authentication code for security issues"
//	Matches: "review" → ["code-review", "quality"]
//	        "code" → ["code-review", "quality"] (duplicate)
//	Output: ["code-review", "quality"]
func ExtractKeywords(prompt string) []string {
	keywords := []string{}

	// Domain-specific pattern matching
	// Each pattern maps to agent capabilities
	domainKeywords := map[string][]string{
		"code|review|quality|bug|issue|fix|check|error|performance|refactor|correct": {"code-review", "quality"},
		"test|coverage|unit-test|mock|integration-test|edge-case":                    {"test-generator", "testing"},
		"architecture|design|pattern|structure|organize|scale|module|boundary":       {"architecture-advisor", "design"},
		"doc|readme|guide|comment|explain|write|api|tutorial":                        {"documentation-writer", "docs"},
	}

	promptLower := strings.ToLower(prompt)
	for pattern, kws := range domainKeywords {
		re := regexp.MustCompile(pattern)
		if re.MatchString(promptLower) {
			// Pattern matched: add associated keywords
			keywords = append(keywords, kws...)
		}
	}

	// Remove duplicates to avoid over-weighting certain agents
	// Multiple patterns may map to the same keywords
	seen := make(map[string]bool)
	unique := []string{}
	for _, kw := range keywords {
		if !seen[kw] {
			unique = append(unique, kw)
			seen[kw] = true
		}
	}

	return unique
}
