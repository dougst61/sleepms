# sleepms - Go Package Documentation

**Version:** 1.0.0 build 12260147
**Generated:** 2025-12-26 01:47:24 UTC

---

## Package Overview

Package main implements a sleep utility that waits for a random duration between
a minimum and maximum value (in milliseconds).

The program generates a random sleep duration, displays it to the user, and then
sleeps for that duration. The sleep can be interrupted at any time by pressing
any key on the keyboard.

For sleep durations longer than 1 minute (60,000ms), a progress bar is displayed
showing the percentage of time elapsed and a visual indicator of progress.

Usage:

    sleepms <min> <max>

Arguments:

    min - Minimum sleep duration in milliseconds
    max - Maximum sleep duration in milliseconds

Example:

    sleepms 1000 5000  # Sleep between 1 and 5 seconds
    sleepms 60000 120000  # Sleep between 1 and 2 minutes (with progress bar)

VARIABLES

var (
	AppName         = "sleepms"
	VersionMajor    = "1"
	VersionMinor    = "0"
	VersionRevision = "0"
	BuildNumber     = "00000000" // Overwritten at build time
	BuildTime       = "unknown"  // Overwritten at build time
)
    Version information variables These are set at build time using -ldflags

