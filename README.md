# sleepms - Random Sleep Utility with Progress Bar

A cross-platform command-line utility that sleeps for a random duration between specified minimum and maximum values (in milliseconds). Features an interactive progress bar for longer sleep durations and supports keyboard interruption.

## Features

- **Random Sleep Duration**: Generates a random sleep time between minimum and maximum values
- **Progress Bar Display**: Shows visual progress for sleep durations over 20 seconds
- **Keyboard Interrupt**: Press any key to interrupt the sleep at any time
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Zero Dependencies**: Uses only Go standard library (plus golang.org/x/term for terminal control)
- **Precise Timing**: Accurate millisecond-level sleep duration
- **Clean Output**: Clear console feedback with percentage completion

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

```bash
sleepms <min> <max>
```

### Arguments

- `min` - Minimum sleep duration in milliseconds (integer)
- `max` - Maximum sleep duration in milliseconds (integer)

### Examples

**Short sleep (1-5 seconds):**
```bash
./sleepms 1000 5000
```
Output:
```
Generated sleep time: 3427ms
Sleep duration complete.
```

**Long sleep with progress bar (1-2 minutes):**
```bash
./sleepms 60000 120000
```
Output:
```
Generated sleep time: 87234ms
Sleeping: [===================>                    ] 47%
```

**Fixed duration (exactly 30 seconds):**
```bash
./sleepms 30000 30000
```

**Check version:**
```bash
./sleepms --version
```
Output:
```
sleepms version 1.0.0 build 11091636
Built: 2025-11-09 16:36:00 UTC
```

## Progress Bar

For sleep durations **20 seconds or longer** (≥20,000ms), sleepms displays a real-time progress bar:

```
Sleeping: [=========>                              ] 25%
```

The progress bar:
- Updates every second
- Shows percentage completion (0-100%)
- Displays visual progress indicator
- 40 characters wide
- Can be interrupted at any time

## Keyboard Interrupt

At any point during the sleep, you can press **any key** to interrupt:

```
Generated sleep time: 45000ms
Sleeping: [======>                                 ] 15%
Program interrupted by keyboard input.
```

The program will immediately exit and restore the terminal to its normal state.

## Console Output

### Normal Operation (Short Duration)
```
Generated sleep time: 2500ms
Sleep duration complete.
```

### Normal Operation (Long Duration)
```
Generated sleep time: 90000ms
Sleeping: [================================>       ] 81%
Sleep duration complete.
```

### Interrupted Operation
```
Generated sleep time: 120000ms
Sleeping: [========>                               ] 21%
Program interrupted by keyboard input.
```

### Error Messages
```
Error: Usage: program <min> <max>
```
```
Error: minimum value cannot be greater than maximum value
```

## Technical Details

### Terminal Mode

The program uses **raw terminal mode** to detect keypresses without requiring Enter:
- Characters are not echoed to the screen
- No line buffering occurs
- Input is detected immediately
- Terminal state is always restored on exit

### Progress Threshold

Progress bar is shown when sleep duration ≥ 20,000ms (20 seconds). This threshold is configurable in the source code:

```go
const progressThreshold = 20000  // milliseconds
```

### Random Number Generation

- Uses `math/rand` with time-based seeding
- Seed: `time.Now().UnixNano()`
- Range: `[min, max]` (inclusive on both ends)
- Distribution: Uniform random

### Build Number Format

Build numbers use the format `MMDDHHMM` in UTC:
- `MM` - Month (01-12)
- `DD` - Day (01-31)
- `HH` - Hour (00-23)
- `MM` - Minute (00-59)

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

### Problem: "strconv.Atoi: parsing error"
**Solution:** Ensure arguments are valid integers (no letters, decimals, or special characters)

### Problem: Program doesn't respond to keypress
**Solution:**
- Ensure terminal supports raw mode (most do)
- Try pressing different keys
- Check that stdin is connected to a terminal (not redirected)

### Problem: "minimum value cannot be greater than maximum value"
**Solution:** Swap the arguments so that `<min>` is less than or equal to `<max>`

### Problem: Build fails with "package golang.org/x/term: not found"
**Solution:** Run `go mod download` to fetch dependencies

## Performance Notes

- **Startup Time**: Negligible (<10ms)
- **Memory Usage**: ~5-10 MB
- **CPU Usage**: Near zero during sleep (polling for keyboard only)
- **Accuracy**: ±1 second for display updates, actual sleep is precise

## Exit Codes

- `0` - Normal completion or keyboard interrupt
- `1` - Error (invalid arguments, terminal setup failure)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines, code standards, and contribution process.

## Documentation

- **README.md** (this file) - User documentation
- **CONTRIBUTING.md** - Development guidelines
- **VERSION_INFO.md** - Versioning system documentation
- **BUILD_TEMPLATE.md** - Build script reuse guide
- **GODOC.md** - Auto-generated API documentation

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
