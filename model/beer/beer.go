package beer

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Model define a quizz question
type Model struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	User       string        `bson:"user"`
	CreatedAt  float64       `bson:"createdAt"`
	ModifiedAt float64       `bson:"modifiedAt"`
	Date       float64       `bson:"date"`
	Steps      []Step        `bson:"steps"`
}

// Step defines a brewing step
type Step struct {
	UUID        string       `bson:"_uuid"`
	Ingredients []Ingredient `bson:"ingredients"`
}

// Ingredient defines a brewing ingredient
type Ingredient struct {
	UUID string `bson:"_uuid"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
	return db.C("quizzs")
}

// Save function save a recipe
func (data Model) Save(db *mgo.Database) (result Model, info *mgo.ChangeInfo, err error) {
	if data.ID != "" {
		info, err = getCollection(db).UpsertId(data.ID, bson.M{"$set": data})
	} else {
		err = getCollection(db).Insert(data)
	}
	result = data
	return
}
