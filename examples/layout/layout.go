package main

import (
	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

var log = logger.NewLogrus("main")
var device launchpad.Controller

var top *launchpad.BasicLayout
var pads [4]*launchpad.BasicLayout
var currentMode = 0
var modeColors = []color.Color{color.Red, color.Blue, color.Green, color.Yellow}

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
	top.SetHandler(launchpad.ModePressed, func(layout *launchpad.BasicLayout, btn button.Button) {
		pads[currentMode].Deactivate()
		currentMode = int(btn - button.Session)
		pads[currentMode].Activate()
		top.SetColor(button.Logo, modeColors[currentMode])
		log.Infof("switched to mode %s", btn)
	})

	for i, _ := range pads {
		id := button.Session + button.Button(i)
		top.SetColor(id, modeColors[i])

		pads[i] = launchpad.NewLayoutPreset(launchpad.MaskPad)
		pads[i].DebugName = id.String()
		pads[i].Connect(device)
		pads[i].SetHandler(launchpad.PadPressed, func(layout *launchpad.BasicLayout, btn button.Button) {
			col := color.FromColor(layout.Color(btn))
			if col.Equal(color.Black) {
				layout.SetColor(btn, common.RandColor())
			} else {
				layout.SetColor(btn, color.Black)
			}
		})
		pads[i].SetColorAll(color.Black)
	}

	top.Connect(device)
	top.Activate()
	pads[currentMode].Activate()
	log.Infof("switched to mode %v", button.Session)
}
