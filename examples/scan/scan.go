package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/draeron/golaunchpad/examples/utils"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
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
	utils.Must(pad.SetColorAll(color.Black))
	defer pad.SetColorAll(color.Black)
	defer pad.Close()

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()

	for _, b := range button.Values() {
		select {
		default:
			pad.SetColor(b, color.White)
			<-time.After(time.Millisecond * 50)
			pad.SetColor(b, color.Black)
		case <-done:
			return
		}
	}
}
