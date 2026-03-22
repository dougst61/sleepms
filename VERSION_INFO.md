# Version Information for sleepms

This document describes the versioning system used in sleepms and how to manage version numbers.

## Table of Contents

1. [Version Format](#version-format)
2. [Semantic Versioning](#semantic-versioning)
3. [Build Number System](#build-number-system)
4. [Version Variables](#version-variables)
5. [Updating Versions](#updating-versions)
6. [Build-Time Injection](#build-time-injection)
7. [Displaying Version](#displaying-version)

## Version Format

sleepms uses a comprehensive version format:

```
<AppName> version <Major>.<Minor>.<Revision> build <BuildNumber>
Built: <BuildTime>
```

**Example:**
```
sleepms version 1.0.0 build 11091636
Built: 2025-11-09 16:36:00 UTC
```

### Components

| Component | Description | Example | Source |
|-----------|-------------|---------|--------|
| AppName | Application name | `sleepms` | Hard-coded |
| Major | Major version | `1` | Manual |
| Minor | Minor version | `0` | Manual |
| Revision | Patch/revision | `0` | Manual |
| BuildNumber | Build timestamp | `11091636` | Auto-generated |
| BuildTime | Human-readable timestamp | `2025-11-09 16:36:00 UTC` | Auto-generated |

## Semantic Versioning

sleepms follows [Semantic Versioning 2.0.0](https://semver.org/) principles:

### Major Version (X.0.0)

Increment when making **incompatible changes**:
- Breaking API changes
- Removing command-line arguments
- Changing output format significantly
- Incompatible behavior changes

**Examples:**
- `1.x.x` → `2.0.0`: Changed from milliseconds to seconds
- `2.x.x` → `3.0.0`: Removed `--version` flag (don't actually do this!)

### Minor Version (1.X.0)

Increment when adding **backwards-compatible functionality**:
- New features
- New command-line options
- Enhanced existing features
- Performance improvements

**Examples:**
- `1.0.x` → `1.1.0`: Added color output to progress bar
- `1.1.x` → `1.2.0`: Added configuration file support

### Revision/Patch Version (1.0.X)

Increment when making **backwards-compatible bug fixes**:
- Bug fixes
- Documentation corrections
- Code refactoring (no behavior change)
- Security patches

**Examples:**
- `1.0.0` → `1.0.1`: Fixed progress bar display bug
- `1.0.1` → `1.0.2`: Fixed terminal restore on Windows

## Build Number System

### Format: MMDDHHMM (UTC)

Build numbers are automatically generated from the current UTC timestamp:

```
MM   DD   HH   MM
│    │    │    └─── Minute (00-59)
│    │    └──────── Hour (00-23)
│    └───────────── Day (01-31)
└────────────────── Month (01-12)
```

**Examples:**

| Build Number | Meaning |
|--------------|---------|
| `11091636` | November 9th, 16:36 UTC |
| `01011200` | January 1st, 12:00 UTC |
| `12312359` | December 31st, 23:59 UTC |

### Build Number Properties

**Advantages:**
1. **Automatic**: No manual tracking needed
2. **Unique**: One per minute (across all developers using UTC)
3. **Human-Readable**: You can decode when the build was made
4. **Sortable**: Lexicographically ordered (mostly)
5. **No Coordination**: No shared build counter required

**Limitations:**
1. **Not Strictly Monotonic**: Wraps at year boundary (December → January)
2. **Collision Possible**: Two builds in the same minute have the same number
3. **Timezone Dependent**: Only works reliably in UTC

**Why UTC?**
- Consistent across all developers worldwide
- No daylight saving time changes
- Standard in software development

## Version Variables

Version information is stored in Go variables at the top of `sleepms.go`:

```go
// Version information variables
// These are set at build time using -ldflags
var (
    AppName         = "sleepms"
    VersionMajor    = "1"
    VersionMinor    = "0"
    VersionRevision = "0"
    BuildNumber     = "00000000"  // Overwritten at build time
    BuildTime       = "unknown"   // Overwritten at build time
)
```

### Variable Details

| Variable | Type | Set By | When |
|----------|------|--------|------|
| AppName | string | Developer | Compile time |
| VersionMajor | string | Developer | Via build script |
| VersionMinor | string | Developer | Via build script |
| VersionRevision | string | Developer | Via build script |
| BuildNumber | string | Build script | Build time |
| BuildTime | string | Build script | Build time |

**Note:** All version components are strings for consistent formatting.

## Updating Versions

### When to Update

Update version numbers in the build scripts before creating a release:

**Major Release:**
1. Increment `VERSION_MAJOR`
2. Reset `VERSION_MINOR` to 0
3. Reset `VERSION_REVISION` to 0
4. Create git tag: `v2.0.0`

**Minor Release:**
1. Increment `VERSION_MINOR`
2. Reset `VERSION_REVISION` to 0
3. Create git tag: `v1.1.0`

**Patch Release:**
1. Increment `VERSION_REVISION`
2. Create git tag: `v1.0.1`

### Where to Update

**File: build.sh (Lines 16-18)**
```bash
# Semantic versioning - update these when releasing new versions
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="0"
```

**File: build.cmd (Lines 17-19)**
```cmd
REM Semantic versioning - update these when releasing new versions
set VERSION_MAJOR=1
set VERSION_MINOR=0
set VERSION_REVISION=0
```

**Important:** Update **both** build scripts to keep them in sync!

### Example Update Process

**Scenario:** Releasing version 1.1.0 (new feature)

1. **Edit build.sh:**
```bash
VERSION_MAJOR="1"
VERSION_MINOR="1"  # Changed from 0
VERSION_REVISION="0"
```

2. **Edit build.cmd:**
```cmd
set VERSION_MAJOR=1
set VERSION_MINOR=1
set VERSION_REVISION=0
```

3. **Commit changes:**
```bash
git add build.sh build.cmd
git commit -m "chore: bump version to 1.1.0"
```

4. **Build and test:**
```bash
./build.sh
./sleepms --version  # Verify version
```

5. **Create git tag:**
```bash
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0
```

## Build-Time Injection

Version information is injected at build time using Go's `-ldflags` mechanism.

### How It Works

The build script constructs a `-ldflags` string that overwrites the Go variables:

**In build.sh:**
```bash
# Generate build number and timestamp
BUILD_NUMBER=$(date -u '+%m%d%H%M')
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')

# Construct linker flags
LDFLAGS="-s -w"
LDFLAGS="${LDFLAGS} -X 'main.AppName=${PROJECT_NAME}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionMajor=${VERSION_MAJOR}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionMinor=${VERSION_MINOR}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionRevision=${VERSION_REVISION}'"
LDFLAGS="${LDFLAGS} -X 'main.BuildNumber=${BUILD_NUMBER}'"
LDFLAGS="${LDFLAGS} -X 'main.BuildTime=${BUILD_TIME}'"

# Build with injected values
go build -ldflags "${LDFLAGS}" -o "${OUTPUT_NAME}" ${MAIN_PACKAGE}
```

### Linker Flag Breakdown

| Flag | Purpose |
|------|---------|
| `-s` | Strip debug symbols (smaller binary) |
| `-w` | Strip DWARF debug info (smaller binary) |
| `-X 'main.VersionMajor=1'` | Set variable value |

**Format:** `-X 'package.Variable=Value'`

**Why Quotes?**
- Single quotes protect spaces in `BUILD_TIME`
- Example: `'2025-11-09 16:36:00 UTC'` contains spaces

### Verification

You can verify the injected values:

```bash
# Build the binary
./build.sh

# Check version
./sleepms --version

# Expected output:
# sleepms version 1.0.0 build 11091636
# Built: 2025-11-09 16:36:00 UTC
```

### Why This Approach?

**Advantages:**
1. **No Source Changes**: Version is set at build time, not in code
2. **Single Source of Truth**: Build script controls versioning
3. **Automatic Build Info**: Timestamp injected automatically
4. **Clean Repository**: No version files to commit
5. **Flexible**: Easy to override for special builds

**Alternatives (Not Used):**
- Embedding version in source: Requires commit for every build
- Version file: Extra file to maintain and commit
- Git tags only: Can't show build timestamp

## Displaying Version

### In Code

The `--version` flag displays version information:

```go
// Check for version flag
if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
    fmt.Printf("%s version %s.%s.%s build %s\n",
        AppName, VersionMajor, VersionMinor, VersionRevision, BuildNumber)
    fmt.Printf("Built: %s\n", BuildTime)
    return
}
```

### User Interface

```bash
$ ./sleepms --version
sleepms version 1.0.0 build 11091636
Built: 2025-11-09 16:36:00 UTC
```

**Supported Flags:**
- `--version` (long form)
- `-v` (short form)

### Version in Documentation

The build script automatically adds version info to generated documentation:

**GODOC.md header:**
```markdown
# sleepms - Go Package Documentation

**Version:** 1.0.0 build 11091636
**Generated:** 2025-11-09 16:36:00 UTC
```

## Best Practices

### Do's

✅ **Update both build scripts** when changing versions
✅ **Follow semantic versioning** principles
✅ **Create git tags** for releases
✅ **Test build** after version changes
✅ **Document** what changed in each version
✅ **Use UTC** for all timestamps

### Don'ts

❌ **Don't** manually edit `BuildNumber` or `BuildTime`
❌ **Don't** commit version changes without testing
❌ **Don't** skip version components (e.g., 1.3 instead of 1.3.0)
❌ **Don't** reuse version numbers
❌ **Don't** forget to update documentation

## Examples

### Example 1: Bug Fix Release

**Current:** `1.0.0`
**Goal:** Fix a bug
**New Version:** `1.0.1`

```bash
# Edit build scripts
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="1"  # Incremented

# Commit, build, tag
git commit -am "fix: terminal restore on interrupt"
./build.sh
git tag -a v1.0.1 -m "Bug fix release"
```

### Example 2: New Feature Release

**Current:** `1.0.1`
**Goal:** Add color output
**New Version:** `1.1.0`

```bash
# Edit build scripts
VERSION_MAJOR="1"
VERSION_MINOR="1"      # Incremented
VERSION_REVISION="0"   # Reset to 0

# Commit, build, tag
git commit -am "feat: add color to progress bar"
./build.sh
git tag -a v1.1.0 -m "Feature release: color output"
```

### Example 3: Breaking Change

**Current:** `1.5.2`
**Goal:** Change to seconds instead of milliseconds
**New Version:** `2.0.0`

```bash
# Edit build scripts
VERSION_MAJOR="2"      # Incremented
VERSION_MINOR="0"      # Reset to 0
VERSION_REVISION="0"   # Reset to 0

# Commit, build, tag
git commit -am "feat!: change to seconds (breaking change)"
./build.sh
git tag -a v2.0.0 -m "Major release: seconds instead of milliseconds"
```

## Troubleshooting

### Problem: Build number shows as "00000000"

**Cause:** Built with `go build` directly instead of build script
**Solution:** Always use `./build.sh` or `build.cmd` to build

### Problem: Version shows wrong timestamp

**Cause:** System clock incorrect or not in UTC
**Solution:** Check system clock and timezone

### Problem: Two builds have same build number

**Cause:** Built within same minute
**Solution:** Wait a minute or manually modify for special case

### Problem: Version mismatch between platforms

**Cause:** Different version numbers in build.sh vs build.cmd
**Solution:** Keep both files in sync

## References

- [Semantic Versioning Specification](https://semver.org/)
- [Go Linker Flags Documentation](https://pkg.go.dev/cmd/link)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Summary

The sleepms versioning system combines:
- **Manual Semantic Versioning** for releases (Major.Minor.Revision)
- **Automatic Build Numbers** from UTC timestamps (MMDDHHMM)
- **Build-Time Injection** via Go linker flags
- **Developer Convenience** through automated build scripts

This provides a robust, maintainable versioning system that requires minimal manual intervention while providing comprehensive version information.
