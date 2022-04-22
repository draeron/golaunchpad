package event

import (
	"github.com/draeron/gopkgs/logger"
)

var log logger.Logger = logger.Dummy{}

func SetLogger(newlogger logger.Logger) {
	log = newlogger
}
