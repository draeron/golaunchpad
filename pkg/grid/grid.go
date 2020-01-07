package grid

import (
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"image/color"
)

type Grid struct {
	launchpad.Layout
	posX, posY int

	grid [][]color.Color
}
