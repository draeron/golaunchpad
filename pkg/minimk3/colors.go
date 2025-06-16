package minimk3

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/mask"
	"github.com/draeron/gopkgs/color"
	"github.com/draeron/gopkgs/color/7bits"

	"go.uber.org/zap/buffer"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/minimk3/cmd"
)

const (
	RgbColorInstruction = 3
)

func (m *Controller) SetColorAll(col color.Color) error {
	return m.SetColorMany(button.Values(), col)
}

func (m *Controller) SetColorPad(x, y int, color color.Color) error {
	btn := button.FromPadXY(x, y)
	return m.SetColor(btn, color)
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

func (m *Controller) SetColorMask(mask mask.Buttons, col color.Color) error {
	return m.SetColorMany(mask.Slice(), col)
}

func (m *Controller) SetColorMaskPreset(preset mask.Preset, col color.Color) error {
	return m.SetColorMask(preset.Mask(), col)
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
		tmp := seven_bits.FromColor(col)
		if _, err := buf.Write([]byte{RgbColorInstruction, btn.Id(), tmp.R, tmp.G, tmp.B}); err != nil {
			return err
		}
	}
	msg := cmd.LedColor.SysEx(buf.Bytes()...)
	return m.device.SendMidi(msg)
}
