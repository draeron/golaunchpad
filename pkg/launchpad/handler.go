package launchpad

import "github.com/draeron/golaunchpad/pkg/launchpad/button"

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(layout *BasicLayout, btn button.Button)

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
