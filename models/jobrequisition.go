package models

import (
	"fmt"
)

type JobRequisition struct {
	ID             int
	Title          string
	JobDescription string
	PostingStatus  bool
	JobReqCountry  Country
	Applicants     []Application
}

var (
	jobReqs   []*JobRequisition
	nextJobID = 1
)

func GetJobRequisitions() []*JobRequisition {
	return jobReqs
}

func GetJobRequisitionByID(id int) (JobRequisition, error) {
	for _, jr := range jobReqs {
		if id == jr.ID {
			return *jr, nil
		}
	}
	return JobRequisition{}, fmt.Errorf("Job Requisition with ID '%v' not found", id)
}

func GetJobRequisitionPosted() []*JobRequisition{
	var postedReqs []*JobRequisition
	for _, jr := range jobReqs {
		if jr.PostingStatus {
			postedReqs = append(postedReqs, jr)
		}
	}
	return postedReqs
}

func AddJobRequisition(jr JobRequisition) (JobRequisition, error) {
	//Validation section
	if jr.ID != 0 {
		return JobRequisition{}, fmt.Errorf("Job Requisition must not contain ID upon creation")
	}

	if !AlreadyExistById(jr.JobReqCountry.ID) {
		return JobRequisition{}, fmt.Errorf("Country inserted for Job Requisition does not exist")
	} else if jr.JobReqCountry.ID == 0 && (jr.JobReqCountry.Name != "" || jr.JobReqCountry.Code != "") {
		jr.JobReqCountry.Name = ""
		jr.JobReqCountry.Code = ""
	}

	if jr.Title == "" || jr.JobDescription == "" {
		return JobRequisition{}, fmt.Errorf("Mandatory fields should be populated upon creating Job Requisition")
	}

	//Create Job Requisition
	jr.ID = nextJobID
	nextJobID++
	jobReqs = append(jobReqs, &jr)
	return jr, nil
}

func UpdateJobRequisition(jr JobRequisition) (JobRequisition, error) {
	//Validation section
	if !AlreadyExistById(jr.JobReqCountry.ID) {
		return JobRequisition{}, fmt.Errorf("Country inserted for Job Requisition does not exist")
	} else if jr.JobReqCountry.ID == 0 && (jr.JobReqCountry.Name != "" || jr.JobReqCountry.Code != "") {
		jr.JobReqCountry.Name = ""
		jr.JobReqCountry.Code = ""
	}

	if jr.Title == "" || jr.JobDescription == "" {
		return JobRequisition{}, fmt.Errorf("All mandatory fields should be populated")
	}

	//Update Job Requisition
	for i, j := range jobReqs {
		if j.ID == jr.ID {
			jobReqs[i] = &jr
			return jr, nil
		}
	}

	//Return Job Req not found
	return JobRequisition{}, fmt.Errorf("Job Requisition with ID '%v' not found", jr.ID)
}

func AddApplicationToJobReq(a Application) error {
	for i, j := range jobReqs {
		if a.JobRequisitionID == j.ID {
			jobReqs[i].Applicants = append(jobReqs[i].Applicants, a)
			return nil
		}
	}
	return fmt.Errorf("Cannot add Application to Job Requisition")
}

func RemoveApplicationFromJobReq(a Application) error {
	for i, j := range jobReqs {
		if a.JobRequisitionID == j.ID {
			for k, apps := range j.Applicants {
				if apps.ID == a.ID {
					jobReqs[i].Applicants = append(jobReqs[i].Applicants[:k], jobReqs[i].Applicants[k+1:]...)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("Cound not remove Application ID '%v' from Job Requisition '%v'", a.ID, a.JobRequisitionID)
}

func DeleteJobRequisition(id int) error {
	for i, j := range jobReqs {
		if id == j.ID {
			for _, app := range j.Applicants {
				RemoveApplicationFromCandidate(app)
			}
			jobReqs = append(jobReqs[:i], jobReqs[i:]...)
			return nil
		}
	}
	return fmt.Errorf("Job Requisition with ID '%v' not found", id)
}

func IsJobReqPosted(id int) (bool, error) {
	jr, err := GetJobRequisitionByID(id)
	if err != nil {
		return false, fmt.Errorf("Counld not find Job Requisition '%v'", id)
	}

	return jr.PostingStatus, nil
}