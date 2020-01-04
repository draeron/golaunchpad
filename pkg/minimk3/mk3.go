package minimk3

import (
  "github.com/draeron/golaunchpad/pkg/device"
  "github.com/draeron/golaunchpad/pkg/device/event"
  "gitlab.com/gomidi/midi/midimessage/sysex"
  "image/color"
)

type MiniMk3 struct {
  device *device.MiniMk3
  mode Mode
}

const (
  Default = Mode(iota)
  DawMode
  ProgrammerMode
)
type Mode int

var (
  DeviceInquiry   = []byte{ 0x7E, 0x7F, 0x06, 0x01 }
  CommandHeader   = []byte{ 0x00, 0x20, 0x29, 0x02, 0x0D }
  EnableDaw       = append(CommandHeader, 0x10)
  SelectLiveMode  = append(CommandHeader, 0x0e)
  SelectLayout    = append(CommandHeader, 0x00)
  BrightnessLevel = append(CommandHeader, 0x08)
  ClearDawState   = append(CommandHeader, 0x12)
  SessionColor    = append(CommandHeader, 0x14)
  TextScrolling   = append(CommandHeader, 0x07)
  LedColor        = append(CommandHeader, 0x03)
)

func Open(mode Mode) (*MiniMk3, error) {
  dev, err := device.Detect()
  if err != nil {
    return nil, err
  }

  err = dev.Open()
  if err != nil {
    return nil, err
  }

  pad := MiniMk3{
    device: dev,
    mode: mode,
  }

  switch mode {
  case ProgrammerMode:
    pad.EnableProgrammerMode(true)
  case DawMode:
    pad.EnableDawMode(true)
  default:
    pad.EnableProgrammerMode(false)
    pad.EnableDawMode(false)
  }

  return &pad, nil
}

func (m *MiniMk3) EnableDebugLogger() {
  m.device.AddListener(func(msg event.Event) {
    log.Debugf("%s", msg.String())
  })
}

func (m *MiniMk3) Close() {
  switch m.mode {
  case ProgrammerMode:
    m.EnableProgrammerMode(false)
  case DawMode:
    m.EnableDawMode(false)
  }
  m.device.Close()
}

func (m *MiniMk3) SetBrightness(level byte) error {
  msg := newSysExMsg(BrightnessLevel, level)
  return m.device.SendMidi(msg)
}

func (m *MiniMk3) SetBtnColor(btn Btn, color color.Color) error {
  r, g, b, _ := color.RGBA()
  id := btn.Id()
  msg := newSysExMsg(LedColor, 3, id, byte(r>>9), byte(g>>9), byte(b>>9))
  return m.device.SendMidi(msg)
}

func (m *MiniMk3) ClearColors(col color.Color) error {
  btns := BtnValues()
  cols := []color.Color{}
  for range btns {
    cols = append(cols, col)
  }
  return m.SetBtnColors(btns, cols)
}

func (m *MiniMk3) SetBtnColors(btns []Btn, colors []color.Color) error {
  buf := append(LedColor)

  for idx, it := range btns {
    if idx > len(colors) {
      break
    }
    r, g, b, _ := colors[idx].RGBA()
        buf = append(buf, 3, it.Id(), byte(r>>9), byte(g>>9), byte(b>>9))
  }
  msg := newSysExMsg(buf)
  return m.device.SendMidi(msg)
}

func (m *MiniMk3) EnableDawMode(enable bool) error {
  msg := sysex.SysEx(EnableDaw)
  if enable {
    log.Infof("enabling daw mode")
    msg = append(msg, 0x01)
  } else {
    log.Infof("disabling daw mode")
    msg = append(msg, 0x00)
  }
  return m.device.SendDaw(msg)
}

func (m *MiniMk3) EnableProgrammerMode(enable bool) error {
  msg := sysex.SysEx(SelectLiveMode)
  if enable {
    log.Infof("enabling live mode")
    msg = append(msg, 0x01)
  } else {
    log.Infof("disabling live mode")
    msg = append(msg, 0x00)
  }
  return m.device.SendMidi(msg)
}

func (m *MiniMk3) SelectLayout(layout Layout) error {
  msg := sysex.SysEx(append(SelectLayout, layout.value()))
  return m.device.SendDaw(msg)
}

func (m *MiniMk3) Diag() error {
  var err error
  msg := sysex.SysEx(DeviceInquiry)
  //err := m.device.SendDaw(msg)
  //if err != nil {
  //  return err
  //}
  err = m.device.SendMidi(msg)
  if err != nil {
    return err
  }
  return err
}

func (m *MiniMk3) String() string {
  return m.device.String()
}

func (m *MiniMk3) Print() {
  m.device.Print()
}

func newSysExMsg(cmd []byte, data... byte) sysex.SysEx {
  return sysex.SysEx(append(cmd, data...))
}