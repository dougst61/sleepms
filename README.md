# sleepms - Random Sleep Utility with Progress Bar

A cross-platform command-line utility that sleeps for a random duration between specified minimum and maximum values. Features an interactive progress bar with time-remaining display for longer sleep durations, keyboard and signal interruption, and distinct exit codes for scripting.

## Features

- **Flexible Duration Input**: Specify durations as plain milliseconds (`5000`) or Go duration strings (`5s`, `1m30s`, `500ms`)
- **Human-Readable Output**: Sleep time displayed as both raw milliseconds and a friendly string (e.g. `1m 27.234s`)
- **Progress Bar with Time Remaining**: Visual progress and countdown for sleep durations ≥ 20 seconds
- **Keyboard Interrupt**: Press any key to cut the sleep short
- **Signal Handling**: SIGINT and SIGTERM restore the terminal and exit cleanly
- **Non-TTY Fallback**: Works in pipelines and scripts where stdin is not a terminal
- **Quiet Mode**: `-q` suppresses all output; use the exit code to detect completion vs. interruption
- **Distinct Exit Codes**: `0` = completed, `1` = interrupted, `2` = error
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Use Cases

1. **Testing Timeouts**: Simulate random delays in testing scripts
2. **Rate Limiting**: Add controlled delays between operations
3. **Simulating Network Latency**: Test application behavior with variable delays
4. **Demo/Presentation**: Create dramatic pauses with visual feedback
5. **Development**: Quick delays during script execution

## Installation

### Prerequisites

- Go 1.21 or higher
- Terminal with ANSI escape code support (for progress bar)

### Building from Source

**Linux/macOS:**
```bash
./build.sh
```

**Windows:**
```cmd
build.cmd
```

The build script will:
- Compile the binary with version information
- Generate documentation (GODOC.md)
- Create a platform-specific binary (e.g., `sleepms-darwin-arm64`)
- Create a symlink/copy for easy access (`sleepms` or `sleepms.exe`)

## Usage

```
sleepms [options] <min> <max>
```

### Arguments

- `min` — minimum sleep duration
- `max` — maximum sleep duration

Durations may be specified as a plain integer (milliseconds) or a Go duration string:

| Input | Interpreted as |
|-------|---------------|
| `5000` | 5000 ms |
| `5s` | 5000 ms |
| `500ms` | 500 ms |
| `1m30s` | 90 000 ms |
| `2m` | 120 000 ms |

### Options

| Flag | Description |
|------|-------------|
| `-q`, `--quiet` | Suppress all output |
| `-v`, `--version` | Print version info and exit |

### Examples

**Short sleep using milliseconds (1–5 seconds):**
```bash
./sleepms 1000 5000
```
```
Sleeping for 3.427s (3427ms)
Sleep complete.
```

**Short sleep using duration strings:**
```bash
./sleepms 1s 5s
```
```
Sleeping for 3.427s (3427ms)
Sleep complete.
```

**Long sleep with progress bar (1–2 minutes):**
```bash
./sleepms 1m 2m
```
```
Sleeping for 1m 27.234s (87234ms)
Sleeping: [===================>                    ] 47%  46s remaining
```

**Fixed duration (exactly 30 seconds):**
```bash
./sleepms 30s 30s
```

**Quiet mode (exit code only):**
```bash
./sleepms -q 1s 5s
if [ $? -eq 1 ]; then echo "interrupted"; fi
```

**Check version:**
```bash
./sleepms --version
```
```
sleepms version 1.1.0 build 11091636
Built: 2025-11-09 16:36:00 UTC
```

## Progress Bar

For sleep durations **20 seconds or longer** (≥ 20,000 ms), sleepms displays a real-time progress bar with time remaining:

```
Sleeping: [=========>                              ] 25%  1m 30s remaining
```

The progress bar:
- Updates every second
- Shows percentage completion and time remaining
- 40 characters wide
- Can be interrupted at any time

## Keyboard Interrupt

At any point during the sleep, press **any key** to interrupt:

```
Sleeping for 45.000s (45000ms)
Sleeping: [======>                                 ] 15%  38s remaining
Sleep interrupted.
```

The program exits immediately with code `1` and restores the terminal to its normal state.

## Signal Handling

`SIGINT` and `SIGTERM` are caught and handled gracefully — the terminal is restored before the process exits. Exit code is `1` when interrupted by a signal.

```bash
kill $sleepms_pid   # terminal restored, exit code 1
```

## Non-TTY / Pipeline Use

