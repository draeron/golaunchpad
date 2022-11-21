package layout

import (
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(layout Layout, btn button.Button)
type HoldHandler func(layout Layout, btn button.Button, first bool)

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

func MaskedHandler(mask launchpad.Mask, handler Handler) Handler {
	return func(layout Layout, btn button.Button) {
		if mask[btn] {
			handler(layout, btn)
		}
	}
}

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
