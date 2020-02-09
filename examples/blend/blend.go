package main

import (
	dcolor "github.com/draeron/golaunchpad/pkg/colors/dynamic"
	"time"

	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var ctrl launchpad.Controller

func main() {
	log.Info("starting blend example")
	defer log.Info("exiting blend example")
	ctrl = common.Setup()
	defer ctrl.Close()

	setup()

	common.WaitExit()
}

func setup() {
	layout := launchpad.NewLayoutPreset(launchpad.MaskAll)
	layout.DebugName = "blend"
	layout.Connect(ctrl)
	layout.Activate()

	common.Must(layout.SetColorAll(color.Black))
	common.Must(layout.SetColor(button.Logo, color.LightGray))
	common.Must(layout.SetColorMany(button.Modes(), color.LightGray))

	for i, b := range button.Rows() {
		common.Must(layout.SetColor(b, color.PaletteColor(i)))
	}

	const blendtime = time.Millisecond * 400

	layout.SetHandler(launchpad.RowPressed, func(layout *launchpad.BasicLayout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := layout.Color(btn)
			if color.Black.Equal(color.FromColor(old)) {
				blend := dcolor.Blend(old, col, blendtime)
				common.Must(layout.SetColor(btn, blend))
			}
		}
	})

	layout.SetHandler(launchpad.RowHold, func(layout *launchpad.BasicLayout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := layout.Color(btn)
			blend := dcolor.Blend(old, col, time.Millisecond * 100)
			common.Must(layout.SetColor(btn, blend))
		}
	})

	// stop blending by assigning a color to the button
	layout.SetHandler(launchpad.PadPressed, func(layout *launchpad.BasicLayout, btn button.Button) {
		layout.SetColor(btn, color.Black)
	})

}

func setPadColor() {

}
