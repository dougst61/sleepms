#!/bin/bash
# =============================================================================
# BUILD SCRIPT FOR sleepms
# =============================================================================
# This script automates the build process for the sleepms application.
# It handles version injection, documentation generation, and cross-platform builds.
#
# Usage: ./build.sh
# Output: sleepms-{OS}-{ARCH} binary and GODOC.md documentation
# =============================================================================

# =============================================================================
# PROJECT CONFIGURATION (modify these for your project)
# =============================================================================
PROJECT_NAME="sleepms"
PROJECT_DESCRIPTION="Random Sleep Utility with Progress Bar"

# Semantic versioning - update these when releasing new versions
VERSION_MAJOR="1"
VERSION_MINOR="0"
VERSION_REVISION="0"

# Source code configuration
SOURCE_FILES="."
MAIN_PACKAGE="."

# Build output configuration
OUTPUT_NAME_TEMPLATE="${PROJECT_NAME}-${OS}-${ARCH}"
SYMLINK_NAME="sleepms"

# Documentation generation
GENERATE_DOCS=true
DOC_OUTPUT="GODOC.md"

# =============================================================================
# BUILD SCRIPT IMPLEMENTATION (generic - no need to modify below this line)
# =============================================================================

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored status messages
print_status() {
    echo -e "${BLUE}[*]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Print build header
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Building ${PROJECT_NAME} - ${PROJECT_DESCRIPTION}${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo

# Check if Go is installed
print_status "Checking Go installation..."
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go and try again."
    exit 1
fi

GO_VERSION=$(go version)
print_success "Go found: ${GO_VERSION}"
echo

# Detect OS and architecture
print_status "Detecting platform..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) print_error "Unsupported OS: $OS"; exit 1 ;;
esac

ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    i386|i686) ARCH="386" ;;
    *) print_error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

print_success "Platform: ${OS}/${ARCH}"
echo

# Generate build number from current UTC time (MMDDHHMM format)
print_status "Generating build number..."
BUILD_NUMBER=$(date -u '+%m%d%H%M')
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
print_success "Build number: ${BUILD_NUMBER}"
print_success "Build time: ${BUILD_TIME}"
echo

# Construct version string
VERSION_STRING="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_REVISION}"
FULL_VERSION="${VERSION_STRING} build ${BUILD_NUMBER}"

print_status "Version: ${FULL_VERSION}"
echo

# Construct linker flags for version injection
LDFLAGS="-s -w"
LDFLAGS="${LDFLAGS} -X 'main.AppName=${PROJECT_NAME}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionMajor=${VERSION_MAJOR}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionMinor=${VERSION_MINOR}'"
LDFLAGS="${LDFLAGS} -X 'main.VersionRevision=${VERSION_REVISION}'"
LDFLAGS="${LDFLAGS} -X 'main.BuildNumber=${BUILD_NUMBER}'"
LDFLAGS="${LDFLAGS} -X 'main.BuildTime=${BUILD_TIME}'"

# Determine output filename
OUTPUT_NAME="${PROJECT_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    OUTPUT_NAME="${OUTPUT_NAME}.exe"
fi

# Build the application
print_status "Building ${PROJECT_NAME}..."
if go build -ldflags "${LDFLAGS}" -o "${OUTPUT_NAME}" ${MAIN_PACKAGE}; then
    print_success "Build successful: ${OUTPUT_NAME}"

    # Get file size
    if [ "$OS" = "darwin" ]; then
        FILE_SIZE=$(stat -f%z "${OUTPUT_NAME}")
    else
        FILE_SIZE=$(stat -c%s "${OUTPUT_NAME}" 2>/dev/null || echo "unknown")
    fi

    if [ "$FILE_SIZE" != "unknown" ]; then
        FILE_SIZE_MB=$(awk "BEGIN {printf \"%.2f\", ${FILE_SIZE}/1024/1024}")
        print_success "Binary size: ${FILE_SIZE_MB} MB"
    fi
else
    print_error "Build failed"
    exit 1
fi
echo

# Create symlink or copy for easy access
print_status "Creating convenience link/copy..."
if [ "$OS" = "windows" ]; then
    SYMLINK_TARGET="${SYMLINK_NAME}.exe"
    if cp "${OUTPUT_NAME}" "${SYMLINK_TARGET}"; then
        print_success "Created copy: ${SYMLINK_TARGET}"
    else
        print_warning "Failed to create copy: ${SYMLINK_TARGET}"
    fi
else
    if ln -sf "${OUTPUT_NAME}" "${SYMLINK_NAME}"; then
        print_success "Created symlink: ${SYMLINK_NAME} -> ${OUTPUT_NAME}"
    else
        print_warning "Failed to create symlink: ${SYMLINK_NAME}"
    fi
fi
echo

# Generate documentation if enabled
if [ "$GENERATE_DOCS" = true ]; then
    print_status "Generating documentation..."

    # Create documentation header
    cat > "${DOC_OUTPUT}" << EOF
# ${PROJECT_NAME} - Go Package Documentation

**Version:** ${FULL_VERSION}
**Generated:** ${BUILD_TIME}

---

## Package Overview

EOF

    # Extract package documentation
    go doc -all >> "${DOC_OUTPUT}" 2>&1

    if [ $? -eq 0 ]; then
        print_success "Documentation generated: ${DOC_OUTPUT}"
    else
        print_warning "Documentation generation had warnings (check ${DOC_OUTPUT})"
    fi
    echo
fi

# Print build summary
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Build Complete!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo
echo -e "  ${BLUE}Project:${NC}     ${PROJECT_NAME}"
echo -e "  ${BLUE}Version:${NC}     ${FULL_VERSION}"
echo -e "  ${BLUE}Platform:${NC}    ${OS}/${ARCH}"
echo -e "  ${BLUE}Binary:${NC}      ${OUTPUT_NAME}"
if [ "$GENERATE_DOCS" = true ]; then
    echo -e "  ${BLUE}Docs:${NC}        ${DOC_OUTPUT}"
fi
echo
echo -e "${YELLOW}Run with:${NC} ./${SYMLINK_NAME} --version"
echo -e "${YELLOW}Usage:${NC}    ./${SYMLINK_NAME} <min> <max>"
echo
