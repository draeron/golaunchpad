package main

import (
	"fmt"
	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var pad launchpad.Controller

func main() {
	log.Info("starting text example")
	defer log.Info("exiting text example")
	pad = common.Setup()
	defer pad.Close()

	setup()

	common.WaitExit()

	pad.DisplayText("", false, 0x40, color.White)
}

func setup() {
	pad.SetColorAll(color.Black)

	for _, b := range button.Rows() {
		pad.SetColor(b, colorFromRow(b))
	}

	pad.SetColor(button.Up, color.LightGray)
	pad.SetColor(button.Down, color.LightGray)
	pad.SetColor(button.User, color.LightGray)

	go func() {
		ch := make(chan event.Event, 10)
		pad.Subscribe(ch)

		lastColor := color.Black.RGB()
		speed := byte(15)
		loop := false

		for e := range ch {
			if e.Type == event.Pressed {
				switch {
				case e.Btn == button.User:
					loop = !loop
					pad.DisplayText("", loop, speed, lastColor)
					if loop {
						pad.SetColor(button.User, color.Green)
					} else {
						pad.SetColor(button.User, color.White)
					}

				case e.Btn.IsRow():
					lastColor = colorFromRow(e.Btn)
					pad.DisplayText("", loop, speed, lastColor)

				case e.Btn.IsArrow():
					switch e.Btn {
					case button.Up:
						speed++
					case button.Down:
						speed--
					}
					pad.DisplayText("", loop, speed, lastColor)

				case e.Btn.IsPad():
					x, y := e.Btn.Coord()
					lastColor = common.RandColor()
					pad.DisplayText(fmt.Sprintf("X: %d, Y: %d", x, y), loop, speed, lastColor)
				}
			}
		}
	}()
}

func colorFromRow(b button.Button) color.RGB {
	return color.Palette[b-button.Row1]
}
