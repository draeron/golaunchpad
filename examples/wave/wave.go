package main

import (
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"time"

	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

var log = logger.New("main")
var pad launchpad.Controller

func main() {
	log.Info("starting wave example")
	defer log.Info("exiting wave example")
	pad = common.Setup()
	defer pad.Close()

	setup()

	common.WaitExit()
}

func setup() {
	common.Must(pad.SetColorAll(color.Black))
	common.Must(pad.SetColor(button.Logo, color.LightGray))

	go func() {
		ch := make(chan event.Event, 10)
		pad.Subscribe(ch)
		for e := range ch {
			if e.Btn.IsPad() {
				if e.Type == event.Pressed {
					col := common.RandColor()
					err := pad.SetColor(e.Btn, col)
					log.LogIfErr(err)
					go wavefx(e.Btn, col)
				} else if e.Type == event.Released {
					err := pad.SetColor(e.Btn, color.Black)
					log.LogIfErr(err)
				}
			}
		}
	}()
}

func wavefx(btn button.Button, col color.Color) {
	x, y := btn.Coord()
	for radius := 1; radius < 8; radius++ {
		btns := []button.Button{
			button.FromXY(x+radius, y),
			button.FromXY(x, y+radius),
			button.FromXY(x-radius, y),
			button.FromXY(x, y-radius),
		}
		pad.SetColorMany(btns, col)
		<-time.After(time.Millisecond * 50)
		pad.SetColorMany(btns, color.Black)
	}
}
