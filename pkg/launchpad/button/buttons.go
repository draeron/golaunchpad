package button

import (
	"log"
	"sort"
)

//go:generate go-enum -f=$GOFILE --noprefix

/*
	Button x ENUM(
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
	Pad21
	Pad31
	Pad41
	Pad51
	Pad61
	Pad71
	Pad81
	Pad12
	Pad22
	Pad32
	Pad42
	Pad52
	Pad62
	Pad72
	Pad82
	Pad13
	Pad23
	Pad33
	Pad43
	Pad53
	Pad63
	Pad73
	Pad83
	Pad14
	Pad24
	Pad34
	Pad44
	Pad54
	Pad64
	Pad74
	Pad84
	Pad15
	Pad25
	Pad35
	Pad45
	Pad55
	Pad65
	Pad75
	Pad85
	Pad16
	Pad26
	Pad36
	Pad46
	Pad56
	Pad66
	Pad76
	Pad86
	Pad17
	Pad27
	Pad37
	Pad47
	Pad57
	Pad67
	Pad77
	Pad87
	Pad18
	Pad28
	Pad38
	Pad48
	Pad58
	Pad68
	Pad78
	Pad88
) */
type Button int

type Buttons []Button

var reverseIdLookup = map[byte]Button{}

func init() {
	for _, b := range Values() {
		reverseIdLookup[b.Id()] = b
	}
}

func (b Button) Id() byte {
	switch {
	case b >= Up && b <= Logo:
		return byte(91 + b - Up)
	case b >= Row1 && b <= StopSoloMute:
		return byte(89 - (b-Row1)*10)
	case b >= Pad11 && b <= Pad88:
		diff := b - Pad11
		return byte(81 + (diff % 8) - (diff / 8 * 10))
	}
	log.Panicf("unknown  value: %s", b)
	return 0
}

func (b Button) IsPad() bool {
	return b >= Pad11 && b <= Pad88
}

func (b Button) IsRow() bool {
	return b >= Row1 && b <= StopSoloMute
}

func (b Button) IsCol() bool {
	return b >= Up && b <= User
}

func (b Button) IsMode() bool {
	return b >= Session && b <= User
}

func (b Button) IsArrow() bool {
	return b >= Up && b <= Right
}

func (b Button) Coord() (x, y int) {
	switch {
	case b.IsCol():
		diff := b - Up
		return int(diff % 9), 0
	case b.IsRow():
		diff := b - Row1
		return 9, int(1 + diff)
	case b.IsPad():
		diff := b - Pad11
		return int(diff % 8), int(1 + diff/8)
	}

	return 0, 9 // default to logo
}

func (b Button) IsValid() bool {
	return b >= Up && b <= Pad88
}

func Pads() (s []Button) {
	for b := Pad11; b <= Pad88; b++ {
		s = append(s, b)
	}
	return
}

func Rows() (s []Button) {
	for b := Row1; b <= StopSoloMute; b++ {
		s = append(s, b)
	}
	return
}

func Columns() (s []Button) {
	for b := Up; b <= Logo; b++ {
		s = append(s, b)
	}
	return
}

func Modes() (s []Button) {
	for b := Session; b <= User; b++ {
		s = append(s, b)
	}
	return
}

func Arrows() (s []Button) {
	for b := Up; b <= Right; b++ {
		s = append(s, b)
	}
	return
}

func FromMidiId(id byte) Button {
	if btn, ok := reverseIdLookup[id]; ok {
		return btn
	} else {
		log.Panicf("invalid button midi id: %v", id)
	}
	return Logo
}

func Values() (s Buttons) {
	for _, b := range _ButtonValue {
		s = append(s, b)
	}
	sort.Sort(s)
	return
}

func FromXY(x, y int) Button {
	if x < 0 || y < 0 || x > 8 || y > 8 {
		return Button(-1)
	}

	if y == 0 { // column
		return Up + Button(x)
	} else if x == 8 { // rows
		return Row1 + Button(y-1)
	} else { // pad
		y-- // remove coord from cols
		return Pad11 + Button(x+y*8)
	}
}

func FromPadXY(x, y int) Button {
	return FromXY(x, y-1)
}

func (b Buttons) Len() int {
	return len(b)
}

func (b Buttons) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b Buttons) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
