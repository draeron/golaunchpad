package minimk3

import (
	"fmt"
	"github.com/draeron/golaunchpad/pkg/device/event"
)

type Event struct {
	Type EventType
	Btn  Btn
}

func (e Event) String() string {
	return fmt.Sprintf("Event: %s - %s", e.Type, e.Btn)
}

func toEvent(evt event.Event) Event {
	e := Event{}

	switch evt.Type {
	case event.NoteOn, event.NoteOff:
		e.Btn = btnFromId(evt.Value)
		if evt.Type == event.NoteOn {
			e.Type = EventTypePressed
		} else if evt.Type == event.NoteOff {
			e.Type = EventTypeReleased
		}
	case event.ControlChange:
		e.Btn = btnFromId(evt.Controller)
		if evt.Value == 127 {
			e.Type = EventTypePressed
		} else {
			e.Type = EventTypeReleased
		}
	}

	return e
}

func (m *Controller) Subscribe(channel chan<- Event) {
	m.subscribers = append(m.subscribers, channel)
}

func (m *Controller) onDeviceEvent(in event.Event) {
	evt := toEvent(in)
	for _, c := range m.subscribers {
		c <- evt
	}
}
