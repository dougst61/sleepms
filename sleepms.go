// Copyright (c) 2026 Doug Stewart

// Package main implements a sleep utility that waits for a random duration
// between a minimum and maximum value.
//
// Durations may be specified as plain integers (milliseconds) or as Go
// duration strings (e.g. 5s, 1m30s, 500ms). For sleep durations of 20
// seconds or longer a real-time progress bar is shown, including time
// remaining. The sleep can be interrupted at any time by pressing any key
// or by sending SIGINT / SIGTERM.
//
// Usage:
//
//	sleepms [options] <min> <max>
//
// Examples:
//
//	sleepms 1000 5000       # between 1 s and 5 s (plain milliseconds)
//	sleepms 1s 5s           # same, using duration strings
//	sleepms 1m 2m           # between 1 and 2 minutes (progress bar shown)
//	sleepms -q 500 2000     # quiet: no output, exit code signals completion
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/term"
)

// Version information — overwritten at build time via -ldflags.
var (
	AppName         = "sleepms"
	VersionMajor    = "1"
	VersionMinor    = "1"
	VersionRevision = "0"
	BuildNumber     = "00000000"
	BuildTime       = "unknown"
)

// Exit codes returned by the program.
const (
	exitSuccess   = 0 // sleep completed normally
	exitInterrupt = 1 // sleep was cut short by keypress or signal
	exitError     = 2 // bad arguments or terminal setup failure
)

// Progress bar configuration.
const (
	progressBarWidth  = 40
	progressThreshold = 20000 // 20 seconds: minimum duration to show a progress bar
	progressUpdateHz  = time.Second
)

// quiet suppresses all printed output when true (-q / --quiet flag).
var quiet bool

func main() {
	os.Exit(run())
}

// run is the real entry point; returning an int lets deferred cleanup run
// before os.Exit is called in main.
func run() int {
	flag.BoolVar(&quiet, "q", false, "suppress all output")
	flag.BoolVar(&quiet, "quiet", false, "suppress all output")
	showVersion := flag.Bool("version", false, "print version and exit")
	showVersionShort := flag.Bool("v", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <min> <max>\n\n", AppName)
		fmt.Fprintf(os.Stderr, "  min, max  Duration as milliseconds (5000) or a Go duration string (5s, 1m30s)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion || *showVersionShort {
		fmt.Printf("%s version %s.%s.%s build %s\n", AppName, VersionMajor, VersionMinor, VersionRevision, BuildNumber)
		fmt.Printf("Built: %s\n", BuildTime)
		return exitSuccess
	}

	minMs, maxMs, err := parsePositionalArgs(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.Usage()
		return exitError
	}

	sleepMs := generateSleepDuration(minMs, maxMs)
	logf("Sleeping for %s (%dms)\n", formatDuration(sleepMs), sleepMs)

	// Attempt raw-mode terminal setup for keypress interrupt. If stdin is not
	// a TTY (e.g. running inside a pipe), skip gracefully.
	isTTY := term.IsTerminal(int(os.Stdin.Fd()))
	var cleanup func()
	if isTTY {
		if cleanup, err = setupTerminal(); err != nil {
			isTTY = false
		}
	}
	if cleanup != nil {
		defer cleanup()
	}

	// Catch SIGINT and SIGTERM so the terminal is restored before we exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if waitWithProgress(sleepMs, isTTY, sigCh) {
		logf("Sleep interrupted.\n")
		return exitInterrupt
	}
	logf("Sleep complete.\n")
	return exitSuccess
}

// logf writes to stdout unless quiet mode is active.
func logf(format string, args ...any) {
	if !quiet {
		fmt.Printf(format, args...)
	}
}

