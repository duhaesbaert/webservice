package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webservice/db"
)

type Tag struct {
	ID		int
	Label 	string
}

var (
	tags 		[]*Tag
	nextTagID 	= updateTagsInMemory()
)

//In Memory: Returns the complete list of tags that has been.
//Returns a hashmap containing the list of tags
func GetTags() []*Tag{
	return tags
}

//In Memory: Finds the corresponding tag on the list of tags.
//Returns the specific tag found, or an error message
func GetTagByLabel(l string) (Tag, error) {
	b,i,e := ExistTagByLabel(l)
	if b {
		return *tags[i], nil
	}

	return Tag{}, e
}

//In DB: Creates a new recod of Tag into the Database.
//Returns the Tag object, and error if not possible to create
func AddTag(t Tag) (Tag, error) {
	//Validation
	if t.ID != 0 {
		return Tag{}, fmt.Errorf("Tag must not contain ID")
	}
	if t.Label == "" {
		return Tag{}, fmt.Errorf("Tag label must not be empty")
	}

	//Test if tag already exists
	b,i,e := ExistTagByLabel(t.Label)
	if b {
		return *tags[i], e
	}

	//Add new tag
	t.ID = nextTagID
	nextTagID++
	tags = append(tags, &t)

	//Validate if able to connect to MongoDB
	client, err := db.OpenConnectionToMongo()
	if err != nil {
		return Tag{}, fmt.Errorf("Could not establish connection to Database")
	}


	coll := client.Database(db.GetDatabaseName()).Collection("Tags")
	doc := bson.D{{"ID", t.ID},{"Label", t.Label}}

	//Insert information into Mongo DB
	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return Tag{}, fmt.Errorf("Could not insert Tag provided into the Database")
	}

	defer db.CloseConnectionToMongo(client)

	updateTagsInMemory()

	return t,nil
}

//Return if a tag exists on the list
func ExistTagByLabel(l string) (bool, int, error) {
	for i,t := range tags {
		if t.Label == l{
			return true,i,nil
		}
	}
	return false, -1, fmt.Errorf("Tag '%v' not found", l)
}

//Updates the hashmap containing all the tags to work with them in memory.
//Return the next ID to be added into the Database
func updateTagsInMemory() int {
	client, err := db.OpenConnectionToMongo()

	if err != nil {
		return -1
	}

	filter := bson.D{}
	projection := bson.D{{"ID", 1}, {"Label" ,1}}
	opts := options.Find().SetProjection(projection)

	coll := client.Database(db.GetDatabaseName()).Collection("Tags")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []bson.D

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	defer db.CloseConnectionToMongo(client)

	biggestId := 1
	for _, v := range results {
		t := bsonToTag(v)
		tags[t.ID] = &t
		if t.ID > biggestId {
			biggestId = t.ID
		}
	}

	return biggestId
}

//Receives a bson object to execute the conversion.
//Returns a Tag object.
func bsonToTag(v bson.D) Tag {
	bsonBytes, _ := bson.Marshal(v)

	var t Tag
	//deconvert the byarray into a struct object
	bson.Unmarshal(bsonBytes, &t)

	return t
}