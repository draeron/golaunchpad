package device

import (
	"bytes"
	"fmt"
	"github.com/draeron/golaunchpad/pkg/device/event"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
)

type MiniMk3 struct {
	dawIn   mid.In
	midiIn  mid.In
	dawOut  mid.Out
	midiOut mid.Out

	dawRdr  *mid.Reader
	midiRdr *mid.Reader

	dawWr   *mid.Writer
	midiWr  *mid.Writer

	channels []chan<- event.Event
}

func (m *MiniMk3) Open() error {
	var err error
	err = m.dawIn.Open()
	if err != nil {
		return err
	}

	err = m.dawOut.Open()
	if err != nil {
		return err
	}

	err = m.midiIn.Open()
	if err != nil {
		return err
	}

	err = m.midiOut.Open()
	if err != nil {
		return err
	}

	m.dawRdr = mid.NewReader(mid.NoLogger())
	m.dawRdr.Msg.Each = m.onMessage
	err = mid.ConnectIn(m.dawIn, m.dawRdr)
	if err != nil {
		return err
	}

	m.midiRdr = mid.NewReader(mid.NoLogger())
	m.midiRdr.Msg.Each = m.onMessage
	err = mid.ConnectIn(m.midiIn, m.midiRdr)
	if err != nil {
		return err
	}

	m.dawWr = mid.ConnectOut(m.dawOut)
	m.midiWr = mid.ConnectOut(m.midiOut)

	return nil
}

func (m *MiniMk3) Subscribe(channel chan<- event.Event) {
	m.channels = append(m.channels, channel)
}

func (m *MiniMk3) onMessage(position *mid.Position, message midi.Message) {
	//log.Debugf("msg: %s", message.String())

	evt, err := event.From(message)
	if err != nil {
		log.Errorf("%s", err.Error())
	} else {
		for _, ch := range m.channels {
			ch <- evt
		}
	}
}

func (m *MiniMk3) receive(data []byte, deltaMicroseconds int64) {
	rd := midireader.New(bytes.NewReader(data), func(message realtime.Message) {
		log.Debugf("realtime msg: %v", message.String())
	})

	mid.NewReader()

	var msg midi.Message
	var err error

	for {
		msg, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		log.Debugf("msg: %s", msg.String())

		msg, err := event.From(msg)
		if err != nil {
			log.Errorf("%s", err.Error())
		} else if msg.IsValid() {
			for _, ch := range m.channels {
				ch <- msg
			}
		}
	}
}

func (m *MiniMk3) SendMidi(msg midi.Message) error {
	return send(m.midiWr, msg)
}

func (m *MiniMk3) SendDaw(msg midi.Message) error {
	return send(m.dawWr, msg)
}

func send(out *mid.Writer, msg midi.Message) error {
	return out.Write(msg)
}

func (m *MiniMk3) Close() {
	if m.dawIn.IsOpen() {
		m.dawIn.Close()
	}

	if m.dawOut.IsOpen() {
		m.dawOut.Close()
	}

	if m.midiIn.IsOpen() {
		m.midiIn.StopListening()
		m.midiIn.Close()
	}

	if m.midiOut.IsOpen() {
		m.midiOut.Close()
	}
}

func (m *MiniMk3) String() string {
	return fmt.Sprintf("Daw In: %s\nDaw Out: %s\nMidi In: %s\nMidi Out: %s\n",
		m.dawIn, m.dawOut, m.midiIn, m.midiOut,
	)
}

func (m *MiniMk3) Print() {
	log.Infof("Daw In: %s", m.dawIn.String())
	log.Infof("Daw Out: %s", m.dawOut.String())
	log.Infof("Midi In: %s", m.midiIn)
	log.Infof("Midi Out: %s", m.midiOut)
}
