package quizz

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/landru29/api-go/helpers/mongo"
)

type Model struct {
	id           string   `bson:"_id"`
	explaination string   `bson:"explaination"`
	image        string   `bson:"image"`
	level        int      `bson:"level"`
	published    bool     `bson:"published"`
	tags         string   `bson:"tags"`
	text         string   `bson:"text"`
	choices      []string `bson:"choices"`
}


func Save(data Model) (info *mgo.ChangeInfo, err error) {
	return GetInstance().UpsertId( data.id, bson.M{ "$set": data} )
}

func GetInstance()  *mgo.Collection {
	return mongo.GetInstance().Quizz
}

func Find(id string) Model {
	result := Model{}
	GetInstance().FindId(id).One(&result)
	return result
}