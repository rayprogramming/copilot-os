// Package cli provides GitHub Copilot CLI invocation capabilities for the
// CopilotOS server.
//
// This package handles:
//   - Agent Invocation: Execute Copilot CLI agents with prompts
//   - Timeout Management: Enforce timeouts on CLI operations
//   - Retry Logic: Automatically retry on transient failures
//   - Result Capture: Capture stdout, stderr, exit codes, and timing
//   - Error Handling: Distinguish between CLI errors and agent failures
//
// # GitHub Copilot CLI Integration
//
// This package wraps the GitHub Copilot CLI (`copilot` command) to invoke
// agents defined in the repository. Each invocation:
//  1. Constructs the CLI command with agent name and prompt
//  2. Executes the command with a timeout context
//  3. Captures the output (stdout/stderr)
//  4. Parses the result and returns structured data
//
// # Command Format
//
// The package executes commands in the format:
//
//	copilot --agent <agent-name> --prompt "<prompt-text>"
//
// The CLI must be installed and authenticated before use:
//
//	copilot auth login
//
// # Timeout Management
//
// Each invocation can have a timeout to prevent hanging operations:
//   - Default timeout: 300 seconds (5 minutes)
//   - Configurable via COPILOT_CLI_TIMEOUT environment variable
//   - Context-aware: respects parent context cancellation
//
// If a timeout occurs:
//   - The CLI process is terminated
//   - An error is returned with timeout details
//   - Partial output (if any) is captured
//
// # Retry Logic
//
// The invoker supports automatic retries for transient failures:
//   - Default retries: 1 (total of 2 attempts)
//   - Retries on: network errors, temporary CLI failures
//   - No retry on: invalid agent names, syntax errors, user cancellation
//
// # Result Structure
//
// Each invocation returns an InvocationResult with:
//   - Agent: Name of the invoked agent
//   - Success: Whether the invocation succeeded
//   - Output: Agent output as JSON (if applicable)
//   - Error: Error message (if failed)
//   - ExitCode: CLI process exit code
//   - Duration: Time taken to execute
//   - Timestamp: When the invocation started
//
// Usage Example
//
//	// Create an invoker
//	invoker := cli.NewInvoker(5*time.Minute, logger)
//
//	// Invoke an agent
//	result, err := invoker.InvokeAgent(ctx, "code-reviewer", "Review auth.go for security issues")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check result
//	if result.Success {
//	    fmt.Printf("Agent output: %s\n", result.Output)
//	} else {
//	    fmt.Printf("Agent failed: %s\n", result.Error)
//	}
//
// # Error Handling
//
// The package distinguishes between different error types:
//
// 1. CLI Execution Errors:
//   - Command not found (copilot CLI not installed)
//   - Permission denied
//   - System errors
//
// 2. Timeout Errors:
//   - Context deadline exceeded
//   - Operation took too long
//
// 3. Agent Errors:
//   - Invalid agent name
//   - Agent execution failed
//   - Agent returned error output
//
// All errors are wrapped with context for easier debugging.
//
// # Context Cancellation
//
// The package respects context cancellation throughout:
//   - Parent context cancellation stops CLI execution
//   - Timeout contexts are created per-invocation
//   - Graceful cleanup on cancellation
//
// # Performance Considerations
//
// Each CLI invocation spawns a new process, which has overhead:
//   - Process creation: ~10-50ms
//   - Copilot CLI startup: ~100-500ms
//   - Agent execution: Variable (seconds to minutes)
//
// For high-throughput scenarios, consider:
//   - Batching prompts when possible
//   - Using shorter timeouts for expected-fast agents
//   - Implementing result caching at a higher level
//
// # Thread Safety
//
// The Invoker type is safe for concurrent use. Multiple goroutines can
// invoke agents simultaneously without additional synchronization.
package cli
