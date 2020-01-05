package device

import (
	"gitlab.com/gomidi/rtmididrv"
	"regexp"
)

func Detect(rxDaw, midiInRx, midiOutRx *regexp.Regexp) (*MiniMk3, error) {
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
		if rxDaw.MatchString(in.String()) {
			pad.dawIn = in
		} else if midiInRx.MatchString(in.String()) {
			pad.midiIn = in
		}
	}

	outs, err := drv.Outs()
	if err != nil {
		return nil, err
	}

	for _, out := range outs {
		if rxDaw.MatchString(out.String()) {
			pad.dawOut = out
		} else if midiOutRx.MatchString(out.String()) {
			pad.midiOut = out
		}
	}

	return &pad, nil
}
