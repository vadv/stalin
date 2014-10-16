package api

import (
	"encoding/json"
	"net/http"
	. "stalin/plugins"
	. "stalin/plugins/store/problem"
	"time"
)

type apiQueryResponse struct {
	Problems  []*Problem   `json:"result"`
	Count     int          `json:"count"`
	SpendTime float64      `json:"time"`
	Query     *ApiQueryReq `json:"query"`
}

func (a *Api) Query(q *ApiQueryReq, w http.ResponseWriter) (err error) {

	// корупт памяти
	var rBytes []byte

	tBeginReq := time.Now()

	result := make([]*Problem, 0)

	for _, item := range a.Store.Items.Items() {
		//todo паралельность
		if q.matchedProblem(item.ProblemObj) {
			result = append(result, item.ProblemObj)
		}
	}

	response := &apiQueryResponse{}
	response.Problems = result
	response.Count = len(result)
	response.SpendTime = time.Now().Sub(tBeginReq).Seconds() * 1000
	response.Query = q

	if q.Pretty {
		rBytes, err = json.MarshalIndent(response, "", "  ")
	} else {
		rBytes, err = json.Marshal(response)
	}

	if err != nil {
		return err
	}

	w.Write(rBytes)

	LogInfo("[MemStore]: Query %v time: %v ms", q.ToJson(), response.SpendTime)

	return nil

}
