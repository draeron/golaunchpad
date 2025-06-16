package button

import (
	"github.com/draeron/gopkgs/color"
	seven_bits "github.com/draeron/gopkgs/color/7bits"
)

type ColorMap map[Button]color.Color

func (c ColorMap) ApplyFrom(other ColorMap) ColorMap {
	for k, v := range other {
		c[k] = v
	}
	return c
}

func isSameColor(a, b color.Color) bool {
	a7 := seven_bits.FromColor(a)
	b7 := seven_bits.FromColor(b)
	return a7.IsSame(b7)
}

func (c ColorMap) DiffFrom(cmap ColorMap) ColorMap {
	out := ColorMap{}

	for btn, col := range c {
		if other, ok := cmap[btn]; ok {
			if !isSameColor(col, other) {
				out[btn] = col
			}
		} else { // missing color considered changed too
			out[btn] = col
		}
	}
	return out
}
