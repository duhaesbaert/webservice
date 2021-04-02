package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"webservice/models"
)

type candidateController struct {
	candidateIDPattern *regexp.Regexp
}

func newCandidateController() *candidateController {
	return &candidateController{
		candidateIDPattern: regexp.MustCompile(`^/candidate/(\d+)/?`),
	}
}

func (c candidateController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/candidate" {
		switch r.Method {
		case http.MethodGet:
			c.getAll(w, r)
		case http.MethodPost:
			c.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := c.candidateIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			c.get(id, w)
		case http.MethodPut:
			c.put(id, w, r)
		case http.MethodDelete:
			c.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (c candidateController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetCandidates(), w)
}

func (c candidateController) get(id int, w http.ResponseWriter) {
	can, err := models.GetCandidateByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(can, w)
}

func (c candidateController) post(w http.ResponseWriter, r *http.Request) {
	can, err := c.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Candidate object"))
		return
	}
	can, err = models.AddCandidate(can)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(can, w)
}

func (c candidateController) put(id int, w http.ResponseWriter, r *http.Request) {
	can, err := c.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Candidate object"))
		return
	}

	if id != can.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in URL"))
		return
	}

	can, err = models.UpdateCandidate(can)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(can, w)
}

func (c candidateController) delete(id int, w http.ResponseWriter) {
	err := models.DeleteCandidate(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c candidateController) parseRequest(r *http.Request) (models.Candidate, error) {
	dec := json.NewDecoder(r.Body)
	var can models.Candidate
	err := dec.Decode(&can)
	if err != nil {
		return models.Candidate{}, err
	}
	return can, nil
}
