package event

import (
	"fmt"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkgs/midi/event"
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
		e.Btn = button.FromMidiId(evt.Controller)
		if evt.Type == event.NoteOn && evt.Value > 0 {
			e.Type = Pressed
		} else if evt.Type == event.NoteOff || evt.Value == 0 {
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
