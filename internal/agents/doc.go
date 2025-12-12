// Package agents provides agent discovery, registration, and selection capabilities
// for the CopilotOS server.
//
// This package handles:
//   - Agent Discovery: Automatically discovers agent definitions from .github/agents/
//   - Agent Registry: Maintains a registry of discovered agents with metadata
//   - Agent Selection: Intelligently selects optimal agents based on keyword matching
//   - Scoring Algorithm: Ranks agents by relevance to user prompts
//
// # Agent Discovery
//
// The Discovery type scans the repository's .github/agents/ directory for Markdown
// files with YAML frontmatter. Each agent file should define:
//   - name: Agent identifier
//   - description: What the agent does
//   - keywords: Capabilities and domains (used for matching)
//
// Example agent file:
//
//	---
//	name: code-reviewer
//	description: Reviews code for quality, bugs, and best practices
//	keywords:
//	  - review
//	  - code quality
//	  - bugs
//	  - best practices
//	---
//	# Code Reviewer Agent Instructions
//	...
//
// # Agent Registry
//
// The Registry maintains discovered agents and provides methods for:
//   - Adding new agents
//   - Retrieving agents by name
//   - Getting all agents in discovery order
//   - Matching agents by keywords
//
// # Agent Selection
//
// Agent selection uses a keyword-based scoring algorithm:
//  1. Extract keywords from user prompt
//  2. Calculate match score for each agent's keywords
//  3. Rank agents by score (higher = better match)
//  4. Return top N agents for execution
//
// The scoring algorithm considers:
//   - Direct keyword matches (highest weight)
//   - Partial keyword matches (lower weight)
//   - Number of matching keywords
//   - Agent keyword coverage
//
// Usage Example
//
//	// Create discovery service
//	discovery := agents.NewDiscovery("/path/to/repo", logger)
//
//	// Discover agents
//	if err := discovery.Discover(); err != nil {
//	    return err
//	}
//
//	// Get registry
//	registry := discovery.Registry()
//
//	// Find agents matching keywords
//	keywords := []string{"code", "review", "quality"}
//	matchedAgents := registry.MatchKeywords(keywords)
//
// # Thread Safety
//
// The Registry is not thread-safe. If concurrent access is needed, external
// synchronization must be used (e.g., sync.RWMutex).
package agents
