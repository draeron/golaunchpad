package minimk3

//go:generate go-enum -f=$GOFILE

/*
	Layout x ENUM(
  Session
  Drum
  Keys
  User
  DawFaders
  Programmer
)
*/
type Layout int

var layoutByteValue = map[Layout]byte{
	LayoutSession:    0x00,
	LayoutDrum:       0x04,
	LayoutKeys:       0x05,
	LayoutUser:       0x06,
	LayoutDawFaders:  0x0d,
	LayoutProgrammer: 0x07f,
}

func (l Layout) value() byte {
	if byt, ok := layoutByteValue[l]; ok {
		return byt
	} else {
		log.Panicf("unsupported layout: %s", l.String())
	}
	return 0
}
