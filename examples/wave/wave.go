package main

import (
	"time"

	"github.com/draeron/golaunchpad/examples/utils"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var pad *minimk3.Controller

func main() {
	log.Info("starting golaunchpad")
	defer log.Info("exiting golaunchpad")

	minimk3.SetLogger(logger.New("minimk3"))

	var err error
	pad, err = minimk3.Open(minimk3.ProgrammerMode)
	utils.Must(err)
	defer pad.Close()
	utils.Must(pad.Diag())

	//pad.EnableDebugLogger()
	setup()

	utils.WaitExit()
}

func setup() {
	utils.Must(pad.SetColorAll(color.Black))
	utils.Must(pad.SetColor(button.Logo, color.LightGray))

	go func() {
		ch := make(chan event.Event, 10)
		pad.Subscribe(ch)
		for e := range ch {
			if e.Btn.IsPad() {
				if e.Type == event.Pressed {
					col := utils.RandColor()
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
		btns := []minimk3.BtnColor{
			{button.FromXY(x+radius, y), col},
			{button.FromXY(x, y+radius), col},
			{button.FromXY(x-radius, y), col},
			{button.FromXY(x, y-radius), col},
		}
		pad.SetColors(btns)
		<-time.After(time.Millisecond * 50)
		btns[0].Color = color.Black
		btns[1].Color = color.Black
		btns[2].Color = color.Black
		btns[3].Color = color.Black
		pad.SetColors(btns)
	}
}

