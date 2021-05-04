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
	applications = make(map[int]*Application)
	nextAppID    = 1
)

func GetApplications() map[int]*Application {
	return applications
}

func GetApplicationByID(id int) (Application, error) {
	if a, found := applications[id]; found {
		return *a, nil
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

	//Job Requisition has not been found
	if err != nil {
		return Application{}, fmt.Errorf("Error getting Job Requisition posting")
	}

	//Candidates cannot apply to requisitions not posted
	if !tf {
		return Application{}, fmt.Errorf("Job Requisition is not posted")
	}

	a.ID = nextAppID

	//Insert record from Applicant on Job Requisition object
	err = AddApplicationToJobReq(a)
	if err != nil {
		return Application{}, err
	}

	//Insert applicant record in Candidate object
	err = AddApplicationToCandidate(a)
	if err != nil {
		RemoveApplicationFromJobReq(a)
		return Application{}, err
	}

	//Add application into application list
	applications[nextAppID] = &a
	nextAppID++
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
	if app, found := applications[id]; found {
		//Remove application from Job Requisition
		err := RemoveApplicationFromJobReq(*app)
		if err != nil {
			return fmt.Errorf("Could not remove Application from Job Requisition")
		}

		//Remove application from Candidate
		err = RemoveApplicationFromCandidate(*app)
		if err != nil {
			AddApplicationToJobReq(*app)
			return fmt.Errorf("Could not remove Application from Candidate")
		}

		delete(applications, id)
		return nil
	}
	return fmt.Errorf("Application with ID '%v' not found", id)
}

func RemoveCandidateApplicationsByID(canID int) error {
	for i, app := range applications {
		if app.CandidateProfileID == canID {
			//Remove the application from Job Requisition
			err := RemoveApplicationFromJobReq(*app)
			if err != nil {
				return fmt.Errorf("Could not remove applications from Job Requisitions")
			}

			delete(applications, i)
		}
	}
	return nil
}
