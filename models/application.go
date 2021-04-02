package models

import (
	"fmt"
)

type Application struct {
	ID                 int
	CandidateProfileID int
	JobRequisitionID   int
}

var (
	applications []*Application
	nextAppID    = 1
)

func GetApplications() []*Application {
	return applications
}

func GetApplicationByID(id int) (Application, error) {
	for _, a := range applications {
		if id == a.ID {
			return *a, nil
		}
	}
	return Application{}, fmt.Errorf("Application with id '%v' not found", id)
}

func AddApplication(a Application) (Application, error) {
	if a.ID != 0 {
		return Application{}, fmt.Errorf("Application must not contain ID upon creation")
	}

	if a.CandidateProfileID == 0 || a.JobRequisitionID == 0 {
		return Application{}, fmt.Errorf("CandidateProfileID and JobRequisitionID are mandatory for submitting application")
	}

	tf, err := IsJobReqPosted(a.JobRequisitionID)
	if err == nil {
		if !tf {
			return Application{}, fmt.Errorf("Job Requisition is not posted")
		}
	} else {
		return Application{}, fmt.Errorf("Error getting Job Requisition posting")
	}

	a.ID = nextAppID

	err = AddApplicationToJobReq(a)
	if err != nil {
		return Application{}, err
	}

	err = AddApplicationToCandidate(a)
	if err != nil {
		RemoveApplicationFromJobReq(a)
		return Application{}, err
	}

	nextAppID++
	applications = append(applications, &a)
	return a, nil
}

/*
func UpdateApplication(a Application) (Application, error) {
	if a.CandidateProfileID != 0 && a.JobRequisitionID != 0 {
		for i, app := range applications {
			if a.ID == app.ID {
				applications[i] = &a
				return a, nil
			}
		}
		return Application{}, fmt.Errorf("Application with ID '%v' not found", a.ID)
	}
	return Application{}, fmt.Errorf("Missing Job Requisition and/or Candidate")
}
*/

func DeleteApplication(id int) error {
	for i, app := range applications {
		if app.ID == id {
			err := RemoveApplicationFromJobReq(*app)
			if err != nil {
				return fmt.Errorf("Could not remove Application from Job Requisition")
			}

			err = RemoveApplicationFromCandidate(*app)
			if err != nil {
				AddApplicationToJobReq(*app)
				return fmt.Errorf("Could not remove Application from Candidate")
			}
			applications = append(applications[:i], applications[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Application with ID '%v' not found", id)
}
