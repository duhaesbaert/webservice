package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webservice/db"
)

type JobRequisition struct {
	ID				int
	Title			string
	JobDescription	string
	PostingStatus	bool
	JrCountryId		int
	JobReqCountry	Country
	Applicants		[]Application
}

var (
	jobReqs		=	make(map[int]*JobRequisition)
	nextJobID	=	updateJobRequisitionInMemory()
)

//In Memory: Returns the complete list of JobRequisition that has been.
//Returns a hashmap containing the list of JobRequisition
func GetJobRequisitions() []*JobRequisition {
	reqArr := make([]*JobRequisition,0)

	for _, v := range jobReqs {
		reqArr = append(reqArr, v)
	}
	return reqArr
	//return jobReqs
}

//In Memory: Searches for a specific JobRequisition on the hashmap.
//Returns a JobRequisition object and an error in case it was not possible to find the record
func GetJobRequisitionByID(id int) (JobRequisition, error) {
	if jr, found := jobReqs[id]; found {
		return *jr, nil
	}

	return JobRequisition{}, fmt.Errorf("Job Requisition with ID '%v' not found", id)
}

//In Memory: Searches for a specific JobRequisition on the hashmap that has posted = true
//Returns a JobRequisition object and an error in case it was not possible to find the record
func GetJobRequisitionPosted() []*JobRequisition{
	//use main method for retrieving job requisitions for updating country values.
	reqs := GetJobRequisitions()

	var postedReqs []*JobRequisition
	for _, jr := range reqs {
		if jr.PostingStatus {
			postedReqs = append(postedReqs, jr)
		}
	}
	return postedReqs
}

//In DB: Creates a new JobRequisition record to the collection and updates the JobRequisition in memory.
//Returns a JobRequisition object and an error in case it was not possible to create the record
func AddJobRequisition(jr JobRequisition) (JobRequisition, error) {
	//Validation section
	if jr.ID != 0 {
		return JobRequisition{}, fmt.Errorf("Job Requisition must not contain ID upon creation")
	}
	if jr.Title == "" || jr.JobDescription == "" {
		return JobRequisition{}, fmt.Errorf("Mandatory fields should be populated upon creating Job Requisition")
	}

	//Add New JobRequisition
	jr.ID = nextJobID

	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return JobRequisition{}, fmt.Errorf("Could not establish conneciton to Database")
	}

	coll := client.Database(db.GetDatabaseName()).Collection("Requisitions")
	doc := bson.D{
		{"ID", jr.ID},
		{"Title", jr.Title},
		{"JobDescription", jr.JobDescription},
		{"PostingStatus", jr.PostingStatus},
		{"JrCountryId", jr.JrCountryId}}

	if _, err = coll.InsertOne(context.TODO(), doc); err != nil {
		return JobRequisition{}, fmt.Errorf("Could not insert Job Requisition provided")
	}

	defer db.CloseConnectionToMongo(client)

	updateJobRequisitionInMemory()
	return GetJobRequisitionByID(jr.ID)
}

//In DB: Updates a JobRequisition record on the collection and updates the JobRequisition in memory.
//Returns a JobRequisition object and an error in case it was not possible to update the record
func UpdateJobRequisition(jr JobRequisition) (JobRequisition, error) {
	if jr.Title == "" || jr.JobDescription == "" {
		return JobRequisition{}, fmt.Errorf("Mandatory fields should be populated upon creating Job Requisition")
	}

	//Update Job Requisition
	if _, found := jobReqs[jr.ID]; found {

		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return JobRequisition{}, fmt.Errorf("Could not establish connection to the Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Requisitions")
		filter := bson.D{{"ID", jr.ID}}

		update := bson.D{{"$set", bson.D{
			{"Title", jr.Title},
			{"JobDescription", jr.JobDescription},
			{"PostingStatus", jr.PostingStatus},
			{"JrCountryId", jr.JrCountryId}}}}

		if _, err := coll.UpdateOne(context.TODO(), filter, update); err != nil {
			return JobRequisition{}, fmt.Errorf("Could not update  Requisition provided")
		}

		defer db.CloseConnectionToMongo(client)

		updateJobRequisitionInMemory()
		return GetJobRequisitionByID(jr.ID)
	}

	//Return Job Req not found
	return JobRequisition{}, fmt.Errorf("Job Requisition with ID '%v' not found", jr.ID)
}

//In DB: Removes a JobRequisition record from the collection and updates the JobRequisition in memory.
//Returns error if failed to complete the deletion on the DB
func DeleteJobRequisition(id int) error {
	if _, found := jobReqs[id]; found {
		DeleteApplicationFromJobReq(id)

		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return fmt.Errorf("Could not connect to Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Requisitions")
		filter := bson.D{{"ID", id}}

		if _, err = coll.DeleteOne(context.TODO(), filter); err != nil {
			return fmt.Errorf("Could not delete requisition wiht ID provided")
		}

		defer db.CloseConnectionToMongo(client)

		updateJobRequisitionInMemory()
		return nil
	}
	return fmt.Errorf("Job Requisition with ID '%v' not found", id)
}

//In Memory: Verify if the JobRequisition id provided is referring to a posted job req.
//Returns a boolean value: True if posted, False if not. Returns an error also in case it does not find the req
func IsJobReqPosted(id int) (bool, error) {
	jr, err := GetJobRequisitionByID(id)
	if err != nil {
		return false, fmt.Errorf("Counld not find Job Requisition '%v'", id)
	}

	return jr.PostingStatus, nil
}

//In Memory: Searches for JobRequisition with Country.
//Return a list of JobRequisition
func GetRequisitionsWithCountry(c int) []JobRequisition {
	ret := make([]JobRequisition, 0)

	for _, v := range GetJobRequisitions() {
		if v.JobReqCountry.ID == c {
			ret = append(ret, *v)
		}
	}

	return ret
}

//Updates the hashmap containing all the JobRequisition to work with them in memory.
//Return the next ID to be added into the Database
func updateJobRequisitionInMemory() int {
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return -1
	}

	filter := bson.D{}
	projection := bson.D{
		{"ID", 1},
		{"Title", 1},
		{"JobDescription", 1},
		{"PostingStatus", 1},
		{"JrCountryId", 1}}
	opts := options.Find().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Requisitions")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	defer db.CloseConnectionToMongo(client)

	biggestId := 1
	jobReqs = make(map[int]*JobRequisition)
	for _, v := range results {
		jr := bsonToJobRequisition(v)

		jr.JobReqCountry, err = GetCountryByID(jr.JrCountryId)

		//Returns all applications of this JobRequisition
		jr.Applicants = GetApplicationsOfJobReq(jr.ID)

		jobReqs[jr.ID] = &jr
		if jr.ID > biggestId {
			biggestId = jr.ID
		}
	}

	return biggestId+1
}

//Receives a bson object to execute the conversion.
//Returns a JobRequisition object.
func bsonToJobRequisition(v bson.D) JobRequisition {
	bsonBytes, _ := bson.Marshal(v)

	var jr JobRequisition
	//deconvert the byarray into a struct object
	bson.Unmarshal(bsonBytes, &jr)

	return jr
}