package launchpad

import "github.com/draeron/golaunchpad/pkg/launchpad/button"

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(layout *Layout, btn button.Button)

/*
	HandlerType x ENUM(
	RowPressed
	RowReleased
	PadPressed
	PadReleased
	ModePressed
	ModeReleased
	ArrowPressed
	ArrowReleased
)
*/
type HandlerType int
