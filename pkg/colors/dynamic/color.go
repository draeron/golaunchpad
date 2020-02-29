package dynamic

import (
	inner "github.com/draeron/gopkgs/color"
	"image/color"
	"sync"
	"time"
)

type DynamicColor struct {
	current [4]int
	target  [4]int
	mutex   sync.Mutex
}

func (c *DynamicColor) RGBA() (r, g, b, A uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return uint32(c.current[0]), uint32(c.current[1]), uint32(c.current[2]), uint32(c.current[3])
}

func (c *DynamicColor) IsEqual(l *DynamicColor) bool {
	return false
}

func Blend(old color.Color, target color.Color, duration time.Duration) *DynamicColor {
	dyn := &DynamicColor{current: toSlice(old)}
	dyn.target = toSlice(target)

	const stepcount = 30

	go func(toupdate *DynamicColor) {
		stepduration := duration / stepcount
		tick := time.NewTicker(stepduration)
		defer tick.Stop()
		_target := toSlice(target)
		_org := toSlice(old)
		cur := toupdate.current
		count := 0

		for count < stepcount {
			for idx := range _org {
				step := (_target[idx] - _org[idx]) / stepcount
				cur[idx] = cur[idx] + step
			}

			toupdate.mutex.Lock()
			for idx := range _org {
				toupdate.current[idx] = cur[idx]
			}
			toupdate.mutex.Unlock()

			<-tick.C

			count++
		}

		toupdate.mutex.Lock()
		for idx := range _target {
			toupdate.current[idx] = _target[idx]
		}
		toupdate.mutex.Unlock()
	}(dyn)

	return dyn
}

func (c *DynamicColor) BlendTo(newTarget color.Color, duration time.Duration) *DynamicColor {
	return Blend(fromSlice(c.current), newTarget, duration)
}

func same(left, right [4]int) bool {
	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}
	return true
}

func fromSlice(rgba [4]int) inner.Color {
	// from 16 bits to 8 bits
	return inner.RGB{
		R: uint8(rgba[0] >> 8),
		G: uint8(rgba[1] >> 8),
		B: uint8(rgba[2] >> 8),
		A: uint8(rgba[3] >> 8),
	}
}

func toSlice(rgba color.Color) [4]int {
	r, g, b, a := rgba.RGBA()
	return [4]int{int(r), int(g), int(b), int(a)}
}
