package quizz

import (
	"github.com/landru29/api-go/helpers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Choice struct {
	Text    string `bson:"text"`
	Scoring int    `bson:"scoring"`
}

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

func (data Model) Save() (result Model, info *mgo.ChangeInfo, err error) {
	if data.Id != "" {
		info, err = getInstance().UpsertId(data.Id, bson.M{"$set": data})
	} else {
		err = getInstance().Insert(data)
	}
	result = data
	return
}

func getInstance() *mgo.Collection {
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
	if level < 0 {
		result, err = getInstance().Find(bson.M{"published": true}).Count()
	} else {
		result, err = getInstance().Find(bson.M{"published": true, "level": level}).Count()
	}
	return
}
