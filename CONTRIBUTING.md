# Contributing to NFS Server Controller

Thank you for your interest in contributing to the NFS Server Controller! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and inclusive in all interactions.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker
- kubectl
- Kind (for local testing)
- Git

### Development Environment Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/nfs-server-controller.git
   cd nfs-server-controller
   ```

2. **Set up your development environment:**
   ```bash
   # Install dependencies
   go mod download
   
   # Install development tools
   make controller-gen kustomize
   ```

3. **Verify your setup:**
   ```bash
   make test
   ```

## Development Workflow

### Making Changes

1. **Create a new branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes following our coding standards:**
   - Write clear, concise commit messages
   - Add tests for new functionality
   - Update documentation as needed
   - Follow Go coding conventions

3. **Run tests locally:**
   ```bash
   # Run unit tests
   make test
   
   # Run linting
   make lint
   
   # Run end-to-end tests
   make test-e2e
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that don't affect code meaning (white-space, formatting, etc.)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to build process or auxiliary tools

Examples:
```
feat: add support for custom NFS export options
fix: resolve PVC binding issue in multi-zone clusters
docs: update installation instructions
test: add unit tests for storage validation
```

## Testing

### Unit Tests

Run unit tests:
```bash
make test
```

Write tests for:
- New functionality
- Bug fixes
- Edge cases
- Error conditions

### End-to-End Tests

Run e2e tests:
```bash
make test-e2e
```

E2e tests validate:
- Complete workflow from resource creation to deletion
- Integration with Kubernetes APIs
- Real-world scenarios

### Test Guidelines

- Tests should be deterministic and repeatable
- Use table-driven tests when appropriate
- Mock external dependencies
- Test both success and failure scenarios
- Ensure good test coverage

## Code Style

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and types

### Linting

Run the linter:
```bash
make lint
```

Fix linting issues:
```bash
make lint-fix
```

## Documentation

### Code Documentation

- Add comments for all exported functions, types, and constants
- Use complete sentences in comments
- Include examples in comments when helpful

### User Documentation

- Update README.md for user-facing changes
- Add or update examples in the `docs/` directory
- Update API documentation for CRD changes

## Submitting Changes

### Pull Request Process

1. **Push your branch to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a pull request:**
   - Use a clear, descriptive title
   - Fill out the pull request template
   - Link related issues
   - Add labels as appropriate

3. **Address review feedback:**
   - Make requested changes
   - Push updates to your branch
   - Respond to reviewer comments

### Pull Request Guidelines

- Keep PRs focused and atomic
- Include tests for new functionality
- Update documentation as needed
- Ensure CI passes
- Rebase your branch if needed

### Review Process

- All PRs require at least one review
- Maintainers will review PRs in a timely manner
- Address feedback constructively
- Be patient and respectful

## Release Process

Releases are managed by maintainers:

1. Version tags follow semantic versioning (v1.2.3)
2. GitHub Actions automatically builds and publishes releases
3. Release notes are generated from commit messages
4. Docker images are pushed to the registry

## Getting Help

- **Questions**: Open a [GitHub Discussion](https://github.com/sharedvolume/nfs-server-controller/discussions)
- **Bugs**: Report on [GitHub Issues](https://github.com/sharedvolume/nfs-server-controller/issues)
- **Security**: Email bilgehan.nal@gmail.com privately

## Recognition

Contributors are recognized in:
- Release notes
- README acknowledgments
- Git commit history

Thank you for contributing to the NFS Server Controller!
