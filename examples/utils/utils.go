package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/draeron/gopkg/color"
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
