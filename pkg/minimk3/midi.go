package minimk3

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi/v2"

	midi2 "github.com/draeron/gopkgs/midi"
	"github.com/draeron/gopkgs/midi/event"
)

type midiDevice struct {
	daw  *midi2.Port
	midi *midi2.Port

	channels []chan<- event.Event
	mutex    sync.Mutex
}

func (m *midiDevice) Open() error {
	var err error

	err = m.daw.Open(m.onMessage)
	if err != nil {
		return err
	}

	err = m.midi.Open(m.onMessage)
	if err != nil {
		return err
	}

	return nil
}

func (m *midiDevice) Subscribe(channel chan<- event.Event) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.channels = append(m.channels, channel)
}

func (m *midiDevice) onMessage(message midi.Message, timestampms int32) {
	// log.Debugf("msg: %s", message.String())

	evt, err := event.From(message)
	if err != nil {
		log.Errorf("%+v", err.Error())
	} else {
		for _, ch := range m.channels {
			ch <- evt
		}
	}
}

func (m *midiDevice) SendMidi(msg midi.Message) error {
	return m.midi.Send(msg)
}

func (m *midiDevice) SendDaw(msg midi.Message) error {
	return m.daw.Send(msg)
}

func (m *midiDevice) Close() {
	m.daw.Close()
	m.midi.Close()
}

func (m *midiDevice) String() string {
	return fmt.Sprintf("Daw In: %s\nDaw Out: %s\nMidi In: %s\nMidi Out: %s\n",
		m.daw.In, m.daw.Out, m.midi.In, m.midi.Out,
	)
}

func (m *midiDevice) Print() {
	log.Infof("Daw In: %s", m.daw.In)
	log.Infof("Daw Out: %s", m.daw.Out)
	log.Infof("Midi In: %s", m.midi.In)
	log.Infof("Midi Out: %s", m.midi.Out)
}
