# Contributing Guide

Thank you for your interest in contributing to CopilotOS! This guide will help you get started.

**📖 Full Contributing Documentation**: For detailed guidelines, see [docs/guides/contributing.md](docs/guides/contributing.md)

## Quick Start for Contributors

### 1. Fork and Clone

```bash
# Visit https://github.com/rayprogramming/copilot-os
# Click "Fork" button
git clone https://github.com/YOUR_USERNAME/copilot-os.git
cd copilot-os
```

### 2. Set Up Development Environment

```bash
# Install dependencies
go mod download
go mod verify

# Authenticate with Copilot CLI
copilot auth login

# Run tests
go test ./...
```

### 3. Create a Feature Branch

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

### 4. Make Your Changes

- Write clear, documented code
- Add tests for new functionality
- Follow Go best practices and idioms
- Update documentation if needed

### 5. Test Your Changes

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Build the project
make build
```

### 6. Submit a Pull Request

```bash
git add .
git commit -m "feat: add awesome feature"
git push origin feature/your-feature-name
```

Then open a Pull Request on GitHub with:
- Clear description of changes
- Reference to any related issues
- Screenshots or examples if applicable

## Development Guidelines

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Write clear, self-documenting code
- Add comments for complex logic
- Use meaningful variable and function names

### Testing

- Write unit tests for new functions
- Maintain or improve test coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

### Documentation

- Update README.md for user-facing changes
- Update API documentation for new features
- Add code comments for complex algorithms
- Include examples in documentation

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`

**Examples**:
- `feat(agents): add new agent selection algorithm`
- `fix(cli): handle timeout errors gracefully`
- `docs: update installation instructions`

## Getting Help

- 📖 Read the documentation in the `docs/` folder
- 🐛 Report bugs via [GitHub Issues](https://github.com/rayprogramming/copilot-os/issues)
- 💬 Ask questions in discussions
- 📧 Contact maintainers for major changes

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions. We aim to foster a welcoming community for all contributors.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
