---
layout: default
title: Contributing Guide
---

# Contributing Guide

Thank you for your interest in contributing to CoPilot OS! This guide will help you get started.

## Getting Started

### 1. Fork the Repository

```bash
# Visit https://github.com/rayprogramming/copilot-agent-chain
# Click "Fork" button
git clone https://github.com/YOUR_USERNAME/copilot-agent-chain.git
cd copilot-agent-chain
```

### 2. Set Up Development Environment

```bash
# Install dependencies
go mod download
go mod verify

# Authenticate with Copilot CLI
copilot auth login

# Verify setup
go test ./...
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

**Branch Naming Convention**:
- `feature/` — New features
- `fix/` — Bug fixes
- `docs/` — Documentation updates
- `refactor/` — Code refactoring
- `test/` — Test improvements
- `perf/` — Performance improvements

### 2. Make Your Changes

Follow Go best practices:

```go
// ✓ Good: Clear, documented code
// Evaluate checks prompt clarity using heuristic rules
func (e *Evaluator) Evaluate(prompt string) *EvaluationResult {
    // Implementation...
}

// ✗ Poor: Unclear and undocumented
func (e *Evaluator) eval(p string) *Result {
    // Code...
}
```

**Code Guidelines**:
- Write clear, self-documenting code
- Add comments for non-obvious logic
- Keep functions small and focused
- Handle errors explicitly
- Use meaningful variable names

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestName -v ./path/to/package

# Check for race conditions
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

**Testing Guidelines**:
- Test normal cases, edge cases, and errors
- Use table-driven tests (Go idiom)
- Aim for at least 80% coverage
- Test error paths thoroughly

### 4. Use Development Agents for Review

Before committing, use the agents to review your code:

```bash
# Code quality review
copilot --agent=code-reviewer --prompt "Review my changes in [package] for correctness, performance, and Go idioms"

# Architecture review
copilot --agent=architecture-advisor --prompt "Does my implementation align with the overall architecture? Any design issues?"

# Test coverage
copilot --agent=test-generator --prompt "Review my tests in [package]. Are there untested code paths or missing edge cases?"

# Documentation
copilot --agent=documentation-writer --prompt "Review my changes. Do they need comments or documentation updates?"
```

### 5. Format and Lint Code

```bash
# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run

# Vet code
go vet ./...
```

**Installing golangci-lint**:
```bash
# macOS
brew install golangci-lint

# Or from source
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 6. Commit Changes

```bash
git add .
git commit -m "feat: add result caching layer"
```

**Commit Message Format**:
```
type: subject (50 characters max)

Body (if needed):
- Explain what changed
- Explain why it changed
- Reference issue numbers: Fixes #123
```

**Types**:
- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `test:` — Test improvements
- `refactor:` — Code restructuring
- `perf:` — Performance improvement
- `chore:` — Build, dependencies, etc.

### 7. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create PR on GitHub
# https://github.com/YOUR_USERNAME/copilot-agent-chain
# Click "New Pull Request"
```

**PR Description Template**:
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] New feature
- [ ] Bug fix
- [ ] Documentation
- [ ] Performance improvement

## Testing
- [ ] Added tests
- [ ] All tests pass
- [ ] Tested manually with command: ...

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] No new warnings generated

## Closes
Fixes #(issue number)
```

## Types of Contributions

### 1. Bug Reports

Report bugs through GitHub Issues:

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
1. Go to '...'
2. Click on '...'
3. Error appears

**Expected behavior**
What you expected to happen.

**Actual behavior**
What actually happened.

**Environment**
- OS: [e.g. Linux 5.10]
- Go version: [e.g. 1.21]
- Copilot CLI version: [e.g. 2.0.1]

**Logs**
Relevant error logs (with LOG_LEVEL=debug)
```

### 2. Feature Requests

Suggest features through GitHub Discussions:

```markdown
**Description**
What feature would you like to see?

**Why**
Why is this important?

**Proposed Implementation**
How would you implement this?

**Alternatives**
Other approaches or features that would solve the problem?
```

### 3. Documentation Improvements

Improve documentation by:

1. **Clarifying existing docs**:
   - Fix typos
   - Improve examples
   - Add missing details

2. **Adding new docs**:
   - Document new features
   - Create how-to guides
   - Add troubleshooting tips

3. **Improving code comments**:
   - Explain complex logic
   - Document function behavior
   - Add usage examples

### 4. Code Improvements

Contribute code by:

1. **Adding new features**:
   - Follow the development workflow above
   - Include tests and documentation
   - Get architecture review before major changes

2. **Fixing bugs**:
   - Include test for the bug
   - Verify fix doesn't break other tests
   - Reference issue in commit message

3. **Performance improvements**:
   - Include benchmarks showing improvement
   - Ensure correctness is maintained
   - Document trade-offs

4. **Refactoring**:
   - Improve code clarity or maintainability
   - Don't change functionality
   - All tests must pass

## Code Style Guide

### Naming

```go
// ✓ Good: Clear, descriptive names
type PromptEvaluator struct { ... }
func (e *PromptEvaluator) Evaluate(prompt string) *EvaluationResult { ... }

// ✗ Poor: Abbreviated or unclear names
type PE struct { ... }
func (e *PE) E(p string) *ER { ... }
```

### Comments

```go
// ✓ Good: Explains WHY, not WHAT
// Confidence starts at 0.7 to account for typical prompt clarity.
// Lower starting values made the evaluator too strict on normal prompts.
confidence := 0.7

