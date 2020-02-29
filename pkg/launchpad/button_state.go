package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"image/color"
	"time"
)

type ButtonState struct {
	Color     color.Color
	pressTime time.Time
}

type ButtonStateMap map[button.Button]*ButtonState

func (bm ButtonStateMap) Press(btn button.Button) {
	if b, ok := bm[btn]; ok {
		b.pressTime = time.Now()
	}
}

func (bm ButtonStateMap) Release(btn button.Button) {
	if b, ok := bm[btn]; ok {
		b.pressTime = time.Time{}
	}
}

func (bm ButtonStateMap) HoldTime(btn button.Button) time.Duration {
	if b, ok := bm[btn]; ok && !b.pressTime.IsZero() {
		return time.Now().Sub(b.pressTime)
	}
	return 0
}

func (bm ButtonStateMap) IsPressed(btn button.Button) bool {
	return bm.HoldTime(btn) > 0
}

func (bm ButtonStateMap) IsHold(btn button.Button, threshold time.Duration) bool {
	return bm.HoldTime(btn) > threshold
}

func (bm ButtonStateMap) SetColor(btn button.Button, col color.Color) {
	val := ButtonState{Color: col}
	if b, ok := bm[btn]; ok {
		val.pressTime = b.pressTime
	}
	bm[btn] = &val
}

func (bm ButtonStateMap) SetColorsMap(mapp button.ColorMap) {
	for k, v := range mapp {
		bm.SetColor(k, v)
	}
}

func (bm ButtonStateMap) SetColors(mask Mask, col color.Color) {
	for k, _ := range mask {
		bm.SetColor(k, col)
	}
}

func (bm ButtonStateMap) Color(button button.Button) color.Color {
	if b, ok := bm[button]; ok {
		return b.Color
	}
	return color.Black
}

func (bm ButtonStateMap) ResetPressed() {
	for _, b := range bm {
		b.pressTime = time.Time{}
	}
}
