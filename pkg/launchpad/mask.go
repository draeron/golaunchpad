package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkgs/color/7bits"
)

//go:generate go-enum -f=$GOFILE --noprefix
/*
	Mask x ENUM(
	MaskPad
	MaskArrows
	MaskRows
	MaskModes
	MaskTop
	MaskAll
)
*/
type MaskPreset int

type Mask map[button.Button]bool

func (m Mask) Intersect(mapp ButtonStateMap) button.ColorMap {
	out := make(button.ColorMap)
	for k, v := range m {
		if v {
			if cl := mapp.Get(k); cl != nil {
				out[k] = seven_bits.FromColor(cl.Color)
			}
		}
	}
	return out
}

func (m Mask) MergePreset(masks ...MaskPreset) Mask {
	out := m
	for _, mask := range masks {
		out.Merge(mask.Mask())
	}
	return out
}

func (m Mask) Merge(masks ...Mask) Mask {
	out := m
	for _, mask := range masks {
		for b, v := range mask {
			out[b] = v
		}
	}
	return out
}

func (mp MaskPreset) Mask() Mask {
	m := Mask{}

	switch mp {
	case MaskAll:
		for _, b := range button.Values() {
			m[b] = true
		}
	case MaskTop:
		for _, b := range button.Columns() {
			m[b] = true
		}
	case MaskArrows:
		for _, b := range button.Arrows() {
			m[b] = true
		}
	case MaskPad:
		for _, b := range button.Pads() {
			m[b] = true
		}
	case MaskRows:
		for _, b := range button.Rows() {
			m[b] = true
		}
	case MaskModes:
		for _, b := range button.Modes() {
			m[b] = true
		}
	}

	return m
}
