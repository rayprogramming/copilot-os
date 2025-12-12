# CI/CD Pipeline Documentation

This document describes the complete continuous integration and continuous deployment (CI/CD) pipeline for the Copilot Agent Chain project.

## Overview

The CI/CD pipeline is built using GitHub Actions and provides:

- **Automated Testing**: Unit tests with race condition detection
- **Code Quality**: Linting and formatting checks
- **Security Analysis**: Vulnerability scanning and code analysis
- **Multi-Platform Builds**: Automated builds for Linux, macOS, and Windows
- **Release Management**: Automated release packaging and GitHub releases
- **Documentation**: Automated GitHub Pages deployment
- **Pre-Commit Checks**: Quality gates before code review

## Workflows

### 1. CI Pipeline (`ci.yml`)

**Trigger**: Push to `main`/`develop` or PR with changes to code/config

**Jobs**:
- **Lint**: Runs `golangci-lint` for code quality
- **Test**: Executes unit tests with coverage tracking
- **Build**: Multi-platform binary builds
- **Security Scan**: Vulnerability scanning with `gosec`
- **Notify**: Final status notification

**Outputs**:
- Build artifacts (binaries)
- Code coverage reports
- Security scan results

### 2. Release Pipeline (`release.yml`)

**Trigger**: Git tag matching `v*.*.*` pattern or manual workflow dispatch

**Jobs**:
- **Validate**: Pre-release verification
- **Build Release**: Multi-platform optimized builds
- **Create Release**: GitHub release with artifacts
- **Docker Build**: Container image creation
- **Notify**: Release completion notification

**Outputs**:
- GitHub Release with binaries and checksums
- Docker image artifacts

### 3. CodeQL Analysis (`codeql.yml`)

**Trigger**: Push to main/develop, PR, or weekly schedule

**Jobs**:
- **Analyze**: SAST analysis using CodeQL

**Outputs**:
- Security vulnerabilities in GitHub Security tab

### 4. Documentation (`docs.yml`)

**Trigger**: Push to `main` with changes to `docs/` or manual

**Jobs**:
- **Validate**: Link checking and config validation
- **Build and Deploy**: Jekyll site build and GitHub Pages deployment

**Outputs**:
- Live documentation on GitHub Pages

### 5. Pre-Commit Checks (`pre-commit.yml`)

**Trigger**: Pull request to main/develop

**Jobs**:
- **Pre-Commit**: Format, imports, vet, and mod tidy checks
- **Dependencies Check**: Vulnerability scanning

**Outputs**:
- Commit message validation
- Dependency security reports

## Local Development

### Using Makefile

The provided `Makefile` supports all CI/CD operations locally:

```bash
# Build binary for current platform
make build

# Build for all platforms
make build-all

# Run tests with coverage
make test

# Generate HTML coverage report
make test-coverage

# Run linters
make lint

# Format code
make fmt

# Build Docker image
make docker-build

# Run full CI pipeline
make ci

# Setup development environment
make dev-setup

# Clean build artifacts
make clean
```

### Setup Development Environment

```bash
# Install required tools
make dev-setup

# Or manually:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go mod tidy
```

## Building Binaries

### Local Build

```bash
# Current platform
make build

# All platforms
make build-all

# Manual for specific platform
GOOS=linux GOARCH=amd64 go build -o dist/copilot-agent-chain-linux-amd64 ./cmd/server
```

### Supported Platforms

| OS      | Architecture | Binary Name                              |
|---------|--------------|------------------------------------------|
| Linux   | x86_64       | `copilot-agent-chain-linux-amd64`       |
| Linux   | ARM64        | `copilot-agent-chain-linux-arm64`       |
| macOS   | Intel        | `copilot-agent-chain-darwin-amd64`      |
| macOS   | Apple Silicon| `copilot-agent-chain-darwin-arm64`      |
| Windows | x86_64       | `copilot-agent-chain-windows-amd64.exe` |

### Build Artifacts

- Location: `dist/` directory
- Binary file + `.sha256` checksum file per platform
- Total size varies by platform (~11MB each)

## Testing

### Running Tests

```bash
# Unit tests with race detection
go test -v -race ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Specific package
go test -v ./internal/orchestrator/...
```

### Coverage Requirements

- Minimum 70% code coverage
- All critical paths must have tests
- Race condition detection enabled

### Local Coverage Reports

```bash
make test-coverage
# Opens coverage.html in coverage/ directory
```

## Security

### Scanning Tools

1. **gosec**: Detects security vulnerabilities in Go code
2. **CodeQL**: SAST (Static Application Security Testing)
3. **Nancy**: Dependency vulnerability scanning
4. **Trivy**: Container image vulnerability scanning

### Running Locally

```bash
# Go security check
gosec ./...

# Dependency check
nancy sleuth

# CodeQL would need GitHub integration
```

## Docker Builds

### Build Image

