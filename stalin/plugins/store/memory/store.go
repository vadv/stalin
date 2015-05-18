package storage

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "stalin/plugins"
	. "stalin/plugins/store/query_language"
)

type Store struct {
	DBPath string
	Items  *Cache
	Log    *Logger
}

func (s *Store) Filename() string {
	return filepath.Join(s.DBPath, "store.gob")
}

func NewStore(dbpath string, janitorTtl time.Duration, log *Logger) (*Store, error) {
	store := &Store{
		Items:  NewCache(janitorTtl),
		DBPath: dbpath,
		Log:    log,
	}
	gob.Register(&Problem{})
	return store, store.Init()
}

func (s *Store) Init() error {
	if fd, err := os.Open(s.Filename()); err != nil {
		return s.save()
	} else {
		defer fd.Close()
		return s.load()
	}
}

func (s *Store) Save() error {
	return s.save()
}

func (s *Store) SetProblem(p2 *Problem) {
	p1, found := s.Items.Get(p2.Key)
	if found {
		p2.UpdateLastFail(p1)
	}
	s.Items.Set(p2.Key, p2, time.Duration(p2.Ttl)*time.Second)
}

func (s *Store) save() error {
	tempfile, err := ioutil.TempFile(s.DBPath, "tmp-")
	if err != nil {
		return err
	}
	savecache := NewCacheFrom(0, s.Items.List())
	if err := savecache.Save(tempfile); err != nil {
		return err
	}
	tempfile.Close()
	savecache = nil
	return os.Rename(tempfile.Name(), s.Filename())
}

func (s *Store) load() error {
	return s.Items.LoadFile(s.Filename())
}
