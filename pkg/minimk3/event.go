package minimk3

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	dev "github.com/draeron/gopkgs/midi/event"
)

func (m *Controller) Subscribe(channel chan<- event.Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.subscribers = append(m.subscribers, channel)
}

func (m *Controller) onDeviceEvent(in dev.Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	evt := event.FromMidiEvent(in)
	for _, c := range m.subscribers {
		c <- evt
	}
}
