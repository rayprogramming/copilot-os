// Package prompt provides prompt evaluation and refinement capabilities for the
// CopilotOS server.
//
// This package handles:
//   - Prompt Clarity Analysis: Detect vague or unclear prompts
//   - Automatic Refinement: Suggest or apply improvements to prompts
//   - Confidence Scoring: Rate prompt quality on a 0-1 scale
//   - Keyword Extraction: Extract relevant keywords for agent selection
//   - Issue Detection: Identify specific problems with prompts
//
// # Prompt Evaluation
//
// The Evaluator analyzes prompts using heuristics to determine clarity:
//  1. Length and word count checks
//  2. Vagueness pattern detection (e.g., "it", "thing", "stuff")
//  3. Specificity indicators (file names, function references, action verbs)
//  4. Context completeness
//
// A confidence score (0.0 to 1.0) is calculated:
//   - 0.0-0.5: Unclear, needs refinement
//   - 0.5-0.7: Somewhat clear, could be improved
//   - 0.7-1.0: Clear and specific
//
// # Evaluation Algorithm
//
// The evaluation algorithm works as follows:
//
//  1. Base Confidence: Start at 0.7
//  2. Apply Penalties:
//     - Too short: -0.2
//     - Too few words: -0.2
//     - Vague terms detected: -0.2
//  3. Apply Bonuses:
//     - Specific file mentioned: +0.1
//     - Function reference: +0.1
//     - Action verb used: +0.1
//  4. Clamp to [0.0, 1.0]
//
// Confidence scores are used by the orchestrator to decide whether to
// automatically refine the prompt before agent selection.
//
// # Vagueness Detection
//
// The package detects vague terms using regex patterns:
//   - Generic terms: "thing", "stuff", "something"
//   - Unclear references: "check it", "fix this", "review that"
//   - Missing context: short prompts without specifics
//
// When vagueness is detected, the evaluator can suggest refinements.
//
// # Automatic Refinement
//
// For unclear prompts, the evaluator can automatically refine them by:
//   - Expanding vague terms with context
//   - Adding specificity based on detected patterns
//   - Suggesting more actionable language
//
// Example refinements:
//   - "fix it" → "fix the authentication error in auth.go"
//   - "review this" → "review the user authentication logic"
//   - "check" → "analyze and identify potential issues"
//
// # Keyword Extraction
//
// The package extracts keywords from prompts for agent selection:
//  1. Tokenize prompt into words
//  2. Remove stop words (common words like "the", "a", "is")
//  3. Normalize to lowercase
//  4. Extract significant terms (nouns, verbs, technical terms)
//  5. Return list of keywords
//
// Extracted keywords are used by the agent registry to match relevant agents.
//
// Usage Example
//
//	// Create an evaluator
//	evaluator := prompt.NewEvaluator()
//
//	// Evaluate a prompt
//	result := evaluator.Evaluate("Review the authentication code")
//
//	// Check clarity
//	if !result.IsClear {
//	    fmt.Printf("Feedback: %s\n", result.Feedback)
//	    fmt.Printf("Suggested: %s\n", result.SuggestedRefinement)
//	}
//
//	// Use refined prompt
//	finalPrompt := result.RefinedPrompt
//
//	// Extract keywords
//	keywords := evaluator.ExtractKeywords(finalPrompt)
//	fmt.Printf("Keywords: %v\n", keywords)
//
// # EvaluationResult Structure
//
// Each evaluation returns a structured result:
//   - IsClear: Boolean indicating if prompt is clear enough
//   - Feedback: Human-readable explanation of issues
//   - SuggestedRefinement: Optional improved version
//   - Confidence: Score from 0.0 to 1.0
//   - DetectedIssues: List of specific problems
//   - RefinedPrompt: Automatically refined version
//   - SuggestedAgentKeywords: Keywords for agent matching
//
// # Configuration
//
// The evaluator can be configured with custom settings:
//   - Minimum prompt length
//   - Minimum word count
//   - Confidence threshold for clarity
//   - Custom vagueness patterns
//
// Currently, configuration is set via code. Future versions may support
// runtime configuration.
//
// # Best Practices
//
// When using this package:
//  1. Always evaluate user prompts before agent selection
//  2. Log evaluation results for debugging and improvement
//  3. Present refinement suggestions to users when appropriate
//  4. Use confidence scores to guide orchestration decisions
//  5. Continuously improve patterns based on real-world usage
//
// # Thread Safety
//
// The Evaluator type is safe for concurrent use. Multiple goroutines can
// evaluate prompts simultaneously without additional synchronization.
package prompt
