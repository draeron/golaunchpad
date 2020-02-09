package minimk3

import "regexp"

var (
	rxDAW     = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 MIDI 2`)
	rxMidiIn  = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 MIDI 1`)
	rxMidiOut = regexp.MustCompile(`Launchpad Mini MK3:Launchpad Mini MK3 MIDI 1`)
)
