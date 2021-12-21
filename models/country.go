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
	countries = make(map[int]*Country)
	nextCountryID = updateCountriesInMemory()
)

//In Memory: Returns the complete list of countries that has been.
//Returns a hashmap containing the list of countries
func GetCountries() map[int]*Country {
	return countries
}

//In Memory: Searches for a specific country on the hashmap.
//Returns a country object and an error in case it was not possible to find the record
func GetCountryByID(id int) (Country, error) {
	//Implementation consulting the databasse.
	//This was replaced by persisting the list in memory and then updating the list when values on the database are updated.
	//Possible only because all insert and updates done on the database pass go this class
	/*client, err := db.OpenConnectionToMongo()

	if err != nil {
		return Country{}, nil
	}

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"ID", id}},
			},
		},
	}
	projection := bson.D{{"ID",1},{"Name", 1},{"Code", 1}}
	opts := options.FindOne().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Countries")

	var result bson.D
	coll.FindOne(context.TODO(),filter,opts).Decode(&result)

	c := bsonToCountry(result)

	if c.Name != "" {
		return c, nil
	}


	*/

	if _, found := countries[id]; found {
		return *countries[id], nil
	}

	return Country{}, fmt.Errorf("Country with ID '%v' not found", id)
}

//In DB: Creates a new country record to the collection and updates the countries in memory.
//Returns a country object and an error in case it was not possible to create the record
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

	//Validate if able to connect to MongoDB
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return Country{}, fmt.Errorf("Could not establish connection to MongoDB")
	}

	coll := client.Database("myFirstDatabase").Collection("Countries")
	doc := bson.D{{"ID", c.ID}, {"Name", c.Name}, {"Code", c.Code}}

	//Insert information into MongoDB
	_ , err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return Country{}, fmt.Errorf("Could not insert Country provided")
	}

	defer db.CloseConnectionToMongo(client)

	updateCountriesInMemory()

	return c, nil
}

//In DB: Updates a country record on the collection and updates the countries in memory.
//Returns a country object and an error in case it was not possible to update the record
func UpdateCountry(c Country) (Country, error) {
	if _, found := countries[c.ID]; found {
		countries[c.ID] = &c

		//establish connection to database
		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return Country{}, fmt.Errorf("Could not establish connection to database")
		}

		//create parameters for updating the values of the country
		coll := client.Database(db.GetDatabaseName()).Collection("Countries")
		filter := bson.D{{"ID", c.ID}}
		update := bson.D{{"$set",bson.D{{"Name", c.Name},{"Code", c.Code}}}}

		//execute update on the database record
		_, err = coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return Country{}, fmt.Errorf("Could not update country on the database")
		}

		defer db.CloseConnectionToMongo(client)

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

		//Update the list stored in memory
		updateCountriesInMemory()
		return c, nil
	}

	return Country{}, fmt.Errorf("Country to be updated not found")
}

//In DB: Removes a country record from the collection and updates the countries in memory.
//Returns error if failed to complete the deletion on the DB
func RemoveCountryByID(id int) error {
	if _, found := countries[id]; found {
		delete(countries, id)

		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return fmt.Errorf("Could not establish connection to database")
		}

		//create parameters for updating the values of the country
		coll := client.Database(db.GetDatabaseName()).Collection("Countries")
		filter := bson.D{{"ID", id}}

		//Execute deletion on the database
		_, err = coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			return fmt.Errorf("Could not delete Country")
		}

		defer db.CloseConnectionToMongo(client)

		//update the list stored in memory
		updateCountriesInMemory()
		return nil
	}

	return fmt.Errorf("Country with ID '%v' not found", id)
}

//In Memory:Validate if the country with given ID already exists on the list.
//Returns true when country exist, and false when it doesn't
func AlreadyExistById(id int) bool {
	_, err := GetCountryByID(id)

	if err != nil {
		return false
	}

	return true
}

//Validate if the country already exists on the list.
//Returns true when the country exists and false when it doesn't
func AlreadyExistByCode(code string) bool {
	for _, c := range countries {
		if code == c.Code {
			return true
		}
	}
	return false
}

//Updates the hashmap containing all the countries to work with them in memory.
//Return the next ID to be added into the Database
func updateCountriesInMemory() int {
	client, err := db.OpenConnectionToMongo()

	if err != nil {
		return -1
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

	biggestId := 1
	for _, v := range results {
		c := bsonToCountry(v)
		//add object into hashmap
		countries[c.ID] = &c
		if c.ID > biggestId{
			biggestId = c.ID
		}
	}

	return biggestId+1
}

//Receives a bson object to execute the conversion.
//Returns a Country object.
func bsonToCountry(v bson.D) Country {
	bsonBytes, _ := bson.Marshal(v)

	var c Country
	//deconvert the byarray into a struct object
	bson.Unmarshal(bsonBytes, &c)

	return c
}