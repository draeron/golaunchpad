package main

import (
	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/grid"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var device launchpad.Controller

func main() {
	log.Info("starting layout example")
	defer log.Info("exiting layout example")
	device = common.Setup()
	defer device.Close()

	setup()

	common.WaitExit()
}

func setup() {
	mask := launchpad.Mask{
		button.User: true,
	}
	gryd := grid.NewGrid(12,12, true, mask)
	gryd.Layout.DebugName = "grid"

	gryd.SetHandler(func(grd *grid.Grid, x, y int, event grid.EventType) {
		if event != grid.Pressed {
			return
		}
		col := color.FromColor(gryd.Color(x,y))
		if col.Equal(color.Black) {
			gryd.SetColor(x,y, common.RandColor())
		} else {
			gryd.SetColor(x,y, color.Black)
		}
	})

	// Add a toggle on the user btn
	wrapped := false
	gryd.Layout.SetColor(button.User, color.Yellow)
	gryd.Layout.SetHandler(launchpad.ModePressed, func(layout *launchpad.BasicLayout, btn button.Button) {
		if btn == button.User {
			wrapped = !wrapped
			gryd.Wrap(wrapped)
			if wrapped {
				gryd.Layout.SetColor(btn, color.Green)
			} else {
				gryd.Layout.SetColor(btn, color.Yellow)
			}
		}
	})

	gryd.Connect(device)
	gryd.Activate()
}
