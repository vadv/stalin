package api

import (
	"net/http"
	"stalin/plugins/store/memory"
)

var API_VERSION = "1"

type Api struct {
	Address string
	Store   *memory.Store
}

func (a *Api) StartApi() error {

	http.HandleFunc("/ping", a.pingHandler)
	http.HandleFunc("/version", a.versionHandler)
	http.HandleFunc("/query", a.queryHandler)
	http.HandleFunc("/save", a.saveHandler) // временная мера

	if err := http.ListenAndServe(a.Address, nil); err != nil {
		return err
	}
	return nil

}

// прием запроса, валидация
func (a *Api) queryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		query, e := NewApiQueryReqFromBody(r)
		if e != nil {
			ApiErrorMsg(w, e.Error())
			return
		}
		a.Query(query, w)
	case "GET":
		query, e := NewApiQueryReqFromUrl(r)
		if e != nil {
			ApiErrorMsg(w, e.Error())
			return
		}
		a.Query(query, w)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func (a *Api) saveHandler(w http.ResponseWriter, r *http.Request) {
	a.Store.Save()
}

func (a *Api) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *Api) versionHandler(w http.ResponseWriter, r *http.Request) {
	ApiVersionResponse(w, API_VERSION)
}
