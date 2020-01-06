package launchpad

import "github.com/draeron/golaunchpad/pkg/launchpad/button"

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

/*
	Remove
*/
func (m Mask) Intersect(mapp button.ColorMap) button.ColorMap {
	out := make(button.ColorMap)
	for k, v := range m {
		if v {
			if cl, ok := mapp[k]; ok {
				out[k] = cl
			}
		}
	}
	return out
}

func (m Mask) Merge(masks... Mask) Mask {
	out := Mask{}
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
