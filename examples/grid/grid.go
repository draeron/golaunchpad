package main

import (
	"time"

	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/grid"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/layout"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

var log = logger.NewLogrus("main")
var device launchpad.Controller

func main() {
	log.Info("starting grid example")
	defer log.Info("exiting grid example")
	device = common.Setup()
	defer device.Close()

	setup()

	common.WaitExit()
}

func setup() {
	mask := launchpad.Mask{
		button.User: true,
	}
	gryd := grid.NewGrid(16, 16, true, mask)
	gryd.Layout.SetName("grid")

	gryd.Layout.SetHoldTimer(layout.ArrowHold, time.Millisecond*20)
	gryd.SetHandler(func(grd *grid.Grid, x, y int, event grid.EventType) {
		if event != grid.Pressed {
			return
		}
		col := color.FromStdColor(gryd.Color(x, y))
		if col.Equal(color.Black) {
			gryd.SetColor(x, y, common.RandColor())
		} else {
			gryd.SetColor(x, y, color.Black)
		}
	})

	// Add a toggle on the user btn
	wrapped := false
	gryd.Layout.SetColor(button.User, color.Yellow)
	gryd.Layout.SetHandler(layout.ModePressed, func(layout layout.Layout, btn button.Button) {
		if btn == button.User {
			wrapped = !wrapped
			gryd.Wrap(wrapped)
			if wrapped {
				gryd.Layout.SetColor(btn, color.Green)
			} else {
				gryd.Layout.SetColor(btn, color.Yellow)
			}
			gryd.UpdateDevice()
		}
	})

	gryd.Connect(device)
	gryd.Activate()
}
