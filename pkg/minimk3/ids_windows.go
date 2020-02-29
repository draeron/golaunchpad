package minimk3

import "regexp"

var (
	rxDAW     = regexp.MustCompile(`LPMiniMK3 MIDI ([[:digit:]])+$`)
	rxMidiIn  = regexp.MustCompile(`MIDIIN([0-9]+) \(LPMiniMK3 MIDI\)`)
	rxMidiOut = regexp.MustCompile(`MIDIOUT([0-9]+) \(LPMiniMK3 MIDI\)`)
)
