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
	mainLayout := layout.NewLayoutPreset(mask.ButtonsAll)
	mainLayout.SetName("blend")
	mainLayout.Connect(ctrl)
	mainLayout.Activate()

	common.Must(mainLayout.SetColorAll(color.Black))
	common.Must(mainLayout.SetColor(button.Logo, color.LightGray))
	common.Must(mainLayout.SetColorMany(button.Modes(), color.LightGray))

	for i, b := range button.Rows() {
		common.Must(mainLayout.SetColor(b, color.PaletteColor(i)))
	}

	const blendtime = time.Millisecond * 400

	mainLayout.SetHandler(layout.RowPressed, func(layout layout.Layout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := mainLayout.Color(btn)
			if color.Black.Equal(color.FromStdColor(old)) {
				blend := dcolor.Blend(old, col, blendtime)
				common.Must(layout.SetColor(btn, blend))
			}
		}
	})

	// mainLayout.SetHandler(mainLayout.RowHold, func(layout *mainLayout.BasicLayout, btn button.Button) {
	mainLayout.SetHandler(layout.RowHold, func(layout layout.Layout, btn button.Button) {
		col := color.PaletteColor(btn - button.Row1)
		for _, btn := range button.Pads() {
			old := mainLayout.Color(btn)
			blend := dcolor.Blend(old, col, time.Millisecond*100)
			common.Must(layout.SetColor(btn, blend))
		}
	})

	// stop blending by assigning a color to the button
	mainLayout.SetHandler(layout.PadPressed, func(layout layout.Layout, btn button.Button) {
		layout.SetColor(btn, color.Black)
	})
}
