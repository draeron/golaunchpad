package main

import (
	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var device launchpad.Controller

var top *launchpad.Layout
var pads [4]*launchpad.Layout
var currentMode = 0
var modeColors = []color.Color{ color.Red, color.Blue, color.Green, color.Yellow }

func main() {
	log.Info("starting layout example")
	defer log.Info("exiting layout example")
	device = common.Setup()
	defer device.Close()

	setup()

	common.WaitExit()
}

func setup() {
	top = launchpad.NewLayoutPreset(launchpad.MaskTop)
	top.DebugName = "top"
	top.SetHandler(launchpad.ModePressed, func(layout *launchpad.Layout, btn button.Button) {
		pads[currentMode].Disable()
		currentMode = int(btn - button.Session)
		pads[currentMode].Enable()
		top.SetColor(button.Logo, modeColors[currentMode])
		log.Infof("switched to mode %s", btn)
	})

	for i, _ := range pads {
		id := button.Session + button.Button(i)
		top.SetColor(id, modeColors[i])

		pads[i] = launchpad.NewLayoutPreset(launchpad.MaskPad)
		pads[i].DebugName = id.String()
		pads[i].Connect(device)
		pads[i].SetHandler(launchpad.PadPressed, func(layout *launchpad.Layout, btn button.Button) {
			if layout.Color(btn) != color.Black {
				layout.SetColor(btn, color.Black)
			} else {
				layout.SetColor(btn, common.RandColor())
			}
		})
		pads[i].SetColorAll(color.Black)
	}

	top.Connect(device)
	top.Enable()
	pads[currentMode].Enable()
	log.Infof("switched to mode %v", button.Session)
}