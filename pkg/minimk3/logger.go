package minimk3

import "github.com/draeron/golaunchpad/pkg/logger"

var log logger.Logger = logger.Dummy{}

func SetLogger(newlogger logger.Logger) {
  log = newlogger
}
