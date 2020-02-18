package common

import (
	"context"
	"fmt"
	"github.com/TheCodeTeam/goodbye"
	"github.com/draeron/golaunchpad/pkg/device"
	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/minimk3"
	"github.com/draeron/gopkgs/logger"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/draeron/gopkgs/color"
)

func WaitExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	ctx := context.Background()
	goodbye.Notify(ctx)
	defer goodbye.Exit(ctx, -1)

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

func StartProfiling() {
	log := logger.New("profiler")

	log.Infof("enabling profiling")

	f, err := os.Create("profile.pprof")
	Must(err)
	err = pprof.StartCPUProfile(f)
	Must(err)

	goodbye.Register(func(ctx context.Context, s os.Signal) {
		pprof.StopCPUProfile()
	})

	go func() {
		tick := time.NewTicker(time.Second * 5)

		goodbye.Register(func(ctx context.Context, s os.Signal) {
			tick.Stop()
		})

		for range tick.C {
			stats := runtime.MemStats{}
			runtime.ReadMemStats(&stats)
			log.Debugf("goroutine count: %d, memory: %d kB", runtime.NumGoroutine(), stats.Alloc/1024)
		}
	}()
}

func Setup() launchpad.Controller {
	device.SetLogger(logger.New("device"))
	minimk3.SetLogger(logger.New("minimk3"))

	//StartProfiling()

	pad, err := minimk3.Open(minimk3.ProgrammerMode)
	Must(err)
	Must(pad.Diag())

	//pad.EnableDebugLogger()

	return pad
}
