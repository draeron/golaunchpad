package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"image/color"
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

	Colorer
}

type Colorer interface {
	SetColorAll(col color.Color) error
	SetColorMany(btns []button.Button, color color.Color) error
	SetColor(btn button.Button, color color.Color) error
	SetColors(sets button.ColorMap) error
}
