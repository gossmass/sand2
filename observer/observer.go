package observer

import "slices"
import "sand2/events"

type Observer interface {
	OnNotify(obj any, event events.Event)
}

type Observerable struct {
	observers []Observer
}

func (b *Observerable) Nofify(event events.Event, args ...any) {
	for _, o := range b.observers {
		o.OnNotify(o, event)
	}
}

func (b *Observerable) AddObserver(o Observer) {
	b.observers = append(b.observers, o)
}

func (b *Observerable) RemoveObserver(o Observer) {
	if idx := slices.Index(b.observers, o); idx != -1 {
		b.observers = slices.Delete(b.observers, idx, idx+1)
	}
}
