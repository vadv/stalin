package message

func NewMessage() *Message {
	return &Message{Events: make([]*Event, 0)}
}

func (m *Message) AddEvent(event *Event) {
	if m != nil && m.Events != nil {
		m.Events = append(m.Events, event)
	}
}

func NewDataMessage(dataType string, data string) *Message {
	msg := NewMessage()
	msg.Data = &data
	msg.DataType = &dataType
	return msg
}
