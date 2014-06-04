package riemann

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNoMetric = errors.New("No metric aviable")

func (m *Event) GetMetric() (float64, error) {
	if m.MetricD != nil {
		return *m.MetricD, nil
	}
	if m.MetricSint64 != nil {
		return float64(*m.MetricSint64), nil
	}
	if m.MetricF != nil {
		return float64(*m.MetricF), nil
	}
	return 0, ErrNoMetric
}

func (m *Event) GetSQLTags() string {
	return fmt.Sprintf("{%s}", strings.Join(m.GetTags(), ","))
}

func reverse(a []string) []string {
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

func (m *Event) GetGraphiteService() string {
	service := strings.Replace(*m.Service, " ", ".", -1)
	service = strings.Replace(service, ",", "", -1) // нечитаемый символ
	// трэшовые hosts
	host := *m.Host
	if strings.Count(host, ".") < 2 {
		host = host + ".please_use_fqdn"
	}
	host = strings.Join(reverse(strings.Split(host, ".")), ".")
	// end hosts
	return strings.Replace(fmt.Sprintf("%s.%s", host, service), "\n", "", -1)
}

func (m *Event) GetHostDomain() (string, string) {
	fqdn := *m.Host
	host := strings.Split(fqdn, ".")[0]
	domain := strings.Replace(fqdn, fmt.Sprintf("%s.", host), "", 1)
	return host, domain
}

func (m *Event) ToGraphite() (string, error) {
	metric, err := m.GetMetric()
	if err != nil {
		return "", ErrNoMetric
	}
	// todo: микрооптимизации по формату metric
	return fmt.Sprintf("%s %v %v\n", m.GetGraphiteService(), metric, m.GetTime()), nil
}

func (m *Event) GetOpentsDBService() string {
	service := strings.Replace(*m.Service, " ", ".", -1)
	service = strings.Replace(service, ",", "", -1) // нечитаемый символ
	return service
}

func (m *Event) ToOpentsDB() (string, error) {
	metric, err := m.GetMetric()
	if err != nil {
		return "", ErrNoMetric
	}
	return fmt.Sprintf("put %s %v %v host=%v\n", m.GetOpentsDBService(), m.GetTime(), metric, m.GetHost()), nil
}
