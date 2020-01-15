package seven_bits

import "image/color"

// color in 7 bits
type SevenColor struct {
	R, G, B uint8
}

func FromColor(color color.Color) SevenColor {
	r,g,b,_ := color.RGBA()
	return SevenColor{
		R: uint8(r >> 10),
		G: uint8(g >> 10),
		B: uint8(b >> 10),
	}
}

/*
	Convert 7 to 16 bits channels
*/
func (bc SevenColor) RGBA() (r, g, b, a uint32) {
	return r << 9, r << 9, r << 9, 0xffff
}

func (bc SevenColor) IsSame(other SevenColor) bool {
	return bc.R == other.R && bc.G == other.G && bc.B == other.B
}
