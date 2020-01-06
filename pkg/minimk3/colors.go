package minimk3

import (
	"image/color"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/minimk3/cmd"
	"go.uber.org/zap/buffer"
)

type BtnColor struct {
	Btn   button.Button
	Color color.Color
}

const (
	RgbColorInstruction = 3
)

func (m *Controller) SetColorAll(col color.Color) error {
	return m.SetColorMany(button.Values(), col)
}

func (m *Controller) SetColorMany(btns []button.Button, color color.Color) error {
	bcs := []BtnColor{}
	for _, btn := range btns {
		bcs = append(bcs, BtnColor{btn, color})
	}
	return m.SetColors(bcs)
}

func (m *Controller) SetColor(btn button.Button, color color.Color) error {
	id := btn.Id()
	r, g, b, _ := color.RGBA()
	msg := cmd.LedColor.SysEx(3, id, byte(r>>9), byte(g>>9), byte(b>>9))
	return m.device.SendDaw(msg)
}

func (m *Controller) SetColors(sets []BtnColor) error {
	buf := buffer.Buffer{}
	for idx, bc := range sets {
		if !bc.Btn.IsValid() {
			continue
		}

		// the launchpad accept up to a maximum of 81 color spec per message
		if idx > 81 {
			log.Warnf("sending too many colors (%d) in a single message", len(sets))
			break
		}
		r, g, b := toColorSpec(bc.Color)
		if _, err := buf.Write([]byte{RgbColorInstruction, bc.Btn.Id(), r, g, b}); err != nil {
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
