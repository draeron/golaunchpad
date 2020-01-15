package minimk3

import (
	seven_bits "github.com/draeron/golaunchpad/pkg/colors/7bits"
	"image/color"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/minimk3/cmd"
	"go.uber.org/zap/buffer"
)

const (
	RgbColorInstruction = 3
)

func (m *Controller) SetColorAll(col color.Color) error {
	return m.SetColorMany(button.Values(), col)
}

func (m *Controller) SetColorMany(btns []button.Button, col color.Color) error {
	mapp := button.ColorMap{}
	t := seven_bits.FromColor(col)
	for _, b := range btns {
		mapp[b] = t
	}
	return m.SetColors(mapp)
}

func (m *Controller) SetColor(btn button.Button, color color.Color) error {
	id := btn.Id()
	rgb := seven_bits.FromColor(color)
	// convert from 16 bits to 7 bits
	msg := cmd.LedColor.SysEx(3, id, rgb.R, rgb.G, rgb.B)
	return m.device.SendDaw(msg)
}

func (m *Controller) SetColors(mapp button.ColorMap) error {
	buf := buffer.Buffer{}
	idx := 0
	for btn, col := range mapp {
		if !btn.IsValid() {
			continue
		}

		idx++

		// the launchpad accept up to a maximum of 81 color spec per message
		if idx > 81 {
			log.Warnf("sending too many colors (%d) in a single message", len(mapp))
			break
		}
		if _, err := buf.Write([]byte{RgbColorInstruction, btn.Id(), col.R, col.G, col.B}); err != nil {
			return err
		}
	}
	msg := cmd.LedColor.SysEx(buf.Bytes()...)
	return m.device.SendMidi(msg)
}
