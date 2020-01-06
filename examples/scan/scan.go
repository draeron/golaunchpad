package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/draeron/golaunchpad/examples/common"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/gopkg/color"
	"github.com/draeron/gopkg/logger"
)

var log = logger.New("main")
var pad launchpad.Controller

func main() {
	log.Info("starting scan example")
	defer log.Info("exiting scan example")
	pad = common.Setup()
	defer pad.Close()

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()

	for {
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
}
