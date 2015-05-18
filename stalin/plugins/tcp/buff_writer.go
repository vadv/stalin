package tcp

import (
	"stalin/events"
)

type buffWriterEvent struct {
	*events.Events
}

func (b *buffWriterEvent) ToBytes() ([]byte, error) {
	return b.Events.ToPbBytes()
}
