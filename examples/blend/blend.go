package main

import (
	"time"

	"github.com/draeron/golaunchpad/pkg/layout"
	dcolor "github.com/draeron/gopkgs/color/dynamic"

	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

var log = logger.NewLogrus("main")
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
	layout := layout.NewLayoutPreset(launchpad.MaskAll)
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

	layout.SetHandler(layout.RowPressed, func(layout *layout.BasicLayout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := layout.Color(btn)
			if color.Black.Equal(color.FromStdColor(old)) {
				blend := dcolor.Blend(old, col, blendtime)
				common.Must(layout.SetColor(btn, blend))
			}
		}
	})

	layout.SetHandler(layout.RowHold, func(layout *layout.BasicLayout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := layout.Color(btn)
			blend := dcolor.Blend(old, col, time.Millisecond*100)
			common.Must(layout.SetColor(btn, blend))
		}
	})

	// stop blending by assigning a color to the button
	layout.SetHandler(layout.PadPressed, func(layout *layout.BasicLayout, btn button.Button) {
		layout.SetColor(btn, color.Black)
	})

}

func setPadColor() {

}
