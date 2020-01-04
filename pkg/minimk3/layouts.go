package minimk3

//go:generate go-enum -f=$GOFILE 

// Layout x ENUM(
/*
  Session
  Drum
  Keys
  User
  DawFaders
  Programmer
*/
// )
type Layout int

func (l Layout) value() byte {
  switch l {
  case LayoutSession:
    return 0x00
  case LayoutDrum:
    return 0x04
  case LayoutKeys:
    return 0x05
  case LayoutUser:
    return 0x06
  case LayoutDawFaders:
    return 0x0d
  case LayoutProgrammer:
    return 0x07f
  default:
    log.DPanicf("unsupported layout: %s", l.String())
    return 0
  }
}