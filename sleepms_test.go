// Copyright (c) 2026 Doug Stewart

package main

import (
	"fmt"
	"testing"
	"time"
)

func TestParseDurationArg(t *testing.T) {
	tests := []struct {
		input   string
		wantMs  int
		wantErr bool
	}{
		{"5000", 5000, false},
		{"0", 0, false},
		{"5s", 5000, false},
		{"500ms", 500, false},
		{"1m", 60000, false},
		{"1m30s", 90000, false},
		{"1.5s", 1500, false},
		{"-1", 0, true},
		{"-1s", 0, true},
		{"abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseDurationArg(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDurationArg(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && got != tt.wantMs {
				t.Errorf("parseDurationArg(%q) = %d, want %d", tt.input, got, tt.wantMs)
			}
		})
	}
}

func TestParsePositionalArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantMin int
		wantMax int
		wantErr bool
	}{
		{"plain ms", []string{"1000", "5000"}, 1000, 5000, false},
		{"duration strings", []string{"1s", "5s"}, 1000, 5000, false},
		{"equal min max", []string{"3000", "3000"}, 3000, 3000, false},
		{"mixed units", []string{"500ms", "2s"}, 500, 2000, false},
		{"min greater than max", []string{"5000", "1000"}, 0, 0, true},
		{"too few args", []string{"1000"}, 0, 0, true},
		{"too many args", []string{"1000", "2000", "3000"}, 0, 0, true},
		{"invalid min", []string{"abc", "2000"}, 0, 0, true},
		{"invalid max", []string{"1000", "abc"}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, err := parsePositionalArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePositionalArgs(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if err == nil {
				if gotMin != tt.wantMin {
					t.Errorf("min = %d, want %d", gotMin, tt.wantMin)
				}
				if gotMax != tt.wantMax {
					t.Errorf("max = %d, want %d", gotMax, tt.wantMax)
				}
			}
		})
	}
}

func TestGenerateSleepDuration(t *testing.T) {
	t.Run("equal min and max", func(t *testing.T) {
		got := generateSleepDuration(3000, 3000)
		if got != 3000 {
			t.Errorf("got %d, want 3000", got)
		}
	})
	t.Run("result within range", func(t *testing.T) {
		for i := 0; i < 200; i++ {
			got := generateSleepDuration(100, 500)
			if got < 100 || got > 500 {
				t.Fatalf("got %d, outside [100, 500]", got)
			}
		}
	})
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms   int
		want string
	}{
		{0, "0ms"},
		{250, "250ms"},
		{999, "999ms"},
		{1000, "1.000s"},
		{3500, "3.500s"},
		{59999, "59.999s"},
		{60000, "1m 0.000s"},
		{90234, "1m 30.234s"},
		{3600000, "60m 0.000s"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%dms", tt.ms), func(t *testing.T) {
			got := formatDuration(tt.ms)
			if got != tt.want {
				t.Errorf("formatDuration(%d) = %q, want %q", tt.ms, got, tt.want)
			}
		})
	}
}

func TestFormatRemaining(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "0s"},
		{500 * time.Millisecond, "1s"}, // rounds up
		{45 * time.Second, "45s"},
		{59 * time.Second, "59s"},
		{60 * time.Second, "1m"},
		{90 * time.Second, "1m 30s"},
		{120 * time.Second, "2m"},
		{3661 * time.Second, "61m 1s"},
	}
	for _, tt := range tests {
		t.Run(tt.d.String(), func(t *testing.T) {
			got := formatRemaining(tt.d)
			if got != tt.want {
				t.Errorf("formatRemaining(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

// assertAllBytes fails the test if any byte in s differs from want.
func assertAllBytes(t *testing.T, s string, want byte) {
	t.Helper()
	for i := range s {
		if s[i] != want {
			t.Errorf("position %d: got %q, want %q", i, s[i], want)
			return
		}
	}
}

func TestRenderProgressBar(t *testing.T) {
	t.Run("length is always progressBarWidth", func(t *testing.T) {
		for _, pct := range []int{0, 25, 50, 75, 100} {
			bar := renderProgressBar(pct)
			if len(bar) != progressBarWidth {
				t.Errorf("pct=%d: len=%d, want %d", pct, len(bar), progressBarWidth)
			}
		}
	})
	t.Run("0 percent shows leading arrow then spaces", func(t *testing.T) {
		bar := renderProgressBar(0)
		if bar[0] != '>' {
			t.Errorf("position 0: got %q, want '>'", bar[0])
		}
		assertAllBytes(t, bar[1:], ' ')
	})
	t.Run("100 percent is all filled", func(t *testing.T) {
		assertAllBytes(t, renderProgressBar(100), '=')
	})
	t.Run("50 percent shows fill then arrow", func(t *testing.T) {
		bar := renderProgressBar(50)
		half := progressBarWidth / 2
		if bar[half-1] != '=' {
			t.Errorf("position %d: got %q, want '='", half-1, bar[half-1])
		}
		if bar[half] != '>' {
			t.Errorf("position %d: got %q, want '>'", half, bar[half])
		}
	})
}
