package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/golaunchpad/pkg/launchpad/mask"
	"github.com/draeron/gopkgs/color"
)

type Controller interface {
	EnableDebugLogger()
	Close()
	SetBrightness(level byte) error
	EnableDawMode() error
	DisableDawMode() error
	EnableProgrammerMode() error
	DisableProgrammerMode() error
	Diag() error
	IsAwake()
	Wake() error
	Sleep() error
	String() string
	Print()
	Subscribe(channel chan<- event.Event)
	DisplayText(text string, loop bool, speed byte, color color.Color) error
	Name() string
	SetName(name string)
	Colorer
}

type Colorer interface {
	SetColorAll(col color.Color) error
	SetColorMany(btns []button.Button, color color.Color) error
	SetColor(btn button.Button, color color.Color) error
	SetColorMask(mask mask.Buttons, col color.Color) error
	SetColorMaskPreset(mask mask.Preset, col color.Color) error
	SetColors(sets button.ColorMap) error
	SetColorPad(x, y int, color color.Color) error
}
