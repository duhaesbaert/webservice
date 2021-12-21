package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webservice/db"
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
	candidates = make(map[int]*Candidate)
	nextCanID  = updateCandidatesInMemory()
)

//In Memory: Returns the complete list of candidates that has been.
//Returns a hashmap containing the list of candidates
func GetCandidates() []*Candidate {
	candArr := make([]*Candidate, 0)

	for _, v := range candidates {
		candArr = append(candArr, v)
	}
	return candArr
	//return candidates
}

//In Memory: Searches for a specific candidate on the hashmap.
//Returns a candidate object and an error in case it was not possible to find the record
func GetCandidateByID(id int) (Candidate, error) {
	if c, found := candidates[id]; found {
		return *c, nil
	}
	return Candidate{}, fmt.Errorf("Candidate with ID '%v' not found", id)
}

//In DB: Creates a new country record to the collection and updates the candidates in memory.
//Returns a country object and an error in case it was not possible to create the record
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

	//Check if tags are part of the candidate creation
	if c.Tags != nil {
		//Validate if tag exists to Add/Reuse
		c.Tags = ValidateTags(c.Tags)
	}

	//Add new Candidate
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return Candidate{}, fmt.Errorf("Could not establish connection to Database")
	}

	coll := client.Database(db.GetDatabaseName()).Collection("Candidates")
	doc := bson.D{
		{"ID", nextCanID},
		{"FirstName", c.FirstName},
		{"LastName", c.LastName},
		{"Email", c.Email},
		{"Address", c.Address},
		{"Tags", c.Tags},
		{"CountryObj", c.CountryObj},
		{"JobsApplied", c.JobsApplied}}

	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return Candidate{}, fmt.Errorf("Could not insert Candidate provided")
	}

	defer db.CloseConnectionToMongo(client)

	c.ID = nextCanID
	updateCandidatesInMemory()

	return c, nil
}

//In DB: Updates a candidate record on the collection and updates the candidates in memory.
//Returns a candidate object and an error in case it was not possible to update the record
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
	if _, found := candidates[c.ID]; found {
		if c.Tags != nil {
			//Validate if tag exist to add/reuse
			c.Tags = ValidateTags(c.Tags)
		}
		//Removes the possibility of editing the JobsApplied when updating candidate
		c.JobsApplied = candidates[c.ID].JobsApplied

		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return Candidate{}, fmt.Errorf("Could not establish connection to Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Candidates")
		filter := bson.D{{"ID", c.ID}}
		update := bson.D{{"$set", bson.D{
			{"FirstName", c.FirstName},
			{"LastName", c.LastName},
			{"Email", c.Email},
			{"Address", c.Address},
			{"Tags", c.Tags},
			{"CountryObj", c.CountryObj},
			{"JobsApplied", c.JobsApplied}}}}

		if _, err = coll.UpdateOne(context.TODO(), filter, update); err != nil {
			return Candidate{}, fmt.Errorf("Could not update candidate provided")
		}

		defer db.CloseConnectionToMongo(client)

		updateCandidatesInMemory()
		return c, nil
	}

	//Return candidate not found
	return Candidate{}, fmt.Errorf("Candidate '%v' was not found", c.FirstName)
}

//In DB: Removes a candidate record from the collection and updates the candidates in memory.
//Returns error if failed to complete the deletion on the DB
func DeleteCandidate(id int) error {
	if _, found := candidates[id]; found {
		for _, app := range candidates[id].JobsApplied {
			DeleteApplication(app.ID)
		}

		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return fmt.Errorf("Could not establish connection to Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Candidates")
		filter := bson.D{{"ID", id}}

		if _, err = coll.DeleteOne(context.TODO(), filter); err != nil {
			return fmt.Errorf("Could not delete Candidate")
		}

		defer db.CloseConnectionToMongo(client)

		updateCandidatesInMemory()
		return nil
	}

	return fmt.Errorf("Candidate with id '%v' not found", id)
}

func AddApplicationToCandidate(a Application) error {
	if _, found := candidates[a.CandidateProfileID]; found {
		candidates[a.CandidateProfileID].JobsApplied = append(candidates[a.CandidateProfileID].JobsApplied, a)
		return nil
	}

	return fmt.Errorf("Cannot add Application to this candidate")
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
	if _, found := candidates[a.CandidateProfileID]; found {
		for i, app := range candidates[a.CandidateProfileID].JobsApplied {
			if app.ID == a.ID {
				candidates[a.CandidateProfileID].JobsApplied = append(candidates[a.CandidateProfileID].JobsApplied[:i], candidates[a.CandidateProfileID].JobsApplied[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("Cannot remove Application from Candidate")
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

func GetCandidatesWithCountry(c int) []Candidate{
	ret := make([]Candidate,0)

	for _, v := range GetCandidates() {
		if v.CountryObj.ID == c {
			ret = append(ret, *v)
		}
	}

	return ret
}

func UpdateApplicationOnCandidates(a Application) {
	for _, v := range GetCandidates() {
		for _, vApp := range v.JobsApplied {
			if vApp.ID == a.ID {
				vApp.SalaryExpectation = a.SalaryExpectation
				vApp.timeOfExperience = a.timeOfExperience
				vApp.ApplicationSource = a.ApplicationSource
			}
		}
	}
}

//Updates the hashmap containing all the countries to work with them in memory.
//Return the next ID to be added into the Database
func updateCandidatesInMemory() int {
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return -1
	}

	filter := bson.D{}
	projection := bson.D{
		{"ID", 1},
		{"FirstName", 1},
		{"LastName", 1},
		{"Email", 1},
		{"Address", 1},
		{"Tags", 1},
		{"CountryObj", 1},
		{"JobsApplied", 1}}
	opts := options.Find().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Candidates")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	defer db.CloseConnectionToMongo(client)

	biggestId := 1
	candidates = make(map[int]*Candidate)
	for _, v := range results {
		c := bsonToCandidate(v)

		candidates[c.ID] = &c
		if c.ID > biggestId {
			biggestId = c.ID
		}
	}

	return biggestId+1
}

//Receives a bson object to execute the conversion.
//Returns a Candidate object.
func bsonToCandidate(v bson.D) Candidate {
	bsonBytes, _ := bson.Marshal(v)

	var c Candidate
	//deconvert the byarray into a struct object
	bson.Unmarshal(bsonBytes, &c)

	return c
}