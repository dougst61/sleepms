# Build Script Template Guide

This document explains how to reuse the sleepms build scripts (`build.sh` and `build.cmd`) in your own Go projects.

## Table of Contents

1. [Overview](#overview)
2. [Build Script Philosophy](#build-script-philosophy)
3. [Quick Start](#quick-start)
4. [Configuration Reference](#configuration-reference)
5. [Customization Guide](#customization-guide)
6. [Advanced Usage](#advanced-usage)
7. [Troubleshooting](#troubleshooting)

## Overview

The sleepms build scripts are designed to be **portable and reusable**. They handle:

- ✅ **Version Injection**: Automatically inject version info at build time
- ✅ **Build Number Generation**: UTC timestamp-based build numbers
- ✅ **Cross-Platform Support**: Linux, macOS, Windows
- ✅ **Documentation Generation**: Auto-generate GODOC.md from code comments
- ✅ **Binary Naming**: Platform-specific naming (e.g., `app-linux-amd64`)
- ✅ **Convenience Links**: Symlinks or copies for easy access
- ✅ **Colored Output**: Clear, informative build status messages

## Build Script Philosophy

### Two-Section Design

Each build script has two distinct sections:

```
┌─────────────────────────────────────┐
│  PROJECT CONFIGURATION              │  ← You modify this
│  (Top ~40 lines)                    │
├─────────────────────────────────────┤
│  BUILD IMPLEMENTATION               │  ← Don't touch this
│  (Rest of file)                     │
└─────────────────────────────────────┘
```

**Key Principle:** You only need to edit the configuration section at the top. The implementation section is generic and works for any Go project.

### Configuration vs Implementation

**Configuration Section** (Project-Specific):
- Project name and description
- Version numbers (semantic versioning)
- Source file locations
- Output naming patterns
- Documentation settings

**Implementation Section** (Generic/Portable):
- Platform detection
- Build number generation
- Go compiler invocation
- Documentation generation
- Status reporting

## Quick Start

### Step 1: Copy Files

Copy both build scripts to your project:

```bash
# From sleepms directory
cp build.sh /path/to/your-project/
cp build.cmd /path/to/your-project/
```

### Step 2: Make Executable (Unix)

```bash
cd /path/to/your-project
chmod +x build.sh
```

### Step 3: Edit Configuration

Open `build.sh` and edit the configuration section (lines 3-40):

```bash
#!/bin/bash
# =============================================================================
# PROJECT CONFIGURATION (modify these for your project)
# =============================================================================
PROJECT_NAME="your-app"              # ← Change this
PROJECT_DESCRIPTION="Your App Desc"  # ← Change this

# Semantic versioning
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="0"

# Source configuration
SOURCE_FILES="."       # Or specific files
MAIN_PACKAGE="."       # Or specific package path

# Build output
OUTPUT_NAME_TEMPLATE="${PROJECT_NAME}-${OS}-${ARCH}"
SYMLINK_NAME="your-app"              # ← Change this

# Documentation generation
GENERATE_DOCS=true
DOC_OUTPUT="GODOC.md"
```

Do the same for `build.cmd` (lines 4-39).

### Step 4: Update Your Go Code

Add version variables to your main package:

```go
// Version information variables
// These are set at build time using -ldflags
var (
    AppName         = "your-app"
    VersionMajor    = "1"
    VersionMinor    = "0"
    VersionRevision = "0"
    BuildNumber     = "00000000"  // Overwritten at build time
    BuildTime       = "unknown"   // Overwritten at build time
)
```

### Step 5: Build

```bash
./build.sh  # Unix
build.cmd   # Windows
```

That's it! You now have a professional build system.

## Configuration Reference

### Project Information

```bash
PROJECT_NAME="sleepms"
PROJECT_DESCRIPTION="Random Sleep Utility with Progress Bar"
```

**PROJECT_NAME:**
- Used in binary naming: `${PROJECT_NAME}-linux-amd64`
- Injected into Go variable: `main.AppName`
- Used in output messages

**PROJECT_DESCRIPTION:**
- Displayed in build header
- Purely informational
- Does not affect binary

### Version Numbers

```bash
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="0"
```

**Rules:**
- Use semantic versioning (see [VERSION_INFO.md](VERSION_INFO.md))
- Always use three components (Major.Minor.Revision)
- Update manually when releasing
- Injected into Go variables at build time

**When to Update:**
- `VERSION_MAJOR`: Breaking changes (1.x.x → 2.0.0)
- `VERSION_MINOR`: New features (1.0.x → 1.1.0)
- `VERSION_REVISION`: Bug fixes (1.0.0 → 1.0.1)

### Source Configuration

```bash
SOURCE_FILES="."
MAIN_PACKAGE="."
```

**SOURCE_FILES:**
- Passed to `go build`
- Usually `.` for current directory
- Can be specific files: `"main.go utils.go"`
- Can be package path: `"./cmd/app"`

**MAIN_PACKAGE:**
- Package containing `main()` function
- Usually `.` for root package
- Example: `"./cmd/server"` for cmd/server structure

### Build Output

```bash
OUTPUT_NAME_TEMPLATE="${PROJECT_NAME}-${OS}-${ARCH}"
SYMLINK_NAME="sleepms"
```

**OUTPUT_NAME_TEMPLATE:**
- Pattern for binary naming
- Available variables: `${PROJECT_NAME}`, `${OS}`, `${ARCH}`
- Example result: `sleepms-darwin-arm64`
- Windows automatically gets `.exe` extension

**SYMLINK_NAME:**
- Simple name for convenience access
- Unix: Creates symlink to versioned binary
- Windows: Creates copy of binary
- Example: `sleepms` → `sleepms-darwin-arm64`

### Documentation Generation

```bash
GENERATE_DOCS=true
DOC_OUTPUT="GODOC.md"
```

**GENERATE_DOCS:**
- `true`: Generate documentation during build
- `false`: Skip documentation generation
- Uses `go doc -all` command

**DOC_OUTPUT:**
- Filename for generated documentation
- Conventional name: `GODOC.md`
- Created/overwritten on each build

## Customization Guide

### Example 1: Multi-Command Project

**Project Structure:**
```
myapp/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
└── build.sh
```

**Option A: Separate Build Scripts**

Create `build-server.sh`:
```bash
PROJECT_NAME="myapp-server"
MAIN_PACKAGE="./cmd/server"
SYMLINK_NAME="server"
```

Create `build-client.sh`:
```bash
PROJECT_NAME="myapp-client"
MAIN_PACKAGE="./cmd/client"
SYMLINK_NAME="client"
```

**Option B: Build All (Advanced)**

Modify the build script to loop through commands:
```bash
# After configuration section, add:
COMMANDS=("server" "client")

for CMD in "${COMMANDS[@]}"; do
    PROJECT_NAME="myapp-${CMD}"
    MAIN_PACKAGE="./cmd/${CMD}"
    # ... build logic ...
done
```

### Example 2: Custom Binary Naming

**Goal:** Include version in binary name

```bash
OUTPUT_NAME_TEMPLATE="${PROJECT_NAME}-v${VERSION_MAJOR}.${VERSION_MINOR}-${OS}-${ARCH}"
```

**Result:** `sleepms-v1.0-darwin-arm64`

**Goal:** Simpler naming (no platform)

```bash
OUTPUT_NAME_TEMPLATE="${PROJECT_NAME}"
```

**Result:** `sleepms` (overwrites on each platform)

### Example 3: Disable Documentation

If your project doesn't need generated docs:

```bash
GENERATE_DOCS=false
```

Or keep it but change the filename:

```bash
DOC_OUTPUT="API.md"  # or "docs/API.md"
```

### Example 4: Specific Go Files

If you have multiple `.go` files but want to build only specific ones:

```bash
SOURCE_FILES="main.go config.go utils.go"
MAIN_PACKAGE="."
```

## Advanced Usage

### Custom Linker Flags

The implementation section constructs these flags automatically:

```bash
LDFLAGS="-s -w"
LDFLAGS="${LDFLAGS} -X 'main.AppName=${PROJECT_NAME}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionMajor=${VERSION_MAJOR}'"
# ... etc ...
```

**To add custom flags**, modify after the linker flags construction:

```bash
# Add after existing LDFLAGS setup
LDFLAGS="${LDFLAGS} -X 'main.BuildUser=${USER}'"
LDFLAGS="${LDFLAGS} -X 'main.GitCommit=$(git rev-parse HEAD)'"
```

Then add variables to your Go code:

```go
var (
    // Existing variables...
    BuildUser  = "unknown"
    GitCommit  = "unknown"
)
```

### Environment-Specific Builds

**Development Build:**
```bash
# Remove -s -w to keep debug symbols
LDFLAGS=""  # Start with empty flags instead of "-s -w"
```

**Production Build:**
```bash
# Default includes -s -w (strip symbols for smaller binary)
LDFLAGS="-s -w"
```

### Cross-Compilation

The build script detects the current platform. To build for other platforms:

**Linux to Windows:**
```bash
GOOS=windows GOARCH=amd64 go build ...
```

**Modify build script for multiple platforms:**
```bash
# After configuration, add:
PLATFORMS=("linux/amd64" "darwin/arm64" "windows/amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    # ... build logic ...
done
```

### Documentation Customization

The documentation header can be customized:

Find this section in the implementation:
```bash
cat > "${DOC_OUTPUT}" << EOF
# ${PROJECT_NAME} - Go Package Documentation

**Version:** ${FULL_VERSION}
**Generated:** ${BUILD_TIME}

---

## Package Overview

EOF
```

Modify to:
```bash
cat > "${DOC_OUTPUT}" << EOF
# ${PROJECT_NAME} API Documentation

Project: ${PROJECT_NAME}
Version: ${FULL_VERSION}
Generated: ${BUILD_TIME}
Platform: ${OS}/${ARCH}

---

EOF
```

## Advanced Scenarios

### Scenario 1: Multi-Module Workspace

**Project Structure:**
```
workspace/
├── module-a/
│   ├── go.mod
│   ├── main.go
│   └── build.sh
└── module-b/
    ├── go.mod
    ├── main.go
    └── build.sh
```

Each module gets its own build script with appropriate `PROJECT_NAME`.

### Scenario 2: Nested Package Structure

**Project Structure:**
```
myapp/
├── go.mod
├── build.sh
├── internal/
│   ├── config/
│   └── utils/
└── cmd/
    └── myapp/
        └── main.go
```

**Configuration:**
```bash
PROJECT_NAME="myapp"
MAIN_PACKAGE="./cmd/myapp"
SOURCE_FILES="."  # Include all packages
```

### Scenario 3: Library (No Main)

If building a library (no executable):

**Option 1:** Skip build, just generate docs
```bash
GENERATE_DOCS=true

# Comment out or skip the go build section
# Just run documentation generation
```

**Option 2:** Build test binary
```bash
go test -c -o test-binary
```

## Troubleshooting

### Problem: "go: command not found"

**Solution:**
- Install Go from https://golang.org/dl/
- Ensure `go` is in your PATH
- Verify: `go version`

### Problem: Wrong version injected

**Check:**
1. Build script configuration variables
2. Go variable names match exactly
3. Package name is correct (`main.VersionMajor` vs `pkg.VersionMajor`)

**Debug:**
```bash
# Build with verbose output
go build -v -ldflags "${LDFLAGS}" ...

# Check what's being injected
echo "${LDFLAGS}"
```

### Problem: Build works but version shows default values

**Cause:** Built with `go build` directly instead of build script

**Solution:** Always use `./build.sh` or `build.cmd`

### Problem: Documentation generation fails

**Check:**
1. Package has godoc comments
2. `go doc` command works: `go doc -all`
3. File permissions for `DOC_OUTPUT`

**Debug:**
```bash
# Test manually
go doc -all > test-doc.md
```

### Problem: Binary not executable

**Unix:**
```bash
chmod +x sleepms-linux-amd64
```

**Windows:** Should not occur (`.exe` files are executable by default)

### Problem: Symlink creation fails

**Unix:**
- Check permissions in directory
- Verify no file exists with same name

**Windows:**
- Copy is used instead of symlink
- Check disk space and permissions

## Best Practices

### Do's

✅ Keep configuration section at top clearly marked
✅ Update both `build.sh` and `build.cmd` together
✅ Test build script after modifications
✅ Use semantic versioning for version numbers
✅ Document any custom modifications
✅ Keep implementation section generic

### Don'ts

❌ Don't modify implementation section unless necessary
❌ Don't hardcode version info in Go source files
❌ Don't skip version updates before releases
❌ Don't use spaces in `PROJECT_NAME`
❌ Don't commit build artifacts (add to `.gitignore`)

## File Checklist

When setting up build scripts in a new project:

- [ ] Copy `build.sh` and `build.cmd`
- [ ] Make `build.sh` executable (`chmod +x`)
- [ ] Update `PROJECT_NAME` in both scripts
- [ ] Update `PROJECT_DESCRIPTION` in both scripts
- [ ] Update `SYMLINK_NAME` in both scripts
- [ ] Adjust `MAIN_PACKAGE` if needed
- [ ] Add version variables to Go code
- [ ] Test build: `./build.sh`
- [ ] Test version flag: `./app --version`
- [ ] Add build artifacts to `.gitignore`
- [ ] Document custom modifications

## Example .gitignore

```gitignore
# Build artifacts
sleepms
sleepms-*
*.exe

# Documentation (optional - you may want to commit this)
GODOC.md

# Go build cache
*.test
*.out
```

## Summary

The sleepms build scripts provide:

1. **Portability**: Copy and configure for any Go project
2. **Automation**: Automatic version injection and doc generation
3. **Flexibility**: Easy to customize via configuration section
4. **Maintainability**: Clean separation between config and implementation
5. **Professional**: Colored output, proper error handling, comprehensive features

By following this guide, you can adopt these build practices in your own Go projects with minimal effort.

## Additional Resources

- [VERSION_INFO.md](VERSION_INFO.md) - Version management details
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development guidelines
- [Go Build Documentation](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Go Linker Flags](https://pkg.go.dev/cmd/link)
