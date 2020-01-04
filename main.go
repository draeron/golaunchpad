package main

import (
  "fmt"
  "github.com/draeron/golaunchpad/pkg/device"
  "github.com/draeron/golaunchpad/pkg/device/event"
  "github.com/draeron/golaunchpad/pkg/minimk3"
  color2 "github.com/draeron/gopkg/color"
  "github.com/draeron/gopkg/logger"
  "image/color"
  "os"
  "os/signal"
)

var log = logger.New("main")

func main() {
  log.Info("starting golaunchpad")
  defer log.Info("exiting golaunchpad")

  minimk3.SetLogger(logger.New("minimk3"))
  event.SetLogger(logger.New("event"))
  device.SetLogger(logger.New("device"))

  pad, err := minimk3.Open(minimk3.ProgrammerMode)
  must(err)
  pad.Print()
  defer pad.Close()
  must(pad.Diag())

  pad.EnableDebugLogger()

  pad.ClearColors(color2.Black)

  pad.SetBtnColors([]minimk3.Btn{
   minimk3.BtnRow2,
   minimk3.BtnPad11,
   minimk3.BtnPad18,
   minimk3.BtnPad21,
   minimk3.BtnPad42,
   minimk3.BtnPad88,
  }, []color.Color{
   color2.Black,
   color2.Blue,
   color2.CyanGreen,
   color2.DarkGray,
   color2.Green,
   color2.MagentaRed,
   color2.Blue,
  })
  //pad.SetBtnColor(minimk3.BtnRow2, color.Black)
  //
  //pad.SetBtnColor(minimk3.BtnPad11, color.Red)
  //pad.SetBtnColor(minimk3.BtnPad18, color.CyanGreen)
  //pad.SetBtnColor(minimk3.BtnPad21, color.DarkGray)
  //pad.SetBtnColor(minimk3.BtnPad42, color.Green)
  //pad.SetBtnColor(minimk3.BtnPad88, color.Blue)

  waitExit()
}

func waitExit() {
  sigs := make(chan os.Signal, 1)
  done := make(chan bool, 1)

  signal.Notify(sigs)
  go func(){
    sig := <-sigs
    fmt.Printf("receive signal %v\n", sig)
    done <- true
  }()

  <- done
}

func must(err error) {
  if err != nil {
    panic(err.Error())
  }
}

