package minimk3

//go:generate go-enum -f=$GOFILE

// Layout x ENUM(
/*
  Up
  Down
  Left
  Right
  Session
  Drums
  Keys
  User
	Logo
  Row1
  Row2
  Row3
  Row4
  Row5
  Row6
  Row7
  StopSoloMute
  Pad11
  Pad12
  Pad13
  Pad14
  Pad15
  Pad16
  Pad17
  Pad18
  Pad21
  Pad22
  Pad23
  Pad24
  Pad25
  Pad26
  Pad27
  Pad28
  Pad31
  Pad32
  Pad33
  Pad34
  Pad35
  Pad36
  Pad37
  Pad38
  Pad41
  Pad42
  Pad43
  Pad44
  Pad45
  Pad46
  Pad47
  Pad48
  Pad51
  Pad52
  Pad53
  Pad54
  Pad55
  Pad56
  Pad57
  Pad58
  Pad61
  Pad62
  Pad63
  Pad64
  Pad65
  Pad66
  Pad67
  Pad68
  Pad71
  Pad72
  Pad73
  Pad74
  Pad75
  Pad76
  Pad77
  Pad78
  Pad81
  Pad82
  Pad83
  Pad84
  Pad85
  Pad86
  Pad87
  Pad88
*/
// )
type Btn int

var reverseBtnIdLookup = map[byte]Btn{}

func init() {
	for _, b := range BtnValues() {
		reverseBtnIdLookup[b.Id()] = b
	}
}

func (b Btn) Id() byte {
	switch {
	case b >= BtnUp && b <= BtnLogo:
		return byte(91 + b - BtnUp)
	case b >= BtnRow1 && b <= BtnStopSoloMute:
		return byte(89 - (b-BtnRow1)*10)
	case b >= BtnPad11 && b <= BtnPad88:
		diff := b - BtnPad11
		return byte(81 + (diff % 8) - (diff / 8 * 10))
	}
	log.Panicf("unknown btn value: %s", b)
	return 0
}

func (b Btn) IsPad() bool {
	return b >= BtnPad11 && b <= BtnPad88
}

func (b Btn) IsRow() bool {
	return b >= BtnRow1 && b <= BtnStopSoloMute
}

func (b Btn) IsCol() bool {
	return b >= BtnUp && b <= BtnUser
}

func (b Btn) IsArrow() bool {
	return b >= BtnUp && b <= BtnRight
}

func btnFromId(id byte) Btn {
	if btn, ok := reverseBtnIdLookup[id]; ok {
		return btn
	} else {
		log.Panicf("invalid button midi id: %v", id)
	}
	return BtnStopSoloMute
}

func BtnValues() (btns []Btn) {
	for _, b := range _BtnValue {
		btns = append(btns, b)
	}
	return
}

/*
func Colors() []PaletteColor {
	var cols []PaletteColor
	for cl, _ := range _PaletteColorMap {
		cols = append(cols, cl)
	}
	return cols
}
*/
