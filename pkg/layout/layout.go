package layout

import (
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/draeron/golaunchpad/pkg/launchpad"
	"github.com/draeron/golaunchpad/pkg/launchpad/button"
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"github.com/draeron/gopkgs/color"
)

type Layout interface {
	Connect(controller launchpad.Controller)
	Disconnect()
	Activate()
	Deactivate()

	SetHandler(htype HandlerType, handler Handler)
	SetHandlerHold(htype HandlerType, handler HoldHandler)
	SetHoldTimer(htype HandlerType, duration time.Duration)
	SetDefaultHoldTimer(duration time.Duration)

	HoldTime(btn button.Button) time.Duration
	IsPressed(btn button.Button) bool
	IsHold(btn button.Button, threshold time.Duration) bool

	launchpad.Colorer

	SetName(s string)
	Name() string
}

type BasicLayout struct {
	name       string
	state      launchpad.ButtonStateMap
	lastColors button.ColorMap
	controler  launchpad.Controller
	handlers   handlersMap
	enabled    atomic.Bool
	eventsCh   chan (event.Event)
	mask       launchpad.Mask
	mutex      sync.RWMutex
	ticker     *time.Ticker

	holdTimer        map[HandlerType]time.Duration
	holdTimerDefault time.Duration
}

type handlersMap map[HandlerType]HoldHandler

const DefaultHoldDuration = time.Millisecond * 250

func NewLayoutPreset(preset launchpad.MaskPreset) *BasicLayout {
	return NewLayout(preset.Mask())
}

func NewLayout(mask launchpad.Mask) *BasicLayout {
	l := &BasicLayout{
		state:            launchpad.NewButtonStateMap(),
		lastColors:       button.ColorMap{},
		handlers:         handlersMap{},
		mask:             mask,
		holdTimerDefault: DefaultHoldDuration,
		holdTimer:        map[HandlerType]time.Duration{},
	}
	l.state.SetColors(mask, color.Black) // allocated state
	return l
}

func (l *BasicLayout) Connect(controller launchpad.Controller) {
	l.mutex.Lock()
	l.controler = controller
	l.mutex.Unlock()

	if l.name != "" {
		log.Infof("connecting layout %s to controller %s", l.name, controller.Name())
	}

	go l.tickEvents()
	go l.tickUpdate()
}

func (l *BasicLayout) Disconnect() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.name != "" {
		log.Infof("disconnecting layout %s from controller %s", l.name, l.controler.Name())
	}

	close(l.eventsCh)
	l.controler = nil
	l.ticker.Stop()
	l.ticker = nil
}

/*
	When enabling a layout, it will transfert it's color state to
*/
func (l *BasicLayout) Activate() {
	l.enabled.Store(true)
}

/*
	When disabling a layout, any pressed state will be deleted
*/
func (l *BasicLayout) Deactivate() {
	l.enabled.Store(false)

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// clear last displayed
	l.lastColors = button.ColorMap{}
	l.state.ResetPressed()
}

/*
	The handler will be
*/
func (l *BasicLayout) SetHandler(htype HandlerType, handler Handler) {
	l.handlers[htype] = func(layout Layout, btn button.Button, first bool) {
		if first {
			handler(layout, btn)
		}
	}
}

func (l *BasicLayout) SetHandlerHold(htype HandlerType, handler HoldHandler) {
	l.handlers[htype] = handler
}

func (l *BasicLayout) SetHoldTimer(htype HandlerType, duration time.Duration) {
	l.holdTimer[htype] = duration
}

func (l *BasicLayout) SetDefaultHoldTimer(duration time.Duration) {
	l.holdTimerDefault = duration
}

func (l *BasicLayout) HoldTime(btn button.Button) time.Duration {
	return l.state.HoldTime(btn)
}

func (l *BasicLayout) IsPressed(btn button.Button) bool {
	return l.state.IsPressed(btn)
}

func (l *BasicLayout) IsHold(btn button.Button, threshold time.Duration) bool {
	return l.state.IsHold(btn, threshold)
}

func (l *BasicLayout) UpdateDevice() error {
	if l.enabled.Load() {
		l.mutex.Lock()
		defer l.mutex.Unlock()

		if l.controler == nil {
			return nil
		}

		colors := l.mask.Intersect(l.state).DiffFrom(l.lastColors)

		if len(colors) > 0 {
			err := l.controler.SetColors(colors)
			l.lastColors = l.lastColors.ApplyFrom(colors)
			return err
		}
	}
	return nil
}

func (l *BasicLayout) Color(btn button.Button) color.Color {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.state.Color(btn)
}

func (l *BasicLayout) SetColorAll(col color.Color) error {
	l.state.SetAllColors(col)
	return nil
}

func (l *BasicLayout) SetColorMask(mask launchpad.MaskPreset, col color.Color) error {
	for b, _ := range mask.Mask() {
		l.state.SetColor(b, col)
	}
	return nil
}

func (l *BasicLayout) SetColorMany(btns []button.Button, color color.Color) error {
	for _, k := range btns {
		l.state.SetColor(k, color)
	}
	return nil
}

func (l *BasicLayout) SetColor(btn button.Button, color color.Color) error {
	l.state.SetColor(btn, color)
	return nil
}

func (l *BasicLayout) SetColors(set button.ColorMap) error {
	l.state.SetColorsMap(set)
	return nil
}

func (l *BasicLayout) tickEvents() {
	l.mutex.Lock()
	l.eventsCh = make(chan event.Event, 20)
	l.controler.Subscribe(l.eventsCh)
	l.mutex.Unlock()

	for e := range l.eventsCh {
		l.dispatch(e)
	}
}

func (l *BasicLayout) tickUpdate() {
	l.mutex.Lock()
	l.ticker = time.NewTicker(time.Second / 60)
	l.mutex.Unlock()

	for range l.ticker.C {
		l.UpdateDevice()
	}
}

func (l *BasicLayout) dispatch(e event.Event) {
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
		l.state.Press(e.Btn)
	}
	l.mutex.Unlock()

	if e.Type == event.Pressed {
		if handle, ok := l.handlers[ht+1]; ok {
			timer := l.holdTimerDefault
			if t, ok := l.holdTimer[ht+1]; ok {
				timer = t
			}
			go func() {
				first := true
				for {
					<-time.After(timer)
					if l.state.IsHold(e.Btn, timer) {
						handle(l, e.Btn, first)
						if first {
							first = false
						}
					} else {
						return
					}
				}
			}()
		}
	}

	if h, ok := l.handlers[ht]; ok {
		h(l, e.Btn, true)
	}

	l.mutex.Lock()
	if e.Type == event.Released {
		l.state.Release(e.Btn)
	}
	l.mutex.Unlock()
}

func (l *BasicLayout) SetName(s string) {
	l.name = s
}

func (l *BasicLayout) Name() string {
	return l.name
}
