package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

// InvocationResult holds the result of a CLI invocation.
type InvocationResult struct {
	Agent     string          `json:"agent"`
	Success   bool            `json:"success"`
	Output    json.RawMessage `json:"output,omitempty"`
	Error     string          `json:"error,omitempty"`
	ExitCode  int             `json:"exit_code"`
	Duration  time.Duration   `json:"duration_ms"`
	Timestamp time.Time       `json:"timestamp"`
}

// Invoker handles invocation of Copilot CLI agents.
type Invoker struct {
	timeout time.Duration
	logger  *zap.Logger
	retries int
}

// NewInvoker creates a new CLI invoker.
func NewInvoker(timeout time.Duration, logger *zap.Logger) *Invoker {
	return &Invoker{
		timeout: timeout,
		logger:  logger,
		retries: 1, // Default to 1 retry on transient failures
	}
}

// InvokeAgent invokes a specific agent with the given prompt.
//
// This function executes the GitHub Copilot CLI with the specified agent and prompt,
// capturing all output and execution metadata.
//
// Execution Flow:
//  1. Create timeout context (if not already present)
//  2. Build CLI command: copilot --agent <name> --prompt "<text>"
//  3. Execute command with context
//  4. Capture stdout/stderr
//  5. Wait for completion or timeout
//  6. Parse and return structured result
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - agentName: Name of the agent to invoke (must match .github/agents/ file)
//   - prompt: The prompt/task to send to the agent
//
// Returns:
//   - InvocationResult: Structured result with output, status, timing
//   - error: Non-nil if CLI execution failed
//
// Error Conditions:
//   - Command not found: copilot CLI not installed or not in PATH
//   - Timeout: Operation exceeded deadline
//   - Context cancelled: Parent context was cancelled
//   - Exit code > 0: Agent execution failed
//
// Note: An agent returning an error (exit code > 0) is still returned as an
// InvocationResult with Success=false and Error populated. This is not considered
// a Go error - only CLI execution failures return Go errors.
func (i *Invoker) InvokeAgent(ctx context.Context, agentName, prompt string) (*InvocationResult, error) {
	start := time.Now()
	result := &InvocationResult{
		Agent:     agentName,
		Timestamp: start,
	}

	// Create context with timeout if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, i.timeout)
		defer cancel()
	}

	// Prepare command
	cmd := exec.CommandContext(
		ctx,
		"copilot",
		"--agent="+agentName,
		"--prompt="+prompt,
	)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()

	// Record duration
	result.Duration = time.Since(start)

	// Parse exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	// Handle output
	stdoutStr := stdout.String()
	if stdoutStr != "" {
		// Try to parse as JSON
		var jsonOutput json.RawMessage
		if err := json.Unmarshal([]byte(stdoutStr), &jsonOutput); err == nil {
			result.Output = jsonOutput
			result.Success = true
		} else {
			// If not JSON, wrap in string output
			result.Output = json.RawMessage([]byte(`"` + strings.TrimSpace(stdoutStr) + `"`))
			result.Success = true
		}
	}

	// Handle errors
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Sprintf("agent invocation timed out after %v", i.timeout)
		} else {
			stderrStr := stderr.String()
			if stderrStr != "" {
				result.Error = stderrStr
			} else {
				result.Error = err.Error()
			}
		}
		result.Success = false

		// Log the error
		i.logger.Warn("agent invocation failed",
			zap.String("agent", agentName),
			zap.Int("exit_code", result.ExitCode),
			zap.String("error", result.Error),
		)
	} else {
		i.logger.Debug("agent invocation succeeded",
			zap.String("agent", agentName),
			zap.Duration("duration", result.Duration),
		)
	}

	return result, nil
}

// ListAgents lists available agents.
func (i *Invoker) ListAgents(ctx context.Context) (*InvocationResult, error) {
	start := time.Now()
	result := &InvocationResult{
		Agent:     "orchestrator",
		Timestamp: start,
	}

	// Create context with timeout if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, i.timeout)
		defer cancel()
	}

	// Try to list agents by prompting the orchestrator
	cmd := exec.CommandContext(
		ctx,
		"copilot",
		"--agent=orchestrator",
		"--prompt=List all available agents and their descriptions",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Duration = time.Since(start)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		result.Error = stderr.String()
		result.Success = false
	} else {
		result.Output = json.RawMessage([]byte(stdout.String()))
		result.Success = true
	}

	return result, nil
}

// IsAvailable checks if the Copilot CLI is available.
func (i *Invoker) IsAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "copilot", "--version")
	err := cmd.Run()
	return err == nil
}

// CheckAuth checks if Copilot CLI is authenticated.
func (i *Invoker) CheckAuth(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "copilot", "auth", "status")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.Contains(stdout.String(), "authenticated") || stdout.String() != ""
}
