package minimk3

import "regexp"

var (
	rxDAW     = regexp.MustCompile(`Launchpad Mini MK3 LPMiniMK3 DAW.*`)
	rxMidiIn  = regexp.MustCompile(`Launchpad Mini MK3 LPMiniMK3 MIDI Out`)
	rxMidiOut = regexp.MustCompile(`Launchpad Mini MK3 LPMiniMK3 MIDI In`)
)
