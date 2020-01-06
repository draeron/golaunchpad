package minimk3

import (
	dev "github.com/draeron/golaunchpad/pkg/device/event"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
)

func (m *Controller) Subscribe(channel chan<- event.Event) {
	m.subscribers = append(m.subscribers, channel)
}

func (m *Controller) onDeviceEvent(in dev.Event) {
	evt := event.FromMidiEvent(in)
	for _, c := range m.subscribers {
		c <- evt
	}
}
