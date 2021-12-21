package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webservice/db"

	"context"
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

func GetCountries() []Country {
	client, err := db.OpenConnectionToMongo()

	if err != nil {
		return []Country{}
	}

	filter := bson.D{}
	projection := bson.D{{"ID",1},{"Name", 1},{"Code", 1}}
	opts := options.Find().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Countries")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	defer db.CloseConnectionToMongo(client)

	var listCountries []Country
	for _, v := range results {
		//turn the bson object into a byte array
		bsonBytes, _ := bson.Marshal(v)

		var c Country
		//deconvert the byarray into a struct object
		bson.Unmarshal(bsonBytes, &c)

		//add object into hashmap
		listCountries = append(listCountries, c)
		//countries[i] = &c
	}

	return listCountries
}

func GetCountryByID(id int) (Country, error) {

	//Old Implementation with values persisted in memory
	//if _, found := countries[id]; found {
	//	return *countries[id], nil
	//}
	return Country{}, fmt.Errorf("Country with ID '%v' not found", id)
}

func AddCountry(c Country) (Country, error) {
	if c.ID != 0 {
		return Country{}, fmt.Errorf("Country must not include ID")
	}

	if AlreadyExistByCode(c.Code) {
		return Country{}, fmt.Errorf("Country with CODE '%v' already exists", c.Code)
	}

	//c.ID = nextCountryID
	//countries[nextCountryID] = &c
	//nextCountryID++

	//Validate if able to connect to MongoDB
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return Country{}, fmt.Errorf("Could not establish connection to MongoDB")
	}

	coll := client.Database("myFirstDatabase").Collection("Countries")
	doc := bson.D{{"ID", c.ID}, {"Name", c.Name}, {"Code", c.Code}}

	//Validate if able to insert information into MongoDB
	_ , err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return Country{}, fmt.Errorf("Could not insert Country provided")
	}

	defer db.CloseConnectionToMongo(client)

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
