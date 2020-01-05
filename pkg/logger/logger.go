package logger

import (
	"fmt"
	"os"
)

/*
  This match the signature of go.uber.org/zap sugarred logger
*/
type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Dummy struct{}

func (d Dummy) Debugf(template string, args ...interface{}) {
}

func (d Dummy) Infof(template string, args ...interface{}) {
}

func (d Dummy) Warnf(template string, args ...interface{}) {
}

func (d Dummy) Errorf(template string, args ...interface{}) {
}

func (d Dummy) DPanicf(template string, args ...interface{}) {
	panic(fmt.Sprintf(template, args))
}

func (d Dummy) Panicf(template string, args ...interface{}) {
	panic(fmt.Sprintf(template, args))
}

func (d Dummy) Fatalf(template string, args ...interface{}) {
	os.Exit(1)
}
