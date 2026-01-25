# Support & Getting Help

## Documentation

Start here for comprehensive guides:

- **[README.md](README.md)** - Project overview and quick start
- **[docs/INSTALL.md](docs/INSTALL.md)** - Installation instructions
- **[docs/INDEX.md](docs/INDEX.md)** - Complete documentation index
- **[docs/WINDOWS_INSTALL.md](docs/WINDOWS_INSTALL.md)** - Windows setup guide
- **[docs/GITHUB_READY.md](docs/GITHUB_READY.md)** - Getting started guide

## Common Issues

### Build Fails

```bash
# Ensure Go 1.22+ is installed
go version

# Clean and rebuild
go clean -cache
go build -o programme
```

### Web Interface Not Working

```bash
# Verify web files exist
ls -la web/

# Check port 8080 availability
netstat -an | grep 8080

# Try different port
./programme web 9000
```

### Performance Issues

- Reduce network size: `--atoms 500` instead of 1000
- Use fast mode: `./programme generate --fast`
- Check available memory: `free -h` (Linux) or Task Manager (Windows)

## Getting Help

### Search First

1. Check existing GitHub issues
2. Review documentation in `/docs`
3. See Common Issues above

### Ask a Question

- Create a GitHub Discussion (preferred)
- Open a GitHub Issue with tag `question`
- Include:
  - IA-ATOMIQUE version
  - Go version
  - Operating system
  - Minimal reproduction steps
  - Error messages or logs

### Report a Bug

Create a GitHub Issue with:

- Clear title describing the bug
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version)
- Relevant logs or screenshots
- Link to related discussions (if any)

## Contributing

Found a solution to an issue? Consider contributing:

- See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines
- Follow conventional commit message format
- Include tests for bug fixes

## Community

- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: Questions and general discussion
- Pull Requests: Code contributions and improvements

## Escalation

For sensitive matters:

- Security issues: See [SECURITY.md](SECURITY.md)
- Code of Conduct violations: See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- Other concerns: Contact maintainers via email

---

**Thank you for using IA-ATOMIQUE!**
