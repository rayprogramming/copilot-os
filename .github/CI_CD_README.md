# GitHub Actions CI/CD Pipeline

Complete automated build, test, and release pipeline for Copilot Agent Chain.

## 📋 Quick Links

- **[Full CI/CD Guide](../../CI_CD_IMPLEMENTATION.md)** - Detailed implementation documentation
- **[Pipeline Guide](../../docs/guides/cicd-pipeline.md)** - Complete usage guide
- **[Makefile Reference](../../Makefile)** - Local build automation
- **[GitHub Actions](workflows/)** - Workflow definitions

## 🚀 Workflows Overview

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **CI** | Push/PR | Build, test, lint, security scan |
| **Release** | Git tag | Multi-platform builds, GitHub Release |
| **CodeQL** | Push/PR/Schedule | SAST security analysis |
| **Docs** | Docs push | GitHub Pages deployment |
| **Pre-Commit** | PR | Code quality gates |

## 🛠️ Local Development

### Quick Commands

```bash
# Setup development environment
make dev-setup

# Run full CI locally
make ci

# Build all platforms
make build-all

# Run tests with coverage
make test-coverage

# Build Docker image
make docker-build

# See all commands
make help
```

### Build Outputs

Binaries are created in `dist/`:
- `copilot-agent-chain-linux-amd64`
- `copilot-agent-chain-linux-arm64`
- `copilot-agent-chain-darwin-amd64`
- `copilot-agent-chain-darwin-arm64`
- `copilot-agent-chain-windows-amd64.exe`

Each binary includes a `.sha256` checksum file.

## 📦 Creating Releases

1. Tag a commit with semantic version:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions automatically:
   - Validates the release
   - Builds all platform binaries
   - Creates GitHub Release with artifacts
   - Generates release notes
   - Builds Docker image

## 🔒 Security

The pipeline includes:
- **gosec** - Go security vulnerability scanning
- **CodeQL** - Static Application Security Testing
- **Nancy** - Dependency vulnerability scanning
- **Trivy** - Container image scanning
- **SARIF** - Upload to GitHub Security tab

## ✅ Quality Gates

All PRs must pass:
- ✅ Code formatting (gofmt)
- ✅ Import organization
- ✅ Module tidiness
- ✅ Unit tests (70%+ coverage)
- ✅ Linting (golangci-lint)
- ✅ Commit message format

## 📊 Workflow Files

### ci.yml
Runs on every push and PR:
- Linting with golangci-lint
- Unit tests with race detection
- Multi-platform binary builds
- Codecov coverage upload
- gosec security scanning

### release.yml
Runs on git tags (v*.*.* format):
- Validation pre-release
- Multi-platform optimized builds
- GitHub Release creation
- Docker image building
- SHA256 checksums

### codeql.yml
Runs on push, PR, and weekly schedule:
- CodeQL SAST analysis
- GitHub Security tab upload
- Automatic code scanning

### docs.yml
Runs when docs/ changes:
- Markdown link validation
- Jekyll site building
- GitHub Pages deployment

### pre-commit.yml
Runs on pull requests:
- Code formatting validation
- Import checking
- Go vet analysis
- Module tidiness
- Dependency scanning
- Commit message validation

## 🐳 Docker

Build and run the Docker image:

```bash
# Build
make docker-build

# Run
docker run -it copilot-agent-chain:latest serve

# With custom config
docker run -e LOG_LEVEL=debug copilot-agent-chain:latest serve
```

The Docker image:
- Uses Alpine Linux (minimal size ~50MB)
- Runs as non-root user for security
- Includes health check
- Multi-stage build for optimization

## 📚 File Structure

```
.github/
├── workflows/
│   ├── ci.yml              # Main CI pipeline
│   ├── release.yml         # Release automation
│   ├── codeql.yml          # Security analysis
│   ├── docs.yml            # Documentation
│   └── pre-commit.yml      # Pre-commit checks
├── CI_CD_README.md         # This file
└── markdown-link-check-config.json

Dockerfile                   # Container build
Makefile                     # Local automation
CI_CD_IMPLEMENTATION.md      # Implementation docs
docs/guides/cicd-pipeline.md # Complete guide
```

## 🎯 Typical Workflow

### Development
```bash
# Make changes
git checkout -b feature/my-feature
make ci                      # Verify locally
git commit -m "feat: add feature"
git push origin feature/my-feature
```

### Pre-Commit (Automated)
GitHub Actions runs:
- Code formatting checks
- Linting
- Unit tests
- Dependency scanning
- Commit message validation

### Merge
Once all checks pass:
- Reviewer approves
- PR gets merged to main
- CI pipeline runs full test suite

### Release
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# GitHub Actions automatically:
# - Builds for all platforms
# - Creates GitHub Release
# - Publishes Docker image
# - Generates release notes
```

## 📊 Performance

Typical build times:
- **Lint**: 2-3 minutes
- **Test**: 2-3 minutes
- **Build**: 1 minute per platform
- **Full CI**: ~10 minutes
- **Release (all platforms)**: ~15 minutes

## 🔧 Configuration

### Go Version
Update in:
- `.github/workflows/*.yml` - `go-version: '1.24.3'`
- `Dockerfile` - `FROM golang:1.24.3-alpine`

### Platform Support
Edit `release.yml` matrix to add/remove platforms.

### Coverage Threshold
Adjust in `Makefile` and workflow files.

## 🚨 Troubleshooting

### Build Failures
1. Check workflow logs in Actions tab
2. Run locally: `make ci`
3. Common fixes:
   - `go fmt ./...`
   - `go mod tidy`
   - `go test ./...`

### Test Coverage Issues
```bash
make test-coverage
open coverage/coverage.html
```

### Docker Build Issues
```bash
docker build -t copilot-agent-chain:test .
```

## 📖 More Information

- **Full Guide**: [CI_CD_IMPLEMENTATION.md](../../CI_CD_IMPLEMENTATION.md)
- **Pipeline Details**: [docs/guides/cicd-pipeline.md](../../docs/guides/cicd-pipeline.md)
- **Makefile Reference**: `make help`

## Support

For issues:
1. Check GitHub Actions logs
2. Review workflow files in `.github/workflows/`
3. Run `make ci` locally to reproduce
4. Check GitHub Security tab for vulnerabilities

---

**Last Updated**: December 7, 2025  
**Status**: ✅ Production Ready