// ✗ Poor: Explains obvious WHAT
// Set confidence to 0.7
confidence := 0.7
```

### Error Handling

```go
// ✓ Good: Explicit error handling
result, err := discoverAgents(repoRoot)
if err != nil {
    return nil, fmt.Errorf("failed to discover agents: %w", err)
}

// ✗ Poor: Silent error
result, _ := discoverAgents(repoRoot)
```

### Function Size

```go
// ✓ Good: Small, focused functions
func selectAgents(keywords []string, registry *Registry) []*Agent {
    matches := registry.MatchKeywords(keywords, 3)
    return matches
}

// ✗ Poor: Large, doing multiple things
func orchestrate(prompt string) ... {
    // 200 lines of code doing: evaluate, select, execute, synthesize
}
```

## Testing Guidelines

### Write Table-Driven Tests

```go
// ✓ Good: Table-driven tests (Go idiom)
func TestEvaluator_Evaluate(t *testing.T) {
    tests := []struct {
        name     string
        prompt   string
        expected bool
    }{
        {"clear prompt", "Review foo.go for performance", true},
        {"vague prompt", "Review the thing", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := evaluator.Evaluate(tt.prompt)
            if result.IsClear != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result.IsClear)
            }
        })
    }
}

// ✗ Poor: Multiple separate tests
func TestEvaluator_ClearPrompt(t *testing.T) { ... }
func TestEvaluator_VaguePrompt(t *testing.T) { ... }
```

### Test Edge Cases

```go
// ✓ Good: Tests edge cases
func TestPromptEvaluator_Evaluate(t *testing.T) {
    tests := []struct {
        name   string
        prompt string
        // ... test cases
    }{
        {"empty prompt", "", ...},
        {"very long prompt", strings.Repeat("a", 10000), ...},
        {"special characters", "Review @#$%^&*()", ...},
        {"non-English", "查看代码", ...},
    }
}
```

## Documentation Style

### Code Comments

```go
// Excellent: Explains WHY, not WHAT
// We start confidence at 0.7 rather than 0.5 because empirically,
// most naturally-phrased prompts without explicit detail markers
// are still actionable. This threshold accounts for:
// - Implicit domain context from file structure
// - Common developer patterns and language
// - The cost of over-refining (user frustration)
// Lower thresholds (e.g., 0.5) made the evaluator reject too many valid prompts.
confidence := 0.7
```

### Documentation Files

```markdown
# Descriptive Heading

Clear opening paragraph explaining what this section covers.

## Subsection

Detailed explanation with:
- Bullet points for clarity
- Code examples where helpful
- Links to related topics

### Examples

Practical, runnable examples with expected output.

## See Also

- [Related Guide](link)
- [Reference](link)
```

## Review Process

### What Reviewers Look For

1. **Code Quality**
   - Follows Go idioms
   - Clear and readable
   - Proper error handling
   - Good test coverage

2. **Functionality**
   - Solves the stated problem
   - Doesn't break existing features
   - Handles edge cases
   - Performance acceptable

3. **Documentation**
   - Code comments explain non-obvious logic
   - Changes documented
   - Examples provided

4. **Testing**
   - Unit tests included
   - Integration tests where needed
   - All tests pass
   - Coverage maintained

### Responding to Reviews

1. **Acknowledge feedback** — Thank reviewers for their input
2. **Ask clarifying questions** — If you don't understand feedback
3. **Make requested changes** — Or explain why you disagree
4. **Push updates** — Reviewer will re-review
5. **Iterate** — Until approval

## Running the Full Test Suite

```bash
# Complete test run with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Check coverage by package
go tool cover -func=coverage.out
```

## Debugging Tips

### Use the Development Agents

```bash
# Ask code-reviewer for feedback on your implementation
copilot --agent=code-reviewer --prompt "Review my implementation in internal/newpackage/main.go. Any issues or improvements?"

# Ask architecture-advisor for design feedback
copilot --agent=architecture-advisor --prompt "Is my approach to [problem] aligned with the overall architecture?"
```

### Enable Debug Logging

```bash
export LOG_LEVEL=debug
export REPO_ROOT=$(pwd)
./copilot-agent-chain &

# Make your requests...
# Look for detailed logs explaining behavior
```

### Run Single Test

```bash
# Run specific test
go test -run TestName -v ./path/to/package

# Run tests matching pattern
go test -run "TestEvaluator" -v ./internal/prompt
```

## Community Guidelines

### Be Respectful
- Treat all contributors with respect
- Welcome diverse perspectives
- Provide constructive feedback
- Assume good intent

### Be Constructive
- Provide specific, actionable feedback
- Offer solutions, not just criticism
- Ask questions to understand intent
- Help each other improve

### Be Patient
- Review takes time
- Maintainers are volunteers
- Not all features will be accepted
- Feedback loops take iterations

## Questions?

- 📖 See the [Development Guide](development.md)
- 🐛 Check [Troubleshooting](troubleshooting.md)
- 💬 Open a [GitHub Discussion](https://github.com/rayprogramming/copilot-agent-chain/discussions)
- 📧 Reach out to maintainers

---

**Thank you for contributing!**
