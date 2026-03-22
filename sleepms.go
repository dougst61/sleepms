// Package main implements a sleep utility that waits for a random duration
// between a minimum and maximum value (in milliseconds).
//
// The program generates a random sleep duration, displays it to the user,
// and then sleeps for that duration. The sleep can be interrupted at any
// time by pressing any key on the keyboard.
//
// For sleep durations longer than 1 minute (60,000ms), a progress bar
// is displayed showing the percentage of time elapsed and a visual
// indicator of progress.
//
// Usage:
//
//	sleepms <min> <max>
//
// Arguments:
//
//	min - Minimum sleep duration in milliseconds
//	max - Maximum sleep duration in milliseconds
//
// Example:
//
//	sleepms 1000 5000  # Sleep between 1 and 5 seconds
//	sleepms 60000 120000  # Sleep between 1 and 2 minutes (with progress bar)
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/term" // Terminal control for raw mode (keyboard input without echo)
)

// Version information variables
// These are set at build time using -ldflags
var (
	AppName         = "sleepms"
	VersionMajor    = "1"
	VersionMinor    = "0"
	VersionRevision = "0"
	BuildNumber     = "00000000" // Overwritten at build time
	BuildTime       = "unknown"  // Overwritten at build time
)

// Constants for progress bar configuration
const (
	progressBarWidth   = 40    // Width of the progress bar in characters
	progressThreshold  = 20000 // Threshold in milliseconds (1 minute) for showing progress
	progressUpdateRate = time.Second
)

// main is the entry point of the program. It orchestrates the sleep utility by:
// 1. Parsing and validating command-line arguments
// 2. Generating a random sleep duration
// 3. Setting up the terminal for keyboard input
// 4. Waiting for the duration (or keyboard interrupt) with optional progress display
func main() {
	// Check for version flag
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("%s version %s.%s.%s build %s\n", AppName, VersionMajor, VersionMinor, VersionRevision, BuildNumber)
		fmt.Printf("Built: %s\n", BuildTime)
		return
	}

	// Parse and validate command-line arguments
	minVal, maxVal, err := parseArguments()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Generate a random sleep duration between min and max
	sleepDuration := generateSleepDuration(minVal, maxVal)
	fmt.Printf("Generated sleep time: %dms\n", sleepDuration)

	// Set up terminal for raw mode input (keyboard without echo)
	cleanup, err := setupTerminal()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer cleanup() // Always restore terminal state on exit

	// Wait for the duration or keyboard interrupt
	waitWithProgress(sleepDuration)
}

// parseArguments parses and validates command-line arguments.
// It expects exactly 2 arguments: minimum and maximum sleep durations.
//
// Returns:
//   - minVal: The minimum sleep duration in milliseconds
//   - maxVal: The maximum sleep duration in milliseconds
//   - error: An error if arguments are invalid or missing
func parseArguments() (int, int, error) {
	// Verify that the user provided exactly 2 arguments (min and max).
	// os.Args[0] is the program name, so we need at least 3 total elements.
	if len(os.Args) < 3 {
		return 0, 0, fmt.Errorf("Usage: program <min> <max>")
	}

	// Parse the minimum value from the first command-line argument.
	// strconv.Atoi converts a string to an integer and returns an error
	// if the string is not a valid integer.
	minVal, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return 0, 0, err
	}

	// Parse the maximum value from the second command-line argument.
	maxVal, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, 0, err
	}

	// Validate that the minimum value is not greater than the maximum value.
	// This is a logical error that would cause issues with random number generation.
	if minVal > maxVal {
		return 0, 0, fmt.Errorf("minimum value cannot be greater than maximum value")
	}

	return minVal, maxVal, nil
}

// generateSleepDuration generates a random sleep duration between min and max (inclusive).
// It uses the current time as a seed to ensure different values on each run.
//
// Parameters:
//   - minVal: The minimum sleep duration in milliseconds
//   - maxVal: The maximum sleep duration in milliseconds
//
// Returns:
//   - A random integer between minVal and maxVal (inclusive)
func generateSleepDuration(minVal, maxVal int) int {
	// Create a new random number generator with a seed based on the current time.
	// Using time.Now().UnixNano() ensures that each run produces different random values.
	// Note: We use rand.New() instead of the deprecated rand.Seed() function.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between minVal and maxVal (inclusive).
	// Start with minVal as the default (handles the case where min == max).
	randomNumber := minVal
	if minVal != maxVal {
		// rand.Intn(n) returns a random number in the range [0, n).
		// To get a range [minVal, maxVal], we:
		// 1. Calculate the range: (maxVal - minVal + 1)
		// 2. Get a random number in [0, range)
		// 3. Add minVal to shift the range to [minVal, maxVal]
		randomNumber = minVal + rng.Intn(maxVal-minVal+1)
	}

	return randomNumber
}

// setupTerminal puts the terminal into raw mode for keyboard input without echo.
// In raw mode:
//   - Characters are not echoed to the screen
//   - No line buffering (we get input immediately)
//   - No special character processing (Ctrl+C, etc.)
//
// Returns:
//   - cleanup: A function to restore the terminal to its original state
//   - error: An error if terminal setup fails
func setupTerminal() (func(), error) {
	// Put the terminal into raw mode so we can detect keypresses without
	// requiring the user to press Enter.
	// We save the old terminal state so we can restore it later.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	// Return a cleanup function that restores the terminal state.
	// This is critical - without this, the terminal would remain in raw mode and be unusable.
	cleanup := func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}

	return cleanup, nil
}

