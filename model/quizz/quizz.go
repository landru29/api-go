package quizz

import (
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type Quizz struct {
	id           string   `bson:"_id"`
	explaination string   `bson:"explaination"`
	image        string   `bson:"image"`
	level        int      `bson:"level"`
	published    bool     `bson:"published"`
	tags         string   `bson:"tags"`
	text         string   `bson:"text"`
	choices      []string `bson:"choices"`
}


