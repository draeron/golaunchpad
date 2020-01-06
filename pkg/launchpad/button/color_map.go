package button

import "image/color"

type ColorMap map[Button]color.Color

func (c ColorMap) ApplyFrom(other ColorMap) ColorMap {
	for k, v := range other {
		c[k] = v
	}
	return c
}

func (c ColorMap) DiffFrom(other ColorMap) ColorMap {
	out := ColorMap{}
	for btn, col := range c {
		if other, ok := other[btn]; ok {
			if other != col {
				out[btn] = col
			}
		} else {
			out[btn] = col
		}
	}
	return out
}
