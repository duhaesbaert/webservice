package models

import (
	"fmt"
)

type Candidate struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Address     string
	Tags		[]Tag
	CountryObj  Country
	JobsApplied []Application
}

var (
	candidates []*Candidate
	nextCanID  = 1
)

func GetCandidates() []*Candidate {
	return candidates
}

func GetCandidateByID(id int) (Candidate, error) {
	for _, c := range candidates {
		if id == c.ID {
			return *c, nil
		}
	}
	return Candidate{}, fmt.Errorf("Candidate with ID '%v' not found", id)
}

func AddCandidate(c Candidate) (Candidate, error) {
	//Validation
	if c.ID != 0 {
		return Candidate{}, fmt.Errorf("Candidate must not include ID")
	}

	if !AlreadyExistById(c.CountryObj.ID) {
		return Candidate{}, fmt.Errorf("Country inserted for candidate does not exist")
	} else if c.CountryObj.ID == 0 && (c.CountryObj.Name != "" || c.CountryObj.Code != "") {
		c.CountryObj.Name = ""
		c.CountryObj.Code = ""
	}

	if c.JobsApplied != nil {
		return Candidate{}, fmt.Errorf("A new candidate cannot have applied to jobs yet")
	}

	b, e := checkRequiredFields(c)
	if b {
		return Candidate{}, fmt.Errorf("Field: %v should be populated", e)
	}

	//Add new Candidate

	//Check if tags are part of the candidate creation
	if c.Tags != nil {
		//Validate if tag exists to Add/Reuse
		c.Tags = ValidateTags(c.Tags)
	}

	c.ID = nextCanID
	nextCanID++
	candidates = append(candidates, &c)
	return c, nil
}

func AddApplicationToCandidate(a Application) error {
	for i, c := range candidates {
		if a.CandidateProfileID == c.ID {
			candidates[i].JobsApplied = append(candidates[i].JobsApplied, a)
			return nil
		}
	}
	return fmt.Errorf("Cannot add Application to this candidate")
}

func UpdateCandidate(c Candidate) (Candidate, error) {
	//Validation section
	if !AlreadyExistById(c.CountryObj.ID) {
		return Candidate{}, fmt.Errorf("Country inserted for candidate does not exist")
	} else if c.CountryObj.ID == 0 && (c.CountryObj.Name != "" || c.CountryObj.Code != "") {
		c.CountryObj.Name = ""
		c.CountryObj.Code = ""
	}

	b, e := checkRequiredFields(c)
	if b {
		return Candidate{}, fmt.Errorf("Field: %v should be populated", e)
	}

	//Update Candidate
	for i, can := range candidates {
		if c.ID == can.ID {
			//Check if tags are on the update payload
			if c.Tags != nil {
				//Validate if tag exist to add/reuse
				c.Tags = ValidateTags(c.Tags)
			}
			//Removes the possibility of editing the JobsApplied when updating candidate
			c.JobsApplied = candidates[i].JobsApplied

			//Updates candidate
			candidates[i] = &c
			return c, nil
		}
	}

	//Return candidate not found
	return Candidate{}, fmt.Errorf("Candidate '%v' was not found", c.FirstName)
}

func checkRequiredFields(c Candidate) (bool, string) {
	retBool := false
	retString := ""

	if c.FirstName == "" {
		retString += "FirstName"
		retBool = true
	}

	if c.LastName == "" {
		if retBool {
			retString += ", "
		} else {
			retBool = true
		}
		retString += "LastName"
	}

	if c.Email == "" {
		if retBool {
			retString += ", "
		} else {
			retBool = true
		}

		retString += "Email"
	}

	return retBool, retString
}

func RemoveApplicationFromCandidate(a Application) error {
	for i, can := range candidates {
		if can.ID == a.CandidateProfileID {
			for j, app := range can.JobsApplied {
				if app.ID == a.ID {
					candidates[i].JobsApplied = append(candidates[i].JobsApplied[:j], candidates[i].JobsApplied[j+1:]...)
					return nil
				}
			}
		}
	}
	return fmt.Errorf("Cannot remove Application from Candidate")
}

func DeleteCandidate(id int) error {
	for i, can := range candidates {
		if can.ID == id {
			for _, app := range can.JobsApplied {
				//The delete application method will remove the application from Candidate and from JobRequisition (slower than specific method)
				DeleteApplication(app.ID)
			}
			candidates = append(candidates[:i], candidates[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Candidate with id '%v' not found", id)
}

func ValidateTags(cTags []Tag) ([]Tag) {
	for _, t := range cTags	{
		b,id,_ := ExistTagByLabel(t.Label)
		if !b {
			AddTag(t)
			continue
		}
		if b {
			t.ID = id
			continue
		}
	}
	return cTags
}