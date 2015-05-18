package pb

func (m *Event) Metric() *float64 {
	if m == nil {
		return nil
	}
	if m.MetricD != nil {
		metric := *m.MetricD
		return &metric
	}
	if m.MetricF != nil {
		metric := float64(*m.MetricF)
		return &metric
	}
	if m.MetricSint64 != nil {
		metric := float64(*m.MetricSint64)
		return &metric
	}
	return nil
}
