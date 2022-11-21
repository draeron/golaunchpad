module github.com/draeron/golaunchpad

go 1.18

require (
	github.com/TheCodeTeam/goodbye v0.0.0-20170927022442-a83968bda2d3
	github.com/draeron/gopkgs v0.4.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	gitlab.com/gomidi/midi/v2 v2.0.25
	go.uber.org/atomic v1.10.0
	go.uber.org/zap v1.23.0
)

require (
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
)

//replace github.com/draeron/gopkgs => ../gopkgs
