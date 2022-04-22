package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"image/color"
	"sync"
	"time"
)

type ButtonState struct {
	Color     color.Color
	pressTime time.Time
}

type ButtonStateMap struct {
	data  map[button.Button]*ButtonState
	mutex *sync.RWMutex
}

func NewButtonStateMap() ButtonStateMap {
	return ButtonStateMap{
		data:  map[button.Button]*ButtonState{},
		mutex: &sync.RWMutex{},
	}
}

func (bm *ButtonStateMap) Get(btn button.Button) *ButtonState {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	if b, ok := bm.data[btn]; ok {
		return b
	} else {
		return nil
	}
}

func (bm *ButtonStateMap) Press(btn button.Button) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if b, ok := bm.data[btn]; ok {
		b.pressTime = time.Now()
	}
}

func (bm *ButtonStateMap) Release(btn button.Button) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if b, ok := bm.data[btn]; ok {
		b.pressTime = time.Time{}
	}
}

func (bm *ButtonStateMap) HoldTime(btn button.Button) time.Duration {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	if b, ok := bm.data[btn]; ok && !b.pressTime.IsZero() {
		return time.Now().Sub(b.pressTime)
	}
	return 0
}

func (bm *ButtonStateMap) IsPressed(btn button.Button) bool {
	return bm.HoldTime(btn) > 0
}

func (bm *ButtonStateMap) IsHold(btn button.Button, threshold time.Duration) bool {
	return bm.HoldTime(btn) > threshold
}

func (bm *ButtonStateMap) SetColor(btn button.Button, col color.Color) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	val := ButtonState{Color: col}
	if b, ok := bm.data[btn]; ok {
		val.pressTime = b.pressTime
	}
	bm.data[btn] = &val
}

func (bm *ButtonStateMap) SetColorsMap(mapp button.ColorMap) {
	for k, v := range mapp {
		bm.SetColor(k, v)
	}
}

func (bm *ButtonStateMap) SetColors(mask Mask, col color.Color) {
	for k, _ := range mask {
		bm.SetColor(k, col)
	}
}

func (bm *ButtonStateMap) SetAllColors(col color.Color) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	for key, val := range bm.data {
		newv := *val // create a copy
		newv.Color = col
		bm.data[key] = &newv
	}
}

func (bm *ButtonStateMap) Color(button button.Button) color.Color {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	if b, ok := bm.data[button]; ok {
		return b.Color
	}
	return color.Black
}

func (bm *ButtonStateMap) ResetPressed() {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	for _, b := range bm.data {
		b.pressTime = time.Time{}
	}
}
