package models

import (
	"fmt"
)

type Country struct {
	ID   int
	Name string
	Code string
}

var (
	countries     = make(map[int]*Country)
	nextCountryID = 1
)

func GetCountries() map[int]*Country {
	return countries
}

func GetCountryByID(id int) (Country, error) {
	if _, found := countries[id]; found {
		return *countries[id], nil
	}

	return Country{}, fmt.Errorf("Country with ID '%v' not found", id)
}

func AddCountry(c Country) (Country, error) {
	if c.ID != 0 {
		return Country{}, fmt.Errorf("Country must not include ID")
	}

	if AlreadyExistByCode(c.Code) {
		return Country{}, fmt.Errorf("Country with CODE '%v' already exists", c.Code)
	}

	c.ID = nextCountryID
	countries[nextCountryID] = &c
	nextCountryID++

	return c, nil
}

func UpdateCountry(c Country) (Country, error) {
	if _, found := countries[c.ID]; found {
		countries[c.ID] = &c
		return c, nil
	}

	//Update Candidate with the new value of country
	for _, v := range GetCandidatesWithCountry(c.ID) {
		v.CountryObj = c
		UpdateCandidate(v)
	}

	//Update Job Requisition with new value of country
	for _, v := range GetRequisitionsWithCountry(c.ID) {
		v.JobReqCountry = c
		UpdateJobRequisition(v)
	}

	return Country{}, fmt.Errorf("Country to be updated not found")
}

func RemoveCountryByID(id int) error {
	if _, found := countries[id]; found {
		delete(countries, id)
		return nil
	}

	return fmt.Errorf("Country with ID '%v' not found", id)
}

//Validate if the country with given ID already exists on the list
//Returns true when country exist, and false when it doesn't
func AlreadyExistById(id int) bool {
	_, found := countries[id]
	return found
}

//Validate if the country already exists on the list
//Returns true when the country exists and false when it doesn't
func AlreadyExistByCode(code string) bool {
	for _, c := range countries {
		if code == c.Code {
			return true
		}
	}
	return false
}
