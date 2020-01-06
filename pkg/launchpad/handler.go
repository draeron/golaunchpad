package launchpad

import "github.com/draeron/golaunchpad/pkg/launchpad/button"

//go:generate go-enum -f=$GOFILE --noprefix

type Handler func(button button.Button)

/*
	HandlerType x ENUM(
	RowPressed
	RowReleased
	PadPressed
	PadReleased
	PagePressed
	PageReleased
	ArrowPressed
	ArrowReleased
)
*/
type HandlerType int
