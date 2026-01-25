# Contributing to IA-ATOMIQUE

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the IA-ATOMIQUE project.

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Issues

If you find a bug or have a suggestion for improvement:

1. Check existing issues to avoid duplicates
2. Create a new issue with a clear title and description
3. Include steps to reproduce (for bugs)
4. Specify your Go version and OS
5. Attach relevant logs or outputs

### Making Changes

1. **Fork the repository** to `nosserb/IA-ATOMIQUE` and create a feature branch from `main`
2. **Write clean, idiomatic Go code** following [Effective Go](https://golang.org/doc/effective_go) guidelines
3. **Add tests** for new functionality
4. **Update documentation** if behavior changes
5. **Commit with conventional commit messages**:
   - `feat:` new feature
   - `fix:` bug fix
   - `docs:` documentation
   - `refactor:` code restructuring
   - `test:` test updates
   - `chore:` maintenance tasks

### Submitting Pull Requests

1. Ensure your branch is up-to-date with `main`
2. Run `go build` to verify compilation
3. Push your branch and create a PR
4. Include a clear description of changes
5. Reference related issues
6. Await review and address feedback

## Development Setup

```bash
# Install Go 1.22+
go version

# Clone and build
git clone https://github.com/yourusername/IA-ATOMIQUE.git
cd IA-ATOMIQUE
go build -o programme

# Run tests
go test ./...
```

## Project Structure

```
IA-ATOMIQUE/
├── internal/          # Core packages
│   ├── commands/      # CLI commands
│   └── tests/         # Test suites
├── docs/              # Documentation
├── web/               # Web interface
├── scripts/           # Build/utility scripts
├── main.go            # Application entry point
└── go.mod             # Module definition
```

## Code Standards

- Use `gofmt` for formatting
- Keep functions small and focused
- Write self-documenting code
- Include comments for complex logic
- Maintain >80% test coverage for critical paths

## Questions?

- Open a discussion in Issues
- Review existing documentation in `/docs`
- Check the main README.md for quick help