// parsePositionalArgs validates the two positional duration arguments.
func parsePositionalArgs(args []string) (minMs, maxMs int, err error) {
	if len(args) != 2 {
		return 0, 0, fmt.Errorf("expected 2 arguments (min max), got %d", len(args))
	}
	minMs, err = parseDurationArg(args[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid min: %w", err)
	}
	maxMs, err = parseDurationArg(args[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid max: %w", err)
	}
	if minMs > maxMs {
		return 0, 0, fmt.Errorf("min (%s) must not exceed max (%s)", formatDuration(minMs), formatDuration(maxMs))
	}
	return minMs, maxMs, nil
}

// parseDurationArg accepts either a plain integer (milliseconds) or a Go
// duration string such as "5s", "1m30s", or "500ms".
func parseDurationArg(s string) (int, error) {
	if n, err := strconv.Atoi(s); err == nil {
		if n < 0 {
			return 0, fmt.Errorf("duration must be non-negative")
		}
		return n, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("use milliseconds (e.g. 5000) or a Go duration string (e.g. 5s, 1m30s, 500ms)")
	}
	if d < 0 {
		return 0, fmt.Errorf("duration must be non-negative")
	}
	return int(d.Milliseconds()), nil
}

// generateSleepDuration returns a uniform random integer in [minMs, maxMs].
func generateSleepDuration(minMs, maxMs int) int {
	if minMs == maxMs {
		return minMs
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return minMs + rng.Intn(maxMs-minMs+1)
}

// formatDuration converts a millisecond count to a human-readable string.
//
// Examples: 250 → "250ms", 3500 → "3.500s", 90234 → "1m 30.234s"
func formatDuration(ms int) string {
	d := time.Duration(ms) * time.Millisecond
	if d < time.Second {
		return fmt.Sprintf("%dms", ms)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.3fs", d.Seconds())
	}
	m := int(d.Minutes())
	s := d.Seconds() - float64(m)*60
	return fmt.Sprintf("%dm %.3fs", m, s)
}

// formatRemaining converts a remaining duration to a compact string for the
// progress bar (rounds to the nearest second).
//
// Examples: 45s → "45s", 90s → "1m 30s", 120s → "2m"
func formatRemaining(d time.Duration) string {
	d = d.Round(time.Second)
	if d <= 0 {
		return "0s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	if s == 0 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dm %ds", m, s)
}

// setupTerminal puts stdin into raw mode so individual keypresses are
// detected without requiring Enter.
func setupTerminal() (cleanup func(), err error) {
	old, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return func() { term.Restore(int(os.Stdin.Fd()), old) }, nil
}

// listenForKeyPress blocks until a byte is read from stdin, then closes ch.
func listenForKeyPress(ch chan<- struct{}) {
	buf := make([]byte, 1)
	os.Stdin.Read(buf) //nolint:errcheck // intentional: any read (including EOF) triggers interrupt
	close(ch)
}

// waitWithProgress waits for the sleep to complete, a keypress, or a signal.
// Returns true if the sleep was interrupted before completion.
func waitWithProgress(durationMs int, isTTY bool, sigCh <-chan os.Signal) bool {
	timeout := time.After(time.Duration(durationMs) * time.Millisecond)
	keyPress := make(chan struct{})
	if isTTY {
		go listenForKeyPress(keyPress)
	}
	if durationMs >= progressThreshold {
		return waitWithProgressBar(timeout, keyPress, sigCh, durationMs)
	}
	return waitSimple(timeout, keyPress, sigCh)
}

// waitSimple handles short sleeps (no progress bar).
func waitSimple(timeout <-chan time.Time, keyPress <-chan struct{}, sigCh <-chan os.Signal) bool {
	select {
	case <-timeout:
		return false
	case <-keyPress:
		return true
	case <-sigCh:
		return true
	}
}

// waitWithProgressBar handles long sleeps, updating the progress bar each second.
func waitWithProgressBar(timeout <-chan time.Time, keyPress <-chan struct{}, sigCh <-chan os.Signal, durationMs int) bool {
	ticker := time.NewTicker(progressUpdateHz)
	defer ticker.Stop()

	start := time.Now()
	total := time.Duration(durationMs) * time.Millisecond

	for {
		select {
		case <-timeout:
			clearProgressLine()
			return false
		case <-keyPress:
			clearProgressLine()
			return true
		case <-sigCh:
			clearProgressLine()
			return true
		case <-ticker.C:
			if !quiet {
				updateProgressBar(time.Since(start), total)
			}
		}
	}
}

// clearProgressLine erases the current terminal line.
func clearProgressLine() {
	if !quiet {
		fmt.Print("\r\033[K")
	}
}

// updateProgressBar renders and prints the current progress.
//
// Example output:
//
//	Sleeping: [===================>                    ] 49%   1m 2s remaining
func updateProgressBar(elapsed, total time.Duration) {
	pct := min(int((float64(elapsed)/float64(total))*100), 100)
	remaining := total - elapsed
	fmt.Printf("\rSleeping: [%s] %3d%%  %s remaining", renderProgressBar(pct), pct, formatRemaining(remaining))
}

// renderProgressBar returns a fixed-width ASCII progress bar string.
func renderProgressBar(pct int) string {
	filled := (pct * progressBarWidth) / 100
	bar := make([]byte, progressBarWidth)
	for i := range bar {
		switch {
		case i < filled:
			bar[i] = '='
		case i == filled && pct < 100:
			bar[i] = '>'
		default:
			bar[i] = ' '
		}
	}
	return string(bar)
}
