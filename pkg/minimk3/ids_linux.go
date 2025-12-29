package minimk3

import "regexp"

var (
	rxDAW     = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 (MIDI 2|LPMiniMK3 DA)`)
	rxMidiIn  = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 (MIDI 1|LPMiniMK3 MI)`)
	rxMidiOut = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 (MIDI 1|LPMiniMK3 MI)`)
)
