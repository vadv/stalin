package memory

import (
	"encoding/gob"
	"os"
	. "stalin/plugins"
	"stalin/plugins/store/cache"
	. "stalin/plugins/store/problem"
	"sync"
	"time"
)

type Store struct {
	FileName string
	Items    *cache.Cache
	mutex    sync.RWMutex
}

func NewStore(defaultTtl, janitorTtl time.Duration) *Store {
	store := &Store{Items: cache.New(defaultTtl, janitorTtl)}
	gob.Register(&Problem{})
	return store
}

func (s *Store) SetProblem(p2 *Problem) {
	p1, found := s.Items.Get(p2.Key)
	if found {
		p2.UpdateLastFail(p1)
		p1 = nil
	}
	s.Items.Set(p2.Key, p2, time.Duration(p2.Ttl)*time.Second)
}

func (s *Store) Init() error {
	if _, err := os.Open(s.FileName); err != nil {
		return s.Save()
	} else {
		return s.Load()
	}
}

func (s *Store) Save() error {
	// сохраняем во временный файл
	tempfile := s.FileName + ".tmp-copy"
	fd, err := os.Create(tempfile)
	if err != nil {
		return err
	}
	tBegin := time.Now()
	LogInfo("[MemStore]: Background save to file: %v, items count: %v", tempfile, len(s.Items.Items()))
	if err = s.Items.Save(fd); err != nil {
		LogErr("[MemStore]: Save in file: %v, spend time: %v. Error: %v", tempfile, time.Now().Sub(tBegin), err)
		return err
	}
	LogInfo("[MemStore]: Succesfull save to file: %v, spend time: %v", tempfile, time.Now().Sub(tBegin))
	return os.Rename(fd.Name(), s.FileName)
}

func (s *Store) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Items.LoadFile(s.FileName)
}
