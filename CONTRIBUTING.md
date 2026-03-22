# Contributing to sleepms

Thank you for your interest in contributing to sleepms! This document provides guidelines and standards for contributing to the project.

## Table of Contents

1. [Development Philosophy](#development-philosophy)
2. [Getting Started](#getting-started)
3. [Code Standards](#code-standards)
4. [Testing Requirements](#testing-requirements)
5. [Documentation Standards](#documentation-standards)
6. [Git Workflow](#git-workflow)
7. [Pull Request Process](#pull-request-process)
8. [Release Checklist](#release-checklist)

## Development Philosophy

### Core Principles

1. **Clarity Over Cleverness**
   - Write code that is easy to read and understand
   - Prefer explicit over implicit behavior
   - Use descriptive names for variables and functions

2. **Documentation First**
   - Every exported function must have godoc comments
   - Comments should explain "why", not just "what"
   - Keep documentation up-to-date with code changes

3. **Fail Fast**
   - Detect errors early and report them clearly
   - Never ignore errors or swallow them silently
   - Provide helpful error messages with context

4. **Cross-Platform Compatibility**
   - Code must work on Linux, macOS, and Windows
   - Test on multiple platforms before submitting
   - Avoid platform-specific assumptions

5. **Minimal Dependencies**
   - Prefer Go standard library when possible
   - Only add dependencies when absolutely necessary
   - Document why each dependency is needed

6. **Production Ready**
   - All code should be production-quality
   - Handle edge cases and error conditions
   - Consider performance and resource usage

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Code editor with Go support (VS Code, GoLand, etc.)

### Setting Up Development Environment

1. Clone the repository:
```bash
git clone <repository-url>
cd sleepms
```

2. Install development tools:
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/lint/golint@latest
```

3. Download dependencies:
```bash
go mod download
```

4. Build the project:
```bash
./build.sh  # Linux/macOS
build.cmd   # Windows
```

5. Test your setup:
```bash
./sleepms 1000 5000
```

## Code Standards

### Formatting

All Go code must be formatted with `gofmt`:

```bash
gofmt -s -w .
```

Use `goimports` for automatic import management:

```bash
goimports -w .
```

### Naming Conventions

- **Exported identifiers**: Use `MixedCaps` (e.g., `ParseArguments`, `AppName`)
- **Unexported identifiers**: Use `mixedCaps` (e.g., `parseArguments`, `appName`)
- **Constants**: Use `MixedCaps` or `SCREAMING_SNAKE_CASE` for constant groups
- **Acronyms**: Keep them uppercase (e.g., `HTTPServer`, not `HttpServer`)

**Good Examples:**
```go
func GenerateSleepDuration(minVal, maxVal int) int
const ProgressBarWidth = 40
var BuildNumber = "00000000"
```

**Bad Examples:**
```go
func generate_sleep_duration(min, max int) int  // Wrong: snake_case
const progress_bar_width = 40                    // Wrong: unexported constant
var build_number = "00000000"                    // Wrong: snake_case
```

### Code Organization

```
sleepms/
├── sleepms.go        # Main source file
├── go.mod            # Module definition
├── go.sum            # Dependency checksums
├── build.sh          # Unix build script
├── build.cmd         # Windows build script
├── README.md         # User documentation
├── CONTRIBUTING.md   # This file
├── VERSION_INFO.md   # Versioning documentation
├── BUILD_TEMPLATE.md # Build script guide
└── GODOC.md          # Auto-generated docs
```

### Comments and Documentation

#### Package-Level Comments

Every package must have a package-level comment:

```go
// Package main implements a sleep utility that waits for a random duration
// between a minimum and maximum value (in milliseconds).
//
// The program generates a random sleep duration, displays it to the user,
// and then sleeps for that duration.
package main
```

#### Function Comments

Every exported function must be documented:

```go
// GenerateSleepDuration generates a random sleep duration between min and max (inclusive).
// It uses the current time as a seed to ensure different values on each run.
//
// Parameters:
//   - minVal: The minimum sleep duration in milliseconds
//   - maxVal: The maximum sleep duration in milliseconds
//
// Returns:
//   - A random integer between minVal and maxVal (inclusive)
func GenerateSleepDuration(minVal, maxVal int) int {
    // Implementation
}
```

#### Inline Comments

Use inline comments to explain complex logic:

```go
// Generate a random number between minVal and maxVal (inclusive).
// rand.Intn(n) returns [0, n), so we adjust to get [minVal, maxVal]
randomNumber := minVal + rng.Intn(maxVal-minVal+1)
```

### Error Handling

#### Always Handle Errors

```go
// Good
oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
    return nil, err
}

// Bad
oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
```

#### Provide Context

Use `fmt.Errorf` with `%w` to wrap errors:

```go
// Good
if err != nil {
    return fmt.Errorf("failed to parse arguments: %w", err)
}

// Bad
if err != nil {
    return err  // Lost context
}
```

#### Error Messages

- Start with lowercase (unless beginning with proper noun)
- Don't end with punctuation
- Be specific and actionable

```go
// Good
return fmt.Errorf("minimum value cannot be greater than maximum value")

// Bad
return fmt.Errorf("Invalid input.")
```

### Concurrency

#### Goroutine Management

```go
// Good - proper goroutine with channel communication
keyPress := make(chan bool)
go listenForKeyPress(keyPress)

// Handle goroutine completion
select {
case <-keyPress:
    fmt.Println("Key pressed")
}
```

#### Channel Best Practices

- Close channels when no more data will be sent
- Use `defer` for cleanup operations
- Avoid sharing memory; communicate via channels

```go
// Good
cleanup := func() {
    term.Restore(int(os.Stdin.Fd()), oldState)
}
defer cleanup()
```

## Testing Requirements

### Writing Tests

Create test files with `_test.go` suffix:

```go
// sleepms_test.go
package main

import "testing"

func TestGenerateSleepDuration(t *testing.T) {
    tests := []struct {
        name    string
        minVal  int
        maxVal  int
        wantMin int
        wantMax int
    }{
        {"equal values", 1000, 1000, 1000, 1000},
        {"range 1-100", 1, 100, 1, 100},
        {"large range", 1000, 10000, 1000, 10000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := GenerateSleepDuration(tt.minVal, tt.maxVal)
            if got < tt.wantMin || got > tt.wantMax {
                t.Errorf("GenerateSleepDuration() = %v, want range [%v, %v]",
                    got, tt.wantMin, tt.wantMax)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...
```

### Coverage Requirements

- Aim for >80% code coverage
- Test both success and failure paths
- Include edge cases and boundary conditions

### Test Categories

1. **Unit Tests**: Test individual functions
2. **Integration Tests**: Test component interactions
3. **Edge Cases**: Test boundary conditions
4. **Error Cases**: Test error handling

## Documentation Standards

### Required Documentation Files

1. **README.md**: User-facing documentation
   - Features and use cases
   - Installation instructions
   - Usage examples
   - Troubleshooting guide

2. **CONTRIBUTING.md**: Development guidelines (this file)
   - Code standards
   - Testing requirements
   - Git workflow

3. **VERSION_INFO.md**: Versioning documentation
   - Version number format
   - Build number generation
   - Release process

4. **BUILD_TEMPLATE.md**: Build script guide
   - How to reuse build scripts
   - Configuration options
   - Customization guidelines

5. **GODOC.md**: Auto-generated API docs
   - Generated by build script
   - Package documentation
   - Function references

### Keeping Documentation Updated

- Update docs in the same commit as code changes
- Review all docs before releasing new versions
- Ensure examples in docs actually work
- Keep versioning info current

## Git Workflow

### Branch Naming

Use descriptive branch names with prefixes:

- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring

Examples:
```
feature/add-color-output
bugfix/fix-terminal-restore
docs/update-readme
refactor/simplify-progress-bar
```

### Commit Message Format

Follow the Conventional Commits specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Maintenance tasks

**Examples:**

```
feat(progress): add color to progress bar

Add ANSI color codes to make the progress bar more visually appealing.
Uses green for the filled portion and gray for empty space.
```

```
fix(terminal): restore terminal state on panic

Ensure terminal is always restored even if program panics.
Uses defer with recover to handle panic cases.

Fixes #123
```

```
docs(readme): add troubleshooting section

Add common issues and solutions based on user feedback.
```

### Commit Best Practices

1. **Atomic Commits**: One logical change per commit
2. **Clear Messages**: Describe what and why, not how
3. **Test Before Commit**: Ensure code builds and tests pass
4. **Small Commits**: Easier to review and revert if needed

## Pull Request Process

### Before Creating a PR

1. **Update your branch** from main:
```bash
git fetch origin
git rebase origin/main
```

2. **Run all checks**:
```bash
# Format code
gofmt -s -w .

# Run tests
go test ./...

# Build project
./build.sh
```

3. **Update documentation** if needed

4. **Review your changes**:
```bash
git diff main
```

### Creating the PR

1. Push your branch:
```bash
git push origin feature/your-feature
```

2. Create PR with clear description:

**Title:** Brief, descriptive title

**Description:**
```markdown
## Summary
Brief description of changes

## Changes Made
- Change 1
- Change 2
- Change 3

## Testing
How you tested the changes

## Related Issues
Fixes #123
```

### PR Review Process

**For Contributors:**
- Address review feedback promptly
- Keep discussions focused and professional
- Update PR based on feedback
- Squash commits if requested

**For Reviewers:**
- Review within 48 hours
- Be constructive and specific
- Test the changes locally
- Approve when all concerns are addressed

### Review Checklist

- [ ] Code follows style guidelines
- [ ] All functions are documented
- [ ] Tests are included and passing
- [ ] No obvious bugs or security issues
- [ ] Error handling is appropriate
- [ ] Code is readable and maintainable
- [ ] Documentation is updated
- [ ] Cross-platform compatibility considered

## Release Checklist

Before creating a release:

- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] No lint warnings (`golint ./...`)
- [ ] Documentation is updated
- [ ] Version numbers updated in build scripts
- [ ] Build scripts generate correct binaries
- [ ] Cross-platform builds tested
- [ ] GODOC.md regenerated
- [ ] README examples verified
- [ ] Git tag created with version number

### Version Number Updates

Update these files before release:

1. **build.sh**: Lines 16-18
```bash
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="0"
```

2. **build.cmd**: Lines 17-19
```cmd
set VERSION_MAJOR=1
set VERSION_MINOR=0
set VERSION_REVISION=0
```

### Creating a Release

```bash
# Update version in build scripts
# Build and test
./build.sh

# Commit version changes
git add build.sh build.cmd
git commit -m "chore: bump version to 1.1.0"

# Create and push tag
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0

# Build release binaries for all platforms
# Upload binaries to release page
```

## Questions or Problems?

- Check existing issues and pull requests
- Review documentation in README.md
- Open a new issue for bugs or feature requests
- Join discussions in existing issues

## Thank You!

Your contributions make this project better for everyone. We appreciate your time and effort!
