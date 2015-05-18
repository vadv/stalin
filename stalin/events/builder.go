package events

import (
	"fmt"
	"github.com/lann/builder"
	. "stalin/events/pb"
	"time"
)

type StalinEvent struct {
	UUID        *string
	ORG         *string
	Host        *string
	Service     *string
	Description *string
	Metric      *float64
	TTL         *int
	Time        *int64
	Tags        []string
	State       *string
}

func pStrCopy(str *string) *string {
	str2 := *str
	return &str2
}

func (s *StalinEvent) ToEvent() *Event {
	ttl := float32(*s.TTL)
	pbEvent := &Event{
		SourceUuid:  pStrCopy(s.UUID),
		OrgUuid:     pStrCopy(s.ORG),
		State:       pStrCopy(s.State),
		Service:     pStrCopy(s.Service),
		Host:        pStrCopy(s.Host),
		Description: pStrCopy(s.Description),
		Tags:        s.Tags,
		Ttl:         &ttl,
	}
	if s.Time != nil {
		pbEvent.Time = s.Time
	} else {
		now := time.Now().Unix()
		pbEvent.Time = &now
	}
	return pbEvent
}

type eventBuilder builder.Builder

func (b eventBuilder) Service(format string, a ...interface{}) eventBuilder {
	str := fmt.Sprintf(format, a...)
	return builder.Set(b, "Service", &str).(eventBuilder)
}

func (b eventBuilder) AddTag(format string, a ...interface{}) eventBuilder {
	return builder.Append(b, "Tags", fmt.Sprintf(format, a...)).(eventBuilder)
}

func (b eventBuilder) Host(format string, a ...interface{}) eventBuilder {
	str := fmt.Sprintf(format, a...)
	return builder.Set(b, "Host", &str).(eventBuilder)
}

func (b eventBuilder) Description(format string, a ...interface{}) eventBuilder {
	str := fmt.Sprintf(format, a...)
	return builder.Set(b, "Description", &str).(eventBuilder)
}

func (b eventBuilder) State(format string, a ...interface{}) eventBuilder {
	str := fmt.Sprintf(format, a...)
	return builder.Set(b, "State", &str).(eventBuilder)
}

func (b eventBuilder) Metric(metric float64) eventBuilder {
	return builder.Set(b, "Metric", &metric).(eventBuilder)
}

func (b eventBuilder) Time(metric int64) eventBuilder {
	return builder.Set(b, "Time", &metric).(eventBuilder)
}

func (b eventBuilder) Ttl(metric int) eventBuilder {
	return builder.Set(b, "TTL", &metric).(eventBuilder)
}

func (b eventBuilder) ToStalinEvent() *StalinEvent {
	event := builder.GetStruct(b).(StalinEvent)
	return &event
}

var EventBuilder = builder.Register(eventBuilder{}, StalinEvent{}).(eventBuilder)
