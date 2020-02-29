package grid

import (
	seven_bits "github.com/draeron/golaunchpad/pkg/colors/7bits"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"image/color"
	"time"
)

/*
	A colors is a layout which is larger than the 8x8. Colors
	are allocated through coordinate system and arrows are
	used to scroll the window

	options exist to setup arrows button automatically and
	the speed which those arrow scroll the colors
*/
type Grid struct {
	Layout       *launchpad.BasicLayout
	posX, posY   int
	sizeX, sizeY int
	colors       [][]color.Color
	wrap         bool
	handler      Handler
	arrows       bool
}

const (
	Pressed = EventType(iota)
	Released
)

type EventType int
type Handler func(grid *Grid, x, y int, event EventType)

func NewGrid(x, y int, arrows bool, mask launchpad.Mask) *Grid {
	mask = launchpad.MaskPad.Mask().Merge(mask)
	if arrows {
		mask = mask.MergePreset(launchpad.MaskArrows)
	}

	grid := Grid{
		colors: make([][]color.Color, x),
		Layout: launchpad.NewLayout(mask),
		sizeX:  x,
		sizeY:  y,
		arrows: arrows,
	}

	for i, _ := range grid.colors {
		grid.colors[i] = make([]color.Color, y)
	}
	grid.SetColorAll(color.Black)

	if arrows {
		grid.Layout.SetHandler(launchpad.ArrowPressed, grid.onArrow)
		grid.Layout.SetHandler(launchpad.ArrowHold, grid.onArrow)
		grid.Layout.SetHoldTimer(launchpad.ArrowHold, time.Millisecond*75)
	}

	grid.Layout.SetHandler(launchpad.PadPressed, func(layout *launchpad.BasicLayout, btn button.Button) {
		grid.onPad(btn, Pressed)
	})
	grid.Layout.SetHandler(launchpad.PadReleased, func(layout *launchpad.BasicLayout, btn button.Button) {
		grid.onPad(btn, Released)
	})

	grid.Layout.SetColorMany(button.Arrows(), color.White)
	grid.updateColors()

	return &grid
}

func (g *Grid) SetHandler(handler Handler) {
	g.handler = handler
}

func (g *Grid) onPad(btn button.Button, eventType EventType) {
	if g.handler != nil {
		x, y := btn.Coord()
		y-- // pads are y+1

		x, y = g.wrapPos(x+g.posX, y+g.posY)
		g.handler(g, x, y, eventType)
	}
}

func (g *Grid) PanUp() bool {
	oldY := g.posY
	if g.CanPanUp() {
		g.posY = (g.posY - 1) % g.sizeY
	}
	return oldY != g.posY
}

func (g *Grid) PanDown() bool {
	oldY := g.posY
	if g.CanPanDown() {
		g.posY = (g.posY + 1) % g.sizeY
	}
	return oldY != g.posY
}

func (g *Grid) PanLeft() bool {
	oldX := g.posX
	if g.CanPanLeft() {
		g.posX = (g.posX - 1) % g.sizeX
	}
	return oldX != g.posX
}

func (g *Grid) PanRight() bool {
	oldX := g.posX
	if g.CanPanRight() {
		g.posX = (g.posX + 1) % g.sizeX
	}
	return oldX != g.posX
}

func (g *Grid) CanPanUp() bool {
	return g.wrap || g.posY > 0
}

func (g *Grid) CanPanDown() bool {
	return g.wrap || g.posY < g.sizeY-8
}

func (g *Grid) CanPanLeft() bool {
	return g.wrap || g.posX > 0
}

func (g *Grid) CanPanRight() bool {
	return g.wrap || g.posX < g.sizeX-8
}

func (g *Grid) onArrow(layout *launchpad.BasicLayout, btn button.Button) {
	switch btn {
	case button.Up:
		g.PanUp()
	case button.Down:
		g.PanDown()
	case button.Left:
		g.PanLeft()
	case button.Right:
		g.PanRight()
	}

	g.updateColors()
}

func (g *Grid) Color(x, y int) color.Color {
	if x < 0 || y < 0 || x >= g.sizeX || y >= g.sizeY {
		return color.Black
	}
	return g.colors[x][y]
}

func (g *Grid) SetColorAll(col color.Color) {
	for x := 0; x < g.sizeX; x++ {
		for y := 0; y < g.sizeY; y++ {
			g.colors[x][y] = col
		}
	}
}

func (g *Grid) SetColor(x, y int, col color.Color) {
	if x < 0 || y < 0 || x >= g.sizeX || y >= g.sizeY {
		return
	}
	g.colors[x][y] = col
	g.updateColors()
}

func (g *Grid) updateColors() {
	mapp := button.ColorMap{}

	if g.arrows {
		for _, btn := range button.Arrows() {
			var canMove func() bool
			switch btn {
			case button.Up:
				canMove = g.CanPanUp
			case button.Down:
				canMove = g.CanPanDown
			case button.Left:
				canMove = g.CanPanLeft
			case button.Right:
				canMove = g.CanPanRight
			}
			if canMove != nil {
				if canMove() {
					mapp[btn] = seven_bits.FromColor(color.White)
				} else {
					mapp[btn] = seven_bits.FromColor(color.RGBA{R: 255})
				}
			}
		}
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			btn := button.FromPadXY(x, y)
			cx, cy := g.wrapPos(x+g.posX, y+g.posY)
			mapp[btn] = seven_bits.FromColor(g.Color(cx, cy))
		}
	}
	g.Layout.SetColors(mapp)
}

func (g *Grid) Connect(controller launchpad.Controller) {
	g.Layout.Connect(controller)
}

func (g *Grid) Disconnect() {
	g.Layout.Disconnect()
}

func (g *Grid) Activate() {
	g.Layout.Activate()
}

func (g *Grid) Deactivate() {
	g.Layout.Deactivate()
}

func (g *Grid) Wrap(b bool) {
	g.wrap = b

	if !g.wrap { // convert back to a valid coordinate
		g.posX = clamp(g.posX, 0, g.sizeX-8)
		g.posY = clamp(g.posY, 0, g.sizeY-8)
	}

	g.updateColors()
}

func (g *Grid) wrapPos(x, y int) (int, int) {
	if g.wrap {
		if x < 0 || x >= g.sizeX {
			x = (x + g.sizeX) % g.sizeX
		}
		if y < 0 || y >= g.sizeY {
			y = (y + g.sizeY) % g.sizeY
		}
	}
	return x, y
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}
