//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package minimk3

import "regexp"

var (
	rxDAW     = regexp.MustCompile(`.*`)
	rxMidiIn  = regexp.MustCompile(`.*`)
	rxMidiOut = regexp.MustCompile(`.*`)
)
