package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"webservice/models"
)

type applicationController struct {
	applicationIDPattern *regexp.Regexp
}

func newApplicationController() *applicationController {
	return &applicationController{
		applicationIDPattern: regexp.MustCompile(`^/application/(\d+)/?`),
	}
}

func (a applicationController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/application" {
		switch r.Method {
		case http.MethodGet:
			a.getAll(w, r)
		case http.MethodPost:
			a.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := a.applicationIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			a.get(id, w)
		case http.MethodPut:
			a.put(id, w, r)
		case http.MethodDelete:
			a.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (a applicationController) parseRequest(r *http.Request) (models.Application, error) {
	dec := json.NewDecoder(r.Body)
	var app models.Application
	err := dec.Decode(&app)
	if err != nil {
		return models.Application{}, err
	}
	return app, nil
}

func (a applicationController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetApplications(), w)
}

func (a applicationController) get(id int, w http.ResponseWriter) {
	app, err := models.GetApplicationByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(app, w)
}

func (a applicationController) post(w http.ResponseWriter, r *http.Request) {
	app, err := a.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse application object"))
		return
	}

	app, err = models.AddApplication(app)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJSON(app, w)
}

func (a applicationController) put(id int, w http.ResponseWriter, r *http.Request) {
	app, err := a.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse application object"))
		return
	}

	if id != app.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted applicant must match ID in the URL"))
		return
	}

	app, err = models.UpdateApplication(app)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(app, w)
	//w.WriteHeader(http.StatusNotImplemented)
}

func (a applicationController) delete(id int, w http.ResponseWriter) {
	err := models.DeleteApplication(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
