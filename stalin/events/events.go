package events

import (
	"bytes"
	"code.google.com/p/gogoprotobuf/proto"
	"net"

	. "stalin/events/pb"
)

type Events struct {
	stalinEvents []*StalinEvent `json:"-"`
	Events       []*Event       `json:"events"`
	DataType     *string        `json:"data_type"`
	Data         *string        `json:"data"`
}

func NewEvents() *Events {
	return &Events{
		Events:       make([]*Event, 0),
		stalinEvents: make([]*StalinEvent, 0),
	}
}

func ReadConn(conn net.Conn) (*Events, error) {
	events := &Message{}
	data, err := Unpack(conn)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(data, events); err != nil {
		return nil, err
	}
	return &Events{
		Events:       events.Events,
		stalinEvents: make([]*StalinEvent, 0),
	}, nil
}

func (e *Events) ToPbBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	msg := &Message{Events: e.List()}
	if e.Data != nil && e.DataType != nil {
		typ := *e.DataType
		data := *e.Data
		msg.Data = &data
		msg.DataType = &typ
	}
	if err := enc.EncodeMsg(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *Events) List() []*Event {
	result := make([]*Event, 0)
	for _, event := range e.stalinEvents {
		e.Events = append(e.Events, event.ToEvent())
	}
	e.stalinEvents = make([]*StalinEvent, 0)
	for _, event := range e.Events {
		if event != nil {
			result = append(result, event)
		}
	}
	e.Events = result
	return e.Events
}

func (e *Events) Add(event interface{}) {
	e.stalinEvents = append(e.stalinEvents, event.(eventBuilder).ToStalinEvent())
}

func (e *Events) AddEvents(events *Events) {
	for _, event := range events.Events {
		e.Events = append(e.Events, event)
	}
	for _, event := range events.stalinEvents {
		e.stalinEvents = append(e.stalinEvents, event)
	}
	return
}
