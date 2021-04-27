package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"webservice/models"
)

type jobReqPosted struct {
	jobReqIDPattern *regexp.Regexp
}

func newJobReqPostedController() *jobReqPosted {
	return &jobReqPosted{
		jobReqIDPattern: regexp.MustCompile(`^/jobrequisition/posted/(\d+)/?`),
	}
}

func (jr jobReqPosted) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/jobrequisition/posted" {
		switch r.Method {
		case http.MethodGet:
			jr.getPosted(w, r)
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
			jr.getIfPosted(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (jr jobReqPosted) getPosted(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetJobRequisitionPosted(), w)
}

func (jr jobReqPosted) getIfPosted(id int, w http.ResponseWriter) {
	j, err := models.IsJobReqPosted(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(j, w)
}