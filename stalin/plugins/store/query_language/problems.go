package ql

type Problems struct {
	Items map[string]*Problem
}

func NewProblems() *Problems {
	return &Problems{
		Items: make(map[string]*Problem, 0),
	}
}

func (p *Problems) List() []*Problem {
	result := make([]*Problem, 0)
	for _, problem := range p.Items {
		result = append(result, problem)
	}
	return result
}

func (q *Problems) Size() int {
	return len(q.Items)
}

func (q *Problems) Exists(p *Problem) bool {
	_, ok := q.Items[p.Key]
	return ok
}

func (q *Problems) Add(p *Problem) {
	q.Items[p.Key] = p
}

func (q1 *Problems) Or(q2 *Problems) (q *Problems) {
	if q1.Size() > q2.Size() {
		q = orFast(q1, q2)
	} else {
		q = orFast(q2, q1)
	}
	return
}

func orFast(big, small *Problems) *Problems {
	for _, p := range small.Items {
		big.Add(p)
	}
	return big
}

func (q1 *Problems) And(q2 *Problems) (q *Problems) {
	if q1.Size() > q2.Size() {
		q = andFast(q1, q2)
	} else {
		q = andFast(q2, q1)
	}
	return
}

func andFast(big, small *Problems) *Problems {
	q := NewProblems()
	for _, p := range small.Items {
		if big.Exists(p) {
			q.Add(p)
		}
	}
	return q
}
