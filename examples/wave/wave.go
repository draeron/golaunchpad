package main

import (
	"fmt"
	"github.com/draeron/golaunchpad/pkg/device"
	devevt "github.com/draeron/golaunchpad/pkg/device/event"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/golaunchpad/pkg/minimk3/button"
	"github.com/draeron/golaunchpad/pkg/minimk3/event"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var log = logger.New("main")
var pad *minimk3.Controller

func main() {
	log.Info("starting golaunchpad")
	defer log.Info("exiting golaunchpad")

	minimk3.SetLogger(logger.New("minimk3"))
	devevt.SetLogger(logger.New("event"))
	device.SetLogger(logger.New("device"))

	var err error
	pad, err = minimk3.Open(minimk3.ProgrammerMode)
	must(err)
	defer pad.Close()
	must(pad.Diag())

	//pad.EnableDebugLogger()
	setup()

	waitExit()
}

func setup() {
	must(pad.ClearColors(color.Black))
	must(pad.SetBtnColor(button.Logo, color.LightGray))

	go func() {
		ch := make(chan event.Event, 10)
		pad.Subscribe(ch)
		for e := range ch {
			if e.Btn.IsPad() {
				if e.Type == event.Pressed {
					col := randColor()
					err := pad.SetBtnColor(e.Btn, col)
					log.LogIfErr(err)
					go wavefx(e.Btn, col)
				} else if e.Type == event.Released {
					err := pad.SetBtnColor(e.Btn, color.Black)
					log.LogIfErr(err)
				}
			}
		}
	}()
}

func wavefx(btn button.Button, col color.Color) {
	x,y := btn.Coord()
	for radius := 1; radius < 8; radius++ {
		btns := []minimk3.BtnColor{
			{button.FromXY(x+radius, y), col},
			{button.FromXY(x, y+radius), col},
			{button.FromXY(x-radius, y), col},
			{button.FromXY(x, y-radius), col},
		}
		pad.SetBtnColors(btns)
		<- time.After(time.Millisecond * 50)
		btns[0].Color = color.Black
		btns[1].Color = color.Black
		btns[2].Color = color.Black
		btns[3].Color = color.Black
		pad.SetBtnColors(btns)
	}
}

func randColor() color.Color {
	col := color.Black.RGB()
	for col == color.Black.RGB() {
		col = color.Palette[rand.Intn(len(color.Palette)-1)]
	}
	return col
}

func waitExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()

	<-done
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
