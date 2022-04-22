package common

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/TheCodeTeam/goodbye"
	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
)

func WaitExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	ctx := context.Background()
	goodbye.Notify(ctx)
	// defer goodbye.Exit(ctx, -1)

	signal.Notify(sigs)
	go func() {
		sig := <-sigs
		fmt.Printf("receive signal %v\n", sig)
		done <- true
	}()

	<-done
}

func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func RandColor() color.RGB {
	col := color.Black.RGB()
	for col == color.Black.RGB() {
		col = color.Palette[rand.Intn(len(color.Palette)-1)]
	}
	return col
}

func Setup() launchpad.Controller {
	color.SetLogger(logger.NewLogrus("color"))
	device.SetLogger(logger.NewLogrus("device"))
	minimk3.SetLogger(logger.NewLogrus("minimk3"))

	// StartProfiling()

	pad, err := minimk3.Open(minimk3.ProgrammerMode)
	Must(err)
	Must(pad.Diag())

	// pad.EnableDebugLogger()

	return pad
}