```bash
# Using Makefile
make docker-build

# Manual
docker build -t copilot-agent-chain:latest .

# With version
docker build --build-arg VERSION=v1.0.0 -t copilot-agent-chain:v1.0.0 .
```

### Image Details

- **Base**: Alpine Linux (minimal size)
- **Non-root user**: App runs as user `app` (UID 1000)
- **Health check**: Included
- **Multi-stage build**: Optimized for size

### Running Container

```bash
docker run -it copilot-agent-chain:latest serve
```

## Releases

### Creating a Release

1. **Tag the commit**:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Validates the release
   - Builds all platform binaries
   - Creates GitHub Release with artifacts
   - Generates release notes
   - Builds Docker image

3. **Manual release**:
   - Trigger `release.yml` workflow manually
   - Specify version (e.g., `v1.0.0`)

### Release Contents

- Multi-platform binaries (5 variants)
- SHA256 checksums for each binary
- Docker image artifact
- Auto-generated release notes with commit history

### Verifying Release

```bash
# Verify binary signature
sha256sum -c copilot-agent-chain-linux-amd64.sha256

# Expected output
copilot-agent-chain-linux-amd64: OK
```

## Deployment Strategies

### Current Pipeline

1. **Push to main/develop** → CI pipeline (tests, lint, build)
2. **Open PR** → Pre-commit checks + PR reviews
3. **Tag release** → Full release pipeline
4. **Push to docs/** → GitHub Pages deployment

### Future Enhancements

Possible additions:
- Automated deployment to staging/production
- Container registry push (Docker Hub, GHCR)
- Slack/Discord notifications
- Performance benchmarking
- Load testing

## Monitoring and Alerts

### GitHub Actions Dashboard

- Visit: Actions tab in repository
- View: Workflow runs, logs, artifacts

### Status Checks

- **Required**: ci.yml, pre-commit.yml pass before merge
- **Recommended**: codeql.yml, docs.yml

### Artifacts

- **Retention**: 5-7 days for CI builds, 7 days for releases
- **Access**: Actions tab → Workflow → Run → Artifacts

## Troubleshooting

### Build Failures

1. Check workflow logs in Actions tab
2. Reproduce locally:
   ```bash
   make ci
   ```
3. Common issues:
   - Code formatting: `go fmt ./...`
   - Dependencies: `go mod tidy`
   - Tests failing: `go test -v ./...`

### Test Coverage Issues

```bash
# Check coverage
make test-coverage

# View report
open coverage/coverage.html
```

### Docker Build Issues

```bash
# Build locally
docker build -t copilot-agent-chain:test .

# Run image
docker run copilot-agent-chain:test serve
```

## Configuration Files

### `.github/workflows/`

- `ci.yml` - Build and test pipeline
- `release.yml` - Release automation
- `codeql.yml` - Security analysis
- `docs.yml` - Documentation deployment
- `pre-commit.yml` - Pre-commit quality gates

### `Makefile`

Multi-platform build automation and development tasks

### `Dockerfile`

Multi-stage production container build

### `.github/markdown-link-check-config.json`

Documentation link validation configuration

## Performance Metrics

### Build Times

- **Lint**: ~2-3 minutes
- **Test**: ~2-3 minutes
- **Build (single platform)**: ~1 minute
- **Full CI pipeline**: ~10 minutes
- **Release (all platforms)**: ~15 minutes

### Resource Usage

- **CI runs**: ~2GB memory, 2 CPUs
- **Docker image size**: ~50MB
- **Binary size**: ~11MB (stripped)

## Best Practices

1. **Always run locally before pushing**:
   ```bash
   make ci
   ```

2. **Follow commit message format**:
   ```
   feat(scope): description
   fix(scope): description
   docs: description
   ```

3. **Write tests for new features**:
   - Minimum 70% coverage
   - Include edge cases
   - Test error paths

4. **Keep dependencies updated**:
   ```bash
   go get -u ./...
   go mod tidy
   ```

5. **Monitor workflow runs**:
   - Check Actions tab regularly
   - Review security scans
   - Keep artifacts clean

## Support and Resources

- **GitHub Actions Documentation**: https://docs.github.com/actions
- **Go Build Documentation**: https://golang.org/cmd/go/
- **Docker Documentation**: https://docs.docker.com/
- **CodeQL Documentation**: https://codeql.github.com/

## Maintenance

### Regular Tasks

- Monthly: Review and update dependencies
- Quarterly: Audit security scanning rules
- Yearly: Review and optimize build times

### Updating Go Version

1. Update in workflows: `.github/workflows/*.yml`
2. Update in Dockerfile: `golang:1.24.3`
3. Update in Makefile comments
4. Test locally before pushing

## Next Steps

To enable additional features:

1. **Container Registry**: Add Docker Hub or GHCR push
2. **Slack Notifications**: Add workflow notifications
3. **Performance Monitoring**: Integrate benchmark tracking
4. **Automated Deployments**: Add production deployment steps
