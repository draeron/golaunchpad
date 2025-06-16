package cmd

import (
	"bytes"
	"log"

	"gitlab.com/gomidi/midi/v2"
)

// go tool github.com/abice/go-enum -f=$GOFILE --noprefix

// Cmd x ENUM(
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
	EnableSession:   0x10,
	SelectLayout:    0x00,
	ProgrammerMode:  0x0e,
	BrightnessLevel: 0x08,
	ClearDawState:   0x12,
	SessionColor:    0x14,
	TextScrolling:   0x07,
	LedColor:        0x03,
	Sleep:           0x09,
}

func (c Cmd) Bytes() []byte {
	switch c {
	case DeviceInquiry:
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

func (c Cmd) SysEx(byts ...byte) midi.Message {
	buf := bytes.Buffer{}
	buf.Write(c.Bytes())
	buf.Write(byts)
	return midi.SysEx(buf.Bytes())
}
