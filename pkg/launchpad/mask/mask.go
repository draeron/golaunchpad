package mask

import (
	"maps"
	"slices"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

type Buttons map[button.Button]bool

func From(btns ...button.Button) Buttons {
	return Buttons{}.Add(btns...)
}

func FromPresets(presets ...Preset) Buttons {
	return Buttons{}.MergePreset(presets...)
}

func (m Buttons) Add(btns ...button.Button) Buttons {
	for _, btn := range btns {
		m[btn] = true
	}
	return m
}

func (m Buttons) MergePreset(masks ...Preset) Buttons {
	out := maps.Clone(m)
	for _, mask := range masks {
		out.Merge(mask.Mask())
	}
	return out
}

func (m Buttons) Merge(masks ...Buttons) Buttons {
	out := m
	for _, mask := range masks {
		for b, v := range mask {
			out[b] = v
		}
	}
	return out
}

func (mp Buttons) Slice() []button.Button {
	keys := slices.Collect(maps.Keys(mp))
	slices.SortFunc(keys, func(a, b button.Button) int {
		return int(a.Id() - b.Id())
	})
	return keys
}
