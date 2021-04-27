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
	jrp := newJobReqPostedController()
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

	//Job Req Posting Controller
	http.Handle("/jobrequisition/posted", *jrp)
	http.Handle("/jobrequisition/posted/", *jrp)

	//Application Controller
	http.Handle("/application", a)
	http.Handle("/application/", a)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