// waitWithProgress waits for the specified duration or until a key is pressed.
// For durations longer than 1 minute (60,000ms), it displays a progress bar.
//
// Parameters:
//   - durationMs: The duration to wait in milliseconds
//
// The function uses a select statement to wait for one of three events:
//  1. The timeout channel fires (sleep duration complete)
//  2. The keyPress channel fires (user interrupted)
//  3. The progress ticker fires (time to update progress bar)
func waitWithProgress(durationMs int) {
	// Create a channel that will receive a signal after the specified duration.
	// time.After() returns a channel that delivers the current time after
	// the specified duration has elapsed.
	timeout := time.After(time.Duration(durationMs) * time.Millisecond)

	// Create a channel to signal when a key is pressed.
	// This allows us to interrupt the sleep at any time.
	keyPress := make(chan bool)

	// Launch a goroutine (concurrent function) to listen for keyboard input.
	// This runs in the background while the main thread waits in the select statement.
	go listenForKeyPress(keyPress)

	// Determine if we should show progress based on duration
	showProgress := durationMs >= progressThreshold

	// If showing progress, set up the progress tracking
	if showProgress {
		waitWithProgressBar(timeout, keyPress, durationMs)
	} else {
		waitSimple(timeout, keyPress)
	}
}

// listenForKeyPress listens for a single key press and signals on the provided channel.
// This function blocks until a key is pressed.
//
// Parameters:
//   - keyPress: Channel to signal when a key is pressed
func listenForKeyPress(keyPress chan bool) {
	// Allocate a 1-byte buffer to read a single character.
	buf := make([]byte, 1)
	// Block until a key is pressed. Since we're in raw mode,
	// this will return immediately when any key is pressed.
	os.Stdin.Read(buf)
	// Send a signal on the keyPress channel to notify the main thread.
	keyPress <- true
}

// waitSimple waits for either a timeout or key press without displaying progress.
// This is used for short durations (less than 1 minute).
//
// Parameters:
//   - timeout: Channel that signals when the duration has elapsed
//   - keyPress: Channel that signals when a key is pressed
func waitSimple(timeout <-chan time.Time, keyPress chan bool) {
	select {
	case <-timeout:
		fmt.Println("Sleep duration complete.")
	case <-keyPress:
		fmt.Println("Program interrupted by keyboard input.")
	}
}

// waitWithProgressBar waits for timeout or key press while displaying a progress bar.
// This is used for long durations (1 minute or more).
//
// Parameters:
//   - timeout: Channel that signals when the duration has elapsed
//   - keyPress: Channel that signals when a key is pressed
//   - durationMs: The total duration in milliseconds
func waitWithProgressBar(timeout <-chan time.Time, keyPress chan bool, durationMs int) {
	// Create a ticker that fires every second to update the progress bar.
	progressTicker := time.NewTicker(progressUpdateRate)
	defer progressTicker.Stop()

	// Record the start time so we can calculate elapsed time.
	startTime := time.Now()
	totalDuration := time.Duration(durationMs) * time.Millisecond

	// Main event loop - wait for timeout, key press, or ticker events
	for {
		select {
		case <-timeout:
			// Sleep duration has elapsed
			clearProgressLine()
			fmt.Println("Sleep duration complete.")
			return

		case <-keyPress:
			// User pressed a key to interrupt
			clearProgressLine()
			fmt.Println("Program interrupted by keyboard input.")
			return

		case <-progressTicker.C:
			// Update the progress bar
			elapsed := time.Since(startTime)
			updateProgressBar(elapsed, totalDuration)
		}
	}
}

// clearProgressLine clears the current line in the terminal.
// Uses ANSI escape codes:
//   - \r = carriage return (move cursor to start of line)
//   - \033[K = clear from cursor to end of line
func clearProgressLine() {
	fmt.Print("\r\033[K")
}

// updateProgressBar calculates and displays the current progress.
//
// Parameters:
//   - elapsed: Time that has elapsed since start
//   - total: Total duration to wait
func updateProgressBar(elapsed, total time.Duration) {
	// Calculate the percentage of completion.
	// We convert to float64 for precise division, then back to int for display.
	percentage := int((float64(elapsed) / float64(total)) * 100)

	// Cap the percentage at 100 to handle any timing edge cases.
	if percentage > 100 {
		percentage = 100
	}

	// Render the progress bar string
	bar := renderProgressBar(percentage)

	// Print the progress bar on the same line (using \r to return to start).
	// This creates an animated effect where the bar grows over time.
	// Example output: "Sleeping: [=========>                              ] 25%"
	fmt.Printf("\rSleeping: [%s] %d%%", bar, percentage)
}

// renderProgressBar creates a visual progress bar string based on the percentage.
// The bar is 40 characters wide and shows:
//   - '=' for completed portions
//   - '>' for the current position indicator
//   - ' ' for remaining portions
//
// Parameters:
//   - percentage: The completion percentage (0-100)
//
// Returns:
//   - A string representing the visual progress bar
//
// Example outputs:
//   - 0%:   "                                        "
//   - 25%:  "==========>                             "
//   - 50%:  "====================>                   "
//   - 100%: "========================================"
func renderProgressBar(percentage int) string {
	// Calculate how many characters should be "filled" based on percentage.
	// For example, at 50%, we'd fill 20 out of 40 characters.
	filledWidth := (percentage * progressBarWidth) / 100

	// Build the progress bar string character by character.
	bar := ""
	for i := 0; i < progressBarWidth; i++ {
		if i < filledWidth {
			// Characters before the filled width show as '=' (filled)
			bar += "="
		} else if i == filledWidth && percentage < 100 {
			// The character at the filled width shows as '>' (progress indicator)
			// Don't show '>' at 100% to keep the bar clean
			bar += ">"
		} else {
			// Characters after the filled width show as ' ' (empty)
			bar += " "
		}
	}

	return bar
}
