// Code generated by go-enum
// DO NOT EDIT!

package event

import (
	"fmt"
)

const (
	// Aftertouch is a Type of type Aftertouch
	Aftertouch Type = iota
	// ControlChange is a Type of type ControlChange
	ControlChange
	// NoteOn is a Type of type NoteOn
	NoteOn
	// NoteOff is a Type of type NoteOff
	NoteOff
	// NoteOffVelocity is a Type of type NoteOffVelocity
	NoteOffVelocity
	// PitchBend is a Type of type PitchBend
	PitchBend
	// PolyAftertouch is a Type of type PolyAftertouch
	PolyAftertouch
	// ProgramChange is a Type of type ProgramChange
	ProgramChange
	// SysEx is a Type of type SysEx
	SysEx
)

const _TypeName = "AftertouchControlChangeNoteOnNoteOffNoteOffVelocityPitchBendPolyAftertouchProgramChangeSysEx"

var _TypeMap = map[Type]string{
	0: _TypeName[0:10],
	1: _TypeName[10:23],
	2: _TypeName[23:29],
	3: _TypeName[29:36],
	4: _TypeName[36:51],
	5: _TypeName[51:60],
	6: _TypeName[60:74],
	7: _TypeName[74:87],
	8: _TypeName[87:92],
}

func (i Type) String() string {
	if str, ok := _TypeMap[i]; ok {
		return str
	}
	return fmt.Sprintf("Type(%d)", i)
}

var _TypeValue = map[string]Type{
	_TypeName[0:10]:  0,
	_TypeName[10:23]: 1,
	_TypeName[23:29]: 2,
	_TypeName[29:36]: 3,
	_TypeName[36:51]: 4,
	_TypeName[51:60]: 5,
	_TypeName[60:74]: 6,
	_TypeName[74:87]: 7,
	_TypeName[87:92]: 8,
}

// ParseType attempts to convert a string to a Type
func ParseType(name string) (Type, error) {
	if x, ok := _TypeValue[name]; ok {
		return Type(x), nil
	}
	return Type(0), fmt.Errorf("%s is not a valid Type", name)
}
