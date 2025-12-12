package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Discovery discovers and loads agents from the repository.
type Discovery struct {
	repoRoot string
	logger   *zap.Logger
	registry *Registry
}

// NewDiscovery creates a new agent discovery service.
func NewDiscovery(repoRoot string, logger *zap.Logger) *Discovery {
	return &Discovery{
		repoRoot: repoRoot,
		logger:   logger,
		registry: NewRegistry(),
	}
}

// Discover scans the repository for agents and populates the registry.
func (d *Discovery) Discover() error {
	agentsDir := filepath.Join(d.repoRoot, ".github", "agents")

	// Check if agents directory exists
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		d.logger.Warn("agents directory not found", zap.String("path", agentsDir))
		return nil
	}

	// Scan for .md files
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return fmt.Errorf("failed to read agents directory: %w", err)
	}

	discoveredCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			filePath := filepath.Join(agentsDir, entry.Name())
			agent, err := d.parseAgentFile(filePath)
			if err != nil {
				d.logger.Warn("failed to parse agent file", zap.String("file", entry.Name()), zap.Error(err))
				continue
			}

			if agent != nil {
				if err := d.registry.Add(agent); err != nil {
					d.logger.Warn("failed to add agent", zap.String("name", agent.Name), zap.Error(err))
				} else {
					discoveredCount++
					d.logger.Debug("discovered agent", zap.String("name", agent.Name))
				}
			}
		}
	}

	d.logger.Info("agent discovery complete", zap.Int("count", discoveredCount))
	return nil
}

// Registry returns the populated registry.
func (d *Discovery) Registry() *Registry {
	return d.registry
}

// parseAgentFile parses a Markdown agent file with YAML frontmatter.
func (d *Discovery) parseAgentFile(filePath string) (*Agent, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract YAML frontmatter (between --- delimiters)
	frontmatter, err := extractFrontmatter(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	if frontmatter == "" {
		return nil, fmt.Errorf("no frontmatter found")
	}

	// Parse YAML frontmatter
	agent := &Agent{
		Keywords: []string{},
	}

	// Simple YAML parsing (handles our use case)
	lines := strings.Split(frontmatter, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			agent.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "description:") {
			agent.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		} else if strings.HasPrefix(line, "keywords:") {
			// Parse keywords array [key1, key2, key3]
			keywordsStr := strings.TrimSpace(strings.TrimPrefix(line, "keywords:"))
			if strings.HasPrefix(keywordsStr, "[") && strings.HasSuffix(keywordsStr, "]") {
				keywordsStr = strings.TrimPrefix(keywordsStr, "[")
				keywordsStr = strings.TrimSuffix(keywordsStr, "]")
				parts := strings.Split(keywordsStr, ",")
				for _, part := range parts {
					kw := strings.TrimSpace(part)
					kw = strings.Trim(kw, "\"'")
					if kw != "" {
						agent.Keywords = append(agent.Keywords, kw)
					}
				}
			}
		}
	}

	// Validate required fields
	if agent.Name == "" {
		return nil, fmt.Errorf("agent name not found in frontmatter")
	}

	return agent, nil
}

// extractFrontmatter extracts YAML frontmatter from content.
// Expects content to start with ---, contain YAML, and end with ---.
//
// Regex Pattern Explanation:
//
//	^---\s*\n    - Start of string, three dashes, optional whitespace, newline
//	([\s\S]*?)  - Capture group: any character (including newlines), non-greedy
//	\n---       - Newline followed by three dashes (end delimiter)
//
// The pattern uses [\s\S]*? instead of .* because:
//   - \s matches whitespace (including \n)
//   - \S matches non-whitespace
//   - Together [\s\S] matches ANY character including newlines
//   - *? is non-greedy: stops at first occurrence of \n---
//
// Example input:
//
//	---
//	name: code-reviewer
//	description: Reviews code
//	---
//	# Agent instructions here
//
// Returns: "name: code-reviewer\ndescription: Reviews code\n"
func extractFrontmatter(content string) (string, error) {
	// Match frontmatter pattern: ---\n<content>\n---
	pattern := `^---\s*\n([\s\S]*?)\n---`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("frontmatter not found")
	}
	return matches[1], nil
}

// ExportAgentsJSON exports the registry as JSON (useful for debugging).
func (d *Discovery) ExportAgentsJSON() (string, error) {
	agents := d.registry.All()
	data, err := json.MarshalIndent(agents, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
