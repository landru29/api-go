package quizz

import (
	"github.com/landru29/api-go/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Choice define a choice in the response
type Choice struct {
	Text    string `bson:"text"`
	Scoring int    `bson:"scoring"`
}

// Model define a quizz question
type Model struct {
	Id           string   `bson:"_id"`
	Explaination string   `bson:"explaination"`
	Image        string   `bson:"image"`
	Level        int      `bson:"level"`
	Published    bool     `bson:"published"`
	Tags         string   `bson:"tags"`
	Text         string   `bson:"text"`
	Choices      []Choice `bson:"choices"`
	CreatedAt    string   `bson:"createdAt"`
	UpdatedAt    string   `bson:"updatedAt"`
}

func getCollection() *mgo.Collection {
	return mongo.GetMongoDatabase().C("quizz")
}

// Save function save a question
func (data Model) Save() (result Model, info *mgo.ChangeInfo, err error) {
	if data.Id != "" {
		info, err = getCollection().UpsertId(data.Id, bson.M{"$set": data})
	} else {
		err = getCollection().Insert(data)
	}
	result = data
	return
}

// Find function find an element
func Find(id string) (result Model, err error) {
	result = Model{}
	err = getCollection().FindId(id).One(&result)
	return
}

// CountAll count all the questions
func CountAll() (result int, err error) {
	result, err = getCollection().Count()
	return
}

// CountPublished count all published questions
func CountPublished(level int) (result int, err error) {
	if level < 0 {
		result, err = getCollection().Find(bson.M{"published": true}).Count()
	} else {
		result, err = getCollection().Find(bson.M{"published": true, "level": level}).Count()
	}
	return
}
