module github.com/draeron/golaunchpad

go 1.16

require (
	github.com/TheCodeTeam/goodbye v0.0.0-20170927022442-a83968bda2d3
	github.com/draeron/gopkgs v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/rakyll/portmidi v0.0.0-20201020180702-d436ceaa537a // indirect
	go.uber.org/atomic v1.5.1
	go.uber.org/zap v1.13.0
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
)

replace gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi => gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi v0.0.0-20220419143954-65c25ed8cc67

replace github.com/draeron/gopkgs => ../gopkgs
