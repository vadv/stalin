package graphite

import (
	"fmt"
	"strings"

	"stalin/events/pb"
)

var errorNoMetric = fmt.Errorf("No metric")

type buffWriterEvent struct {
	*pb.Event
}

func (b *buffWriterEvent) toGraphiteService() string {
	service := strings.Replace(*b.Event.Service, " ", ".", -1)
	service = strings.Replace(service, ",", "", -1)
	host := *b.Event.Host
	host = strings.Join(reverseStringSlice(strings.Split(host, ".")), ".")
	return strings.Replace(fmt.Sprintf("%s.%s", host, service), "\n", "", -1)
}

func (b *buffWriterEvent) ToBytes() ([]byte, error) {
	if b.Event.Metric() == nil {
		return []byte{}, errorNoMetric
	}
	return []byte(fmt.Sprintf("%s %f %d\n", b.toGraphiteService(), *b.Event.Metric(), *b.Event.Time)), nil
}

func reverseStringSlice(a []string) []string {
	if len(a) == 4 {
		a[0], a[1], a[2], a[3] = a[3], a[2], a[1], a[0]
		return a
	}
	if len(a) == 3 {
		a[0], a[1], a[2] = a[2], a[1], a[0]
		return a
	}
	if len(a) == 2 {
		a[0], a[1] = a[1], a[0]
		return a
	}
	lp, rp := 0, len(a)-1
	for (lp != rp) && (lp < rp) {
		a[lp], a[rp] = a[rp], a[lp]
		lp++
		rp--
	}
	return a
}
