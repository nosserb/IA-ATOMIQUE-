# Security Policy

## Supported Versions

We maintain security updates for the following versions:

| Version | Supported          | EOL Date   |
|---------|--------------------|------------|
| 5.2.x   | Yes (current)      | 2027-01    |
| 5.1.x   | Yes                | 2026-06    |
| 5.0.x   | Limited            | 2026-03    |
| < 5.0   | No                 | -          |

## Reporting Vulnerabilities

**Do not** create public issues for security vulnerabilities.

Instead:

1. Email: `security@ia-atomique.local` with details:
   - Type of vulnerability
   - Location in code (file, line number)
   - Description and impact
   - Proof of concept (if applicable)

2. **Expected response time**: 24-48 hours
3. **Responsible disclosure**: We aim to provide a fix within 7 days
4. **Credit**: Your responsible disclosure will be acknowledged (if desired)

## Security Features

### Current Implementation

- AES-256 encryption for sensitive data
- Input validation and sanitization
- No external API dependencies (local processing)
- Memory-safe Go implementation
- Regular dependency updates

### Known Limitations

- Pattern database stored in plain text (not encrypted by default)
- Lexicon files accessible to process user
- No built-in authentication (assumes trusted environment)

### Recommendations

For production deployment:

1. Run in isolated network environment
2. Restrict file system access
3. Use process isolation (Docker/containers)
4. Monitor resource usage (CPU, memory)
5. Keep Go toolchain updated

## Security Updates Process

1. Vulnerability confirmed by security team
2. Fix developed in private branch
3. Testing on all supported versions
4. Security advisory prepared
5. Release published with advisory
6. Announcement to known downstream users

## Dependencies

All Go dependencies are managed through `go.mod` and kept up-to-date.
Run `go mod audit` to check for known vulnerabilities:

```bash
go list -json -m all | nancy sleuth
```

## Bug Bounty

Currently, we do not operate a formal bug bounty program. However, we greatly appreciate responsible disclosure of security issues.