When stdin is not a terminal (e.g. the process is running inside a pipe or a CI script), sleepms detects this automatically and skips raw terminal mode. The sleep runs normally without keyboard-interrupt support. Progress is still displayed unless `-q` is used.

```bash
echo "" | ./sleepms 1s 5s   # works without error
```

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Sleep completed normally |
| `1` | Sleep interrupted (keypress or signal) |
| `2` | Error (invalid arguments, terminal setup failure) |

## Console Output

### Normal Operation (Short Duration)
```
Sleeping for 2.500s (2500ms)
Sleep complete.
```

### Normal Operation (Long Duration)
```
Sleeping for 1m 30.000s (90000ms)
Sleeping: [================================>       ] 81%  17s remaining
Sleep complete.
```

### Interrupted Operation
```
Sleeping for 2m 0.000s (120000ms)
Sleeping: [========>                               ] 21%  1m 34s remaining
Sleep interrupted.
```

### Quiet Mode
```
(no output)
```

### Error Messages
```
Error: expected 2 arguments (min max), got 1
Error: invalid min: use milliseconds (e.g. 5000) or a Go duration string (e.g. 5s, 1m30s, 500ms)
Error: min (5.000s) must not exceed max (1.000s)
```

## Technical Details

### Terminal Mode

The program uses **raw terminal mode** to detect keypresses without requiring Enter:
- Characters are not echoed to the screen
- No line buffering occurs
- Input is detected immediately
- Terminal state is always restored on exit (including on SIGTERM)

### Progress Threshold

Progress bar is shown when sleep duration ≥ 20,000 ms (20 seconds).

### Random Number Generation

- Uses `math/rand` with time-based seeding
- Range: `[min, max]` (inclusive on both ends)
- Distribution: Uniform random

### Build Number Format

Build numbers use the format `MMDDHHMM` in UTC:
- `MM` — Month (01–12)
- `DD` — Day (01–31)
- `HH` — Hour (00–23)
- `MM` — Minute (00–59)

Example: `11091636` = November 9th, 16:36 UTC

## Platform Support

| Platform | Status | Binary Name |
|----------|--------|-------------|
| Linux (amd64) | ✓ Supported | sleepms-linux-amd64 |
| Linux (arm64) | ✓ Supported | sleepms-linux-arm64 |
| macOS (amd64) | ✓ Supported | sleepms-darwin-amd64 |
| macOS (arm64) | ✓ Supported | sleepms-darwin-arm64 |
| Windows (amd64) | ✓ Supported | sleepms-windows-amd64.exe |
| Windows (arm64) | ✓ Supported | sleepms-windows-arm64.exe |

## Troubleshooting

### Problem: "go: command not found"
**Solution:** Install Go from https://golang.org/dl/

### Problem: Progress bar shows garbled characters
**Solution:** Use a terminal that supports ANSI escape codes (most modern terminals do)

### Problem: Invalid duration error
**Solution:** Use a plain integer (milliseconds) or a valid Go duration string: `5s`, `500ms`, `1m30s`

### Problem: Program doesn't respond to keypress
**Solution:**
- Ensure terminal supports raw mode (most do)
- Try pressing different keys
- Check that stdin is connected to a terminal (not redirected)

### Problem: Build fails with "package golang.org/x/term: not found"
**Solution:** Run `go mod download` to fetch dependencies

## Performance Notes

- **Startup Time**: Negligible (< 10 ms)
- **Memory Usage**: ~5–10 MB
- **CPU Usage**: Near zero during sleep (polling for keyboard only)
- **Accuracy**: ± 1 second for display updates; actual sleep is precise

## Documentation

- **README.md** (this file) — User documentation
- **CONTRIBUTING.md** — Development guidelines
- **VERSION_INFO.md** — Versioning system documentation
- **BUILD_TEMPLATE.md** — Build script reuse guide
- **GODOC.md** — Auto-generated API documentation

## Version Information

Version format: `Major.Minor.Revision build BuildNumber`

- **Major**: Breaking changes
- **Minor**: New features (backwards compatible)
- **Revision**: Bug fixes (backwards compatible)
- **BuildNumber**: Auto-generated from build timestamp (MMDDHHMM UTC)

See [VERSION_INFO.md](VERSION_INFO.md) for complete versioning documentation.

## License

This project is provided as-is for educational and practical use.

## Author

Created as a demonstration of clean Go programming practices with professional documentation and build automation.

---
Copyright (c) 2026 Doug Stewart
