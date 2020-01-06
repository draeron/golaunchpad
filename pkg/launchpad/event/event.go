package event

import (
	"fmt"

	"github.com/draeron/golaunchpad/pkg/device/event"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

type Event struct {
	Type Type
	Btn  button.Button
}

func (e Event) String() string {
	return fmt.Sprintf("Event: %s - %s", e.Type, e.Btn)
}

func FromMidiEvent(evt event.Event) Event {
	e := Event{}

	switch evt.Type {
	case event.NoteOn, event.NoteOff:
		e.Btn = button.FromMidiId(evt.Value)
		if evt.Type == event.NoteOn {
			e.Type = Pressed
		} else if evt.Type == event.NoteOff {
			e.Type = Released
		}
	case event.ControlChange:
		e.Btn = button.FromMidiId(evt.Controller)
		if evt.Value == 127 {
			e.Type = Pressed
		} else {
			e.Type = Released
		}
	}

	return e
}
