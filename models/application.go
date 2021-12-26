package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webservice/db"
)

type Application struct {
	ID                 int
	CandidateProfileID int
	JobRequisitionID   int
	SalaryExpectation  string
	ApplicationSource  string
	TimeOfExperience   int
}

var (
	applications = make(map[int]*Application)
	nextAppID    = updateApplicantsInMemory()
)

//In Memory: Returns the complete list of Application that has been.
//Returns a hashmap containing the list of Application
func GetApplications() []*Application {
	appArr := make([]*Application,0)
	for _, v := range applications {
		appArr = append(appArr, v)
	}
	return appArr
	//return applications
}

//In Memory: Searches for a specific Application on the hashmap.
//Returns a Application object and an error in case it was not possible to find the record
func GetApplicationByID(id int) (Application, error) {
	if a, found := applications[id]; found {
		return *a, nil
	}

	return Application{}, fmt.Errorf("Application with id '%v' not found", id)
}

//In Memory: Searches for Application that belong to the candidate with id received as parameter on the hashmap.
//Returns a list of Application object.
func GetApplicationsOfCandidate(id int) []Application {
	appArr := make([]Application,0)

	for _, v := range applications {
		if v.CandidateProfileID == id {
			appArr = append(appArr, *v)
		}
	}

	return appArr
}

//In Memory: Searches for Application done to the JobRequisition with id received as parameter on the hashmap.
//Returns a list of Application object.
func GetApplicationsOfJobReq(id int) []Application {
	appArr := make([]Application,0)

	for _, v := range applications {
		if v.JobRequisitionID == id {
			appArr = append(appArr, *v)
		}
	}

	return appArr
}

//In DB: Creates a new Application record to the collection and updates the Application in memory.
//Returns a Application object and an error in case it was not possible to create the record
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

	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return Application{}, fmt.Errorf("Could not establish connection to Database")
	}

	coll := client.Database(db.GetDatabaseName()).Collection("Applications")
	doc := bson.D{
		{"ID", a.ID},
		{"CandidateProfileID", a.CandidateProfileID},
		{"JobRequisitionID", a.JobRequisitionID},
		{"SalaryExpectation", a.SalaryExpectation},
		{"ApplicationSource", a.ApplicationSource},
		{"TimeOfExperience", a.TimeOfExperience}}

	if _, err = coll.InsertOne(context.TODO(), doc); err != nil {
		return Application{}, fmt.Errorf("Could not insert application provided")
	}

	defer db.CloseConnectionToMongo(client)

	updateApplicantsInMemory()
	updateJobRequisitionInMemory()
	updateCandidatesInMemory()
	return a, nil
}

//In DB: Updates a Application record on the collection and updates the Application in memory.
//Returns a Application object and an error in case it was not possible to update the record
func UpdateApplication(a Application) (Application, error) {
	if a.CandidateProfileID == 0 || a.JobRequisitionID == 0{
		return Application{}, fmt.Errorf("Missing Job Requisition ID and/or Candidate ID")
	}

	if _, found := applications[a.ID]; found {
		client, err := db.OpenConnectionToMongo()
		if err != nil {
			return Application{}, fmt.Errorf("Could not establish connection to Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Applications")
		filter := bson.D{{"ID", a.ID}}
		update := bson.D{{"$set", bson.D{
			{"CandidateProfileID", a.CandidateProfileID},
			{"JobRequisitionID", a.JobRequisitionID},
			{"SalaryExpectation", a.SalaryExpectation},
			{"ApplicationSource", a.ApplicationSource},
			{"TimeOfExperience", a.TimeOfExperience}}}}

		if _, err = coll.UpdateOne(context.TODO(), filter, update); err != nil {
			return Application{}, fmt.Errorf("Could not update application provided")
		}

		defer db.CloseConnectionToMongo(client)

		updateApplicantsInMemory()
		updateCandidatesInMemory()
		updateJobRequisitionInMemory()
		return a, nil
	} else {
		return Application{}, fmt.Errorf("Application with ID '%v' not found", a.ID)
	}
}

//In DB: Removes a Application record from the collection and updates the Application in memory.
//Returns error if failed to complete the deletion on the DB
func DeleteApplication(id int) error {
	if _, found := applications[id]; found {
		client, err := db.OpenConnectionToMongo()

		if err != nil {
			return fmt.Errorf("Could not establish connection to Database")
		}

		coll := client.Database(db.GetDatabaseName()).Collection("Applications")
		filter := bson.D{{"ID", id}}

		if _, err = coll.DeleteOne(context.TODO(),filter); err != nil {
			return fmt.Errorf("Could not delete Application with id provided")
		}

		defer db.CloseConnectionToMongo(client)

		updateApplicantsInMemory()
		updateJobRequisitionInMemory()
		updateCandidatesInMemory()
		return nil
	}
	return fmt.Errorf("Application with ID '%v' not found", id)
}

//In DB: Removes all Application from a specified Candidate.
func DeleteApplicationFromCandidate(id int) {
	for _, v := range applications {
		if v.CandidateProfileID == id {
			DeleteApplication(v.ID)
		}
	}
}

//In DB: Removes all Application from a specified JobRequisition.
func DeleteApplicationFromJobReq(id int) {
	for _, v := range applications {
		if v.JobRequisitionID == id {
			DeleteApplication(v.ID)
		}
	}
}

//Updates the hashmap containing all the Application to work with them in memory.
//Return the next ID to be added into the Database
func updateApplicantsInMemory() int {
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return -1
	}

	filter := bson.D{}
	projection := bson.D{
		{"ID", 1},
		{"CandidateProfileID", 1},
		{"JobRequisitionID", 1},
		{"SalaryExpectation", 1},
		{"ApplicationSource", 1},
		{"TimeOfExperience", 1},
	}
	opts := options.Find().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Applications")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	defer db.CloseConnectionToMongo(client)

	biggestId := 1
	applications = make(map[int]*Application)
	for _, v := range results {
		a := bsonToApplicant(v)

		applications[a.ID] = &a
		if a.ID > biggestId {
			biggestId = a.ID
		}
	}

	return biggestId+1
}

//Receives a bson object to execute the conversion.
//Returns a Application object.
func bsonToApplicant(v bson.D) Application {
	bsonBytes, _ := bson.Marshal(v)

	var a Application
	//deconvert the byarray into a struct object
	bson.Unmarshal(bsonBytes, &a)

	return a
}