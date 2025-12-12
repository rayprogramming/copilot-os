package agents

import "fmt"

// Agent represents a discovered agent with metadata.
type Agent struct {
	Name        string
	Description string
	Keywords    []string
}

// Registry holds discovered agents.
type Registry struct {
	agents map[string]*Agent
	order  []string // Maintain discovery order
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]*Agent),
		order:  []string{},
	}
}

// Add adds an agent to the registry.
func (r *Registry) Add(agent *Agent) error {
	if agent.Name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	if _, exists := r.agents[agent.Name]; exists {
		return fmt.Errorf("agent %q already registered", agent.Name)
	}
	r.agents[agent.Name] = agent
	r.order = append(r.order, agent.Name)
	return nil
}

// Get retrieves an agent by name.
func (r *Registry) Get(name string) *Agent {
	return r.agents[name]
}

// All returns all registered agents.
func (r *Registry) All() []*Agent {
	agents := make([]*Agent, len(r.order))
	for i, name := range r.order {
		agents[i] = r.agents[name]
	}
	return agents
}

// MatchKeywords finds agents matching the given keywords.
// Returns agents ranked by match score (highest first).
func (r *Registry) MatchKeywords(keywords []string) []*Agent {
	type scored struct {
		agent *Agent
		score float64
	}

	// Create a set of keywords for faster lookup
	keywordSet := make(map[string]bool)
	for _, kw := range keywords {
		keywordSet[kw] = true
	}

	// Score each agent
	scores := make([]scored, 0)
	for _, agent := range r.All() {
		score := calculateMatchScore(agent.Keywords, keywordSet)
		if score > 0 {
			scores = append(scores, scored{agent, score})
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Extract agents
	result := make([]*Agent, len(scores))
	for i, s := range scores {
		result[i] = s.agent
	}
	return result
}

// calculateMatchScore computes a match score between agent keywords and search keywords.
// Returns 0 if no match, otherwise returns score based on number and type of matches.
//
// Scoring Algorithm:
//   - Direct keyword match: +2.0 points
//   - No normalization: raw score returned
//
// The scoring is intentionally simple to provide predictable agent selection.
// Higher scores indicate better matches. Agents with score > 0 are returned,
// sorted by score descending.
//
// Example:
//
//	Agent keywords: ["code", "review", "quality"]
//	Search keywords: {"code": true, "review": true}
//	Score: 4.0 (2 matches × 2.0 points each)
//
// Future enhancements could include:
//   - Partial/fuzzy matching (e.g., "reviewing" matches "review")
//   - Weighted keywords (e.g., primary vs secondary capabilities)
//   - Normalized scores (0.0 to 1.0 range)
func calculateMatchScore(agentKeywords []string, searchKeywords map[string]bool) float64 {
	if len(agentKeywords) == 0 || len(searchKeywords) == 0 {
		return 0
	}

	score := 0.0
	for _, kw := range agentKeywords {
		if searchKeywords[kw] {
			// Exact match: highest weight
			// Using 2.0 instead of 1.0 to allow room for future partial matching
			score += 2.0
		}
	}

	return score
}
