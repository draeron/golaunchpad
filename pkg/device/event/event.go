package event

import (
	"errors"
	"fmt"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/sysex"
)

type Event struct {
	Type           Type
	Channel        channel.Channel
	Controller     uint8
	ControllerName string
	Value          uint8
	valid          bool
}

func (e Event) IsValid() bool {
	return e.valid
}

func From(msg midi.Message) (Event, error) {

	var err error
	evt := Event{
		valid: true,
	}

	switch v := msg.(type) {
	case sysex.Message:
		evt.Type = SysEx
		log.Infof(v.String())

	case channel.ControlChange:
		evt.Type = ControlChange
		evt.Controller = v.Controller()
		evt.Channel = channel.Channel(v.Channel())
		evt.Value = v.Value()

	case channel.NoteOn:
		evt.Type = NoteOn
		evt.Channel = channel.Channel(v.Channel())
		evt.Value = v.Key()

	case channel.NoteOff:
		evt.Type = NoteOff
		evt.Channel = channel.Channel(v.Channel())
		evt.Value = v.Key()

	case channel.NoteOffVelocity:
		evt.Type = NoteOffVelocity
		evt.Channel = channel.Channel(v.Channel())

	case channel.Aftertouch:
		evt.Type = Aftertouch
		evt.Channel = channel.Channel(v.Channel())
		evt.Value = v.Pressure()

	case channel.PolyAftertouch:
		evt.Type = PolyAftertouch
		evt.Channel = channel.Channel(v.Channel())
		evt.Value = v.Pressure()

	case channel.Pitchbend:
		evt.Type = PitchBend
		evt.Channel = channel.Channel(v.Channel())

	default:
		log.Warnf("unsupported message: %s\n", v.String())
		err = errors.New("unsupported message")
		evt.valid = false
	}

	return evt, err
}

func (e Event) String() string {
	return fmt.Sprintf("%s - Channel: %d, Controller: %s (%d), Value: %d",
		e.Type,
		e.Channel,
		e.ControllerName,
		e.Controller,
		e.Value,
	)
}
