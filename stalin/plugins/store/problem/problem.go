package problem

import (
	"fmt"
	"stalin/message"
	"strings"
)

type Problem struct {
	Hostname    string   `json:"host"`
	Domain      string   `json:"domain"`
	Fqdn        string   `json:"-"`
	Service     string   `json:"service"`
	State       string   `json:"state"`
	Description string   `json:"description"`
	Metric      *float64 `json:"metric"`
	Tags        []string `json:"tags"`
	Time        int64    `json:"time"`
	Ttl         float32  `json:"ttl"`
	LastFail    int64    `json:"last_fail"`
	Key         string   `json:"-"`
	Failed      bool     `json:"-"`
}

func NewProblemFromEvent(e *message.Event, defTtl int64) (*Problem, error) {
	if e.GetHost() == "" || e.GetService() == "" || e.GetService() == "" {
		return nil, fmt.Errorf("Empty event")
	}
	p := &Problem{}
	p.Service = e.GetService()
	p.State = e.GetState()
	p.Description = e.GetDescription()
	p.Tags = e.GetTags()
	p.Ttl = e.GetTtl()
	if e.GetTtl() == 0 {
		p.Ttl = float32(defTtl)
	} else {
		p.Ttl = e.GetTtl()
	}
	p.Time = e.GetTime()
	p.SetFailed()
	if metric, err := e.GetMetric(); err == nil {
		p.Metric = &metric
	}
	p.Domain, p.Hostname = splitDomainHost(e.GetHost())
	p.Fqdn = e.GetHost()
	p.Key = p.Hostname + p.Domain + p.Service
	return p, nil
}

// название last_fail - устоявщееся.
// -если стейт "ok" указывает на время когда последний раз проблемы закончились
// -если стейт не "ok" указывает на время когда проблемы начались
func (p2 *Problem) UpdateLastFail(p1 *Problem) {

	//p2 - new, p1 - old.
	if p1 == nil || p2 == nil {
		return
	}

	// стейт не поменялся, самый частый use-case
	if p2.State == p1.State {
		p2.LastFail = p1.LastFail
		return
	}

	// вдруг стало все плохо
	if p2.Failed && !p1.Failed {
		p2.LastFail = p2.Time
		return
	}

	// вдруг стало все хорошо
	if !p2.Failed && p1.Failed {
		p2.LastFail = p1.Time
		return
	}

}

func (p *Problem) SetFailed() {
	if p.State != "ok" || p.State != "" {
		p.Failed = true
		p.LastFail = p.Time
	}
}

func splitDomainHost(hostname string) (domain, host string) {
	host = strings.Split(hostname, ".")[0]
	domain = strings.Replace(hostname, host+".", "", 1)
	return
}
