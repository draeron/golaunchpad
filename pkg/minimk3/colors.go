package minimk3

import (
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
	mapp := map[button.Button]color.Color{}
	for _, b := range btns {
		mapp[b] = col
	}
	return m.SetColors(mapp)
}

func (m *Controller) SetColor(btn button.Button, color color.Color) error {
	id := btn.Id()
	r, g, b, _ := color.RGBA()
	msg := cmd.LedColor.SysEx(3, id, byte(r>>9), byte(g>>9), byte(b>>9))
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
		r, g, b := toColorSpec(col)
		if _, err := buf.Write([]byte{RgbColorInstruction, btn.Id(), r, g, b}); err != nil {
			return err
		}
	}
	msg := cmd.LedColor.SysEx(buf.Bytes()...)
	return m.device.SendMidi(msg)
}

func toColorSpec(color color.Color) (byte, byte, byte) {
	r, g, b, _ := color.RGBA()
	return byte(r >> 9), byte(g >> 9), byte(b >> 9)
}
