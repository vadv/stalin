package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	. "stalin/plugins/store/problem"
	"strconv"
)

type ApiQueryReq struct {
	Domain        *string `json:"domain"`
	Fqdn          *string `json:"fqdn"`
	fqdnRegexp    *regexp.Regexp
	Host          *string  `json:"host"`
	State         []string `json:"state"`
	Tags          []string `json:"tags"`
	Description   *string  `json:"description"`
	descRegexp    *regexp.Regexp
	Service       *string `json:"service"`
	serviceRegexp *regexp.Regexp
	Pretty        bool `json:"pretty"`
}

func (a *ApiQueryReq) empty() bool {
	if a.Domain == nil && a.Fqdn == nil && a.Service == nil &&
		a.Description == nil && len(a.State) == 0 && len(a.Tags) == 0 {
		return true
	}
	return false
}

func (a *ApiQueryReq) setRegexp() error {

	var myRegexp *regexp.Regexp
	var err error

	if a.Service != nil {
		if myRegexp, err = regexp.CompilePOSIX(*a.Service); err != nil {
			return err
		}
		a.serviceRegexp = myRegexp
	}

	if a.Fqdn != nil {
		if myRegexp, err = regexp.CompilePOSIX(*a.Fqdn); err != nil {
			return err
		}
		a.fqdnRegexp = myRegexp
	}

	if a.Description != nil {
		if myRegexp, err = regexp.CompilePOSIX(*a.Description); err != nil {
			return err
		}
		a.descRegexp = myRegexp
	}

	return nil

}

func missParams(r *http.Request) []string {
	params := r.URL.Query()
	all_keys := []string{}
	for key, _ := range params {
		if key != "domain" && key != "fqdn" &&
			key != "host" && key != "state" &&
			key != "tags" && key != "service" &&
			key != "pretty" && key != "description" {
			all_keys = append(all_keys, key)
		}
	}
	return all_keys
}

// todo: miss params
func NewApiQueryReqFromBody(r *http.Request) (*ApiQueryReq, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	query := &ApiQueryReq{}
	if err := json.Unmarshal(body, query); err != nil {
		return nil, err
	}
	if err := query.setRegexp(); err != nil {
		return nil, err
	}
	return query, nil
}

func NewApiQueryReqFromUrl(r *http.Request) (*ApiQueryReq, error) {

	var err error

	if len(missParams(r)) > 0 {
		// не частый use-case считаем еще раз
		return nil, fmt.Errorf("Unknown params: %#v", missParams(r))
	}

	params := r.URL.Query()
	query := &ApiQueryReq{}

	if len(params["domain"]) > 0 {
		if len(params["domain"]) > 1 {
			return nil, fmt.Errorf("multiple 'domain' is not supported")
		}
		query.Domain = &params["domain"][0]
	}
	if len(params["host"]) > 0 {
		if len(params["host"]) > 1 {
			return nil, fmt.Errorf("multiple 'host' is not supported")
		}
		query.Host = &params["host"][0]
	}
	if len(params["fqdn"]) > 0 {
		if len(params["fqdn"]) > 1 {
			return nil, fmt.Errorf("multiple 'fqdn' is not supported")
		}
		query.Fqdn = &params["fqdn"][0]
	}
	if len(params["description"]) > 0 {
		if len(params["description"]) > 1 {
			return nil, fmt.Errorf("multiple 'description' is not supported")
		}
		query.Description = &params["description"][0]
	}
	if len(params["service"]) > 0 {
		if len(params["service"]) > 1 {
			return nil, fmt.Errorf("multiple 'service' is not supported")
		}
		query.Service = &params["service"][0]
	}
	if len(params["tags"]) > 0 {
		query.Tags = params["tags"]
	}
	if len(params["state"]) > 0 {
		query.State = params["state"]
	}
	if len(params["pretty"]) > 0 {
		query.Pretty, err = strconv.ParseBool(params["pretty"][0])
		if err != nil {
			return nil, fmt.Errorf("Error convert pretty: %v", err)
		}
	}
	if query.empty() {
		return nil, fmt.Errorf("Query is empty")
	}
	if err := query.setRegexp(); err != nil {
		return nil, err
	}
	return query, nil
}

func (q *ApiQueryReq) matchedProblem(problem *Problem) bool {

	// если стейт один из указанных
	if len(q.State) > 0 {
		if !stringInSlice(problem.State, q.State) {
			return false
		}
	}

	// фильтруем domain (четкий поиск)
	if q.Domain != nil && *q.Domain != problem.Domain {
		return false
	}

	// фильтруем host (четкий поиск)
	if q.Host != nil && *q.Host != problem.Hostname {
		return false
	}

	// фильтруем fqdn
	if q.Fqdn != nil {
		if !q.fqdnRegexp.MatchString(problem.Fqdn) {
			return false
		}
	}

	// сервис
	if q.Service != nil {
		if !q.serviceRegexp.MatchString(problem.Service) {
			return false
		}
	}

	// дескрипшен
	if q.Description != nil {
		if !q.descRegexp.MatchString(problem.Description) {
			return false
		}
	}

	// если попадает хоть один таг
	if len(q.Tags) > 0 {
		for _, tag := range q.Tags {
			if stringInSlice(tag, problem.Tags) {
				return true
			}
		}
		// таг указан, но не попал
		return false
	}

	return true
}

func (q *ApiQueryReq) ToJson() string {
	b, _ := json.Marshal(q)
	return string(b)
}
