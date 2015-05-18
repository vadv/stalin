package store

import (
	"net/http"
	"runtime"
	"time"

	"stalin/events"
	. "stalin/plugins"
	store "stalin/plugins/store/memory"
	ql "stalin/plugins/store/query_language"
)

func init() {
	Plugin.Creator(newStorage).Type("Storage").Description("Persistent storage with search api.").Register()
}

type Storage struct {
	Log             *Logger `json:"-"`
	Address         string  `json:"address" description:"Listen address"`
	DBPath          string  `json:"dbpath" description:"Database storage path"`
	CleanUpInterval int     `json:"cleanup_interval" description:"Clear old event period, in sec"`
	SaveInterval    int     `json:"save_interval" description:"Save period for storage, in sec"`
	Store           *store.Store
}

func newStorage(name string) PluginInterface {
	return &Storage{
		Log:             NewLog(name),
		Address:         "127.0.0.1:7777",
		DBPath:          "/var/lib/stalin",
		CleanUpInterval: 60,
		SaveInterval:    60 * 5,
	}
}

func (s *Storage) Inject(msg *events.Events) error {
	if s.Store == nil {
		return nil
	}
	for _, event := range msg.List() {
		s.Store.SetProblem(ql.NewProblemFromEvent(event))
	}
	return nil
}

func (s *Storage) Start() error {
	begin := time.Now()
	s.Log.Info("Load db from path: %v", s.DBPath)
	db, err := store.NewStore(s.DBPath, time.Duration(s.CleanUpInterval)*time.Second, s.Log)
	if err != nil {
		return err
	}
	s.Store = db
	s.Log.Info("DB %v loaded (items: %d, spend: %v)", s.Store.Filename(), s.Store.Items.ItemCount(), time.Now().Sub(begin))
	s.Log.Info("Starting http-api at address: %v", s.Address)
	go s.Tik()
	http.HandleFunc("/query", s.Query)
	if err := http.ListenAndServe(s.Address, nil); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Tik() {
	saveStorage := time.Tick(time.Duration(s.SaveInterval) * time.Second)
	for {
		select {
		case <-saveStorage:
			begin := time.Now()
			if err := s.Store.Save(); err != nil {
				s.Log.Error("Can't save storage: error %v (spend: %v)", err.Error(), time.Now().Sub(begin))
			} else {
				s.Log.Info("Store saved to %s (items: %d, spend: %v)", s.Store.Filename(), s.Store.Items.ItemCount(), time.Now().Sub(begin))
			}
			runtime.GC()
		}
	}
}
