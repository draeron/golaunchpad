package device

import (
	"gitlab.com/gomidi/rtmididrv"
	"regexp"
)

var (
	DAW_ID      = regexp.MustCompile(`LPMiniMK3 MIDI ([[:digit:]])+$`)
	MIDI_IN_ID  = regexp.MustCompile(`MIDIIN([0-9]+) \(LPMiniMK3 MIDI\)`)
	MIDI_OUT_ID = regexp.MustCompile(`MIDIOUT([0-9]+) \(LPMiniMK3 MIDI\)`)
)

func Detect() (*MiniMk3, error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, err
	}

	pad := MiniMk3{}

	ins, err := drv.Ins()
	if err != nil {
		return nil, err
	}

	for _, in := range ins {
		if DAW_ID.MatchString(in.String()) {
			pad.dawIn = in
		} else if MIDI_IN_ID.MatchString(in.String()) {
			pad.midiIn = in
		}
	}

	outs, err := drv.Outs()
	if err != nil {
		return nil, err
	}

	for _, out := range outs {
		if DAW_ID.MatchString(out.String()) {
			pad.dawOut = out
		} else if MIDI_OUT_ID.MatchString(out.String()) {
			pad.midiOut = out
		}
	}

	return &pad, nil
}
