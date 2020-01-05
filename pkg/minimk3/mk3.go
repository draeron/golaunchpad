package minimk3

import (
	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/device/event"
)

type Controller struct {
	device      *device.MiniMk3
	mode        Mode
	subscribers []chan<- Event
	eventsChan  chan event.Event
}

const (
	Default = Mode(iota)
	DawMode
	ProgrammerMode
)

type Mode int

func Open(mode Mode) (*Controller, error) {
	dev, err := device.Detect()
	if err != nil {
		return nil, err
	}

	err = dev.Open()
	if err != nil {
		return nil, err
	}

	pad := Controller{
		device: dev,
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
	pad.eventsChan = make(chan event.Event)
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
		ch := make(chan event.Event)
		m.device.Subscribe(ch)
		for evt := range ch {
			//log.Debugf("%s", evt.String())
			log.Debugf("%s", toEvent(evt).String())
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

func (m *Controller) SetBrightness(level byte) error {
	msg := CmdBrightnessLevel.SysEx(level)
	return m.device.SendMidi(msg)
}

func (m *Controller) EnableDawMode() error {
	log.Infof("enabling daw mode")
	msg := CmdEnableSession.SysEx(1)
	return m.device.SendDaw(msg)
}

func (m *Controller) DisableDawMode() error {
	log.Infof("disabling daw mode")
	msg := CmdEnableSession.SysEx(0)
	return m.device.SendDaw(msg)
}

func (m *Controller) EnableProgrammerMode() error {
	log.Infof("enabling programmer mode")
	msg := CmdProgrammerMode.SysEx(1)
	return m.device.SendMidi(msg)
}

func (m *Controller) DisableProgrammerMode() error {
	log.Infof("disabling programmer mode")
	msg := CmdProgrammerMode.SysEx(0)
	return m.device.SendMidi(msg)
}


func (m *Controller) SelectLayout(layout Layout) error {
	msg := CmdSelectLayout.SysEx(layout.value())
	return m.device.SendDaw(msg)
}

func (m *Controller) Diag() error {
	var err error
	msg := CmdDeviceInquiry.SysEx()
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
	msg := CmdSleep.SysEx()
	m.device.SendMidi(msg)
}

func (m *Controller) Wake() error {
	log.Infof("sending wake message")
	msg := CmdSleep.SysEx(1)
	return m.device.SendMidi(msg)
}

func (m *Controller) Sleep() error {
	log.Infof("sending sleep message")
	msg := CmdSleep.SysEx(0)
	return m.device.SendMidi(msg)
}

func (m *Controller) String() string {
	return m.device.String()
}

func (m *Controller) Print() {
	m.device.Print()
}
