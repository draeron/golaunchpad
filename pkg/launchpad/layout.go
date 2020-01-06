package launchpad

import (
	"github.com/draeron/golaunchpad/pkg/launchpad/event"
	"image/color"
	"time"

	"github.com/draeron/golaunchpad/pkg/launchpad/button"
)

type Layout struct {
	colors    map[button.Button]color.Color
	pressed   map[button.Button]time.Time
	controler Controller
	handlers  map[HandlerType]Handler
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

func (l *Layout) IsHold(btn button.Button, threshold time.Duration) bool {
	return l.HoldTime(btn) > threshold
}

func (l *Layout) dispatch(e event.Event) {
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
			ht = PagePressed
		} else {
			ht = PageReleased
		}
	}

	if e.Type == event.Pressed {
		l.pressed[e.Btn] = time.Now()
	} else {
		delete(l.pressed, e.Btn)
	}

	if h, ok := l.handlers[ht]; ok {
		h(e.Btn)
	}
}
