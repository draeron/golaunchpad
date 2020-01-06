package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/atomic"
	"image/color"
	"time"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

type Layout struct {
	DebugName  string
	colors     button.ColorMap
	lastColors button.ColorMap
	pressed    map[button.Button]time.Time
	controler  Controller
	handlers   handlersMap
	enabled    atomic.Bool
	eventsCh   chan (event.Event)
	mask       Mask
	mutex      deadlock.RWMutex
	ticker     *time.Ticker
}

type handlersMap map[HandlerType]Handler

func NewLayoutPreset(preset MaskPreset) *Layout {
	return NewLayout(preset.Mask())
}

func NewLayout(mask Mask) *Layout {
	return &Layout{
		colors:     button.ColorMap{},
		lastColors: button.ColorMap{},
		pressed:    map[button.Button]time.Time{},
		handlers:   handlersMap{},
		mask:       mask,
	}
}

func (l *Layout) Connect(controller Controller) {
	l.mutex.Lock()
	l.controler = controller
	l.mutex.Unlock()

	go l.tickEvents()
	go l.tickUpdate()
}

func (l *Layout) Disconnect() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	close(l.eventsCh)
	l.controler = nil
	l.ticker.Stop()
	l.ticker = nil
}

/*
	When enabling a layout, it will transfert it's color state to
*/
func (l *Layout) Enable() {
	l.enabled.Store(true)
}

/*
	When disabling a layout, any pressed state will be deleted
*/
func (l *Layout) Disable() {
	l.enabled.Store(false)

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// clear last displayed
	l.lastColors = button.ColorMap{}
	// clear pressed timers
	l.pressed = map[button.Button]time.Time{}
}

func (l *Layout) SetHandler(htype HandlerType, handler Handler) {
	l.handlers[htype] = handler
}

func (l *Layout) HoldTime(btn button.Button) time.Duration {
	if start, ok := l.pressed[btn]; ok {
		return time.Now().Sub(start)
	}
	return 0
}

func (l *Layout) IsPressed(btn button.Button) bool {
	return l.HoldTime(btn) > 0
}

func (l *Layout) IsHold(btn button.Button, threshold time.Duration) bool {
	return l.HoldTime(btn) > threshold
}

func (l *Layout) UpdateDevice() error {
	if l.enabled.Load() {
		l.mutex.Lock()
		defer l.mutex.Unlock()

		if l.controler == nil {
			return nil
		}

		colors := l.mask.Intersect(l.colors).DiffFrom(l.lastColors)

		if len(colors) > 0 {
			err := l.controler.SetColors(colors)
			l.lastColors = l.lastColors.ApplyFrom(colors)
			return err
		}
	}
	return nil
}

func (l *Layout) Color(btn button.Button) color.Color {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if c, ok := l.colors[btn]; ok {
		return c
	} else {
		return color.Black
	}
}

func (l *Layout) SetColorAll(col color.Color) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for k, _ := range l.mask {
		l.colors[k] = col
	}
	return nil
}

func (l *Layout) SetColorMany(btns []button.Button, color color.Color) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, k := range btns {
		l.colors[k] = color
	}
	return nil
}

func (l *Layout) SetColor(btn button.Button, color color.Color) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.colors[btn] = color
	return nil
}

func (l *Layout) SetColors(sets button.ColorMap) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for k, v := range sets {
		l.colors[k] = v
	}
	return nil
}

func (l *Layout) tickEvents() {
	l.mutex.Lock()
	l.eventsCh = make(chan event.Event, 20)
	l.controler.Subscribe(l.eventsCh)
	l.mutex.Unlock()

	for e := range l.eventsCh {
		l.dispatch(e)
	}
}

func (l *Layout) tickUpdate() {
	l.mutex.Lock()
	l.ticker = time.NewTicker(time.Second/60)
	l.mutex.Unlock()

	for range l.ticker.C {
		l.UpdateDevice()
	}
}

func (l *Layout) dispatch(e event.Event) {
	if !l.enabled.Load() || !l.mask[e.Btn] {
		return
	}
	var ht HandlerType

	switch {
	case e.Btn.IsPad():
		if e.Type == event.Pressed {
			ht = PadPressed
		} else {
			ht = PadReleased
		}

	case e.Btn.IsRow():
		if e.Type == event.Pressed {
			ht = RowPressed
		} else {
			ht = RowReleased
		}

	case e.Btn.IsArrow():
		if e.Type == event.Pressed {
			ht = ArrowPressed
		} else {
			ht = ArrowReleased
		}

	case e.Btn.IsMode():
		if e.Type == event.Pressed {
			ht = ModePressed
		} else {
			ht = ModeReleased
		}
	}

	l.mutex.Lock()
	if e.Type == event.Pressed {
		l.pressed[e.Btn] = time.Now()
	} else {
		delete(l.pressed, e.Btn)
	}
	l.mutex.Unlock()

	if h, ok := l.handlers[ht]; ok {
		h(l, e.Btn)
	}
}
