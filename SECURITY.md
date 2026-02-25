# Security Policy

## Supported Scope

This repository is actively maintained on the `main` branch.

## Reporting a Vulnerability

Please do not open public issues for sensitive vulnerabilities.

Report privately by email to the repository owner with:

- affected component/path
- impact summary
- reproduction steps
- suggested mitigation (if available)

## Response Targets

- initial triage: within 5 business days
- status update: within 10 business days

## Security Baseline

- dependency checks in CI (`govulncheck`, npm audit)
- secret scanning in CI (`gitleaks`)
- request input limits and transport validation
