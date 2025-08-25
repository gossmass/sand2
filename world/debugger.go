package world

import (
	"sand2/events"
)

type Debugger struct {
}

func NewDebugger() *Debugger {
	return &Debugger{}
}

func (d *Debugger) OnNotify(obj any, event events.Event) {
	switch event {
	case events.EVENT_CHUNK_UPDATE:
		break
	}
}
