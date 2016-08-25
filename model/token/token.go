package token

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Model define a user
type Model struct {
	ID           bson.ObjectId `bson:"_id"`
	TokenType    string        `bson:"tokenType"`
	AccessToken  string        `bson:"accressToken"`
	RefreshToken string        `bson:"refreshToken"`
	Origin       string        `bson:"origin"`
	Email        string        `bson:"email"`
	Expiry       int64         `bson:"expiry"`
}

func getCollection(db *mgo.Database) *mgo.Collection {
	return db.C("token")
}

// FindToken find a user
func FindToken(db *mgo.Database, ID string) (result Model, err error) {
	err = getCollection(db).FindId(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&result)
	return
}

// Save save a user
func (data Model) Save(db *mgo.Database) (result Model, err error) {
	var count int
	result = Model{}
	count, err = getCollection(db).Find(bson.M{"email": data.Email}).Limit(1).Count()
	fmt.Println(err)
	if err != nil {
		return
	}
	if count == 0 {
		data.ID = bson.NewObjectId()
		err = getCollection(db).Insert(data)
	} else {
		data.ID = bson.NewObjectId()
		err = getCollection(db).Update(bson.M{"email": data.Email}, data)
	}
	fmt.Println(err)
	result = data
	return
}
