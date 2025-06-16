package state

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/mask"
	"github.com/draeron/gopkgs/color"
	seven_bits "github.com/draeron/gopkgs/color/7bits"

	"sync"
	"time"
)

type ButtonState struct {
	Color     color.Color
	pressTime time.Time
}

type Map struct {
	data  map[button.Button]*ButtonState
	mutex *sync.RWMutex
}

func NewButtonStateMap() Map {
	return Map{
		data:  map[button.Button]*ButtonState{},
		mutex: &sync.RWMutex{},
	}
}

func (m *Map) Get(btn button.Button) *ButtonState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if b, ok := m.data[btn]; ok {
		return b
	} else {
		return nil
	}
}

func (m *Map) Press(btn button.Button) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if b, ok := m.data[btn]; ok {
		b.pressTime = time.Now()
	}
}

func (m *Map) Release(btn button.Button) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if b, ok := m.data[btn]; ok {
		b.pressTime = time.Time{}
	}
}

func (m *Map) HoldTime(btn button.Button) time.Duration {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if b, ok := m.data[btn]; ok && !b.pressTime.IsZero() {
		return time.Now().Sub(b.pressTime)
	}
	return 0
}

func (m *Map) IsPressed(btn button.Button) bool {
	return m.HoldTime(btn) > 0
}

func (m *Map) IsHold(btn button.Button, threshold time.Duration) bool {
	return m.HoldTime(btn) > threshold
}

func (m *Map) SetColor(btn button.Button, col color.Color) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val := ButtonState{Color: col}
	if b, ok := m.data[btn]; ok {
		val.pressTime = b.pressTime
	}
	m.data[btn] = &val
}

func (m *Map) Intersect(mask mask.Buttons) button.ColorMap {
	out := make(button.ColorMap)
	for k, v := range mask {
		if v {
			if cl := m.Get(k); cl != nil {
				out[k] = seven_bits.FromColor(cl.Color)
			}
		}
	}
	return out
}

func (m *Map) SetColorsMap(mapp button.ColorMap) {
	for k, v := range mapp {
		m.SetColor(k, v)
	}
}

func (m *Map) SetColors(mask mask.Buttons, col color.Color) {
	for k, _ := range mask {
		m.SetColor(k, col)
	}
}

func (m *Map) SetAllColors(col color.Color) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for key, val := range m.data {
		newv := *val // create a copy
		newv.Color = col
		m.data[key] = &newv
	}
}

func (m *Map) Color(button button.Button) color.Color {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if b, ok := m.data[button]; ok {
		return b.Color
	}
	return color.Black
}

func (m *Map) ResetPressed() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, b := range m.data {
		b.pressTime = time.Time{}
	}
}
