package mask

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

//go:generate go tool github.com/abice/go-enum -f=$GOFILE --noprefix --values
/*
	Mask x ENUM(
	None
	Pad = 0x01
	Arrows = 0x02
	Modes = 0x04
	Top = 0x06
	Rows = 0x10
	MaskAll = 0xFF
	)
*/
type Preset int

var presetButtonMapping = map[Preset][]button.Button{
	Top:    append(button.Arrows(), button.Modes()...),
	Arrows: button.Arrows(),
	Rows:   button.Rows(),
	Modes:  button.Modes(),
	Pad:    button.Pads(),
	// MaskAll:    button.Values(),
}

func (mp Preset) Mask() Buttons {
	m := Buttons{}

	if mp == MaskAll {
		for _, b := range button.Values() {
			m[b] = true
		}
		return m
	}

	for _, preset := range PresetValues() {
		if mp.Has(preset) {
			for _, btn := range presetButtonMapping[preset] {
				m[btn] = true
			}
		}
	}

	return m
}

func (m Preset) Has(preset Preset) bool {
	return (m & preset) == preset
}
