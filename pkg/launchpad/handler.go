package launchpad

import "github.com/draeron/golaunchpad/pkg/launchpad/button"

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(layout *BasicLayout, btn button.Button)
type HoldHandler func(layout *BasicLayout, btn button.Button, first bool)

/*
	HandlerType x ENUM(
	RowPressed
  RowHold
	RowReleased
	PadPressed
  PadHold
	PadReleased
	ModePressed
	ModeHold
	ModeReleased
	ArrowPressed
	ArrowHold
	ArrowReleased
)
*/
type HandlerType int

func (h HandlerType) IsPressed() bool {
	switch h {
	case RowPressed, PadPressed, ModePressed, ArrowPressed:
		return true
	default:
		return false
	}
}

func (h HandlerType) IsReleased() bool {
	switch h {
	case RowReleased, PadReleased, ModeReleased, ArrowReleased:
		return true
	default:
		return false
	}
}

func (h HandlerType) IsHold() bool {
	switch h {
	case RowHold, PadHold, ModeHold, ArrowHold:
		return true
	default:
		return false
	}
}
