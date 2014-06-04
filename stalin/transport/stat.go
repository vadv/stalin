package transport

import (
	"log"
	//"os"
	"sync"
	"time"
)

type Stat struct {
	Input     int
	qInput    int
	Graphite  int
	qGraphite int
	Postgres  int
	qPostgres int
	OpentsDB  int
	qOpentsDB int
	Resend    int
	qResend   int
	mutex     sync.RWMutex
	StatTime  int
}

func (s *Stat) inInput() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.qInput++
}

func (s *Stat) doneInput() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Input++
	s.qInput--
}

func (s *Stat) inPostgres() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.qPostgres++
}

func (s *Stat) donePostgres() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Postgres++
	s.qPostgres--
}

func (s *Stat) inGraphite() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.qGraphite++
}

func (s *Stat) doneGraphite() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Graphite++
	s.qGraphite--
}

func (s *Stat) inOpentsDB() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.qOpentsDB++
}

func (s *Stat) doneOpentsDB() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.OpentsDB++
	s.qOpentsDB--
}

func (s *Stat) inResend() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Resend++
	s.qResend--
}

func (s *Stat) doneResend() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.qResend++
}

func (s *Stat) clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Postgres = 0
	s.Graphite = 0
	s.OpentsDB = 0
	s.Resend = 0
	s.Input = 0
}

func (s *Stat) print() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	log.Printf("[Statistics] Input: [%6.2f/%v], Resend: [%6.2f/%v], Graphite: [%6.2f/%v], OpentsDB:[%6.2f/%v], Pg: [%6.2f/%v]",
		float64(s.Input)/float64(s.StatTime), s.qInput, float64(s.Resend)/float64(s.StatTime), s.qResend,
		float64(s.Graphite)/float64(s.StatTime), s.qGraphite, float64(s.OpentsDB)/float64(s.StatTime), s.qOpentsDB, float64(s.Postgres)/float64(s.StatTime), s.qPostgres)
}

// todo: send to graphite metric
//func (s *Stat) GraphiteReport() []byte {
//s.mutex.RLock()
//defer s.mutex.RUnlock()
//result := []string{}
//append(result, "")
//}

func (s *Stat) Ticks() {
	c := time.Tick(time.Duration(s.StatTime) * time.Second)
	for {
		select {
		case <-c:
			s.print()
			s.clear()
		}
	}
}
