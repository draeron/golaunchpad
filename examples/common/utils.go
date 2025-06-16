package common

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/logger"
	"github.com/draeron/gopkgs/midi"
)

func WaitExit() {
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

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
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
	logrus.SetLevel(logrus.DebugLevel)
	color.SetLogger(logger.NewLogrus("color"))
	midi.SetLogger(logger.NewLogrus("midi"))
	minimk3.SetLogger(logger.NewLogrus("minimk3"))
	launchpad.SetLogger(logger.NewLogrus("launchpad"))

	// StartProfiling()

	pad, err := minimk3.Open(minimk3.ProgrammerMode)
	Must(err)
	Must(pad.Diag())

	err = pad.SetColorAll(color.Black)
	if err != nil {
		log.Fatalf("could not reset colors: %+v", err)
	}

	pad.EnableDebugLogger()

	return pad
}
