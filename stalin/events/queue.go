package events

import (
	"sync"
)

type queueEvent struct {
	mutex  *sync.Mutex
	events []*StalinEvent
	count  int
}

func NewQueue() *queueEvent {
	return &queueEvent{mutex: &sync.Mutex{}, events: make([]*StalinEvent, 0)}
}

func (q *queueEvent) Push(e *StalinEvent) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.events = append(q.events, e)
	q.count++
}

func (q *queueEvent) Pop() *StalinEvent {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.count == 0 {
		return nil
	}
	q.count--
	return q.events[q.count]
}
