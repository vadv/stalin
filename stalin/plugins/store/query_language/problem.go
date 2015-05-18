package ql

import (
	"fmt"
	"strings"
	"time"

	. "stalin/events/pb"
)

type Problem struct {
	Fqdn            string   `json:"fqdn"`
	Hostname        string   `json:"host"`
	Domain          string   `json:"domain"`
	Service         string   `json:"service,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	TsdbTags        []string `json:"tsdb_tags,omitempty"`
	TsdbService     string   `json:"tsdb_service,omitempty"`
	GraphiteService string   `json:"graphite_service,omitempty"`
	State           string   `json:"state,omitempty"`
	Description     string   `json:"description"`
	Metric          *float64 `json:"metric"`
	Time            int64    `json:"time"`
	Ttl             float32  `json:"ttl"`
	LastFail        int64    `json:"last_fail"`
	Expiration      int64    `json:"-"`
	Key             string   `json:"-"`
	Failed          bool     `json:"-"`
}

func NewProblemFromEvent(e *Event) *Problem {
	if e == nil {
		return nil
	}
	p := &Problem{
		Fqdn:        e.GetHost(),
		Service:     e.GetService(),
		Tags:        e.GetTags(),
		TsdbTags:    e.GetTsdbTags(),
		State:       e.GetState(),
		Description: e.GetDescription(),
		Metric:      e.Metric(),
		Time:        e.GetTime(),
		Ttl:         e.GetTtl(),
	}
	p.Hostname, p.Domain = splitDomainHost(e.GetHost())
	p.SetKey()
	p.SetFailed()
	p.SetGraphiteService()
	p.SetExpiration()
	return p
}

func (p *Problem) SetGraphiteService() {
	service := strings.Replace(p.Service, " ", ".", -1)
	service = strings.Replace(service, ",", "", -1)
	host := strings.Join(reverseStringSlice(strings.Split(p.Fqdn, ".")), ".")
	p.GraphiteService = strings.Replace(fmt.Sprintf("%s.%s", host, service), "\n", "", -1)
}

func (p *Problem) SetKey() {
	p.Key = p.Fqdn + p.Service
}

func (p *Problem) SetFailed() {
	if p.State != "ok" || p.State != "" {
		p.Failed = true
		p.LastFail = p.Time
	}
}

func (p *Problem) SetExpiration() {
	if p.Ttl == 0 {
		p.Expiration = p.Time + (60 * 15)
	} else {
		p.Expiration = p.Time + int64(p.Ttl)
	}
}

func (p2 *Problem) UpdateLastFail(p1 *Problem) {
	if p1 == nil || p2 == nil {
		return
	}
	if p2.State == p1.State {
		p2.LastFail = p1.LastFail
		return
	}
	if p2.Failed && !p1.Failed {
		p2.LastFail = p2.Time
		return
	}
	if !p2.Failed && p1.Failed {
		p2.LastFail = p1.Time
		return
	}
}

func (p1 *Problem) Eq(p2 *Problem) bool {
	return p1.Key == p2.Key
}

func (problem *Problem) GetTextValue(val string) string {
	switch val {
	case "state":
		return problem.State
	case "domain":
		return problem.Domain
	case "fqdn":
		return problem.Fqdn
	case "hostname", "host":
		return problem.Hostname
	case "service":
		return problem.Service
	default:
		return ""
	}
}

func (p *Problem) Expired() bool {
	return time.Now().Unix() > p.Expiration
}

func splitDomainHost(hostname string) (host, domain string) {
	host = strings.Split(hostname, ".")[0]
	domain = strings.Replace(hostname, host+".", "", 1)
	return
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
