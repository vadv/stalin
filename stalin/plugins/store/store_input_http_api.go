package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	ql "stalin/plugins/store/query_language"
	parser "stalin/plugins/store/query_language/parser"
)

type qlResponse struct {
	Problems  []*ql.Problem `json:"result"`
	Count     int           `json:"count"`
	SpendTime float64       `json:"time"`
	Query     *qlReq        `json:"query"`
}

type qlReq struct {
	Query  string `json:"query"`
	Pretty bool   `json:"pretty"`
}

func writeError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (s *Storage) Query(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if s.Store == nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	tBegin := time.Now()

	query := fmt.Sprintf("select * from problem where %s", string(body))
	parser.New(query)
	ast, err := parser.AST()
	if err != nil {
		writeError(err, w)
		return
	}

	sel := ast.Child.(*parser.NodeSelect)
	where := sel.Where

	result, err := ql.Walk(where, s.Store.Items.List())
	if err != nil {
		writeError(err, w)
		return
	}

	response := &qlResponse{}
	response.Problems = result.List()
	response.Count = result.Size()
	response.SpendTime = time.Now().Sub(tBegin).Seconds() * 1000
	response.Query = &qlReq{Query: query}

	var rBytes []byte
	rBytes, err = json.Marshal(response)
	if err != nil {
		writeError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rBytes)

	rBytes = nil
	result = nil
	response = nil

	s.Log.Info("Query: '%s' time: %v", query, time.Now().Sub(tBegin))

}
