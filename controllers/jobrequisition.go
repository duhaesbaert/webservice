package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"webservice/models"
)

type jobRequisitionController struct {
	jobReqIDPattern *regexp.Regexp
}

func newJobRequisitionController() *jobRequisitionController {
	return &jobRequisitionController{
		jobReqIDPattern: regexp.MustCompile(`^/jobrequisition/(\d+)/?`),
	}
}

func (jr jobRequisitionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/jobrequisition" {
		switch r.Method {
		case http.MethodGet:
			jr.getAll(w, r)
		case http.MethodPost:
			jr.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := jr.jobReqIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			jr.get(id, w)
		case http.MethodPut:
			jr.put(id, w, r)
		case http.MethodDelete:
			jr.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (jr jobRequisitionController) parseRequest(r *http.Request) (models.JobRequisition, error) {
	dec := json.NewDecoder(r.Body)
	var can models.JobRequisition
	err := dec.Decode(&can)
	if err != nil {
		return models.JobRequisition{}, err
	}
	return can, nil
}

func (jr jobRequisitionController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetJobRequisitions(), w)
}

func (jr jobRequisitionController) get(id int, w http.ResponseWriter) {
	j, err := models.GetJobRequisitionByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(j, w)
}

func (jr jobRequisitionController) post(w http.ResponseWriter, r *http.Request) {
	j, err := jr.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing object Job Requisition"))
		return
	}

	j, err = models.AddJobRequisition(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(j, w)
}

func (jr jobRequisitionController) put(id int, w http.ResponseWriter, r *http.Request) {
	j, err := jr.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing object Job Requisition"))
		return
	}

	if id != j.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in URL"))
		return
	}

	j, err = models.UpdateJobRequisition(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(j, w)
}

func (jr jobRequisitionController) delete(id int, w http.ResponseWriter) {
	err := models.DeleteJobRequisition(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
