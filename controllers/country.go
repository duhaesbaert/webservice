package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"webservice/models"
)

type countryController struct {
	countryIDPattern *regexp.Regexp
}

func newCountryController() *countryController {
	return &countryController{
		countryIDPattern: regexp.MustCompile(`^/country/(\d+)/?`),
	}
}

func (cntC countryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/country" {
		switch r.Method {
		case http.MethodGet:
			cntC.getAll(w, r)
		case http.MethodPost:
			cntC.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := cntC.countryIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		switch r.Method {
		case http.MethodGet:
			cntC.get(id, w)
		case http.MethodPut:
			cntC.put(id, w, r)
		case http.MethodDelete:
			cntC.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (cntC countryController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetCountries(), w)
}

func (cntC countryController) get(id int, w http.ResponseWriter) {
	c, err := models.GetCountryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(c, w)
}

func (cntC countryController) post(w http.ResponseWriter, r *http.Request) {
	c, err := cntC.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Country object"))
		return
	}

	c, err = models.AddCountry(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(c, w)
}

func (cntC countryController) put(id int, w http.ResponseWriter, r *http.Request) {
	c, err := cntC.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse Country object"))
		return
	}

	if id != c.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in URL"))
		return
	}

	c, err = models.UpdateCountry(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(c, w)
}

func (cntC countryController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveCountryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (cntC countryController) parseRequest(r *http.Request) (models.Country, error) {
	dec := json.NewDecoder(r.Body)
	var c models.Country
	err := dec.Decode(&c)
	if err != nil {
		return models.Country{}, err
	}
	return c, nil
}
