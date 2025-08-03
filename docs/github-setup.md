# GitHub Repository Setup

This document outlines the recommended GitHub repository settings for the NFS Server Controller project.

## Repository Settings

### General Settings

**Repository Name**: `nfs-server-controller`
**Description**: "A Kubernetes operator for managing NFS servers as custom resources"
**Website**: `https://github.com/sharedvolume/nfs-server-controller`
**Topics**: `kubernetes`, `operator`, `nfs`, `storage`, `controller`, `kubebuilder`, `go`, `crd`

### Features

- ✅ **Issues**: Enable for bug reports and feature requests
- ✅ **Projects**: Enable for project management
- ✅ **Wiki**: Disable (use docs/ directory instead)
- ✅ **Discussions**: Enable for community Q&A
- ✅ **Sponsorships**: Enable if planning to accept sponsorships

### Pull Requests

- ✅ **Allow merge commits**: Enable
- ✅ **Allow squash merging**: Enable (recommended default)
- ✅ **Allow rebase merging**: Enable
- ✅ **Always suggest updating pull request branches**: Enable
- ✅ **Allow auto-merge**: Enable
- ✅ **Automatically delete head branches**: Enable

### Branch Protection Rules

#### Main Branch (`main`)

**Branch name pattern**: `main`

**Protect matching branches**:
- ✅ Require a pull request before merging
  - ✅ Require approvals: 1
  - ✅ Dismiss stale PR approvals when new commits are pushed
  - ✅ Require review from code owners
- ✅ Require status checks to pass before merging
  - ✅ Require branches to be up to date before merging
  - Required status checks:
    - `test`
    - `lint`
    - `build`
    - `e2e` (if stable)
- ✅ Require conversation resolution before merging
- ✅ Require signed commits
- ✅ Include administrators
- ✅ Restrict pushes that create files that have a path matching a pattern: `bin/**/*`

## Security Settings

### Security Features

- ✅ **Dependency graph**: Enable
- ✅ **Dependabot alerts**: Enable
- ✅ **Dependabot security updates**: Enable
- ✅ **Code scanning**: Enable with CodeQL
- ✅ **Secret scanning**: Enable
- ✅ **Secret scanning push protection**: Enable

### Dependabot Configuration

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  # Go modules
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    reviewers:
      - "maintainer-username"
    assignees:
      - "maintainer-username"
    commit-message:
      prefix: "deps"
      include: "scope"

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "maintainer-username"
    assignees:
      - "maintainer-username"
    commit-message:
      prefix: "ci"
      include: "scope"

  # Docker
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "maintainer-username"
    assignees:
      - "maintainer-username"
    commit-message:
      prefix: "docker"
      include: "scope"
```

## Actions and Secrets

### Required Secrets

**For CI/CD**:
- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub access token
- `GITHUB_TOKEN`: Automatically provided

**For Release**:
- `GPG_PRIVATE_KEY`: For signing releases (optional)
- `GPG_PASSPHRASE`: GPG key passphrase (optional)

### Environments

Create environments for deployment:

**Development**:
- No protection rules
- Used for feature branch testing

**Staging**:
- Required reviewers: maintainers
- Used for pre-release testing

**Production**:
- Required reviewers: 2+ maintainers
- Deployment branch rule: `main` only
- Used for stable releases

## Labels

### Issue Labels

**Type**:
- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements or additions to documentation
- `question` - Further information is requested
- `help wanted` - Extra attention is needed
- `good first issue` - Good for newcomers

**Priority**:
- `priority/critical` - Critical priority
- `priority/high` - High priority
- `priority/medium` - Medium priority
- `priority/low` - Low priority

**Component**:
- `area/api` - API related
- `area/controller` - Controller logic
- `area/documentation` - Documentation
- `area/testing` - Testing related
- `area/ci` - Continuous Integration

**Status**:
- `status/needs-review` - Needs review
- `status/in-progress` - Work in progress
- `status/blocked` - Blocked by external dependency
- `status/duplicate` - Duplicate issue
- `status/wontfix` - Won't fix

### Pull Request Labels

**Size**:
- `size/XS` - Extra small PR
- `size/S` - Small PR
- `size/M` - Medium PR
- `size/L` - Large PR
- `size/XL` - Extra large PR

**Type**:
- `type/feature` - New feature
- `type/bugfix` - Bug fix
- `type/refactor` - Code refactoring
- `type/docs` - Documentation changes
- `type/tests` - Test changes

## Milestones

Create milestones for releases:
- `v0.1.0` - Initial stable release
- `v0.2.0` - Feature enhancement release
- `v1.0.0` - Production ready release

## Projects

Create GitHub Projects for:
- **Roadmap**: High-level feature planning
- **Release Planning**: Specific release milestones
- **Bug Triage**: Bug tracking and prioritization

## Code Owners

Create `.github/CODEOWNERS`:

```
# Global owners
* @maintainer-username

# Documentation
*.md @maintainer-username
/docs/ @maintainer-username

# API changes
/api/ @maintainer-username @api-reviewer

# Controller code
/internal/controller/ @maintainer-username @controller-reviewer

# CI/CD
/.github/ @maintainer-username @devops-reviewer
/Makefile @maintainer-username @devops-reviewer
/Dockerfile @maintainer-username @devops-reviewer

# Security
/SECURITY.md @maintainer-username @security-reviewer
```

## Release Process

### Pre-release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in appropriate files
- [ ] Security scan completed
- [ ] Performance testing completed

### Release Steps

1. Create release branch from `main`
2. Update version numbers
3. Update CHANGELOG.md
4. Create pull request for release
5. After merge, create git tag
6. GitHub Actions automatically builds and publishes
7. Create GitHub release with notes
8. Announce release in discussions

## Repository Maintenance

### Regular Tasks

**Weekly**:
- Review and merge Dependabot PRs
- Triage new issues
- Review open PRs

**Monthly**:
- Update documentation
- Review and update labels
- Clean up stale branches
- Security review

**Quarterly**:
- Review repository settings
- Update contributing guidelines
- Performance and security audit
- Roadmap planning

### Automation

**Stale Issues/PRs**:
```yaml
# .github/stale.yml
daysUntilStale: 60
daysUntilClose: 7
staleLabel: 'status/stale'
markComment: >
  This issue has been automatically marked as stale because it has not had
  recent activity. It will be closed if no further activity occurs.
closeComment: >
  This issue has been automatically closed due to inactivity.
```

## Community Guidelines

### README.md Badge Requirements

```markdown
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/sharedvolume/nfs-server-controller)](https://goreportcard.com/report/github.com/sharedvolume/nfs-server-controller)
[![codecov](https://codecov.io/gh/sharedvolume/nfs-server-controller/branch/main/graph/badge.svg)](https://codecov.io/gh/sharedvolume/nfs-server-controller)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/sharedvolume/nfs-server-controller/badge)](https://securityscorecards.dev/viewer/?uri=github.com/sharedvolume/nfs-server-controller)
```

### Links and References

- **Documentation**: Link to GitHub Pages or docs site
- **Issues**: Link to issue templates
- **Discussions**: Link to community guidelines
- **Security**: Link to security policy
- **Contributing**: Link to contributing guidelines

This configuration ensures a professional, secure, and well-maintained open source project.
