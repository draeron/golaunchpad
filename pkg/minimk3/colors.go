package minimk3

import (
	"github.com/draeron/golaunchpad/pkg/minimk3/button"
	"github.com/draeron/golaunchpad/pkg/minimk3/cmd"
	"go.uber.org/zap/buffer"
	"image/color"
)

type BtnColor struct {
	Btn   button.Button
	Color color.Color
}

const (
	RgbColorInstruction = 3
)

func (m *Controller) SetBtnColor(btn button.Button, color color.Color) error {
	id := btn.Id()
	r, g, b, _ := color.RGBA()
	msg := cmd.LedColor.SysEx(3, id, byte(r>>9), byte(g>>9), byte(b>>9))
	return m.device.SendDaw(msg)
}

func (m *Controller) ClearColors(col color.Color) error {
	btns := button.Values()
	bcs := []BtnColor{}
	for _, btn := range btns {
		bcs = append(bcs, BtnColor{btn, col})
	}
	return m.SetBtnColors(bcs)
}

func (m *Controller) SetBtnColors(sets []BtnColor) error {
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
		r, g, b, _ := bc.Color.RGBA()
		if _, err := buf.Write([]byte{RgbColorInstruction, bc.Btn.Id(), byte(r >> 9), byte(g >> 9), byte(b >> 9)}); err != nil {
			return err
		}
	}
	msg := cmd.LedColor.SysEx(buf.Bytes()...)
	return m.device.SendMidi(msg)
}
