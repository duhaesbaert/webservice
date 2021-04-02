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
	countries     []*Country
	nextCountryID = 1
)

func GetCountries() []*Country {
	return countries
}

func GetCountryByID(id int) (Country, error) {
	for _, c := range countries {
		if id == c.ID {
			return *c, nil
		}
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
	nextCountryID++
	countries = append(countries, &c)
	return c, nil
}

func UpdateCountry(c Country) (Country, error) {
	for i, cnt := range countries {
		if c.ID == cnt.ID {
			countries[i] = &c
			return c, nil
		}
	}
	return Country{}, fmt.Errorf("Country to be updated not found")
}

func RemoveCountryByID(id int) error {
	for i, c := range countries {
		if id == c.ID {
			countries = append(countries[:i], countries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Country with ID '%v' not found", id)
}

//Validate if the country with given ID already exists on the list
//Returns true when country exist, and false when it doesn't
func AlreadyExistById(id int) bool {
	if id == 0 {
		return true
	}

	for _, c := range countries {
		if id == c.ID {
			return true
		}
	}

	return false
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
