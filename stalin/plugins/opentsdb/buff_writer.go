package opentsdb

import (
	"fmt"

	"stalin/events/pb"
)

var (
	errorNoMetric = fmt.Errorf("No metric")
)

type buffWriterEvent struct {
	*pb.Event
}

func (b *buffWriterEvent) toTsdb() []byte {
	service := *(b.Event.TsdbService)
	metric := b.Event.Metric()
	time := *(b.Event.Time)
	tags := ""
	for _, tag := range b.Event.TsdbTags {
		tags = fmt.Sprintf("%s %s", tags, tag)
	}
	return []byte(fmt.Sprintf("put %s %v %v %s\n", service, time, metric, tags))
}

func (b *buffWriterEvent) ToBytes() ([]byte, error) {
	if b.Event.TsdbService == nil || b.Event.Metric() == nil {
		return []byte{}, errorNoMetric
	}
	return b.toTsdb(), nil
}
