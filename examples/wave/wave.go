package main

import (
	"fmt"
	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/device/event"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
	"math/rand"
	"os"
	"os/signal"
)

var log = logger.New("main")
var pad *minimk3.Controller

func main() {
	log.Info("starting golaunchpad")
	defer log.Info("exiting golaunchpad")

	minimk3.SetLogger(logger.New("minimk3"))
	event.SetLogger(logger.New("event"))
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

	btns := []minimk3.BtnColor{}
	for i := minimk3.BtnRow1; i <= minimk3.BtnStopSoloMute; i++ {
		btns = append(btns, minimk3.BtnColor{i, color.Palette[i-minimk3.BtnRow1]})
	}
	must(pad.SetBtnColors(btns))

	must(pad.SetBtnColor(minimk3.BtnLogo, color.LightGray))

	go func() {
		ch := make(chan minimk3.Event, 4)
		pad.Subscribe(ch)
		for e := range ch {
			if e.Btn.IsPad() {
				if e.Type == minimk3.EventTypePressed {
					col := color.Palette[rand.Intn(len(color.Palette)-1)]
					err := pad.SetBtnColor(e.Btn, col)
					log.LogIfErr(err)
				} else if e.Type == minimk3.EventTypeReleased {
					err := pad.SetBtnColor(e.Btn, color.Black)
					log.LogIfErr(err)
					//time.AfterFunc(time.Millisecond * 100, func() {
					//})
				}
			}
		}
	}()
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
