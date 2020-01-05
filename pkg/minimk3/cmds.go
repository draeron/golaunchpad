package minimk3

import (
	"bytes"
	"gitlab.com/gomidi/midi/midimessage/sysex"
)

//go:generate go-enum -f=$GOFILE

// Layout x ENUM(
/*
  DeviceInquiry
	EnableSession
	ProgrammerMode
	SelectLayout
	BrightnessLevel
	ClearDawState
	SessionColor
	TextScrolling
	LedColor
	Sleep
*/
// )
type Cmd int

var CommandHeader = []byte{0x00, 0x20, 0x29, 0x02, 0x0D}

var cmdByte = map[Cmd]byte{
	CmdEnableSession:   0x10,
	CmdSelectLayout:    0x00,
	CmdProgrammerMode:  0x0e,
	CmdBrightnessLevel: 0x08,
	CmdClearDawState:   0x12,
	CmdSessionColor:    0x14,
	CmdTextScrolling:   0x07,
	CmdLedColor:        0x03,
	CmdSleep:           0x09,
}

func (c Cmd) Bytes() []byte {
	switch c {
	case CmdDeviceInquiry:
		return []byte{0x7E, 0x7F, 0x06, 0x01}
	default:
		if byt, ok := cmdByte[c]; ok {
			return append(CommandHeader, byt)
		} else {
			log.Panicf("missing byte command value for command %s", c.String())
		}
	}
	return nil
}

func (c Cmd) SysEx(byts ...byte) sysex.SysEx {
	buf := bytes.Buffer{}
	buf.Write(c.Bytes())
	buf.Write(byts)
	return sysex.SysEx(buf.Bytes())
}
