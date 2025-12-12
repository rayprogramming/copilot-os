package agents

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestDiscovery_ParseAgentFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		expectName  string
		expectDesc  string
		expectKW    []string
	}{
		{
			name: "valid agent file",
			content: `---
name: test-agent
description: A test agent
keywords: [testing, example, demo]
---

# Test Agent

This is test content.
`,
			expectError: false,
			expectName:  "test-agent",
			expectDesc:  "A test agent",
			expectKW:    []string{"testing", "example", "demo"},
		},
		{
			name: "missing frontmatter",
			content: `# Test Agent

No frontmatter here.
`,
			expectError: true,
		},

		{
			name: "missing required fields",
			content: `---
name: test-agent
---
`,
			expectError: false,
			expectName:  "test-agent",
			expectDesc:  "",
			expectKW:    nil,
		},
	}

	logger := zap.NewNop()
	discovery := NewDiscovery(".", logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test-agent.md")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			agent, err := discovery.parseAgentFile(tmpFile)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if agent.Name != tt.expectName {
				t.Errorf("expected name %q, got %q", tt.expectName, agent.Name)
			}

			if agent.Description != tt.expectDesc {
				t.Errorf("expected description %q, got %q", tt.expectDesc, agent.Description)
			}

			if len(agent.Keywords) != len(tt.expectKW) {
				t.Errorf("expected %d keywords, got %d", len(tt.expectKW), len(agent.Keywords))
			}

			for i, kw := range tt.expectKW {
				if i >= len(agent.Keywords) {
					break
				}
				if agent.Keywords[i] != kw {
					t.Errorf("expected keyword %q at position %d, got %q", kw, i, agent.Keywords[i])
				}
			}
		})
	}
}

func TestDiscovery_Discover(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, ".github", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test agent files
	agents := []struct {
		filename string
		content  string
	}{
		{
			filename: "agent1.md",
			content: `---
name: agent1
description: First agent
keywords: [one, first]
---
`,
		},
		{
			filename: "agent2.md",
			content: `---
name: agent2
description: Second agent
keywords: [two, second]
---
`,
		},
		{
			filename: "invalid.md",
			content:  "No frontmatter",
		},
	}

	for _, a := range agents {
		path := filepath.Join(agentsDir, a.filename)
		if err := os.WriteFile(path, []byte(a.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	logger := zap.NewNop()
	discovery := NewDiscovery(tmpDir, logger)

	err := discovery.Discover()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	registry := discovery.Registry()
	all := registry.All()

	// Should have discovered 2 valid agents (invalid.md should be skipped)
	if len(all) != 2 {
		t.Errorf("expected 2 agents, got %d", len(all))
	}

	// Verify agents are in registry
	agent1 := registry.Get("agent1")
	if agent1 == nil {
		t.Error("expected to find agent1")
	}

	agent2 := registry.Get("agent2")
	if agent2 == nil {
		t.Error("expected to find agent2")
	}
}

func TestDiscovery_Discover_NoAgentsDir(t *testing.T) {
	tmpDir := t.TempDir()

	logger := zap.NewNop()
	discovery := NewDiscovery(tmpDir, logger)

	err := discovery.Discover()
	// Should not error if directory doesn't exist, just log warning
	if err != nil {
		t.Errorf("unexpected error when agents dir missing: %v", err)
	}

	registry := discovery.Registry()
	if len(registry.All()) != 0 {
		t.Error("expected empty registry when no agents dir")
	}
}

func TestDiscovery_ExportAgentsJSON(t *testing.T) {
	logger := zap.NewNop()
	discovery := NewDiscovery(".", logger)

	registry := discovery.Registry()
	registry.Add(&Agent{
		Name:        "test",
		Description: "Test agent",
		Keywords:    []string{"test"},
	})

	json, err := discovery.ExportAgentsJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(json) == 0 {
		t.Error("expected non-empty JSON export")
	}

	// Basic validation that it's valid JSON
	if json[0] != '[' {
		t.Error("expected JSON to start with '['")
	}
}

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectYAML  string
		expectError bool
	}{
		{
			name: "valid frontmatter",
			content: `---
name: test
---
Content`,
			expectYAML:  "name: test",
			expectError: false,
		},
		{
			name: "no frontmatter",
			content: `# Heading
Content`,
			expectYAML:  "",
			expectError: true,
		},
		{
			name: "incomplete frontmatter",
			content: `---
name: test
Content without closing`,
			expectYAML:  "",
			expectError: true,
		},
		{
			name: "frontmatter with extra content",
			content: `---
name: test
description: multi
  line
---

# More content
`,
			expectYAML:  "name: test\ndescription: multi\n  line",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yaml, err := extractFrontmatter(tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if yaml != tt.expectYAML {
				t.Errorf("expected YAML:\n%s\ngot:\n%s", tt.expectYAML, yaml)
			}
		})
	}
}
