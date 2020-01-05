package main

import (
	"fmt"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/golaunchpad/pkg/minimk3/button"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
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

	var err error
	pad, err = minimk3.Open(minimk3.ProgrammerMode)
	must(err)
	must(pad.ClearColors(color.Black))
	defer pad.ClearColors(color.Black)
	defer pad.Close()

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()


	for _, b := range button.Values(){
		select {
		default:
			pad.SetBtnColor(b, color.White)
			<- time.After(time.Millisecond * 50)
			pad.SetBtnColor(b, color.Black)
		case <-done:
			return
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
