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

var data Model


func GetInstance()  *mgo.Collection {
	return mongo.GetInstance().clQuizz
}

func Save() (info *mgo.ChangeInfo, err error) {
	return GetInstance().UpsertId( data.id, bson.M{ "$set": data} )
}