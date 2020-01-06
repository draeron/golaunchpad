package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/golaunchpad/pkg/minimk3"
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
	SelectLayout(layout minimk3.Layout) error
	Diag() error
	IsAwake()
	Wake() error
	Sleep() error
	String() string
	Print()
	Subscribe(channel chan<- event.Event)
	SetColorAll(col color.Color) error
	SetColorMany(btns []button.Button, color color.Color) error
	SetColor(btn button.Button, color color.Color) error
	SetColors(sets []minimk3.BtnColor) error
}
