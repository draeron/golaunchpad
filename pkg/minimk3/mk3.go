package minimk3

import (
	"image/color"
	"sync"

	"github.com/pkg/errors"

	seven_bits "github.com/draeron/golaunchpad/pkg/colors/7bits"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/golaunchpad/pkg/minimk3/cmd"
	"github.com/draeron/gopkgs/midi"
	midiEvt "github.com/draeron/gopkgs/midi/event"
)

type Controller struct {
	device      *midiDevice
	mode        Mode
	subscribers []chan<- event.Event
	eventsChan  chan midiEvt.Event
	mutex       sync.RWMutex
}

const (
	Default = Mode(iota)
	DawMode
	ProgrammerMode
)

type Mode int

func Open(mode Mode) (*Controller, error) {
	var err error

	mk3 := midiDevice{}

	mk3.daw, err = midi.DetectPort(rxDAW, rxDAW)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open launchpad DAW ports")
	}
	mk3.midi, err = midi.DetectPort(rxMidiIn, rxMidiOut)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open launchpad MIDI ports")
	}

	err = mk3.Open()
	if err != nil {
		return nil, err
	}

	pad := Controller{
		device: &mk3,
		mode:   mode,
	}

	switch mode {
	case ProgrammerMode:
		err = pad.EnableProgrammerMode()
	case DawMode:
		err = pad.EnableDawMode()
	default:
		err = pad.EnableProgrammerMode()
		err = pad.DisableDawMode()
	}

	pad.IsAwake()
	pad.Wake()
	pad.eventsChan = make(chan midiEvt.Event, 20)
	pad.device.Subscribe(pad.eventsChan)
	pad.Diag()

	go func() {
		for e := range pad.eventsChan {
			pad.onDeviceEvent(e)
		}
	}()

	return &pad, err
}

func (m *Controller) EnableDebugLogger() {
	go func() {
		log.Debugf("enable debug logging of events")
		ch := make(chan midiEvt.Event, 20)
		m.device.Subscribe(ch)
		for evt := range ch {
			// log.Debugf("%s", evt.String())
			log.Debugf("%s", event.FromMidiEvent(evt).String())
		}
	}()
}

func (m *Controller) Close() {
	switch m.mode {
	case ProgrammerMode:
		m.DisableProgrammerMode()
	case DawMode:
		m.DisableDawMode()
	}
	m.device.Close()
}

func (m *Controller) Name() string {
	// todo: add index to name?
	return "MidiDevice"
}

func (m *Controller) SetBrightness(level byte) error {
	msg := cmd.BrightnessLevel.SysEx(level)
	return m.device.SendMidi(msg)
}

func (m *Controller) EnableDawMode() error {
	log.Infof("enabling daw mode")
	msg := cmd.EnableSession.SysEx(1)
	return m.device.SendDaw(msg)
}

func (m *Controller) DisableDawMode() error {
	log.Infof("disabling daw mode")
	msg := cmd.EnableSession.SysEx(0)
	return m.device.SendDaw(msg)
}

func (m *Controller) EnableProgrammerMode() error {
	log.Infof("enabling programmer mode")
	msg := cmd.ProgrammerMode.SysEx(1)
	return m.device.SendMidi(msg)
}

func (m *Controller) DisableProgrammerMode() error {
	log.Infof("disabling programmer mode")
	msg := cmd.ProgrammerMode.SysEx(0)
	return m.device.SendMidi(msg)
}

func (m *Controller) DisplayText(text string, loop bool, speed byte, color color.Color) error {
	log.Debugf("sending text '%s', loop: %v, speed: %d pad/s, color: %v", text, loop, speed, color)
	loopbyte := byte(0)
	if loop {
		loopbyte = 1
	}

	rgb := seven_bits.FromColor(color)
	buf := append([]byte{}, loopbyte, speed, 0x1, rgb.R, rgb.G, rgb.B)
	buf = append(buf, []byte(text)...)
	msg := cmd.TextScrolling.SysEx(buf...)
	return m.device.SendMidi(msg)
}

func (m *Controller) SelectLayout(layout Layout) error {
	msg := cmd.SelectLayout.SysEx(layout.value())
	return m.device.SendDaw(msg)
}

func (m *Controller) Diag() error {
	var err error
	msg := cmd.DeviceInquiry.SysEx()
	err = m.device.SendDaw(msg)
	if err != nil {
		return err
	}
	err = m.device.SendMidi(msg)
	if err != nil {
		return err
	}
	return err
}

func (m *Controller) IsAwake() {
	msg := cmd.Sleep.SysEx()
	m.device.SendMidi(msg)
}

func (m *Controller) Wake() error {
	log.Infof("sending wake message")
	msg := cmd.Sleep.SysEx(1)
	return m.device.SendMidi(msg)
}

func (m *Controller) Sleep() error {
	log.Infof("sending sleep message")
	msg := cmd.Sleep.SysEx(0)
	return m.device.SendMidi(msg)
}

func (m *Controller) String() string {
	return m.device.String()
}

func (m *Controller) Print() {
	m.device.Print()
}
