package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	c := newCandidateController()
	cntC := newCountryController()
	jr := newJobRequisitionController()
	a := newApplicationController()

	//Candidate controller
	http.Handle("/candidate", *c)
	http.Handle("/candidate/", *c)

	//Country controller
	http.Handle("/country", *cntC)
	http.Handle("/country/", *cntC)

	//Job Req Controller
	http.Handle("/jobrequisition", *jr)
	http.Handle("/jobrequisition/", *jr)

	http.Handle("/application", a)
	http.Handle("/application/", a)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
