package minimk3

import (
	dev "github.com/draeron/golaunchpad/pkg/device/event"
	mk3 "github.com/draeron/golaunchpad/pkg/minimk3/event"
)

func (m *Controller) Subscribe(channel chan<- mk3.Event) {
	m.subscribers = append(m.subscribers, channel)
}

func (m *Controller) onDeviceEvent(in dev.Event) {
	evt := mk3.FromMidiEvent(in)
	for _, c := range m.subscribers {
		c <- evt
	}
}
