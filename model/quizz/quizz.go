package quizz

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/landru29/api-go/helpers/mongo"
)

type Choice struct {
	text    string `bson:"text"`
	scoring int    `bson:"scoring"`
}

type Model struct {
	id           string   `bson:"_id"`
	explaination string   `bson:"explaination"`
	image        string   `bson:"image"`
	level        int      `bson:"level"`
	published    bool     `bson:"published"`
	tags         string   `bson:"tags"`
	text         string   `bson:"text"`
	choices      []Choice `bson:"choices"`
	createdAt    string   `bson:"createdAt"`
	updatedAt    string   `bson:"updatedAt"`
}


func (data Model) Save() (result Model, info *mgo.ChangeInfo, err error) {
	if (data.id != "") {
		info, err =  getInstance().UpsertId( data.id, bson.M{ "$set": data} )
	} else {
		err =  getInstance().Insert( data )
	}
	result = data
	return
}

func getInstance()  *mgo.Collection {
	return mongo.GetInstance().Quizz
}

func Find(id string) (result Model, err error) {
	result = Model{}
	err = getInstance().FindId(id).One(&result)
	return
}

func CountAll() (result int, err error) {
	result, err = getInstance().Count()
	return
}

func CountPublished(level int) (result int, err error) {
	if (level <0) {
		result, err = getInstance().Find(bson.M{ "published": true}).Count()
	} else {
		result, err = getInstance().Find(bson.M{ "published": true, "level": level}).Count()
	}
	return
}