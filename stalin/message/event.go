package message

func (e *Event) SetHost(host string) {
	if e != nil {
		e.Host = &host
	}
}

func (e *Event) SetTime(mytime int64) {
	if e != nil {
		e.Time = &mytime
	}
}

func (e *Event) SetService(service string) {
	if e != nil {
		e.Service = &service
	}
}

func (e *Event) SetTsdbService(service string) {
	if e != nil {
		e.TsdbService = &service
	}
}

func (e *Event) SetTags(tags []string) {
	if e != nil {
		e.Tags = tags
	}
}

func (e *Event) SetTsdbTags(tags []string) {
	if e != nil {
		e.TsdbTags = tags
	}
}

func (e *Event) SetDescription(description string) {
	if e != nil {
		e.Description = &description
	}
}

func (e *Event) SetTtl(ttl float32) {
	if e != nil {
		e.Ttl = &ttl
	}
}

func (e *Event) SetMetric(metric float64) {
	if e != nil {
		e.MetricD = &metric
	}
}
